package stemcells_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cloudops/omen/internal/stemcells"
	"github.com/pivotal-cloudops/omen/internal/stemcells/stemcellsfakes"
	"time"
	"encoding/json"
	"github.com/onsi/ginkgo/extensions/table"
)

var _ = Describe("Stemcells", func() {

	Describe("Report available stemcell/product_id upgrades", func() {
		var client *stemcellsfakes.FakeClient

		Describe("No new stemcells are available", func() {
			It("returns an empty list", func() {
				client = &stemcellsfakes.FakeClient{
					GetStub: func(endpoint string, timeout time.Duration) ([]byte, error) {
						return []byte(noNewStemcellsAvailableResponse), nil
					},
				}

				reporter := stemcells.NewUpdateReporter(client)
				report, err := reporter.Report()

				Expect(err).ToNot(HaveOccurred())
				Expect(client.GetCallCount()).To(Equal(1))
				Expect(report.StemcellUpdates).To(BeEmpty())
			})
		})

		table.DescribeTable(
			"Stemcell update(s)",

			func(response, expectedOutput string) {
				client = &stemcellsfakes.FakeClient{
					GetStub: func(endpoint string, timeout time.Duration) ([]byte, error) {
						return []byte(response), nil
					},
				}

				reporter := stemcells.NewUpdateReporter(client)

				report, err := reporter.Report()

				Expect(err).ToNot(HaveOccurred())
				jsonReport, _ := json.Marshal(report)
				Expect(jsonReport).To(MatchJSON(expectedOutput))

			},

			table.Entry("is available for one product",
				newStemcellForOneProductResponse, newStemcellForOneProductOutput),
			table.Entry("is available for multiple products",
				newStemcellForMultipleProductsResponse, newStemcellForMultipleProductsOutput),
			table.Entry("are available for one product",
				newStemcellsForOneProductResponse, newStemcellsForOneProductOutput),
			table.Entry("are available for multiple products",
				newStemcellsForMultipleProductsResponse, newStemcellsForMultipleProductsOutput))

	})

})
