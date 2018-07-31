package cmd

import (
	"errors"
	"fmt"

	"github.com/pivotal-cf/om/api"
	"github.com/pivotal-cloudops/omen/internal/errands"
	"github.com/pivotal-cloudops/omen/internal/tile"
	"github.com/spf13/cobra"
)

var errandProducts []string

var errandsCmd = &cobra.Command{
	Use:   "errands",
	Short: "list the errands and their state",
	Long:  "Display a list of errands, optionally filtered by the product name",
	Run:   errandsFunc,
}

func init() {
	errandsCmd.Flags().StringSliceVar(&errandProducts, "products", []string{},
		`(Optional) A comma-delimited list of products for errand updates. When omitted, all products will be affected.`)
}

var errandsFunc = func(*cobra.Command, []string) {
	c := setupOpsmanClient()
	es := api.NewErrandsService(c)
	et := errands.NewErrandReporter(es, tr)
	tl := tile.NewTilesLoader(c)

	if len(errandProducts) > 0 {
		guids, err := mapGuid(tl, errandProducts)

		if err != nil {
			rp.Fail(err)
		}

		err = et.Execute(guids)

		if err != nil {
			rp.Fail(err)
		}
	} else {
		reportAllErrands(tl, et)
	}
}

func mapGuid(tl tile.Loader, productNames []string) ([]string, error) {
	var guids []string
	deployedProducts, err := tl.LoadDeployed(false)
	for _, product := range productNames {
		_tile, err := deployedProducts.FindBySlug(product)
		if err != nil {
			return nil, err
		}
		guids = append(guids, _tile.GUID)
	}
	return guids, err
}

func reportAllErrands(tl tile.Loader, er errands.ErrandReporter) {
	deployedProducts, err := tl.LoadDeployed(false)
	if err != nil {
		rp.Fail(errors.New(fmt.Sprintf("Unable to fetch deployed products:\n%#v", err)))
	}
	for _, product := range deployedProducts.Data {
		err := er.Execute([]string{product.GUID})
		if err != nil {
			rp.Fail(err)
		}
	}
}
