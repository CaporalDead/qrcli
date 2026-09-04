package render

import "testing"

// FuzzRenderers exercises both renderers for panics across arbitrary
// bitmaps, quiet zones and polarities. Adopted from the 2026-09-05 audit
// (issue #22), where it ran ~1M execs clean; CI runs a short smoke on the
// weekly security lane.
func FuzzRenderers(f *testing.F) {
	f.Add([]byte{1, 0, 1, 1}, 2, false)
	f.Add([]byte{}, 0, true)
	f.Fuzz(func(_ *testing.T, data []byte, quiet int, invert bool) {
		if quiet < -8 || quiet > 32 {
			quiet %= 32
		}
		n := 0
		for n*n < len(data) {
			n++
		}
		bm := make([][]bool, n)
		i := 0
		for y := 0; y < n; y++ {
			bm[y] = make([]bool, n)
			for x := 0; x < n && i < len(data); x++ {
				bm[y][x] = data[i]&1 == 1
				i++
			}
		}
		opts := Options{QuietZone: quiet, Invert: invert}
		_ = Terminal(bm, opts)
		_ = ASCII(bm, opts)
	})
}
