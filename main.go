package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
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
