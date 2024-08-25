package printer

import (
	"sync"
)

type StubPrinter struct {
	RecordRows   bool
	RowsReceived [][]string
}

func (p *StubPrinter) PrintFrom(inputChan <-chan []string, wg *sync.WaitGroup) {
	defer wg.Done()

	for row := range inputChan {
		if p.RecordRows {
			p.RowsReceived = append(p.RowsReceived, row)
		}
	}
}

func NewStubPrinter(recordRows bool) *StubPrinter {
	return &StubPrinter{
		RecordRows: recordRows,
	}
}
