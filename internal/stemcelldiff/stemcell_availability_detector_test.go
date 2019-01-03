package stemcelldiff_test

import (
	"io/ioutil"
	"net/http"
	"strings"

	"fmt"
	. "github.com/Benjamintf1/Expanded-Unmarshalled-Matchers"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cloudops/omen/internal/stemcelldiff"
	"github.com/pivotal-cloudops/omen/internal/stemcelldiff/stemcelldifffakes"
)

const (
	/* Decisions made in this fixture
		- ensured 2 types of stemcells are specified
		- p-healthwatch doesn't need to be updated
		- cf is using xenial (look at the assignments var)
		- p-redis is using trusty (again look at the assignments var)
	 */
	stemcellUpdatesSomeProducts = `{
    "stemcell_updates": [
      {
        "stemcell_version": "3468.46",
        "release_id": 106153,
        "products": [
          {
            "product_id": "p-redis-a4de4d5a4bad5"
          }
        ]
      },
      {
        "stemcell_version": "170.15",
        "release_id": 106151,
        "products": [
          {
            "product_id": "cf-97c6b6c7f53d2124"
          }
        ]
      }
    ]
  }`

	/* Decisions made in this fixture
		- Simplified it by removing all fields that are not used
		- Made sure both xenial and trusty are represented
		- Used a reasonable number of products to cover most cases
	 */
	stemcellAssignments = `{
    "products": [
      {
        "guid": "p-bosh-7d6f7d6b6c2d3b2a3",
        "required_stemcell_os": "ubuntu-trusty"
      },
      {
        "guid": "p-redis-a4de4d5a4bad5",
        "required_stemcell_os": "ubuntu-trusty"
      },
      {
        "guid": "p-healthwatch-a4de4d5a4bad5",
        "required_stemcell_os": "ubuntu-trusty"
      },
      {
        "guid": "cf-97c6b6c7f53d2124",
        "required_stemcell_os": "ubuntu-xenial"
      }
    ]
  }`

	expectedDiff = `{
  "stemcell_updates": [
    {
      "stemcell_version": "3468.46",
	  "stemcell_os": "ubuntu-trusty",
	  "release_id": 106153,
      "products": [
        "p-redis-a4de4d5a4bad5"
      ]
    },
    {
      "stemcell_version": "170.15",
	  "stemcell_os": "ubuntu-xenial",
	  "release_id": 106151,
      "products": [
        "cf-97c6b6c7f53d2124"
      ]
    }
  ]
}`
)

var _ = Describe("StemcellAvailabilityDetector", func() {

	table.DescribeTable("Stemcell reporting", func(stemcells, assignments, report string) {
		client := stemcelldifffakes.FakeHttpClient{}
		rep := stemcelldifffakes.FakeReporter{}
		client.DoStub = func(request *http.Request) (*http.Response, error) {
			var response *http.Response
			if strings.HasSuffix(request.URL.Path, "/stemcell_updates") {
				response = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(strings.NewReader(stemcells))}
			} else if strings.HasSuffix(request.URL.Path, "/stemcell_assignments") {
				response = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(strings.NewReader(assignments))}
			} else {
				return nil, fmt.Errorf("unexpected URL %#v", request.URL)
			}
			return response, nil
		}

		detector := stemcelldiff.NewStemcellUpdateDetector(&client, &rep)
		err := detector.DetectMissingStemcells()
		Expect(err).NotTo(HaveOccurred())
		output := rep.PrintReportArgsForCall(0)
		Expect(output).To(MatchUnorderedJSON(report))
	},
		table.Entry("no products need a stemcell update",
			`{"stemcell_updates": []}`, nil, `{"stemcell_updates": []}`),
		table.Entry("some products need a stemcell update",
			stemcellUpdatesSomeProducts, stemcellAssignments, expectedDiff),
	)
})
