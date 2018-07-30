package errands_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-cloudops/omen/internal/errands"
	"github.com/pivotal-cloudops/omen/internal/errands/errandsfakes"
	"github.com/pivotal-cf/om/api"
	"errors"
)

var (
	subject ErrandReporter
	es      *errandsfakes.FakeErrandService
	rp      *errandsfakes.FakeReporter
)

var _ = Describe("ErrandReporter", func() {

	BeforeEach(func() {
		es = &errandsfakes.FakeErrandService{}
		rp = &errandsfakes.FakeReporter{}
		subject = NewErrandReporter(es, rp)
	})

	It("can be executed", func() {
		Expect(subject.Execute([]string{"blah"})).NotTo(HaveOccurred())
	})

	It("fetches errands for a given product", func() {
		subject.Execute([]string{"sainthood"})

		productId := es.ListArgsForCall(0)
		Expect(productId).To(Equal("sainthood"))
	})

	It("fetches errands for multiple products", func() {
		subject.Execute([]string{"an", "interstellar", "space"})

		Expect(es.ListCallCount()).To(Equal(3))
		Expect(es.ListArgsForCall(0)).To(Equal("an"))
		Expect(es.ListArgsForCall(1)).To(Equal("interstellar"))
		Expect(es.ListArgsForCall(2)).To(Equal("space"))
	})

	It("Propagates the errand service errors", func() {
		es.ListReturns(api.ErrandsListOutput{}, errors.New("oh"))

		err := subject.Execute([]string{"a"})

		Expect(err).To(HaveOccurred())

	})

	It("Reports no errands", func() {
		subject.Execute([]string{"a"})

		Expect(rp.PrintReportCallCount()).To(Equal(3))
		Expect(rp.PrintReportArgsForCall(0)).To(HavePrefix("Listing errands for product: a"))
		Expect(rp.PrintReportArgsForCall(1)).To(Equal("No errands defined"))
		Expect(rp.PrintReportArgsForCall(2)).To(HavePrefix("----------------"))
	})

	It("reports the errands appropriately", func() {
		errands := api.ErrandsListOutput{
			Errands: []api.Errand{
				{
					Name:       "as-i-have-remembered",
					PreDelete:  false,
					PostDeploy: true,
				},
				{
					Name:       "rhubarb",
					PreDelete:  true,
					PostDeploy: false,
				},
			},
		}

		es.ListReturns(errands, nil)

		subject.Execute([]string{"jam"})

		Expect(rp.PrintReportCallCount()).To(Equal(4))
		Expect(rp.PrintReportArgsForCall(1)).To(HavePrefix("Errand name: as-i-have-remembered; Post-deploy enabled: yes; Pre-delete enabled: no"))
		Expect(rp.PrintReportArgsForCall(2)).To(HavePrefix("Errand name: rhubarb; Post-deploy enabled: no; Pre-delete enabled: yes"))
		Expect(rp.PrintReportArgsForCall(3)).To(HavePrefix("---------------"))
	})

	It("reports the default state", func() {
		errands := api.ErrandsListOutput{
			Errands: []api.Errand{
				{
					Name: "spatii",
				},
			},
		}
		es.ListReturns(errands, nil)

		subject.Execute([]string{"nuclear"})
		Expect(rp.PrintReportArgsForCall(1)).To(HavePrefix("Errand name: spatii; Post-deploy enabled: default; Pre-delete enabled: default"))
	})

	It("reports custom state", func() {
		errands := api.ErrandsListOutput{
			Errands: []api.Errand{
				{
					Name: "cow",
					PostDeploy: "dogma",
					PreDelete: "lotus",
				},
			},
		}
		es.ListReturns(errands, nil)

		subject.Execute([]string{"jazz"})

		Expect(rp.PrintReportArgsForCall(1)).To(HavePrefix("Errand name: cow; Post-deploy enabled: dogma; Pre-delete enabled: lotus"))

	})
})
