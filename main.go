package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/robertoseba/csv_parser/pkg/app"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		app.Run("./data.csv", "col2,col3", "col2=row2col2")
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
