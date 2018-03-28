package errands_test

import (
	"errors"

	"fmt"

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

		Context("valid product output", func() {
			errandServiceResponse := api.ErrandsListOutput{
				Errands: []api.Errand{
					{
						Name:       "errand1",
						PostDeploy: true,
						PreDelete:  true,
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
						PostDeploy: false,
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
			BeforeEach(func() {
				es.ListReturns(errandServiceResponse, nil)
			})

			Describe("disable", func() {
				It("outputs current and desired state for post-deploy errands of affected products", func() {
					et.Disable().Execute([]string{"PEANUTS-and-butter"})
					output := ""
					for i := 0; i < rp.PrintReportCallCount(); i++ {
						output += rp.PrintReportArgsForCall(i)
					}

					Expect(output).To(ContainSubstring("Errands for PEANUTS-and-butter"))
					Expect(output).To(MatchRegexp("errand1\\s+enabled => disabled"))
					Expect(output).To(MatchRegexp("errand2\\s+disabled\\n"))
					Expect(output).To(MatchRegexp("errand3\\s+when-changed => disabled"))
					Expect(output).To(MatchRegexp("errand4\\s+disabled"))
					Expect(output).NotTo(ContainSubstring("errand5"))
				})

				It("only disables post-deploy errands not at desired state", func() {
					et.Disable().Execute([]string{"PEANUTS-and-butter"})
					output := ""
					for i := 0; i < rp.PrintReportCallCount(); i++ {
						output += rp.PrintReportArgsForCall(i)
					}

					Expect(es.SetStateCallCount()).To(Equal(2))

					for i := range []int{0, 1} {
						productName, errandName, postDeployState, preDeleteState := es.SetStateArgsForCall(i)
						Expect(productName).To(Equal("PEANUTS-and-butter"))
						Expect(postDeployState).To(BeFalse())

						Expect(errandName).ToNot(Equal("errand2"))
						Expect(errandName).ToNot(Equal("errand5"))

						if errandName == "errand1" {
							Expect(preDeleteState).To(BeTrue())
						} else {
							Expect(preDeleteState).To(BeNil())
						}

						expectedOut := fmt.Sprintf("updating %s to disabled", errandName)
						Expect(output).To(ContainSubstring(expectedOut))
					}
				})
			})

			Describe("enable", func() {
				It("outputs current and desired state for post-deploy errands of affected products", func() {
					et.Enable().Execute([]string{"PEANUTS-and-butter"})
					output := ""
					for i := 0; i < rp.PrintReportCallCount(); i++ {
						output += rp.PrintReportArgsForCall(i)
					}

					Expect(output).To(ContainSubstring("Errands for PEANUTS-and-butter"))
					Expect(output).To(MatchRegexp("errand1\\s+enabled\\n"))
					Expect(output).To(MatchRegexp("errand2\\s+disabled => enabled"))
					Expect(output).To(MatchRegexp("errand3\\s+when-changed => enabled"))
					Expect(output).To(MatchRegexp("errand4\\s+disabled => enabled"))
					Expect(output).NotTo(ContainSubstring("errand5"))
				})

				It("only enables post-deploy errands not at desired state", func() {
					et.Enable().Execute([]string{"PEANUTS-and-butter"})
					output := ""
					for i := 0; i < rp.PrintReportCallCount(); i++ {
						output += rp.PrintReportArgsForCall(i)
					}

					Expect(es.SetStateCallCount()).To(Equal(3))

					for i := range []int{0, 1, 2} {
						productName, errandName, postDeployState, preDeleteState := es.SetStateArgsForCall(i)
						Expect(productName).To(Equal("PEANUTS-and-butter"))
						Expect(postDeployState).To(BeTrue())

						Expect(errandName).ToNot(Equal("errand1"))
						Expect(errandName).ToNot(Equal("errand5"))

						Expect(preDeleteState).To(BeNil())

						expectedOut := fmt.Sprintf("updating %s to enabled", errandName)
						Expect(output).To(ContainSubstring(expectedOut))
					}
				})
			})

			Describe("default", func() {
				It("outputs current and desired state for post-deploy errands of affected products", func() {
					et.Default().Execute([]string{"PEANUTS-and-butter"})
					output := ""
					for i := 0; i < rp.PrintReportCallCount(); i++ {
						output += rp.PrintReportArgsForCall(i)
					}

					Expect(output).To(ContainSubstring("Errands for PEANUTS-and-butter"))
					Expect(output).To(MatchRegexp("errand1\\s+enabled => default"))
					Expect(output).To(MatchRegexp("errand2\\s+disabled => default"))
					Expect(output).To(MatchRegexp("errand3\\s+when-changed => default"))
					Expect(output).To(MatchRegexp("errand4\\s+disabled => default"))
					Expect(output).NotTo(ContainSubstring("errand5"))
				})

				It("only enables post-deploy errands not at desired state", func() {
					et.Default().Execute([]string{"PEANUTS-and-butter"})
					output := ""
					for i := 0; i < rp.PrintReportCallCount(); i++ {
						output += rp.PrintReportArgsForCall(i)
					}

					Expect(es.SetStateCallCount()).To(Equal(4))

					for i := range []int{0, 1, 2, 3} {
						productName, errandName, postDeployState, preDeleteState := es.SetStateArgsForCall(i)
						Expect(productName).To(Equal("PEANUTS-and-butter"))
						Expect(postDeployState).To(Equal("default"))

						Expect(errandName).To(Equal(fmt.Sprintf("errand%d", i+1)))

						if errandName == "errand1" {
							Expect(preDeleteState).To(BeTrue())
						} else {
							Expect(preDeleteState).To(BeNil())
						}

						expectedOut := fmt.Sprintf("updating %s to default", errandName)
						Expect(output).To(ContainSubstring(expectedOut))
					}
				})

				It("propagates the error from the errand service", func() {
					es.SetStateReturns(errors.New("blah"))

					err := et.Default().Execute([]string{"PEANUTS-and-butter"})
					Expect(err).To(HaveOccurred())
				})
			})
		})
	})
})
