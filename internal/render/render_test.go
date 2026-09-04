package render

import (
	"strings"
	"testing"
)

func TestTerminal(t *testing.T) {
	tests := []struct {
		name   string
		bitmap [][]bool
		opts   Options
		want   string
	}{
		{
			name:   "empty bitmap",
			bitmap: [][]bool{},
			opts:   Options{},
			want:   "",
		},
		{
			// Single dark module: top half is dark (terminal background),
			// the half-row below the code is light.
			name:   "single dark module",
			bitmap: [][]bool{{true}},
			opts:   Options{},
			want:   "▄\n",
		},
		{
			name:   "single dark module inverted",
			bitmap: [][]bool{{true}},
			opts:   Options{Invert: true},
			want:   "▀\n",
		},
		{
			name:   "2x2 checkerboard",
			bitmap: [][]bool{{true, false}, {false, true}},
			opts:   Options{},
			want:   "▄▀\n",
		},
		{
			name:   "2x2 checkerboard inverted",
			bitmap: [][]bool{{true, false}, {false, true}},
			opts:   Options{Invert: true},
			want:   "▀▄\n",
		},
		{
			// Odd height: the last pair of rows extends past the bitmap;
			// the out-of-range half must render as light (quiet).
			name: "3x3 diamond, odd height",
			bitmap: [][]bool{
				{true, false, true},
				{false, true, false},
				{true, false, true},
			},
			opts: Options{},
			want: "▄▀▄\n▄█▄\n",
		},
		{
			name:   "quiet zone wraps the code in light modules",
			bitmap: [][]bool{{true}},
			opts:   Options{QuietZone: 1},
			want:   "█▀█\n███\n",
		},
		{
			name:   "negative quiet zone is clamped to zero",
			bitmap: [][]bool{{true}},
			opts:   Options{QuietZone: -3},
			want:   "▄\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Terminal(tt.bitmap, tt.opts)
			if got != tt.want {
				t.Errorf("Terminal() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTerminalGeometry(t *testing.T) {
	// A 21x21 bitmap (QR version 1) with quiet zone 2 must render as
	// 25-column lines, 13 of them (ceil(25/2)), each newline-terminated.
	size := 21
	bitmap := make([][]bool, size)
	for i := range bitmap {
		bitmap[i] = make([]bool, size)
		for j := range bitmap[i] {
			bitmap[i][j] = (i+j)%2 == 0
		}
	}

	out := Terminal(bitmap, Options{QuietZone: 2})
	if !strings.HasSuffix(out, "\n") {
		t.Fatal("output must end with a newline")
	}
	lines := strings.Split(strings.TrimSuffix(out, "\n"), "\n")
	wantLines := (size + 4 + 1) / 2
	if len(lines) != wantLines {
		t.Errorf("got %d lines, want %d", len(lines), wantLines)
	}
	for i, line := range lines {
		if n := len([]rune(line)); n != size+4 {
			t.Errorf("line %d has %d cells, want %d", i, n, size+4)
		}
	}
}

func TestASCII(t *testing.T) {
	tests := []struct {
		name   string
		bitmap [][]bool
		opts   Options
		want   string
	}{
		{
			name:   "empty bitmap",
			bitmap: [][]bool{},
			opts:   Options{},
			want:   "",
		},
		{
			name:   "single dark module",
			bitmap: [][]bool{{true}},
			opts:   Options{},
			want:   "  \n",
		},
		{
			name:   "single dark module inverted",
			bitmap: [][]bool{{true}},
			opts:   Options{Invert: true},
			want:   "##\n",
		},
		{
			name:   "2x2 checkerboard",
			bitmap: [][]bool{{true, false}, {false, true}},
			opts:   Options{},
			want:   "  ##\n##  \n",
		},
		{
			name:   "quiet zone wraps the code in lit modules",
			bitmap: [][]bool{{true}},
			opts:   Options{QuietZone: 1},
			want:   "######\n##  ##\n######\n",
		},
		{
			name:   "negative quiet zone is clamped to zero",
			bitmap: [][]bool{{true}},
			opts:   Options{QuietZone: -3},
			want:   "  \n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ASCII(tt.bitmap, tt.opts)
			if got != tt.want {
				t.Errorf("ASCII() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestASCIIIsSevenBitOnly(t *testing.T) {
	bitmap := [][]bool{{true, false}, {false, true}}
	for _, opts := range []Options{{}, {Invert: true}, {QuietZone: 2}} {
		for _, r := range ASCII(bitmap, opts) {
			if r > 127 {
				t.Fatalf("ASCII() emitted non-ASCII rune %q with opts %+v", r, opts)
			}
		}
	}
}

func TestASCIIGeometry(t *testing.T) {
	// One row per module, two columns per module: a 21x21 bitmap with
	// quiet zone 2 renders as 25 lines of 50 characters.
	size := 21
	bitmap := make([][]bool, size)
	for i := range bitmap {
		bitmap[i] = make([]bool, size)
	}

	out := ASCII(bitmap, Options{QuietZone: 2})
	lines := strings.Split(strings.TrimSuffix(out, "\n"), "\n")
	if len(lines) != size+4 {
		t.Errorf("got %d lines, want %d", len(lines), size+4)
	}
	for i, line := range lines {
		if len(line) != 2*(size+4) {
			t.Errorf("line %d has %d chars, want %d", i, len(line), 2*(size+4))
		}
	}
}

func TestBlock(t *testing.T) {
	cases := map[[2]bool]string{
		{true, true}:   "█",
		{true, false}:  "▀",
		{false, true}:  "▄",
		{false, false}: " ",
	}
	for in, want := range cases {
		if got := block(in[0], in[1]); got != want {
			t.Errorf("block(%v, %v) = %q, want %q", in[0], in[1], got, want)
		}
	}
}
