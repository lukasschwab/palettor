package palettor

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPalette(t *testing.T) {
	colorWeights := map[color.Color]float64{
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

	assert.Equal(t, len(colorWeights), palette.Count())
	assert.Equal(t, converged, palette.Converged())
	assert.Equal(t, iterations, palette.Iterations())

	assert.Equal(t, 0.75, palette.Weight(black), "wrong weight for black")
	assert.Equal(t, 0.00, palette.Weight(red), "wrong weight for unknown color")

	for _, color := range palette.Colors() {
		assert.Contains(t, colorWeights, color)
	}

	// ensure entries are sorted by weight
	expectedEntries := []Entry{
		{white, 0.25},
		{black, 0.75},
	}
	assert.Equal(t, expectedEntries, palette.Entries())
}
