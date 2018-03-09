package fakes

type FakeReportPrinter struct {
	FakeReportFunc func(string, error)
}

func (rp FakeReportPrinter) PrintReport(report string, err error) {
	rp.FakeReportFunc(report, err)
}