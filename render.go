package gofrac

import (
	"github.com/lucasb-eyer/go-colorful"
	"image/color"
	"math"
	"math/cmplx"
)

// colorPalette stores a numerically indexed lookup table of colors
type colorPalette [maxIterations]color.Color

// newSpectralPalette creates a spectral color map, but with the final element colored black
func newSpectralPalette(sweep float64, offset float64) (p colorPalette) {
	for i := 0; i < maxIterations-1; i++ {
		t := float64(i) / float64(maxIterations-1)
		p[i] = colorful.Hsv(sweep*t+offset, 1.0, 1.0).Clamped()
	}
	p[len(p)-1] = color.Black
	return p
}

// newRainbowPalette creates the familiar (and !!not colorblind-friendly!!) rainbow color map, but with the final element colored black
func newRainbowPalette() colorPalette {
	return newSpectralPalette(210.0, 135.0)
}

// bitmap stores a 2D field of color.Color that can be used to generate images
type bitmap [][]color.Color

// newBitmap initializes a bitmap
func newBitmap(r int, c int) bitmap {
	b := make(bitmap, r)
	for r := range b {
		b[r] = make([]color.Color, c)
	}
	return b
}

func RenderEscapeTime(results ResultsReader) bitmap {
	rows, cols := results.Dimensions()
	bitmap := newBitmap(rows, cols)
	palette := newRainbowPalette()
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			_, _, iterations := results.At(row, col)
			bitmap[row][col] = palette[iterations]
		}
	}

	return bitmap
}

func RenderSmoothedEscapeTime(results ResultsReader) bitmap {
	rows, cols := results.Dimensions()
	bitmap := newBitmap(rows, cols)
	palette := newRainbowPalette()
	log2 := math.Log(2.0)
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			z, _, iterations := results.At(row, col)
			if iterations < maxIterations-1 {
				lz := math.Log(cmplx.Abs(z))
				nu := math.Log(lz/log2) / log2

				f := float64(iterations+1) - nu
				flo := int(math.Floor(f))
				//c1 := palette[flo]
				//c2 := palette[int(math.Ceil(f))]
				bitmap[row][col] = palette[flo]
			} else {
				bitmap[row][col] = palette[iterations]
			}
		}
	}
	return bitmap
}
