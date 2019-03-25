package stemcelldiff

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

type stemcellProduct struct {
	Guid                      string   `json:"guid"`
	Identifier                string   `json:"identifier"`
	Label                     string   `json:"label"`
	ProductVersion            string   `json:"product_version"`
	StagedStemcellVersion     string   `json:"staged_stemcell_version"`
	DeployedStemcellVersion   string   `json:"deployed_stemcell_version"`
	IsStagedForDeletion       bool     `json:"is_staged_for_deletion"`
	AvailableStemcellVersions []string `json:"available_stemcell_versions"`
	RequiredStemcellVersion   string   `json:"required_stemcell_version"`
	RequiredStemcellOs        string   `json:"required_stemcell_os"`
}

type omStemcellAssignments struct {
	Products []stemcellProduct `json:"products"`
}

type omProductEntry struct {
	ProductId string `json:"product_id"`
}

type omStemcellUpdateEntry struct {
	StemcellVersion string           `json:"stemcell_version"`
	ReleaseId       int32            `json:"release_id"`
	Products        []omProductEntry `json:"products"`
}

type omStemcellUpdates struct {
	StemcellUpdates []omStemcellUpdateEntry `json:"stemcell_updates"`
}

type availableStemcellProduct struct {
	GUID string `json:"guid"`
	Slug string `json:"slug"`
}

type availableStemcellEntry struct {
	StemcellVersion string                     `json:"stemcell_version"`
	StemcellOS      string                     `json:"stemcell_os"`
	ReleaseId       int32                      `json:"release_id"`
	Products        []availableStemcellProduct `json:"products"`
}

type availableStemcellUpdates struct {
	StemcellUpdates []availableStemcellEntry `json:"stemcell_updates"`
}

type availableStemcells struct {
	AvailableStemcells []availableStemcellEntry
}

func (o *availableStemcells) register(stemcellVersion string, stemcellOS string, products []availableStemcellProduct, releaseId int32) {
	if o.AvailableStemcells == nil {
		o.AvailableStemcells = []availableStemcellEntry{}
	}
	o.AvailableStemcells = append(o.AvailableStemcells, availableStemcellEntry{
		StemcellVersion: stemcellVersion,
		StemcellOS:      stemcellOS,
		Products:        products,
		ReleaseId:       releaseId,
	})
}

//go:generate counterfeiter . httpClient
type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

//go:generate counterfeiter . reporter
type reporter interface {
	PrintReport(report string)
}

type StemcellUpdateDetector struct {
	Client   httpClient
	Reporter reporter
}

func NewStemcellUpdateDetector(client httpClient, r reporter) StemcellUpdateDetector {
	return StemcellUpdateDetector{Client: client, Reporter: r}
}

func (s *StemcellUpdateDetector) DetectMissingStemcells() error {

	updates, err := s.getStemcellUpdates()
	if err != nil {
		return err
	}

	var assignments omStemcellAssignments
	if len(updates.StemcellUpdates) > 0 {
		assignments, err = s.getStemcellAssignments()
		if err != nil {
			return err
		}
	}

	stemcells := enhanceStemcellUpgrades(updates, assignments)
	outputBytes, err := json.Marshal(
		&availableStemcellUpdates{StemcellUpdates: stemcells.AvailableStemcells},
	)
	if err != nil {
		return err
	}

	s.Reporter.PrintReport(string(outputBytes))

	return nil
}

func enhanceStemcellUpgrades(omUpdates omStemcellUpdates, assignments omStemcellAssignments) availableStemcells {
	stemcells := availableStemcells{AvailableStemcells: []availableStemcellEntry{}}

	for _, updateEntry := range omUpdates.StemcellUpdates {
		stemcells.register(
			updateEntry.StemcellVersion,
			assignments.findStemcellOS(updateEntry.Products[0].ProductId),
			products(updateEntry, assignments),
			updateEntry.ReleaseId,
		)
	}

	return stemcells
}

func products(updateEntry omStemcellUpdateEntry, assignments omStemcellAssignments) []availableStemcellProduct {
	unupdatedProducts := []availableStemcellProduct{}
	for _, updateProduct := range updateEntry.Products {
		unupdatedProducts = append(unupdatedProducts, availableStemcellProduct{
			GUID: updateProduct.ProductId,
			Slug: assignments.findProductSlug(updateProduct.ProductId),
		})
	}
	return unupdatedProducts
}

func (s *StemcellUpdateDetector) getStemcellUpdates() (omStemcellUpdates, error) {
	availableStemcellsPath := "/api/v0/pivotal_network/stemcell_updates"
	availableStemcells, err := s.getContentForOmPath(availableStemcellsPath)
	if err != nil {
		return omStemcellUpdates{}, err
	}

	var latestStemcells omStemcellUpdates
	err = json.Unmarshal(availableStemcells, &latestStemcells)
	if err != nil {
		return omStemcellUpdates{}, err
	}

	return latestStemcells, nil
}

func (s *StemcellUpdateDetector) getStemcellAssignments() (omStemcellAssignments, error) {
	availableStemcellsPath := "/api/v0/stemcell_assignments"
	availableStemcells, err := s.getContentForOmPath(availableStemcellsPath)
	if err != nil {
		return omStemcellAssignments{}, err
	}
	var latestStemcells omStemcellAssignments
	err = json.Unmarshal(availableStemcells, &latestStemcells)
	if err != nil {
		return omStemcellAssignments{}, err
	}

	return latestStemcells, nil

}

func (s *StemcellUpdateDetector) getContentForOmPath(path string) ([]byte, error) {
	req, err := http.NewRequest("GET", path, nil)

	if err != nil {
		return nil, err
	}

	response, err := s.Client.Do(req)

	if err != nil {
		return nil, err
	}

	reply, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	return reply, nil
}

func (s *omStemcellAssignments) findStemcellOS(productId string) string {
	return s.findField(productId, "RequiredStemcellOs")

}

func (s *omStemcellAssignments) findProductSlug(productId string) string {
	return s.findField(productId, "Identifier")
}

func (s *omStemcellAssignments) findField(productId string, fieldName string) string {
	for _, product := range s.Products {
		if product.Guid == productId {
			return reflect.Indirect(reflect.ValueOf(product)).FieldByName(fieldName).String()
		}
	}
	return fmt.Sprintf("undefined_%s", strings.ToLower(fieldName))
}
