package printer

import (
	"fmt"
	"strings"
	"sync"
)

type PlainPrinter struct{}

func (p *PlainPrinter) PrintFrom(inputChan <-chan []string, wg *sync.WaitGroup) {
	defer wg.Done()

	for row := range inputChan {
		fmt.Printf("%s\n", strings.Join(row, ","))
	}
}

func newPlainPrinter() *PlainPrinter {
	return &PlainPrinter{}
}
