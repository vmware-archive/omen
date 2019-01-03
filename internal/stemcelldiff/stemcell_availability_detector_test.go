package stemcelldiff_test

import (
	"io/ioutil"
	"net/http"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/Benjamintf1/Expanded-Unmarshalled-Matchers"
	"github.com/onsi/ginkgo/extensions/table"
	"github.com/pivotal-cloudops/omen/internal/stemcelldiff"
	"github.com/pivotal-cloudops/omen/internal/stemcelldiff/stemcelldifffakes"
	"fmt"
)

const (
	stemcellUpdates = `{
    "stemcell_updates": [
      {
        "stemcell_version": "3468.46",
        "release_id": 106153,
        "products": [
          {
            "product_id": "p-redis-a4de4d5a4bad5"
          },
          {
            "product_id": "p-healthwatch-876a87d6b6c8f7"
          },
          {
            "product_id": "stackdriver-nozzle-develop-89d98f67c76d7a97e6"
          },
          {
            "product_id": "p-event-alerts-78a765d675f55b5c65a6d5"
          }
        ]
      },
      {
        "stemcell_version": "3445.48",
        "release_id": 106150,
        "products": [
          {
            "product_id": "pivotal-mysql-6a5d5d6c6c65b5b5d"
          }
        ]
      },
      {
        "stemcell_version": "3541.30",
        "release_id": 106151,
        "products": [
          {
            "product_id": "cf-97c6b6c7f53d2124"
          }
        ]
      }
    ]
  }`
	oldAssignments = `{
    "products": [
      {
        "guid": "p-bosh-7d6f7d6b6c2d3b2a3",
        "identifier": "p-bosh",
        "label": "BOSH Director",
        "product_version": "2.1-build.214",
        "staged_stemcell_version": "3541.12",
        "deployed_stemcell_version": "3541.12",
        "is_staged_for_deletion": false,
        "available_stemcell_versions": [
          "3541.12"
        ],
        "required_stemcell_version": "3541.12",
        "required_stemcell_os": "ubuntu-trusty"
      },
      {
        "guid": "p-redis-a4de4d5a4bad5",
        "identifier": "p-redis",
        "label": "Redis",
        "product_version": "1.12.0",
        "staged_stemcell_version": "3468.45",
        "deployed_stemcell_version": "3468.45",
        "is_staged_for_deletion": false,
        "available_stemcell_versions": [
          "3468.46"
        ],
        "required_stemcell_version": "3468",
        "required_stemcell_os": "ubuntu-trusty"
      },
      {
        "guid": "p-healthwatch-a4de4d5a4bad5",
        "identifier": "p-healthwatch",
        "label": "PCF Healthwatch",
        "product_version": "1.2.0-build.41",
        "staged_stemcell_version": "3468.45",
        "deployed_stemcell_version": "3468.45",
        "is_staged_for_deletion": false,
        "available_stemcell_versions": [
          "3468.46"
        ],
        "required_stemcell_version": "3468",
        "required_stemcell_os": "ubuntu-trusty"
      },
      {
        "guid": "stackdriver-nozzle-develop-89d98f67c76d7a97e6",
        "identifier": "stackdriver-nozzle-develop",
        "label": "Stackdriver Nozzle (develop)",
        "product_version": "0.0.145",
        "staged_stemcell_version": "3468.45",
        "deployed_stemcell_version": "3468.45",
        "is_staged_for_deletion": false,
        "available_stemcell_versions": [
          "3468.46"
        ],
        "required_stemcell_version": "3468",
        "required_stemcell_os": "ubuntu-trusty"
      },
      {
        "guid": "pivotal-mysql-6a5d5d6c6c65b5b5d",
        "identifier": "pivotal-mysql",
        "label": "MySQL for Pivotal Cloud Foundry v2",
        "product_version": "2.2.5-build.25",
        "staged_stemcell_version": "3445.47",
        "deployed_stemcell_version": "3445.47",
        "is_staged_for_deletion": false,
        "available_stemcell_versions": [
          "3445.48"
        ],
        "required_stemcell_version": "3445.42",
        "required_stemcell_os": "ubuntu-trusty"
      },
      {
        "guid": "p-event-alerts-78a765d675f55b5c65a6d5",
        "identifier": "p-event-alerts",
        "label": "PCF Event Alerts",
        "product_version": "1.1.3",
        "staged_stemcell_version": "3468.45",
        "deployed_stemcell_version": "3468.45",
        "is_staged_for_deletion": false,
        "available_stemcell_versions": [
          "3468.46"
        ],
        "required_stemcell_version": "3468.42",
        "required_stemcell_os": "ubuntu-trusty"
      },
      {
        "guid": "cf-97c6b6c7f53d2124",
        "identifier": "cf",
        "label": "Pivotal Application Service",
        "product_version": "2.1.5",
        "staged_stemcell_version": "3541.29",
        "deployed_stemcell_version": "3541.29",
        "is_staged_for_deletion": false,
        "available_stemcell_versions": [
          "3541.30"
        ],
        "required_stemcell_version": "3541.25",
        "required_stemcell_os": "ubuntu-trusty"
      }
    ],
    "stemcell_library": [
      {
        "version": "3445.48",
        "os": "ubuntu-trusty",
        "infrastructure": "google",
        "hypervisor": "kvm",
        "light": true
      },
      {
        "version": "3468.42",
        "os": "ubuntu-trusty",
        "infrastructure": "google",
        "hypervisor": "kvm",
        "light": true
      },
      {
        "version": "3468.46",
        "os": "ubuntu-trusty",
        "infrastructure": "google",
        "hypervisor": "kvm",
        "light": true
      },
      {
        "version": "3541.12",
        "os": "ubuntu-trusty",
        "infrastructure": "google",
        "hypervisor": "kvm",
        "light": true
      },
      {
        "version": "3541.25",
        "os": "ubuntu-trusty",
        "infrastructure": "google",
        "hypervisor": "kvm",
        "light": true
      },
      {
        "version": "3541.30",
        "os": "ubuntu-trusty",
        "infrastructure": "google",
        "hypervisor": "kvm",
        "light": true
      }
    ]
  }`
	mixedAssignments = `{
    "products": [
      {
        "guid": "p-bosh-7d6f7d6b6c2d3b2a3",
        "identifier": "p-bosh",
        "label": "BOSH Director",
        "product_version": "2.1-build.214",
        "staged_stemcell_version": "3541.12",
        "deployed_stemcell_version": "3541.12",
        "is_staged_for_deletion": false,
        "available_stemcell_versions": [
          "3541.12"
        ],
        "required_stemcell_version": "3541.12",
        "required_stemcell_os": "ubuntu-trusty"
      },
      {
        "guid": "p-redis-a4de4d5a4bad5",
        "identifier": "p-redis",
        "label": "Redis",
        "product_version": "1.12.0",
        "staged_stemcell_version": "3468.45",
        "deployed_stemcell_version": "3468.45",
        "is_staged_for_deletion": false,
        "available_stemcell_versions": [
          "3468.46"
        ],
        "required_stemcell_version": "3468",
        "required_stemcell_os": "ubuntu-trusty"
      },
      {
        "guid": "p-healthwatch-a4de4d5a4bad5",
        "identifier": "p-healthwatch",
        "label": "PCF Healthwatch",
        "product_version": "1.2.0-build.41",
        "staged_stemcell_version": "3468.46",
        "deployed_stemcell_version": "3468.46",
        "is_staged_for_deletion": false,
        "available_stemcell_versions": [
          "3468.46"
        ],
        "required_stemcell_version": "3468",
        "required_stemcell_os": "ubuntu-trusty"
      },
      {
        "guid": "stackdriver-nozzle-develop-89d98f67c76d7a97e6",
        "identifier": "stackdriver-nozzle-develop",
        "label": "Stackdriver Nozzle (develop)",
        "product_version": "0.0.145",
        "staged_stemcell_version": "3468.46",
        "deployed_stemcell_version": "3468.46",
        "is_staged_for_deletion": false,
        "available_stemcell_versions": [
          "3468.46"
        ],
        "required_stemcell_version": "3468",
        "required_stemcell_os": "ubuntu-trusty"
      },
      {
        "guid": "pivotal-mysql-6a5d5d6c6c65b5b5d",
        "identifier": "pivotal-mysql",
        "label": "MySQL for Pivotal Cloud Foundry v2",
        "product_version": "2.2.5-build.25",
        "staged_stemcell_version": "3445.48",
        "deployed_stemcell_version": "3445.48",
        "is_staged_for_deletion": false,
        "available_stemcell_versions": [
          "3445.48"
        ],
        "required_stemcell_version": "3445.42",
        "required_stemcell_os": "ubuntu-trusty"
      },
      {
        "guid": "p-event-alerts-78a765d675f55b5c65a6d5",
        "identifier": "p-event-alerts",
        "label": "PCF Event Alerts",
        "product_version": "1.1.3",
        "staged_stemcell_version": "3468.46",
        "deployed_stemcell_version": "3468.46",
        "is_staged_for_deletion": false,
        "available_stemcell_versions": [
          "3468.46"
        ],
        "required_stemcell_version": "3468.42",
        "required_stemcell_os": "ubuntu-trusty"
      },
      {
        "guid": "cf-97c6b6c7f53d2124",
        "identifier": "cf",
        "label": "Pivotal Application Service",
        "product_version": "2.1.5",
        "staged_stemcell_version": "3541.30",
        "deployed_stemcell_version": "3541.30",
        "is_staged_for_deletion": false,
        "available_stemcell_versions": [
          "3541.30"
        ],
        "required_stemcell_version": "3541.25",
        "required_stemcell_os": "ubuntu-trusty"
      }
    ],
    "stemcell_library": [
      {
        "version": "3445.48",
        "os": "ubuntu-trusty",
        "infrastructure": "google",
        "hypervisor": "kvm",
        "light": true
      },
      {
        "version": "3468.42",
        "os": "ubuntu-trusty",
        "infrastructure": "google",
        "hypervisor": "kvm",
        "light": true
      },
      {
        "version": "3468.46",
        "os": "ubuntu-trusty",
        "infrastructure": "google",
        "hypervisor": "kvm",
        "light": true
      },
      {
        "version": "3541.12",
        "os": "ubuntu-trusty",
        "infrastructure": "google",
        "hypervisor": "kvm",
        "light": true
      },
      {
        "version": "3541.25",
        "os": "ubuntu-trusty",
        "infrastructure": "google",
        "hypervisor": "kvm",
        "light": true
      },
      {
        "version": "3541.30",
        "os": "ubuntu-trusty",
        "infrastructure": "google",
        "hypervisor": "kvm",
        "light": true
      }
    ]
  }`
	newAssignments = `{
    "products": [
      {
        "guid": "p-bosh-7d6f7d6b6c2d3b2a3",
        "identifier": "p-bosh",
        "label": "BOSH Director",
        "product_version": "2.1-build.214",
        "staged_stemcell_version": "3541.12",
        "deployed_stemcell_version": "3541.12",
        "is_staged_for_deletion": false,
        "available_stemcell_versions": [
          "3541.12"
        ],
        "required_stemcell_version": "3541.12",
        "required_stemcell_os": "ubuntu-trusty"
      },
      {
        "guid": "p-redis-a4de4d5a4bad5",
        "identifier": "p-redis",
        "label": "Redis",
        "product_version": "1.12.0",
        "staged_stemcell_version": "3468.46",
        "deployed_stemcell_version": "3468.46",
        "is_staged_for_deletion": false,
        "available_stemcell_versions": [
          "3468.46"
        ],
        "required_stemcell_version": "3468",
        "required_stemcell_os": "ubuntu-trusty"
      },
      {
        "guid": "p-healthwatch-a4de4d5a4bad5",
        "identifier": "p-healthwatch",
        "label": "PCF Healthwatch",
        "product_version": "1.2.0-build.41",
        "staged_stemcell_version": "3468.46",
        "deployed_stemcell_version": "3468.46",
        "is_staged_for_deletion": false,
        "available_stemcell_versions": [
          "3468.46"
        ],
        "required_stemcell_version": "3468",
        "required_stemcell_os": "ubuntu-trusty"
      },
      {
        "guid": "stackdriver-nozzle-develop-89d98f67c76d7a97e6",
        "identifier": "stackdriver-nozzle-develop",
        "label": "Stackdriver Nozzle (develop)",
        "product_version": "0.0.145",
        "staged_stemcell_version": "3468.46",
        "deployed_stemcell_version": "3468.46",
        "is_staged_for_deletion": false,
        "available_stemcell_versions": [
          "3468.46"
        ],
        "required_stemcell_version": "3468",
        "required_stemcell_os": "ubuntu-trusty"
      },
      {
        "guid": "pivotal-mysql-6a5d5d6c6c65b5b5d",
        "identifier": "pivotal-mysql",
        "label": "MySQL for Pivotal Cloud Foundry v2",
        "product_version": "2.2.5-build.25",
        "staged_stemcell_version": "3445.48",
        "deployed_stemcell_version": "3445.48",
        "is_staged_for_deletion": false,
        "available_stemcell_versions": [
          "3445.48"
        ],
        "required_stemcell_version": "3445.42",
        "required_stemcell_os": "ubuntu-trusty"
      },
      {
        "guid": "p-event-alerts-78a765d675f55b5c65a6d5",
        "identifier": "p-event-alerts",
        "label": "PCF Event Alerts",
        "product_version": "1.1.3",
        "staged_stemcell_version": "3468.46",
        "deployed_stemcell_version": "3468.46",
        "is_staged_for_deletion": false,
        "available_stemcell_versions": [
          "3468.46"
        ],
        "required_stemcell_version": "3468.42",
        "required_stemcell_os": "ubuntu-trusty"
      },
      {
        "guid": "cf-97c6b6c7f53d2124",
        "identifier": "cf",
        "label": "Pivotal Application Service",
        "product_version": "2.1.5",
        "staged_stemcell_version": "3541.30",
        "deployed_stemcell_version": "3541.30",
        "is_staged_for_deletion": false,
        "available_stemcell_versions": [
          "3541.30"
        ],
        "required_stemcell_version": "3541.25",
        "required_stemcell_os": "ubuntu-trusty"
      }
    ],
    "stemcell_library": [
      {
        "version": "3445.48",
        "os": "ubuntu-trusty",
        "infrastructure": "google",
        "hypervisor": "kvm",
        "light": true
      },
      {
        "version": "3468.42",
        "os": "ubuntu-trusty",
        "infrastructure": "google",
        "hypervisor": "kvm",
        "light": true
      },
      {
        "version": "3468.46",
        "os": "ubuntu-trusty",
        "infrastructure": "google",
        "hypervisor": "kvm",
        "light": true
      },
      {
        "version": "3541.12",
        "os": "ubuntu-trusty",
        "infrastructure": "google",
        "hypervisor": "kvm",
        "light": true
      },
      {
        "version": "3541.25",
        "os": "ubuntu-trusty",
        "infrastructure": "google",
        "hypervisor": "kvm",
        "light": true
      },
      {
        "version": "3541.30",
        "os": "ubuntu-trusty",
        "infrastructure": "google",
        "hypervisor": "kvm",
        "light": true
      }
    ]
  }`

	fullDiff = `{
  "stemcell_updates": [
    {
      "stemcell_version": "3445.48",
	  "stemcell_os": "ubuntu-trusty",
	  "release_id": 106150,
      "products": [
        "pivotal-mysql-6a5d5d6c6c65b5b5d"
      ]
    },
    {
      "stemcell_version": "3468.46",
	  "stemcell_os": "ubuntu-trusty",
	  "release_id": 106153,
      "products": [
        "p-redis-a4de4d5a4bad5",
        "p-healthwatch-876a87d6b6c8f7",
        "stackdriver-nozzle-develop-89d98f67c76d7a97e6",
        "p-event-alerts-78a765d675f55b5c65a6d5"
      ]
    },
    {
      "stemcell_version": "3541.30",
	  "stemcell_os": "ubuntu-trusty",
	  "release_id": 106151,
      "products": [
        "cf-97c6b6c7f53d2124"
      ]
    }
  ]
}`

	mixedDiff = `{
  "stemcell_updates": [
    {
      "stemcell_version": "3468.46",
	  "stemcell_os": "ubuntu-trusty",
	  "release_id": 106153,
      "products": [
        "p-redis-a4de4d5a4bad5",
      	"p-healthwatch-876a87d6b6c8f7"
      ]
    }
  ]
}`

	emptyDiff = `{"stemcell_updates": []}`
)

var _ = Describe("StemcellAvailabilityDetector", func() {

	table.DescribeTable("Stemcell reporting", func(assignments, report string) {
		client := stemcelldifffakes.FakeHttpClient{}
		rep := stemcelldifffakes.FakeReporter{}
		client.DoStub = func(request *http.Request) (*http.Response, error) {
			var response *http.Response
			if strings.HasSuffix(request.URL.Path, "/stemcell_updates") {
				response = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(strings.NewReader(stemcellUpdates))}
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
		table.Entry("all stemcells are installed", newAssignments, emptyDiff),
		table.Entry("some stemcells are installed", mixedAssignments, mixedDiff),
		table.Entry("no stemcells are installed", oldAssignments, fullDiff))
})