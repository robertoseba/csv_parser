// Helper for generating CSV files
package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

func generateCSV(filename string, columnCount int, rowCount int) {
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

func randomData(columns int, rows int) [][]string {
	data := make([][]string, rows)

	data[0] = make([]string, columns)

	for i := 0; i < columns; i++ {
		data[0][i] = "col" + fmt.Sprintf("%d", i+1)
	}

	for i := range data {
		if i == 0 {
			continue
		}

		data[i] = make([]string, columns)

		for j := 0; j < columns; j++ {
			data[i][j] = "row" + fmt.Sprintf("%d", i) + "col" + fmt.Sprintf("%d", j+1)
		}
	}
	return data
}
