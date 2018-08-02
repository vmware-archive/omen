package cmd

import (
	"github.com/pivotal-cloudops/omen/internal/tile"
	"github.com/pivotal-cloudops/omen/internal/userio"
	"github.com/spf13/cobra"
)

var listTilesCmd = &cobra.Command{
	Use:   "list-tiles",
	Short: "list all the deployed tiles",
	Long:  "Display a list of all the deployed tiles in a foundation",
	Run:   listTiles,
}

func listTiles(_ *cobra.Command, _ []string) {
	client := setupOpsmanClient()
	tileLoader := tile.NewTilesLoader(client)
	reporter := userio.NewTableReporter()
	tileLister := tile.NewTileLister(tileLoader, reporter)

	err := tileLister.Execute()
	if err != nil {
		rp.Fail(err)
	}
}
