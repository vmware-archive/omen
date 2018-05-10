package stemcells

import (
	"time"
	"encoding/json"
	"strconv"
)

//go:generate counterfeiter . Client
type Client interface {
	Get(endpoint string, timeout time.Duration) ([]byte, error)
}

type UpdateReporter struct {
	client Client
}

type product struct {
	Guid                      string   `json:"guid"`
	DeployedStemcellVersion   string   `json:"deployed_stemcell_version"`
	AvailableStemcellVersions []string `json:"available_stemcell_versions"`
}

type stemcellInfo struct {
}

type assignmentResponse struct {
	Products        []product      `json:"products"`
	StemcellLibrary []stemcellInfo `json:"stemcell_library"`
}

type reportProduct struct {
	ProductId string `json:"product_id"`
}

type stemcellUpdate struct {
	StemcellVersion string          `json:"stemcell_version"`
	Products        []reportProduct `json:"products"`
}

type Report struct {
	StemcellUpdates []stemcellUpdate `json:"stemcell_updates"`
}

func NewUpdateReporter(client Client) UpdateReporter {
	return UpdateReporter{
		client: client,
	}
}

func (u UpdateReporter) Report() (Report, error) {
	data, _ := u.client.Get("/api/v0/stemcell_assignments", time.Minute)
	reply := assignmentResponse{}
	json.Unmarshal(data, &reply)

	return reply.asReport()
}

func (r assignmentResponse) asReport() (Report, error) {
	report := Report{}
	versionMap := map[string][]reportProduct{}

	for _, product := range r.Products {
		newStemcell := product.maxAvailableStemcell()
		if (newStemcell != "") && (newStemcell != product.DeployedStemcellVersion) {
			versionMap[newStemcell] = append(versionMap[newStemcell], reportProduct{product.Guid})
		}
	}

	for stemcell, products := range versionMap {
		report.StemcellUpdates =
			append(report.StemcellUpdates, stemcellUpdate{StemcellVersion: stemcell, Products: products})
	}

	return report, nil
}

func (p product) maxAvailableStemcell() string {
	var maxVersion = ""
	for _, v := range p.AvailableStemcellVersions {
		if maxVersion == "" {
			maxVersion = v
		} else {
			maxVersionFloat, _ := strconv.ParseFloat(maxVersion, 32)
			vFloat, _ := strconv.ParseFloat(v, 32)
			if vFloat > maxVersionFloat {
				maxVersion = v
			}
		}
	}
	return maxVersion
}
