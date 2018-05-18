package cmd

import (
	"fmt"

	"strings"

	"github.com/pivotal-cloudops/omen/internal/applychanges"
	"github.com/pivotal-cloudops/omen/internal/manifest"
	"github.com/pivotal-cloudops/omen/internal/tile"
	"github.com/spf13/cobra"
)

var nonInteractive bool
var products string

var applyChangesCmd = &cobra.Command{
	Use:   "apply-changes",
	Short: "apply any staged changes",
	Long:  "Produces a diff of staged versus deployed changes and then applies those staged changes",
	Run:   applyChangesFunc,
}

func init() {
	applyChangesCmd.Flags().StringVarP(&products, "products", "P", "",
		`Optional flag to set the products to apply changes for (e.g. "product-1" or "product-1,product-2")`)

	applyChangesCmd.Flags().BoolVarP(&nonInteractive, "non-interactive", "n", false,
		"Set to true to skip user confirmation for apply change")
}

var applyChangesFunc = func(cmd *cobra.Command, args []string) {
	c := setupOpsmanClient()
	tl := tile.NewTilesLoader(c)
	ml := manifest.NewManifestsLoader(c, tl)

	var guids []string
	if len(products) == 0 {
		fmt.Println("Applying changes to all products")
	} else {
		fmt.Println("Applying changes to these products:", products)
		products = strings.TrimSpace(products)
		for _, s := range strings.Split(products, ",") {
			guids = append(guids, strings.TrimSpace(s))
		}
	}
	options := applychanges.ApplyChangesOptions{TileSlugs: guids, NonInteractive: nonInteractive}
	err := applychanges.Execute(ml, tl, c, rp, options)
	if err != nil {
		rp.Fail(err)
	}
}
