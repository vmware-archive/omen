package applychanges

import (
	"fmt"
	"time"

	"github.com/pivotal-cloudops/omen/internal/diff"
	"github.com/pivotal-cloudops/omen/internal/manifest"
	"github.com/pivotal-cloudops/omen/internal/userio"
	"strings"
)

const APPLY_CHANGES_BODY = `{
    "ignore_warnings": true,
    "deploy_products": "%s"
}`

type manifestsLoader interface {
	LoadAll(status manifest.ProductStatus) (manifest.Manifests, error)
	Load(status manifest.ProductStatus, tileGuids []string) (manifest.Manifests, error)
}

type reportPrinter interface {
	PrintReport(report string, err error)
}

type opsmanClient interface {
	Post(endpoint, data string, timeout time.Duration) ([]byte, error)
}

func Execute(ml manifestsLoader, c opsmanClient, tileGuids []string, nonInteractive bool, rp reportPrinter) {
	mDiff, err := printDiff(ml, tileGuids, rp)

	if err != nil {
		fmt.Println(err)
		return
	}

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

	var body string
	if len(tileGuids) == 0 {
		body = fmt.Sprintf(APPLY_CHANGES_BODY, "all")
	} else {
		body = fmt.Sprintf(APPLY_CHANGES_BODY, strings.Join(tileGuids, ","))
	}

	resp, err := c.Post("/api/v0/installations", body, 10*time.Minute)
	if err != nil {
		fmt.Printf("An error occurred applying changes: %v \n", err)
		return
	}

	fmt.Printf("Successfully applied changes: %s \n", string(resp))
}

func printDiff(ml manifestsLoader, tileGuids []string, rp reportPrinter) (string, error) {
	var (
		manifestA manifest.Manifests
		manifestB manifest.Manifests
		err       error
	)

	if len(tileGuids) == 0 {
		manifestA, err = ml.LoadAll(manifest.DEPLOYED)
	} else {
		manifestA, err = ml.Load(manifest.DEPLOYED, tileGuids)
	}

	if err != nil {
		return "", err
	}

	if len(tileGuids) == 0 {
		manifestB, err = ml.LoadAll(manifest.STAGED)
	} else {
		manifestB, err = ml.Load(manifest.STAGED, tileGuids)
	}

	if err != nil {
		return "", err
	}

	d, err := diff.FlatDiff(manifestA, manifestB)

	if err != nil {
		return "", err
	}

	rp.PrintReport(d, nil)

	return d, err
}
