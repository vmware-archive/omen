package cmd

import (
	"fmt"

	"github.com/pivotal-cloudops/omen/internal/applychanges"
	"github.com/pivotal-cloudops/omen/internal/manifest"
	"github.com/pivotal-cloudops/omen/internal/tile"
	"github.com/spf13/cobra"
)

var non_interactive bool
var products string

var applyChangesCmd = &cobra.Command{
	Use:   "apply-changes",
	Short: "apply any staged changes",
	Long:  "Produces a diff of staged versus deployed changes and then applies those staged changes",
	Run:   applyChangesFunc,
}

func init() {
	applyChangesCmd.Flags().StringVarP(&products, "products", "P", "all",
		`Optional flag to set the products to apply changes for (e.g. "product-1" or "product-1,product-2")`)

	applyChangesCmd.Flags().BoolVarP(&non_interactive, "non-interactive", "n", false,
		"Set to true to skip user confirmation for apply change")
}

var applyChangesFunc = func(cmd *cobra.Command, args []string) {
	c := getOpsmanClient()
	tl := tile.NewTilesLoader(c)
	ml := manifest.NewManifestsLoader(c, tl)

	fmt.Println("Applying changes to these products:", products)
	applychanges.Execute(ml, c, products, non_interactive)
}
