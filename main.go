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
	var orderedFlag = flag.Bool("ordered", false, "Print the rows in the same order as they appear in the file")
	flag.Parse()

	filename := flag.Arg(0)

	app.Run(filename, *colFilterFlag, *colRulesFlag, *orderedFlag)
}
