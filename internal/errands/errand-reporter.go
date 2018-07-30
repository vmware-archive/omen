package errands

import (
	"fmt"
	"github.com/pivotal-cf/om/api"
)

const divider = "---------------------------------\n"

type ErrandReporter interface {
	Execute(products []string) error
}

type errandReporter struct {
	errandService errandService
	reporter      reporter
}

func NewErrandReporter(es errandService, rp reporter) ErrandReporter {
	return &errandReporter{errandService: es, reporter: rp}
}

func (er *errandReporter) Execute(products []string) error {
	for _, product := range products {
		header := fmt.Sprintf("Listing errands for product: %s\n", product)
		er.reporter.PrintReport(header + divider)
		output, err := er.errandService.List(product)
		if err != nil {
			return err
		}

		if len(output.Errands) == 0 {
			er.reporter.PrintReport("No errands defined")
		} else {
			for _, errand := range output.Errands {
				er.reporter.PrintReport(formatErrand(errand))
			}
		}
		er.reporter.PrintReport(divider)

	}
	return nil
}

func formatErrand(errand api.Errand) string {
	return fmt.Sprintf("Errand name: %s; Post-deploy enabled: %s; Pre-delete enabled: %s\n",
		errand.Name, boolStringFromType(errand.PostDeploy), boolStringFromType(errand.PreDelete))
}

func boolStringFromType(object interface{}) string {
	switch p := object.(type) {
	case string:
		return p
	case bool:
		if object.(bool) {
			return "yes"
		} else {
			return "no"
		}
	default:
		return "default"
	}
}
