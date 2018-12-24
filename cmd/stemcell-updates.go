package cmd

import (
	"github.com/spf13/cobra"
	"github.com/pivotal-cloudops/omen/internal/stemcelldiff"
)

var stemcellUpdatesCmd = &cobra.Command{
	Use:   "stemcell-updates",
	Short: "display available stemcell updates",
	Long:  "List all the stemcell versions that can be updated and the affected products. " +
		"Ops Manager must have a PivNet token installed.",
	Run:   stemcellUpdatesFunc,
}

var stemcellUpdatesFunc = func(*cobra.Command, []string) {
	c := setupOpsmanClient()
	sd := stemcelldiff.NewStemcellUpdateDetector(c, rp)
	err := sd.DetectMissingStemcells()
	if err != nil {
		rp.Fail(err)
	}
}
