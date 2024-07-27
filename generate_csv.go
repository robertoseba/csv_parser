// Helper for generating CSV files
package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
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
			if j == 0 {
				data[i][j] = "row_" + fmt.Sprintf("%d", i)
				continue
			}
			data[i][j] = fmt.Sprintf("%d", rand.Intn(1000))

		}
	}
	return data
}
