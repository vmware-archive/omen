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
	errands, err := et.client.List(product)
	if err != nil {
		return err
	}

	report := fmt.Sprintf("Errands for %s\n", product)
	et.rp.PrintReport(report, nil)

	for _, errand := range errands.Errands {
		et.reportErrand(errand)
	}
	return nil
}

func (et errandToggler) reportErrand(errand api.Errand) {
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
	default:
		return
	}

	report := ""
	if errandState != et.action {
		report = fmt.Sprintf("%s\t%s => %s\n", errand.Name, errandState, et.action)
	} else {
		report = fmt.Sprintf("%s\t%s\n", errand.Name, errandState)
	}

	et.rp.PrintReport(report, nil)
}
