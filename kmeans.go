package palettor

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/lucasb-eyer/go-colorful"
)

// clusterColors finds k clusters in the given colors using the "standard"
// k-means clustering algorithm. It returns a Palette, after running the
// algorithm up to maxIterations times.
//
// Note: in terms of the standard algorithm[1], an observation in this
// implementation is simply a color, and we use the RGB channels as Euclidean
// coordinates for the purposes of finding the distance between two colors.
//
// [1]: https://en.wikipedia.org/wiki/K-means_clustering#Standard_algorithm
func clusterColors(k, maxIterations int, outerColors []colorful.Color) (*Palette, error) {
	// Convert once to HCL space.
	colors := make([]hclColor, len(outerColors))
	for i, color := range outerColors {
		colors[i] = toHCL(color)
	}

	colorCount := len(colors)
	if colorCount < k {
		return nil, fmt.Errorf("too few colors for k (%d < %d)", colorCount, k)
	}

	centroids := initializeStep(k, colors)
	var clusters map[hclColor][]hclColor
	var converged bool

	// The algorithm isn't guaranteed to converge, so we put a limit on the
	// number of attempts we will make.
	var iterations int
	for iterations = 0; iterations < maxIterations; iterations++ {
		clusters = assignmentStep(centroids, colors)
		converged, centroids = updateStep(clusters)
		if converged {
			break
		}
	}

	// Convert back to colorful RGB space for output.
	//
	// NOTE: this conversion may really be unnecessary; it could be ideal just
	// to use it as an intermediary (encapsulated with `hclColor`) to convert
	// back to somthing implementing `color.Color`.
	//
	// Alternatively, `hclColor` could implement `color.Color` directly.
	clusterWeights := make(map[colorful.Color]float64, k)
	for centroid, cluster := range clusters {
		clusterWeights[centroid.toColorfulColor()] = float64(len(cluster)) / float64(colorCount)
	}
	return &Palette{
		colorWeights: clusterWeights,
		iterations:   iterations,
		converged:    converged,
	}, nil
}

// Generate the initial list of k centroids from the given list of colors.
//
// TODO: Try other initialization methods?
// https://en.wikipedia.org/wiki/K-means_clustering#Initialization_methods
func initializeStep(k int, colors []hclColor) []hclColor {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	centroids := make([]hclColor, k)
	colorCount := len(colors)

	// Track random indexes we've used to avoid picking the same index for
	// multiple centroids in the case len(colors) is close to k.
	usedIndexes := make(map[int]struct{}, k)
	var index int
	for i := 0; i < k; i++ {
		for {
			index = r.Intn(colorCount)
			if _, used := usedIndexes[index]; !used {
				usedIndexes[index] = struct{}{}
				break
			}
		}
		centroids[i] = colors[index]
	}
	return centroids
}

// Assign each color to the cluster of the closest centroid.
func assignmentStep(centroids, colors []hclColor) map[hclColor][]hclColor {
	clusters := make(map[hclColor][]hclColor)
	for _, x := range colors {
		centroid := nearest(x, centroids)
		cluster, found := clusters[centroid]
		if !found {
			// allocate slice w/ maximum possible capacity to avoid possible
			// allocations per-append below
			cluster = make([]hclColor, 0, len(colors))
		}
		clusters[centroid] = append(cluster, x)
	}
	return clusters
}

// Pick new centroids from each cluster. If none of the centroids change, the
// clusters have stabilized and the algorithm has converged.
func updateStep(clusters map[hclColor][]hclColor) (bool, []hclColor) {
	converged := true
	newCentroids := make([]hclColor, 0, len(clusters))
	for centroid, cluster := range clusters {
		newCentroid := findCentroid(cluster)
		if newCentroid != centroid {
			converged = false
		}
		newCentroids = append(newCentroids, newCentroid)
	}
	return converged, newCentroids
}

// Find the color closest to the mean of the given colors.
//
// Note: I think this is a departure from the "standard" algorithm, which seems
// to instead use the actual mean of the given colors (which is likely
// not actually present in those colors).
func findCentroid(colors []hclColor) hclColor {
	center := meanColor(colors)
	return nearest(center, colors)
}

// Find the average color in a list of colors.
func meanColor(colors []hclColor) hclColor {
	var hSum, cSum, lSum float64
	for _, color := range colors {
		hSum += color.h
		cSum += color.c
		lSum += color.l
	}
	count := float64(len(colors))
	return hclColor{hSum / count, cSum / count, lSum / count}
}

// Find the item in the haystack to which the needle is closest.
func nearest(needle hclColor, haystack []hclColor) hclColor {
	var minDist float64
	var result hclColor
	for i, candidate := range haystack {
		dist := distanceSquared(needle, candidate)
		if i == 0 || dist < minDist {
			minDist = dist
			result = candidate
		}
	}
	return result
}

// Calculate the square of the Euclidean distance between two colors, ignoring
// the alpha channel.
func distanceSquared(a, b hclColor) float64 {
	// NOTE: consider using one of the (colorful.Color).DistanceX functions
	// rather than homebrewing euclidean distance.
	//
	// + `DistanceCIEDE2000` is the most accurate, but slow.
	// + `DistanceCIEDE2000klch` is the same, but allows specifying non-1 weights.
	// + `DistanceLab` is essentially this code but with an additional square
	// 	 root; That probably shouldn't matter for k-means.
	//
	// I suspect they're *all* slower than the euclidean distance here, but they
	// may produce better results.
	return a.distanceSquared(b)
}
