package errands

import (
	"github.com/pivotal-cf/om/api"
)

//go:generate counterfeiter . errandService
type errandService interface {
	List(productID string) (api.ErrandsListOutput, error)
	SetState(productID string, errandName string, postDeployState interface{}, preDeleteState interface{}) error
}

//go:generate counterfeiter . reporter
type reporter interface {
	PrintReport(report string)
}

//go:generate counterfeiter . tableReporter
type tableReporter interface {
	Write([]byte) (int, error)
	Flush() error
}
