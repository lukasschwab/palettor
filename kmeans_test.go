package palettor

import (
	"image"
	"image/color"
	_ "image/jpeg"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/stretchr/testify/assert"
)

var (
	r         = rand.New(rand.NewSource(time.Now().UnixNano()))
	black     = newColor(0, 0, 0, 255)
	white     = newColor(255, 255, 255, 255)
	red       = newColor(255, 0, 0, 255)
	green     = newColor(0, 255, 0, 255)
	blue      = newColor(0, 0, 255, 255)
	darkGray  = newColor(1, 1, 1, 255)
	mostlyRed = newColor(200, 0, 0, 255)

	// FIXME: Bodge; restandardize this interface.
	hclBlack     = toHCL(black)
	hclWhite     = toHCL(white)
	hclRed       = toHCL(red)
	hclGreen     = toHCL(green)
	hclBlue      = toHCL(blue)
	hclDarkGrey  = toHCL(darkGray)
	hclMostlyRed = toHCL(mostlyRed)
)

func randomColor() colorful.Color {
	return colorful.Hcl(r.Float64()*360, r.Float64(), r.Float64())
}

func newColor(r, g, b, a int) colorful.Color {
	// Bodge: keep constants.
	color, ok := colorful.MakeColor(&color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: uint8(a),
	})
	if !ok {
		panic("Color fixtures must have nonzero A-channel")
	}

	return color
}

func TestDistanceSquared(t *testing.T) {
	a := toHCL(newColor(0, 0, 0, 255))
	b := toHCL(newColor(255, 255, 255, 255))

	assert.InDelta(t, 1, distanceSquared(a, b), .0001, "distance should be square of Euclidean distance")

	a = toHCL(newColor(0, 0, 0, 1))
	b = toHCL(newColor(0, 0, 0, 255))
	assert.Equal(t, 0.00, distanceSquared(a, b), "alpha channel should be ignored for the purpose of distance")

	c := toHCL(randomColor())
	assert.Equal(t, 0.00, distanceSquared(c, c), "distance from between identical colors should be 0")
}

func TestNearest(t *testing.T) {
	var haystack = []hclColor{hclBlack, hclWhite, hclRed, hclGreen, hclBlue}

	assert.Equal(t, hclBlack, nearest(hclBlack, haystack), "nearest color to self should be self")
	assert.Equal(t, hclBlack, nearest(hclDarkGrey, haystack), "dark gray should be nearest to black")
	assert.Equal(t, hclRed, nearest(hclMostlyRed, haystack), "mostly-red should be nearest to red")
}

func TestFindCentroid(t *testing.T) {
	var cluster = []hclColor{hclBlack, hclWhite, hclRed, hclMostlyRed}
	centroid := findCentroid(cluster)

	assert.Contains(t, cluster, centroid, "centroid should be a member of the cluster")
}

func TestCluster(t *testing.T) {
	var colors = []colorful.Color{black, white, red}

	k := 4
	_, err := clusterColors(k, 100, colors)
	assert.Error(t, err, "too few colors should result in an error")

	k = 3
	palette, err := clusterColors(k, 100, colors)
	assert.NoError(t, err)
	assert.Equal(t, k, palette.Count(), "got unexpected number of clusters")

	k = 2
	colors = []colorful.Color{black, white}
	t.Logf("colors: %+v", colors)

	palette, _ = clusterColors(k, 100, colors)

	t.Logf("Palette: %+v", palette)

	assert.Equal(t, 0.5, palette.Weight(black), "expected weight of black cluster to be 0.5")
	assert.Equal(t, 0.5, palette.Weight(white), "expected weight of white cluster to be 0.5")

	// If there are not enough unique colors to cluster, it's okay for the size
	// of the extracted palette to be < k
	k = 3
	palette, _ = clusterColors(k, 100, []colorful.Color{black, black, black, black, black, white})
	assert.LessOrEqual(t, palette.Count(), 2, "actual palette can be smaller than k")
}

func BenchmarkClusterColors200x200(b *testing.B) {
	reader, err := os.Open("testdata/resized.jpg")
	if err != nil {
		b.Fatal(err)
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		b.Fatal(err)
	}

	colors, err := getColors(img)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := clusterColors(4, 100, colors); err != nil {
			b.Error(err)
		}
	}
}
