package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var (
	errandAction         string
	errandType           string
	errandProducts       []string
	errandNonInteractive bool
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

func toggleErrandsFunc(cmd *cobra.Command, args []string) {
	validateFlags()

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
	switch action {
	case "enable":
	case "disable":
	case "default":
		return true
	}

	return false
}
