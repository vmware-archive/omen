package errands

import (
	"fmt"

	"github.com/pivotal-cf/om/api"
)

const (
	errandStateEnabled  = "enabled"
	errandStateDisabled = "disabled"
	errandStateDefault  = "default"

	defaultTogglerAction = errandStateDisabled
)

type ErrandToggler interface {
	Execute(products []string) error
	Disable() ErrandToggler
	Default() ErrandToggler
	Enable() ErrandToggler
}

type errandToggler struct {
	errandService errandService
	interactive   bool
	action        string
	reporter      reporter
}

func NewErrandToggler(es errandService, rp reporter) ErrandToggler {
	return errandToggler{
		errandService: es,
		interactive:   false,
		action:        defaultTogglerAction,
		reporter:      rp,
	}
}

func (et errandToggler) Disable() ErrandToggler {
	return et.errandTogglerWithState(errandStateDisabled)
}

func (et errandToggler) Enable() ErrandToggler {
	return et.errandTogglerWithState(errandStateEnabled)
}

func (et errandToggler) Default() ErrandToggler {
	return et.errandTogglerWithState(errandStateDefault)
}

func (et errandToggler) errandTogglerWithState(state string) ErrandToggler {
	return errandToggler{
		errandService: et.errandService,
		interactive:   et.interactive,
		action:        state,
		reporter:      et.reporter,
	}
}

func (et errandToggler) Execute(products []string) error {
	for _, p := range products {
		et.reporter.PrintReport("\n")
		err := et.updateErrandsForProduct(p)
		if err != nil {
			return err
		}
	}
	return nil
}

func (et errandToggler) updateErrandsForProduct(product string) error {
	errandsList, err := et.errandService.ListStagedProductErrands(product)
	if err != nil {
		return err
	}

	report := fmt.Sprintf("Errands for %s\n", product)
	et.reporter.PrintReport(report)
	et.reporter.PrintReport("---------------------------------")

	transitioningErrands := et.getTransitioningErrands(errandsList.Errands)
	for _, errand := range errandsList.Errands {
		if errand.PostDeploy == nil {
			continue
		}

		if _, ok := transitioningErrands[errand]; ok {
			report = fmt.Sprintf("%s\t%s => %s\n", errand.Name, getErrandStateString(errand), et.action)
		} else {
			report = fmt.Sprintf("%s\t%s\n", errand.Name, getErrandStateString(errand))
		}
		et.reporter.PrintReport(report)
	}

	et.reporter.PrintReport("")
	et.reporter.PrintReport("---------------------------------")

	for errand := range transitioningErrands {
		report := fmt.Sprintf("updating %s to %s", errand.Name, et.action)
		et.reporter.PrintReport(report)

		err := et.errandService.UpdateStagedProductErrands(product, errand.Name, et.getErrandStateFlag(), errand.PreDelete)
		if err != nil {
			return err
		}
	}
	return nil
}

func (et errandToggler) getTransitioningErrands(errands []api.Errand) map[api.Errand]interface{} {
	result := make(map[api.Errand]interface{})

	for _, errand := range errands {
		errandState := getErrandStateString(errand)
		if errandState == "" {
			continue
		}
		if errandState != et.action {
			result[errand] = nil
		}
	}

	return result
}

func (et errandToggler) getErrandStateFlag() interface{} {
	switch et.action {
	case errandStateDefault:
		return "default"
	case errandStateDisabled:
		return false
	case errandStateEnabled:
		return true
	default:
		return nil
	}
}

func getErrandStateString(errand api.Errand) string {
	errandState := ""

	switch errand.PostDeploy.(type) {
	case bool:
		if errand.PostDeploy == true {
			errandState = errandStateEnabled
		} else {
			errandState = errandStateDisabled
		}
	case string:
		errandState = errand.PostDeploy.(string)
	}

	return errandState
}
