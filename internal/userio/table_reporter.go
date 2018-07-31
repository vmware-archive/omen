package userio

import (
	"os"
	"text/tabwriter"
)

type TableReporter interface {
	Write([]byte) (int, error)
	Flush() error
}

func NewTableReporter() TableReporter {
	return tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
}
