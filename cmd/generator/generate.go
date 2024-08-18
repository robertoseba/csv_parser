// Helper for generating CSV files
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
)

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

func main() {
	var cols = flag.Uint("cols", 2, "Number of columns to be generated")
	var rows = flag.Uint("rows", 10, "Number of rows to be generated")
	flag.Parse()

	filename := flag.Arg(0)

	if filename == "" {
		generateCSV("./test.csv", *cols, *rows)
		os.Exit(0)
	}

	generateCSV(filename, *cols, *rows)
}
