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

// go:generate counterfeiter . errandService
type errandService interface {
	List(productID string) (api.ErrandsListOutput, error)
	SetState(productID string, errandName string, postDeployState interface{}, preDeleteState interface{}) error
}

// go:generate counterfeiter . reporter
type reporter interface {
	PrintReport(report string, err error)
}

type ErrandToggler interface {
	Execute(products []string) error
	Disable() ErrandToggler
	Default() ErrandToggler
	Enable() ErrandToggler
}

type errandToggler struct {
	client      errandService
	interactive bool
	action      string
	rp          reporter
}

func NewErrandToggler(client errandService, rp reporter) ErrandToggler {
	return errandToggler{
		client:      client,
		interactive: false,
		action:      defaultTogglerAction,
		rp:          rp,
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
		client:      et.client,
		interactive: et.interactive,
		action:      state,
		rp:          et.rp,
	}
}

func (et errandToggler) Execute(products []string) error {
	for _, p := range products {
		err := et.updateErrandsForProduct(p)
		if err != nil {
			return err
		}
	}
	return nil
}

func (et errandToggler) updateErrandsForProduct(product string) error {
	errandsList, err := et.client.List(product)
	if err != nil {
		return err
	}

	report := fmt.Sprintf("Errands for %s\n", product)
	et.rp.PrintReport(report, nil)

	transitioningErrands := et.getTransitioningErrands(errandsList.Errands)
	for _, errand := range errandsList.Errands {
		if errand.PostDeploy == nil {
			continue
		}

		if _, ok := transitioningErrands[errand]; ok {
			et.client.SetState(product, errand.Name, et.getErrandStateFlag(), errand.PreDelete)
			report = fmt.Sprintf("%s\t%s => %s\n", errand.Name, getErrandStateString(errand), et.action)
		} else {
			report = fmt.Sprintf("%s\t%s\n", errand.Name, getErrandStateString(errand))
		}
		et.rp.PrintReport(report, nil)
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
