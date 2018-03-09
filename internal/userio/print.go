package userio

import (
	"fmt"
	"os"
)

type ReportPrinter struct {}

func (rp ReportPrinter) PrintReport(report string, err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println(report)
}
