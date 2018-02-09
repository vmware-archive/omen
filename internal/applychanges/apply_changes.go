package applychanges

import (
	"fmt"
	"time"

	"github.com/pivotal-cloudops/omen/internal/diff"
	"github.com/pivotal-cloudops/omen/internal/manifest"
	"github.com/pivotal-cloudops/omen/internal/userio"
)

const APPLY_CHANGES_BODY = `{
    "ignore_warnings": true,
    "deploy_products": "all"
}`

type manifestsLoader interface {
	LoadStaged() (manifest.Manifests, error)
	LoadDeployed() (manifest.Manifests, error)
}

type opsmanClient interface {
	Post(endpoint, data string, timeout time.Duration) ([]byte, error)
}

func Execute(ml manifestsLoader, c opsmanClient, nonInteractive bool) {
	mDiff := printDiff(ml)

	if len(mDiff) <= 0 {
		fmt.Println("Warning: Opsman has detected no pending changes")
	}

	if nonInteractive == false {
		proceed := userio.GetConfirmation("Do you wish to continue (y/n)?")

		if proceed == false {
			fmt.Println("Cancelled apply changes")
			return
		}

		fmt.Println("Applying changes")
	}

	resp, err := c.Post("/api/v0/installations", APPLY_CHANGES_BODY, 10*time.Minute)
	if err != nil {
		fmt.Printf("An error occurred applying changes: %v \n", err)
		return
	}

	fmt.Printf("Successfully applied changes: %v \n", resp)
}

func printDiff(ml manifestsLoader) string {
	manifestA, err := ml.LoadDeployed()
	if err != nil {
		userio.PrintReport("", err)
	}

	manifestB, err := ml.LoadStaged()
	if err != nil {
		userio.PrintReport("", err)
	}

	d, err := diff.FlatDiff(manifestA, manifestB)
	if err != nil {
		userio.PrintReport("", err)
	}

	userio.PrintReport(d, nil)

	return d
}
