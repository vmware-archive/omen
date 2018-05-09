package stemcells_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cloudops/omen/internal/stemcells"
	"github.com/pivotal-cloudops/omen/internal/stemcells/stemcellsfakes"
)

var _ = Describe("Stemcells", func() {

	Describe("Report available stemcell/product_id upgrades", func() {

		It("calls OpsMan API", func() {

			client := &stemcellsfakes.FakeClient{}

			reporter := stemcells.NewUpdateReporter(client)

			report, err := reporter.Report()

			Expect(err).ToNot(HaveOccurred())
			Expect(report).ToNot(BeEmpty())
			//Expect(report).To(MatchJSON())

		})

	})

})
