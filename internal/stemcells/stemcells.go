package stemcells

import "time"

//go:generate counterfeiter . Client
type Client interface {
	Get(endpoint string, timeout time.Duration) ([]byte, error)
}

type UpdateReporter struct {
	client Client
}

type Report struct {
	product_guid                string   `json:"guid,omitempty"`
	available_stemcell_versions []string `json:"available_stemcell_versions,omitempty"`
}

func NewUpdateReporter(client Client) UpdateReporter {
	return UpdateReporter{
		client: client,
	}
}

func (u UpdateReporter) Report()
