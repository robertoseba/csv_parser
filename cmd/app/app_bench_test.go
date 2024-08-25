package app

import (
	"encoding/csv"
	"fmt"
	"os"
	"testing"

	printer "github.com/robertoseba/csv_parser/internal/printer/test"
	"golang.org/x/exp/rand"
)

func BenchmarkRun(b *testing.B) {
	filename := "test.csv"

	generateCSV(filename, 10, 10000)

	input := &InputOptions{
		Filename:    filename,
		FilterInput: "col1,col2",
		RulesInput:  "col1:eq(row_22)||eq(row_200)",
	}

	stubPrinter := printer.NewStubPrinter(false)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Run(input, stubPrinter)
	}

	b.StopTimer()
	os.Remove(filename)
}

func generateCSV(filename string, columnCount uint, rowCount uint) {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	data := randomData(columnCount, rowCount)

	for _, value := range data {
		err := writer.Write(value)
		if err != nil {
			panic(err)
		}
	}
}

func randomData(columns uint, rows uint) [][]string {
	data := make([][]string, rows+1)

	data[0] = make([]string, columns)

	for i := 0; i < int(columns); i++ {
		data[0][i] = "col" + fmt.Sprintf("%d", i+1)
	}

	for i := range data {
		if i == 0 {
			continue
		}

		data[i] = make([]string, columns)

		for j := 0; j < int(columns); j++ {
			if j == 0 {
				data[i][j] = "row_" + fmt.Sprintf("%d", i)
				continue
			}
			data[i][j] = fmt.Sprintf("%d", rand.Intn(1000))

		}
	}
	return data
}
