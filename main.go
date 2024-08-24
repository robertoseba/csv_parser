package main

import (
	"flag"
	"os"
	"runtime/pprof"

	"github.com/robertoseba/csv_parser/cmd/app"
)

func main() {
	f, _ := os.Create("cpu.pprof")
	defer f.Close()
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	var colFilterFlag = flag.String("filter", "", "Filter the CSV file by the specified columns")
	var colRulesFlag = flag.String("rules", "", "Apply rules to the specified columns. Ex: -rules \"col1:eq(100)\"")
	flag.Parse()

	filename := flag.Arg(0)

	app.Run(filename, *colFilterFlag, *colRulesFlag)
}

//TODO: Fix printer for when output columns are greater than screen size
//TODO: print total records showing / total records processed

//TODO: Create readme
//TODO: Create CI/CD
//TODO: Publish project
