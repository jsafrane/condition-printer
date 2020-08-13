package main

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/logrusorgru/aurora/v3"
	"golang.org/x/crypto/ssh/terminal"
)

// Quick and dirty implementation of table printer.
// I was not able to find an existing one that supports *both* printing
// colorful cells and wrapping long message lines into cells.

// Table. All rows must have the same nr. of columns. All columns except for
// the last one are expected to be short enough to fit any screen (they're not
// wrapped). The last column is wrapped when its text is too long.
type table []row

// A single row of a table.
type row []item

// A single cell of a table. Has text and color.
type item struct {
	color aurora.Color
	text  string
}

func (t table) Sort() {
	sort.Slice(t, func(i, j int) bool {
		return t[i][0].text < t[j][0].text
	})
}

func (t table) Print() {
	if len(t) == 0 {
		return
	}

	au := aurora.NewAurora(true)
	width, _, err := terminal.GetSize(1 /* stdout */)
	if err != nil {
		// no terminal present, don't wrap long lines and don't use colors
		width = math.MaxInt32
		au = aurora.NewAurora(false)
	}

	colSizes := t.calcColumnSizes()

	for i := range t {
		start := 0
		for j := range t[i] {
			if j != 0 {
				fmt.Print(" | ")
			}
			if j < len(colSizes)-1 {
				t[i][j].printAligned(au, colSizes[j])
			} else {
				t[i][j].printWrapped(au, start, width)
			}
			start = start + colSizes[j] + 3
		}
		fmt.Print("\n")
	}
}

func (t table) calcColumnSizes() []int {
	header := t[0]
	colSizes := make([]int, len(header))
	for i := range t {
		for j := range t[i] {
			item := t[i][j]
			if len(item.text) > colSizes[j] {
				colSizes[j] = len(item.text)
			}
		}
	}
	return colSizes
}

func (it item) printAligned(au aurora.Aurora, columnSize int) {
	fmt.Print(au.Colorize(it.text, it.color))
	numSpaces := columnSize - len(it.text)
	fmt.Print(strings.Repeat(" ", numSpaces))
}

func (it item) printWrapped(au aurora.Aurora, skip int, lineLength int) {
	text := it.text
	start := 0
	for {
		// Find the next \n or the max. string we can print
		var i int
		for i = start; i < len(text) && skip+i-start < lineLength && text[i] != '\n'; i++ {
		}
		if i == len(text) {
			// Found end of the text
			fmt.Print(au.Colorize(text[start:], it.color))
			return
		}
		// The text continues. Print this line and start a new one.
		fmt.Println(au.Colorize(text[start:i], it.color))
		if text[i] == '\n' {
			// Found end of line, consume it
			i = i + 1
		}
		start = i
		fmt.Print(strings.Repeat(" ", skip))
	}
}
