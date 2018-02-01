package cmd

import (
	"fmt"
	"time"

	"github.com/pivotal-cloudops/omen/internal/diff"
	"github.com/pivotal-cloudops/omen/internal/manifest"
	"github.com/pivotal-cloudops/omen/internal/tile"
	"github.com/pivotal-cloudops/omen/internal/user_io"
	"github.com/spf13/cobra"
)

var non_interactive bool

const APPLY_CHANGES_BODY = `{
    "ignore_warnings": true,
    "deploy_products": "all"
}`

type manifestsLoader interface {
	LoadStaged() (manifest.Manifests, error)
	LoadDeployed() (manifest.Manifests, error)
}

var applyChangesCmd = &cobra.Command{
	Use:   "apply-changes",
	Short: "apply any staged changes",
	Long:  "Produces a diff of staged versus deployed changes and then applies those staged changes",
	Run:   applyChangesFunc,
}

func init() {
	applyChangesCmd.Flags().BoolVarP(&non_interactive, "non-interactive", "n", false,
		"Set to true to skip user confirmation for apply change")
}

var applyChangesFunc = func(cmd *cobra.Command, args []string) {
	client := getOpsmanClient()
	tileLoader := tile.NewTilesLoader(client)
	manifestLoader := manifest.NewManifestsLoader(client, tileLoader)

	mDiff := printDiff(manifestLoader)

	if len(mDiff) <= 0 {
		fmt.Println("Warning: Opsman has detected no pending changes")
	}

	if non_interactive == false {
		proceed := user_io.GetConfirmation("Do you wish to continue (y/n)?")

		if proceed == false {
			fmt.Println("Cancelled apply changes")
			return
		}

		fmt.Println("Applying changes")
	}

	resp, err := client.Post("/api/v0/installations", APPLY_CHANGES_BODY, 10*time.Minute)
	if err != nil {
		fmt.Printf("An error occurred applying changes: %v \n", err)
		return
	}

	fmt.Printf("Successfully applied changes: %v \n", resp)
}

func printDiff(ml manifestsLoader) string {
	manifestA, err := ml.LoadDeployed()
	if err != nil {
		printReport("", err)
	}

	manifestB, err := ml.LoadStaged()
	if err != nil {
		printReport("", err)
	}

	d, err := diff.FlatDiff(manifestA, manifestB)
	if err != nil {
		printReport("", err)
	}

	printReport(d, nil)

	return d
}
