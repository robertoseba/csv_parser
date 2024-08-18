package app

import (
	"encoding/csv"
	"fmt"
	"os"
	"syscall"
	"testing"

	"golang.org/x/exp/rand"
)

func BenchmarkEntrypoint(b *testing.B) {
	filename := "test.csv"

	defer func(stdout *os.File) {
		os.Stdout = stdout
	}(os.Stdout)
	os.Stdout = os.NewFile(uintptr(syscall.Stdin), os.DevNull)

	generateCSV(filename, 10, 1000)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Run(filename, "col1,col2", "col1:eq(row_22)")
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
