package tile

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"github.com/pkg/errors"
	"fmt"
)

type Tiles struct {
	Data []*Tile `json:"tiles"`
}

type Tile struct {
	InstallationName string                 `json:"installation_name,omitempty"`
	GUID             string                 `json:"guid,omitempty"`
	Type             string                 `json:"type,omitempty"`
	Networks         map[string]interface{} `json:"networks_and_azs,omitempty"`
	Errands          map[string]interface{} `json:"errands,omitempty"`
	Properties       map[string]interface{} `json:"properties,omitempty"`
	Resources        map[string]interface{} `json:"resources,omitempty"`
}

func (t Tiles) FindBySlug(slug string) (Tile, error) {
	for _, t := range t.Data {
		if t.Type == slug {
			return *t, nil
		}
	}

	return Tile{}, errors.New(fmt.Sprintf("slug %s not found", slug))
}

func (t Tiles) Write(path string) error {
	for _, t := range t.Data {
		dir := path + "/" + t.Type
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
		files := []struct {
			filename string
			member   map[string]interface{}
		}{
			{"networks_and_azs.json", t.Networks},
			{"errands.json", t.Errands},
			{"properties.json", t.Properties},
			{"resources.json", t.Resources},
		}
		for _, f := range files {
			b, err := json.MarshalIndent(f.member, "", "  ")
			if err != nil {
				return err
			}
			err = ioutil.WriteFile(dir+"/"+f.filename, b, 0644)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
