package manifest

import (
	"encoding/json"
	"fmt"

	"time"

	"github.com/pivotal-cloudops/omen/internal/tile"
)

type omClient interface {
	Get(endpoint string, timeout time.Duration) ([]byte, error)
}

type tilesLoader interface {
	LoadStaged(bool) (tile.Tiles, error)
	LoadDeployed(bool) (tile.Tiles, error)
}

type Manifest struct {
	Name           string      `json:"name,omitempty"`
	Releases       interface{} `json:"releases,omitempty"`
	Stemcells      interface{} `json:"stemcells,omitempty"`
	InstanceGroups interface{} `json:"instance_groups,omitempty"`
	Update         interface{} `json:"update,omitempty"`
	Variables      interface{} `json:"variables,omitempty"`
}

type Manifests struct {
	Data        []Manifest  `json:"manifests"`
	CloudConfig interface{} `json:"cloud_config"`
}

type Loader struct {
	client omClient
	tl     tilesLoader
}

func NewManifestsLoader(omClient omClient, tl tilesLoader) Loader {
	return Loader{client: omClient, tl: tl}
}

func (l Loader) LoadStaged() (Manifests, error) {
	return l.load("staged")
}

func (l Loader) LoadDeployed() (Manifests, error) {
	return l.load("deployed")
}

func (l Loader) load(status string) (Manifests, error) {
	var (
		tiles tile.Tiles
		err   error
	)

	if status == "deployed" {
		tiles, err = l.tl.LoadDeployed(false)
	} else {
		tiles, err = l.tl.LoadStaged(false)
	}

	if err != nil {
		return Manifests{}, err
	}

	manifests, err := l.loadManifests(tiles, status)
	if err != nil {
		return Manifests{}, err
	}

	cloudConfig, err := l.loadCloudConfig(status)
	if err != nil {
		return Manifests{}, err
	}

	return Manifests{manifests, cloudConfig}, nil
}

func (l Loader) loadCloudConfig(status string) (interface{}, error) {
	response, err := l.client.Get(fmt.Sprintf("/api/v0/%s/cloud_config", status), 10*time.Minute)
	if err != nil {
		return nil, err
	}
	var cloudConfig map[string]interface{}
	err = json.Unmarshal(response, &cloudConfig)
	return cloudConfig["cloud_config"], err
}

func getEndpoint(tile tile.Tile, status string) string {
	if tile.Type == "p-bosh" {
		return fmt.Sprintf("/api/v0/%s/director/manifest", status)
	}
	return fmt.Sprintf("/api/v0/%s/products/%s/manifest", status, tile.GUID)
}

func (l Loader) loadManifests(tiles tile.Tiles, status string) ([]Manifest, error) {
	var manifests []Manifest
	for _, t := range tiles.Data {
		data, err := l.client.Get(getEndpoint(*t, status), 10*time.Minute)
		if err != nil {
			return nil, err
		}

		var (
			temp map[string]Manifest
			m    Manifest
		)

		m = Manifest{}

		if status == "deployed" {
			err = json.Unmarshal(data, &m)
			if err != nil {
				return nil, err
			}
		} else {
			temp = make(map[string]Manifest)
			err = json.Unmarshal(data, &temp)
			if err != nil {
				return nil, err
			}
			m = temp["manifest"]
		}

		manifests = append(manifests, m)
	}
	return manifests, nil
}
