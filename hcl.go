package palettor

import (
	"image/color"

	"github.com/lucasb-eyer/go-colorful"
)

type hcl struct {
	h, c, l float64
}

func (c hcl) toColor() color.Color {
	// Bodge: squash floating point error to simplify testing with expected
	// output palettes.
	r, g, b := colorful.Hcl(c.h, c.c, c.l).Clamped().RGB255()
	return colorful.Color{
		R: float64(r / 255),
		G: float64(g / 255),
		B: float64(b / 255),
	}
}

// Calculate the square of the Euclidean distance between two colors, ignoring
// the alpha channel.
func (c hcl) distanceSquared(other hcl) float64 {
	dh := c.h - other.h
	dc := c.c - other.c
	dl := c.l - other.l
	return dh*dh + dc*dc + dl*dl
}

func toHCL(color colorful.Color) hcl {
	h, c, l := color.Hcl()
	return hcl{h, c, l}
}
