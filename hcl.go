package palettor

import "github.com/lucasb-eyer/go-colorful"

type hclColor struct {
	h, c, l float64
}

func (c hclColor) toColorfulColor() colorful.Color {
	return colorful.Hcl(c.h, c.c, c.l)
}

func (c hclColor) distanceSquared(other hclColor) float64 {
	dh := c.h - other.h
	dc := c.c - other.c
	dl := c.l - other.l
	return dh*dh + dc*dc + dl*dl
}

func toHCL(color colorful.Color) hclColor {
	h, c, l := color.Hcl()
	return hclColor{h, c, l}
}
