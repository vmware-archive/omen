package errands_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/om/api"
	"github.com/pivotal-cloudops/omen/internal/errands"
	"github.com/pivotal-cloudops/omen/internal/errands/errandsfakes"
)

var _ = Describe("Toggle Errands", func() {
	var (
		es *errandsfakes.FakeErrandService
		et errands.ErrandToggler
		rp *errandsfakes.FakeReporter
	)

	BeforeEach(func() {
		es = &errandsfakes.FakeErrandService{}
		rp = &errandsfakes.FakeReporter{}
		et = errands.NewErrandToggler(es, rp)
	})

	It("can be constructed", func() {
		Expect(et).ToNot(BeNil())
	})

	Describe("specific products", func() {
		errandServiceResponse := api.ErrandsListOutput{
			Errands: []api.Errand{
				{
					Name:       "errand1",
					PostDeploy: true,
				},
				{
					Name:       "errand2",
					PostDeploy: false,
				},
				{
					Name:       "errand3",
					PostDeploy: "when-changed",
				},
				{
					Name:       "errand4",
					PostDeploy: "default",
				},
				{
					Name:      "errand5",
					PreDelete: true,
				},
				{
					Name:      "errand6",
					PreDelete: false,
				},
			},
		}

		It("retrieves errand state for only specified products", func() {
			et.Execute([]string{"PEANUTS-and-butter", "almond-butter"})

			Expect(es.ListCallCount()).To(Equal(2))
			product1Id := es.ListArgsForCall(0)
			Expect(product1Id).To(Equal("PEANUTS-and-butter"))

			product2Id := es.ListArgsForCall(1)
			Expect(product2Id).To(Equal("almond-butter"))
		})

		It("fails with error if an invalid product is specified", func() {
			es.ListReturns(api.ErrandsListOutput{}, errors.New("product not found"))

			err := et.Execute([]string{"PEANUTS-and-butter"})

			Expect(err).To(HaveOccurred())
		})

		Describe("disable", func() {
			It("outputs current and desired state for post-deploy errands of affected products", func() {
				es.ListReturns(errandServiceResponse, nil)

				et.Disable().Execute([]string{"PEANUTS-and-butter"})
				output := ""
				for i := 0; i < rp.PrintReportCallCount(); i++ {
					text, err := rp.PrintReportArgsForCall(i)
					Expect(err).NotTo(HaveOccurred())
					output += text
				}

				Expect(output).To(ContainSubstring("Errands for PEANUTS-and-butter"))
				Expect(output).To(MatchRegexp("errand1\\s+enabled => disabled"))
				Expect(output).To(MatchRegexp("errand2\\s+disabled\\n"))
				Expect(output).To(MatchRegexp("errand3\\s+when-changed => disabled"))
				Expect(output).To(MatchRegexp("errand4\\s+default => disabled"))
				Expect(output).NotTo(ContainSubstring("errand5"))
			})

			It("only enables post-deploy errands not at desired state", func() {

			})
		})

		Describe("enable", func() {
			It("outputs current and desired state for post-deploy errands of affected products", func() {
				es.ListReturns(errandServiceResponse, nil)

				et.Enable().Execute([]string{"PEANUTS-and-butter"})
				output := ""
				for i := 0; i < rp.PrintReportCallCount(); i++ {
					text, err := rp.PrintReportArgsForCall(i)
					Expect(err).NotTo(HaveOccurred())
					output += text
				}

				Expect(output).To(ContainSubstring("Errands for PEANUTS-and-butter"))
				Expect(output).To(MatchRegexp("errand1\\s+enabled\\n"))
				Expect(output).To(MatchRegexp("errand2\\s+disabled => enabled"))
				Expect(output).To(MatchRegexp("errand3\\s+when-changed => enabled"))
				Expect(output).To(MatchRegexp("errand4\\s+default => enabled"))
				Expect(output).NotTo(ContainSubstring("errand5"))
			})

			It("only enables post-deploy errands not at desired state", func() {

			})
		})

		Describe("default", func() {
			It("outputs current and desired state for post-deploy errands of affected products", func() {
				es.ListReturns(errandServiceResponse, nil)

				et.Default().Execute([]string{"PEANUTS-and-butter"})
				output := ""
				for i := 0; i < rp.PrintReportCallCount(); i++ {
					text, err := rp.PrintReportArgsForCall(i)
					Expect(err).NotTo(HaveOccurred())
					output += text
				}

				Expect(output).To(ContainSubstring("Errands for PEANUTS-and-butter"))
				Expect(output).To(MatchRegexp("errand1\\s+enabled => default"))
				Expect(output).To(MatchRegexp("errand2\\s+disabled => default"))
				Expect(output).To(MatchRegexp("errand3\\s+when-changed => default"))
				Expect(output).To(MatchRegexp("errand4\\s+default"))
				Expect(output).NotTo(ContainSubstring("errand5"))
			})

			It("only enables post-deploy errands not at desired state", func() {

			})
		})
	})

	Describe("all errands", func() {
		It("retrieves errand state for all products", func() {

		})
	})
})
