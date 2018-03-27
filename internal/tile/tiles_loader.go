package tile

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/pivotal-cloudops/omen/internal/common"
)

type omClient interface {
	Get(endpoint string, timeout time.Duration) ([]byte, error)
}

type Loader struct {
	client omClient
}

func NewTilesLoader(omClient omClient) Loader {
	return Loader{client: omClient}
}

func (l Loader) LoadStaged(fetchTileMetadata bool) (Tiles, error) {
	return l.load(fetchTileMetadata, common.STAGED)
}

func (l Loader) LoadDeployed(fetchTileMetadata bool) (Tiles, error) {
	return l.load(fetchTileMetadata, common.DEPLOYED)
}

func (l Loader) load(fetchTileMetadata bool, status common.ProductStatus) (Tiles, error) {
	b, err := l.client.Get(fmt.Sprintf("/api/v0/%s/products", status), 10*time.Minute)

	if err != nil {
		return Tiles{}, err
	}

	var data []*Tile
	err = json.Unmarshal(b, &data)
	if err != nil {
		return Tiles{}, err
	}

	if fetchTileMetadata {
		for _, d := range data {
			err = l.loadTileMetadata(d, status)
			if err != nil {
				return Tiles{}, err
			}
		}
	}

	return Tiles{data}, nil
}

func (l Loader) loadTileMetadata(t *Tile, status common.ProductStatus) error {
	urlsToPointer := []struct {
		url     string
		pointer *map[string]interface{}
	}{
		{"/api/v0/%s/products/%s/networks_and_azs", &t.Networks},
		{"/api/v0/%s/products/%s/errands", &t.Errands},
		{"/api/v0/%s/products/%s/resources", &t.Resources},
		{"/api/v0/%s/products/%s/properties", &t.Properties},
	}

	for _, up := range urlsToPointer {
		data, err := l.client.Get(fmt.Sprintf(up.url, status, t.GUID), 10*time.Minute)
		if err != nil {
			return err
		}
		err = json.Unmarshal(data, up.pointer)
		if err != nil {
			return err
		}
	}

	return nil
}
