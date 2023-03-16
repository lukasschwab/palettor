package palettor

import (
	"reflect"
	"testing"

	"github.com/lucasb-eyer/go-colorful"
)

func TestPalette(t *testing.T) {
	colorWeights := map[colorful.Color]float64{
		black: 0.75,
		white: 0.25,
	}
	iterations := 1
	converged := true
	palette := &Palette{
		colorWeights: colorWeights,
		converged:    converged,
		iterations:   iterations,
	}

	if palette.Count() != len(colorWeights) {
		t.Errorf("wrong number of colors in palette")
	}

	if palette.Converged() != converged {
		t.Errorf("wrong value for converged in palette")
	}

	if palette.Iterations() != iterations {
		t.Errorf("wrong number of iterations in palette")
	}

	if palette.Weight(black) != 0.75 {
		t.Errorf("wrong weight for black")
	}
	if palette.Weight(red) != 0 {
		t.Errorf("wrong weight for unknown color")
	}

	for _, color := range palette.Colors() {
		found := false
		for inputColor := range colorWeights {
			if color == inputColor {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("missing color %v from palette", color)
		}
	}

	// ensure entries are sorted by weight
	expectedEntries := []Entry{
		{white, 0.25},
		{black, 0.75},
	}
	if entries := palette.Entries(); !reflect.DeepEqual(entries, expectedEntries) {
		t.Errorf("expected entries %v, got %v", expectedEntries, entries)
	}
}
