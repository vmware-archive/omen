package userio

import (
	"fmt"
	"os"
)

func PrintReport(report string, err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println(report)
}
