package userio

import (
	"fmt"
	"os"
	"strings"
)

type ReportPrinter struct{}

func (rp ReportPrinter) PrintReport(report string) {
	if strings.HasSuffix(report, "\n") {
		fmt.Print(report)
	} else {
		fmt.Println(report)
	}
}

func (rp ReportPrinter) Fail(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
