package errands

import (
	"fmt"
	"strings"

	"github.com/pivotal-cf/om/api"
)

type ErrandReporter interface {
	Execute(products []string) error
}

type errandReporter struct {
	errandService errandService
	reporter      tableReporter
}

func NewErrandReporter(es errandService, rp tableReporter) ErrandReporter {
	return &errandReporter{errandService: es, reporter: rp}
}

func (er *errandReporter) Execute(products []string) error {
	for _, product := range products {
		header := fmt.Sprintf("%s\n%s\n\n", product, strings.Repeat("=", len(product)))

		er.reporter.Write([]byte(header))
		output, err := er.errandService.List(product)
		if err != nil {
			return err
		}

		if len(output.Errands) == 0 {
			er.reporter.Write([]byte("No errands defined\n\n"))
		} else {
			er.reporter.Write([]byte("Name\tPost-deploy\tPre-delete\n"))
			er.reporter.Write([]byte("----\t-----------\t----------\n"))
			for _, errand := range output.Errands {
				er.reporter.Write([]byte(formatErrand(errand)))
			}
			er.reporter.Write([]byte("\n\n"))
		}
		er.reporter.Flush()
	}
	return nil
}

func formatErrand(errand api.Errand) string {
	return fmt.Sprintf("%s\t%s\t%s\n",
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
		return "~"
	}
}
