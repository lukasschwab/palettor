# Palettor

Yet another way to extract dominant colors from an image using [k-means clustering][1].

[![Build Status](https://travis-ci.org/mccutchen/palettor.svg?branch=master)](http://travis-ci.org/mccutchen/palettor)
[![GoDoc](https://godoc.org/github.com/mccutchen/palettor?status.svg)](https://godoc.org/github.com/mccutchen/palettor)
[![Coverage](http://gocover.io/_badge/github.com/mccutchen/palettor?0)](http://gocover.io/github.com/mccutchen/palettor)


## Tests

### Unit tests

```
make test
```

### Benchmarks

```
make benchmark
```


## Example

```go
package main

import (
    "image"
    _ "image/gif"
    _ "image/jpeg"
    _ "image/png"
    "log"
    "os"

    "github.com/mccutchen/palettor"
    "github.com/nfnt/resize"
)

func main() {
    // Read an image from STDIN
    originalImg, _, err := image.Decode(os.Stdin)
    if err != nil {
        log.Fatal(err)
    }

    // Reduce it to a manageable size
    img := resize.Thumbnail(200, 200, originalImg, resize.Lanczos3)

    // Extract the 3 most dominant colors, halting the clustering algorithm
    // after 100 iterations if the clusters have not yet converged.
    k := 3
    maxIterations := 100
    palette, err := palettor.Extract(k, maxIterations, img)

    // Err will only be non-nil if k is larger than the number of pixels in the
    // input image.
    if err != nil {
        log.Fatalf("image too small")
    }

    // Palette is a mapping from color to the weight of that color's cluster,
    // which can be used as an approximation for that color's relative
    // dominance
    for _, color := range palette.Colors() {
        log.Printf("color: %v; weight: %v", color, palette.Weight(color))
    }

    // Example output:
    // 2015/07/19 10:27:52 color: {44 120 135}; weight: 0.17482142857142857
    // 2015/07/19 10:27:52 color: {140 103 150}; weight: 0.39558035714285716
    // 2015/07/19 10:27:52 color: {189 144 118}; weight: 0.42959821428571426
}
```


[1]: https://en.wikipedia.org/wiki/K-means_clustering#Standard_algorithm
