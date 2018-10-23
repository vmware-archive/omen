package tile

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

type Tiles struct {
	Data []*Tile `json:"tiles"`
}

type Tile struct {
	InstallationName string                 `json:"installation_name,omitempty"`
	GUID             string                 `json:"guid,omitempty"`
	Type             string                 `json:"type,omitempty"`
	ProductVersion   string                 `json:"product_version,omitempty"`
	Networks         map[string]interface{} `json:"networks_and_azs,omitempty"`
	Errands          map[string]interface{} `json:"errands,omitempty"`
	Properties       map[string]interface{} `json:"properties,omitempty"`
	Resources        map[string]interface{} `json:"resources,omitempty"`
}

type tileFinders struct {
	main  func(string) (Tile, error)
	other func(string) (Tile, error)
}

func (t Tiles) FindBySlug(slug string) (Tile, error) {
	for _, t := range t.Data {
		if t.Type == slug {
			return *t, nil
		}
	}

	return Tile{}, errors.New(fmt.Sprintf("product %s not found", slug))
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

func (t Tiles) FindBySlugsOrGUIDs(products []string) ([]*Tile, error) {
	if len(products) == 0 {
		return []*Tile{}, nil
	}

	tiles := make([]*Tile, 0)
	finders, err := t.getFindersForProduct(products[0])
	if err != nil {
		return []*Tile{}, err
	}

	for _, product := range products {
		tile, err := finders.main(product)
		if err != nil {
			_, err = finders.other(product)
			if err == nil {
				return nil, errors.New("input contains a mix of GUIDs and names")
			}
			return []*Tile{}, errorProductNotFound(product)
		}
		tiles = append(tiles, &tile)
	}
	return tiles, nil
}

func (t Tiles) FindByGuid(guid string) (Tile, error) {
	for _, t := range t.Data {
		if t.GUID == guid {
			return *t, nil
		}
	}
	return Tile{}, errors.New(fmt.Sprintf("product guid %s not found", guid))
}

func (t Tiles) getFindersForProduct(product string) (tileFinders, error) {
	_, err := t.FindBySlug(product)
	if err != nil {
		_, err := t.FindByGuid(product)
		if err != nil {
			return tileFinders{}, errorProductNotFound(product)
		}
		return tileFinders{t.FindByGuid, t.FindBySlug}, nil
	} else {
		return tileFinders{t.FindBySlug, t.FindByGuid}, nil
	}

}

func errorProductNotFound(product string) error {
	return errors.New(fmt.Sprintf("product %s is not found", product))
}
