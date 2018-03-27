package cmd

import (
	"errors"

	"fmt"

	"github.com/pivotal-cf/om/api"
	"github.com/pivotal-cloudops/omen/internal/errands"
	"github.com/pivotal-cloudops/omen/internal/tile"
	"github.com/spf13/cobra"
)

var (
	errandAction         string
	errandType           string
	errandProducts       []string
	errandNonInteractive bool

	actionEnable  = "enable"
	actionDisable = "disable"
	actionDefault = "default"
)

var toggleErrandsCmd = &cobra.Command{
	Use:   "toggle-errands",
	Short: "toggle the errand state for products",
	Long:  "Set the errand state for a list of products or all products",
	Run:   toggleErrandsFunc,
}

func init() {
	toggleErrandsCmd.Flags().StringVar(&errandAction, "action", "",
		`Set the toggle errand action. Valid values are: enable, disable, default`)

	toggleErrandsCmd.Flags().StringVar(&errandType, "errand-type", "",
		`Set to true to skip user confirmation for apply change`)

	toggleErrandsCmd.Flags().StringSliceVar(&errandProducts, "products", []string{},
		`(Optional) A comma-delimited list of products for errand updates. When omitted, all products will be affected.`)

	toggleErrandsCmd.Flags().BoolVarP(&errandNonInteractive, "non-interactive", "n", false,
		`Set to true to skip user confirmation for apply change`)
}

var toggleErrandsFunc = func(*cobra.Command, []string) {
	validateFlags()
	c := getOpsmanClient()
	es := api.NewErrandsService(c)
	et := errands.NewErrandToggler(es, rp)

	switch errandAction {
	case actionDefault:
		et = et.Default()
	case actionEnable:
		et = et.Enable()
	}
	if len(errandProducts) > 0 {
		et.Execute(errandProducts)
	} else {
		tl := tile.NewTilesLoader(c)
		toggleAllErrands(tl, et)
	}
}

func toggleAllErrands(tl tile.Loader, et errands.ErrandToggler) {
	deployedProducts, err := tl.LoadDeployed(false)
	if err != nil {
		rp.PrintReport("", errors.New(fmt.Sprintf("Unable to fetch deployed products:\n%#v", err)))
	}
	for _, product := range deployedProducts.Data {
		et.Execute([]string{product.GUID})
	}
}

func validateFlags() {
	if !isErrandActionValid(errandAction) {
		rp.PrintReport("", errors.New("invalid value specified for mandatory flag 'action'"))
	}

	if errandType != "post-deploy" {
		rp.PrintReport("", errors.New("invalid value specified for mandatory flag 'errand-type'"))
	}
}

func isErrandActionValid(action string) bool {
	_, ok := map[string]interface{}{
		actionEnable:  nil,
		actionDisable: nil,
		actionDefault: nil,
	}[action]
	return ok
}
