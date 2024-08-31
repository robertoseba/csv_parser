package printer

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

type PrettyPrinter struct {
	maxColWidth      []int
	minCellWidth     int
	colOverflowAtIdx int
	maxHeight        int
	maxWidth         int
	lineNumber       int
	style            lipgloss.Style
	headers          []string
}

func newPrettyPrinter() *PrettyPrinter {
	var width, height int

	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width = 100
		height = 80
	}

	re := lipgloss.NewRenderer(os.Stdout)
	baseStyle := re.NewStyle().
		Padding(0, 2).
		TabWidth(4)

	return &PrettyPrinter{
		minCellWidth: 6,
		maxHeight:    height,
		maxWidth:     width,
		lineNumber:   0,
		style:        baseStyle,
	}
}

func (p *PrettyPrinter) PrintFrom(inputChan <-chan []string, wg *sync.WaitGroup) {
	startTime := time.Now()
	defer wg.Done()
	defer p.wrapUpPrint(startTime)

	p.headers = <-inputChan
	// Since the headers can have column names that are longer
	// than the actual data, we use the first row data to calculate
	// the max width of each column.
	firstRow := <-inputChan

	if len(firstRow) == 0 {
		p.calcMaxColWidth(p.headers)
	} else {
		p.calcMaxColWidth(firstRow)
	}

	p.printHeader(p.headers)
	p.printRow(firstRow)

	for row := range inputChan {
		p.printRow(row)

		if !p.shouldKeepPrinting() {
			break
		}
	}
}

func (p *PrettyPrinter) printHeader(headers []string) {
	headerStyle := p.style.Foreground(lipgloss.Color("#000000")).Background(lipgloss.Color("#D7FF87"))
	p.print(headerStyle, headers, true)
}

func (p *PrettyPrinter) printRow(line []string) {
	rowStyle := p.style.Foreground(lipgloss.Color("#D7FF87")).Faint(p.lineNumber%2 == 0)
	p.print(rowStyle, line, false)
}

func (p *PrettyPrinter) print(style lipgloss.Style, row []string, isHeaders bool) {
	if len(row) == 0 {
		return
	}

	if isHeaders {
		fmt.Print(style.Render("Line#"))
	} else {
		fmt.Print(style.Foreground(lipgloss.Color("#6F6F6F")).Render(fmt.Sprintf("%5d", p.lineNumber)))
	}

	for idx, cell := range row {
		if idx > p.colOverflowAtIdx {
			fmt.Print(style.Foreground(lipgloss.Color("#6F6F6F")).Padding(0).Render("..."))
			break
		}
		fmt.Print(style.Render(resizeCell(cell, p.maxColWidth[idx])))
	}

	fmt.Println()
	p.lineNumber++
}

func (p *PrettyPrinter) calcMaxColWidth(row []string) {
	p.maxColWidth = make([]int, 0, len(row))
	p.colOverflowAtIdx = len(row)

	hPadding := p.style.GetHorizontalPadding()
	totalWidth := lipgloss.Width("Line#") + hPadding

	for idx, cell := range row {
		textWidth := lipgloss.Width(cell)
		if textWidth < p.minCellWidth {
			textWidth = p.minCellWidth
		}

		if totalWidth+textWidth+hPadding > p.maxWidth-lipgloss.Width("...")-hPadding {
			p.colOverflowAtIdx = idx
			textWidth = p.maxWidth - totalWidth - hPadding - lipgloss.Width("...")
			if textWidth < 0 {
				textWidth = 0
			}

			p.maxColWidth = append(p.maxColWidth, textWidth)
			break
		}
		totalWidth += textWidth + hPadding
		p.maxColWidth = append(p.maxColWidth, textWidth)
	}

	if totalWidth < p.maxWidth && p.colOverflowAtIdx == len(row) {
		extraSpacePerCell := (p.maxWidth - totalWidth) / len(p.maxColWidth)
		for idx, cellWidth := range p.maxColWidth {

			p.maxColWidth[idx] = cellWidth + extraSpacePerCell
		}
	}
}

func resizeCell(cell string, maxWidth int) string {
	if len(cell) > maxWidth {
		if maxWidth < lipgloss.Width("...") {
			return ""
		}
		return fmt.Sprintf("%s...", cell[:maxWidth-lipgloss.Width("...")])
	}
	return cell + strings.Repeat(" ", maxWidth-len(cell))
}

func (p *PrettyPrinter) shouldKeepPrinting() bool {
	if p.lineNumber%(p.maxHeight-2) == 0 {
		fmt.Println("Press Enter to continue or type q to exit")
		var input string
		fmt.Scanln(&input)

		if strings.ToLower(input) == "q" {
			return false
		}
		p.printHeader(p.headers)
	}
	return true
}

func (p *PrettyPrinter) wrapUpPrint(start time.Time) {
	fmt.Println()
	fmt.Printf("%d lines printed\n", p.lineNumber-1)
	fmt.Println()
}
