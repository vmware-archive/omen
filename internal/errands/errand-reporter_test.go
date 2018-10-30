package errands_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/om/api"
	. "github.com/pivotal-cloudops/omen/internal/errands"
	"github.com/pivotal-cloudops/omen/internal/errands/errandsfakes"
	"errors"
)

var (
	subject ErrandReporter
	es      *errandsfakes.FakeErrandService
	rp      *errandsfakes.FakeTableReporter
)

var _ = Describe("ErrandReporter", func() {

	BeforeEach(func() {
		es = &errandsfakes.FakeErrandService{}
		rp = &errandsfakes.FakeTableReporter{}
		subject = NewErrandReporter(es, rp)
	})

	It("can be executed", func() {
		Expect(subject.Execute([]string{"blah"})).NotTo(HaveOccurred())
	})

	It("fetches errands for a given product", func() {
		subject.Execute([]string{"sainthood"})

		productId := es.ListStagedProductErrandsArgsForCall(0)
		Expect(productId).To(Equal("sainthood"))
	})

	It("fetches errands for multiple products", func() {
		subject.Execute([]string{"an", "interstellar", "space"})

		Expect(es.ListStagedProductErrandsCallCount()).To(Equal(3))
		Expect(es.ListStagedProductErrandsArgsForCall(0)).To(Equal("an"))
		Expect(es.ListStagedProductErrandsArgsForCall(1)).To(Equal("interstellar"))
		Expect(es.ListStagedProductErrandsArgsForCall(2)).To(Equal("space"))
	})

	It("Propagates the errand service errors", func() {
		es.ListStagedProductErrandsReturns(api.ErrandsListOutput{}, errors.New("oh"))

		err := subject.Execute([]string{"a"})

		Expect(err).To(HaveOccurred())

	})

	It("Reports no errands", func() {
		subject.Execute([]string{"a"})

		Expect(rp.WriteCallCount()).To(Equal(2))
		Expect(string(rp.WriteArgsForCall(0))).To(HavePrefix("a"))
		Expect(string(rp.WriteArgsForCall(1))).To(HavePrefix("No errands defined"))
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

		es.ListStagedProductErrandsReturns(errands, nil)

		subject.Execute([]string{"jam"})

		Expect(rp.WriteCallCount()).To(Equal(6))
		Expect(string(rp.WriteArgsForCall(0))).To(HavePrefix("jam"))
		Expect(string(rp.WriteArgsForCall(1))).To(HavePrefix("Name\tPost-deploy\tPre-delete"))
		Expect(string(rp.WriteArgsForCall(3))).To(HavePrefix("as-i-have-remembered\tyes\tno"))
		Expect(string(rp.WriteArgsForCall(4))).To(HavePrefix("rhubarb\tno\tyes"))
	})

	It("flushes the reporter at the end", func() {
		rp.FlushStub = func() error {
			Expect(rp.WriteCallCount()).To(Equal(2))
			return nil
		}

		subject.Execute([]string{"blah"})

		Expect(rp.FlushCallCount()).To(Equal(1))
	})

	It("reports the undefined state", func() {
		errands := api.ErrandsListOutput{
			Errands: []api.Errand{
				{
					Name: "spatii",
				},
			},
		}
		es.ListStagedProductErrandsReturns(errands, nil)

		subject.Execute([]string{"nuclear"})
		Expect(string(rp.WriteArgsForCall(3))).To(HavePrefix("spatii\t~\t~"))
	})

	It("reports custom state", func() {
		errands := api.ErrandsListOutput{
			Errands: []api.Errand{
				{
					Name:       "cow",
					PostDeploy: "dogma",
					PreDelete:  "lotus",
				},
			},
		}
		es.ListStagedProductErrandsReturns(errands, nil)

		subject.Execute([]string{"jazz"})

		Expect(string(rp.WriteArgsForCall(3))).To(HavePrefix("cow\tdogma\tlotus"))

	})
})
