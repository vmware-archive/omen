package applychanges

import (
	"fmt"
	"time"

	"strings"

	"github.com/pivotal-cloudops/omen/internal/common"
	"github.com/pivotal-cloudops/omen/internal/diff"
	"github.com/pivotal-cloudops/omen/internal/manifest"
	"github.com/pivotal-cloudops/omen/internal/tile"
	"github.com/pivotal-cloudops/omen/internal/userio"
)

const applyChangesBody = `{
    "ignore_warnings": true,
    "deploy_products": "%s"
}`

type manifestsLoader interface {
	LoadAll(status common.ProductStatus) (manifest.Manifests, error)
	Load(status common.ProductStatus, tileGuids []string) (manifest.Manifests, error)
}

type tilesLoader interface {
	LoadStaged(bool) (tile.Tiles, error)
	LoadDeployed(bool) (tile.Tiles, error)
}

//go:generate counterfeiter . reportPrinter
type reportPrinter interface {
	PrintReport(string)
}

type opsmanClient interface {
	Post(endpoint, data string, timeout time.Duration) ([]byte, error)
}

func Execute(ml manifestsLoader, tl tilesLoader, c opsmanClient, tileSlugs []string, nonInteractive bool, rp reportPrinter) error {
	tileGuids, err := slugsToGuids(tileSlugs, tl)
	if err != nil {
		return err
	}

	mDiff, err := printDiff(ml, tileGuids, rp)

	if err != nil {
		fmt.Println(err)
		return err
	}

	if len(mDiff) <= 0 {
		fmt.Println("Warning: Opsman has detected no pending changes")
	}

	if nonInteractive == false {
		proceed := userio.GetConfirmation("Do you wish to continue (y/n)?")

		if proceed == false {
			fmt.Println("Cancelled apply changes")
			return nil
		}

		fmt.Println("Applying changes")
	}

	var body string
	if len(tileGuids) == 0 {
		body = fmt.Sprintf(applyChangesBody, "all")
	} else {
		body = fmt.Sprintf(applyChangesBody, strings.Join(tileGuids, ","))
	}

	resp, err := c.Post("/api/v0/installations", body, 10*time.Minute)
	if err != nil {
		fmt.Printf("An error occurred applying changes: %v \n", err)
		return err
	}

	fmt.Printf("Successfully applied changes: %s \n", string(resp))
	return nil
}

func slugsToGuids(slugs []string, tl tilesLoader) ([]string, error) {
	if len(slugs) == 0 {
		return []string{}, nil
	}

	tiles, err := tl.LoadStaged(false)
	if err != nil {
		return nil, err
	}
	var resp []string
	for _, s := range slugs {
		t, err := tiles.FindBySlug(s)
		if err != nil {
			return nil, err
		}
		resp = append(resp, t.GUID)
	}
	return resp, nil
}

func printDiff(ml manifestsLoader, tileGuids []string, rp reportPrinter) (string, error) {
	var (
		manifestA manifest.Manifests
		manifestB manifest.Manifests
		err       error
	)

	if len(tileGuids) == 0 {
		manifestA, err = ml.LoadAll(common.DEPLOYED)
	} else {
		manifestA, err = ml.Load(common.DEPLOYED, tileGuids)
	}

	if err != nil {
		return "", err
	}

	if len(tileGuids) == 0 {
		manifestB, err = ml.LoadAll(common.STAGED)
	} else {
		manifestB, err = ml.Load(common.STAGED, tileGuids)
	}

	if err != nil {
		return "", err
	}

	d, err := diff.FlatDiff(manifestA, manifestB)

	if err != nil {
		return "", err
	}

	rp.PrintReport(d)

	return d, err
}
