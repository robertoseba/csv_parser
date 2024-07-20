package main

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/robertoseba/csv_parser/csv_parser"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		file, err := os.Open("./data.csv")

		if err != nil {
			fmt.Println("Failed to open file")
			panic(err)
		}

		defer file.Close()

		colFilters := []string{"col1","col4"}

		csvConfig := &csv_parser.CsvConfig{
			Separator: ',',
			ColFilters: colFilters, 
		}

		reader, err := csv_parser.New(file,	csvConfig) 
		
		if err != nil {
			fmt.Println("Failed to create reader")
			panic(err)
		}
		
		fmt.Println(reader.FilteredHeaders().Str())

		for {
			row, err := reader.ReadLine()

			if err == io.EOF{
				break
			}

			if err != nil {
				fmt.Println("Failed to read line")
				panic(err)
			}

			fmt.Println(row.Str())
		}
		return
	}

	if args[0] == "--generate-csv" {
		fmt.Printf("Generating CSV file for %s columns and %s rows...\n", args[1], args[2])

		rowCount, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println("Invalid row count")
			panic(err)
		}

		columnCount, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Invalid column count")
			panic(err)
		}

		generateCSV("./data.csv", columnCount, rowCount)

		fmt.Println("CSV file generated successfully")
	}
}
