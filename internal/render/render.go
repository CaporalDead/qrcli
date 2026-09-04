// Package render converts QR code bitmaps into plain-text strings for a
// terminal. It uses Unicode half-block characters so that two vertically
// stacked QR modules fit into a single character cell, which keeps the
// code compact and roughly square on screen.
//
// The output contains no ANSI escape sequences: it renders correctly when
// piped, redirected to a file, or displayed in CI logs.
package render

import "strings"

// Options controls how a bitmap is rendered.
type Options struct {
	// Invert swaps dark and light cells. By default the light modules
	// (background and quiet zone) are drawn with filled blocks, which
	// produces a scannable dark-on-light code on dark terminals. Set
	// Invert when the terminal background is light.
	Invert bool

	// QuietZone is the number of light modules added around the code.
	// The QR spec asks for 4; terminals conventionally use 2, which
	// scanners handle fine. Negative values are treated as 0.
	QuietZone int
}

// Terminal renders a square QR bitmap (true = dark module) as lines of
// half-block characters. Rows are consumed two at a time: the upper module
// maps to the top half of the cell, the lower module to the bottom half.
// Every line ends with a newline. An empty bitmap renders as an empty string.
func Terminal(bitmap [][]bool, opts Options) string {
	size := len(bitmap)
	if size == 0 {
		return ""
	}

	quiet := opts.QuietZone
	if quiet < 0 {
		quiet = 0
	}
	total := size + 2*quiet

	dark := func(x, y int) bool {
		x, y = x-quiet, y-quiet
		if x < 0 || y < 0 || x >= size || y >= size {
			return false // outside the code everything is light
		}
		return bitmap[y][x]
	}
	lit := func(x, y int) bool {
		if opts.Invert {
			return dark(x, y)
		}
		return !dark(x, y)
	}

	var b strings.Builder
	b.Grow((total*3 + 1) * (total/2 + 1))
	for y := 0; y < total; y += 2 {
		for x := 0; x < total; x++ {
			b.WriteString(block(lit(x, y), lit(x, y+1)))
		}
		b.WriteByte('\n')
	}
	return b.String()
}

// block picks the half-block character lighting the requested cell halves.
func block(top, bottom bool) string {
	switch {
	case top && bottom:
		return "█"
	case top:
		return "▀"
	case bottom:
		return "▄"
	default:
		return " "
	}
}
