// Package render converts QR code bitmaps into plain-text strings for a
// terminal. The default renderer uses Unicode half-block characters so that
// two vertically stacked QR modules fit into a single character cell, which
// keeps the code compact and roughly square on screen. An ASCII renderer is
// provided as a fallback for terminals without Unicode block glyphs.
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

// grid wraps a bitmap with the quiet-zone and polarity logic shared by all
// renderers, so each of them only decides how to draw a lit cell.
type grid struct {
	bitmap [][]bool
	size   int
	quiet  int
	invert bool
}

func newGrid(bitmap [][]bool, opts Options) grid {
	quiet := opts.QuietZone
	if quiet < 0 {
		quiet = 0
	}
	return grid{bitmap: bitmap, size: len(bitmap), quiet: quiet, invert: opts.Invert}
}

// total is the rendered width/height in modules, quiet zone included.
func (g grid) total() int {
	return g.size + 2*g.quiet
}

// lit reports whether the module at (x, y) — quiet zone included, out of
// range treated as light — must be drawn with a filled character.
func (g grid) lit(x, y int) bool {
	x, y = x-g.quiet, y-g.quiet
	dark := x >= 0 && y >= 0 && x < g.size && y < g.size && g.bitmap[y][x]
	if g.invert {
		return dark
	}
	return !dark
}

// Terminal renders a square QR bitmap (true = dark module) as lines of
// half-block characters. Rows are consumed two at a time: the upper module
// maps to the top half of the cell, the lower module to the bottom half.
// Every line ends with a newline. An empty bitmap renders as an empty string.
func Terminal(bitmap [][]bool, opts Options) string {
	if len(bitmap) == 0 {
		return ""
	}
	g := newGrid(bitmap, opts)
	total := g.total()

	var b strings.Builder
	b.Grow((total*3 + 1) * (total/2 + 1))
	for y := 0; y < total; y += 2 {
		for x := 0; x < total; x++ {
			b.WriteString(block(g.lit(x, y), g.lit(x, y+1)))
		}
		b.WriteByte('\n')
	}
	return b.String()
}

// ASCII renders a square QR bitmap using only 7-bit ASCII: "##" for lit
// modules, two spaces otherwise, one row per module. Twice as large as
// Terminal output in both directions, but survives any terminal that can
// print '#'. Same quiet-zone and invert semantics.
func ASCII(bitmap [][]bool, opts Options) string {
	if len(bitmap) == 0 {
		return ""
	}
	g := newGrid(bitmap, opts)
	total := g.total()

	var b strings.Builder
	b.Grow((total*2 + 1) * total)
	for y := 0; y < total; y++ {
		for x := 0; x < total; x++ {
			if g.lit(x, y) {
				b.WriteString("##")
			} else {
				b.WriteString("  ")
			}
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
