package printer

import (
	"sync"
)

type IPrinter interface {
	PrintFrom(inputChan <-chan []string, wg *sync.WaitGroup)
}

func NewPrinter(pretty bool) IPrinter {
	if pretty {
		return newPrettyPrinter()
	}
	return newPlainPrinter()

}
