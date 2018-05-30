package stemcelldiff

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
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

type stemcellLibraryEntry struct {
	Version        string `json:"version"`
	Os             string `json:"os"`
	Infrastructure string `json:"infrastructure"`
	Hypervisor     string `json:"hypervisor"`
	Light          bool   `json:"light"`
}

type stemcellAssignments struct {
	Products        []stemcellProduct      `json:"products"`
	StemcellLibrary []stemcellLibraryEntry `json:"stemcell_library"`
}

type stemcellUpdateProductEntry struct {
	ProductId string `json:"product_id"`
}

type stemcellUpdateEntry struct {
	StemcellVersion string                       `json:"stemcell_version"`
	ReleaseId       int32                        `json:"release_id"`
	Products        []stemcellUpdateProductEntry `json:"products"`
}

type stemcellUpdates struct {
	StemcellUpdates []stemcellUpdateEntry `json:"stemcell_updates"`
}

type availableStemcellEntry struct {
	StemcellVersion string   `json:"stemcell_version"`
	Products        []string `json:"products"`
}

type availableStemcellOutput struct {
	AvailableStemcells map[string][]string `json:"available_stemcells"`
}

func (o *availableStemcellOutput) register(stemcell string, product string) {
	if o.AvailableStemcells[stemcell] == nil {
		o.AvailableStemcells[stemcell] = []string{product}
	} else {
		o.AvailableStemcells[stemcell] = append(o.AvailableStemcells[stemcell], product)
	}
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
	assignments, err := s.getStemcellAssignments()
	if err != nil {
		return err
	}

	updates, err := s.getStemcellUpdates()
	if err != nil {
		return err
	}

	output := &availableStemcellOutput{}
	output.AvailableStemcells = map[string][]string{}

	for _, updateEntry := range updates.StemcellUpdates {
		for _, updateProduct := range updateEntry.Products {
			if assignments.isStemcellDeployedForProduct(updateEntry.StemcellVersion, updateProduct.ProductId) {
				break
			}
			output.register(updateEntry.StemcellVersion, updateProduct.ProductId)
		}
	}
	outputBytes, err := json.Marshal(output)
	if err != nil {
		return err
	}

	s.Reporter.PrintReport(string(outputBytes))

	return nil
}

func (s *StemcellUpdateDetector) getStemcellUpdates() (stemcellUpdates, error) {
	availableStemcellsPath := "/api/v0/pivotal_network/stemcell_updates"
	availableStemcells, err := s.getContentForOmPath(availableStemcellsPath)
	if err != nil {
		return stemcellUpdates{}, err
	}

	var latestStemcells stemcellUpdates
	err = json.Unmarshal(availableStemcells, &latestStemcells)
	if err != nil {
		return stemcellUpdates{}, err
	}

	return latestStemcells, nil
}

func (s *StemcellUpdateDetector) getStemcellAssignments() (stemcellAssignments, error) {
	availableStemcellsPath := "/api/v0/stemcell_assignments"
	availableStemcells, err := s.getContentForOmPath(availableStemcellsPath)
	if err != nil {
		return stemcellAssignments{}, err
	}
	var latestStemcells stemcellAssignments
	err = json.Unmarshal(availableStemcells, &latestStemcells)
	if err != nil {
		return stemcellAssignments{}, err
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

func (s *stemcellAssignments) isStemcellDeployedForProduct(stemcellVersion string, productId string) bool {
	for _, product := range s.Products {
		if product.DeployedStemcellVersion == stemcellVersion && product.Guid == productId {
			return true
		}
	}
	return false
}
