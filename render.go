package gofrac

import (
	"github.com/lucasb-eyer/go-colorful"
	"image/color"
	"math"
	"math/cmplx"
)

type ColorSampler interface {
	SampleColor(val float64) color.Color
}

type SpectralPalette struct {
	Sweep  float64
	Offset float64
}

func (p *SpectralPalette) SampleColor(val float64) color.Color {
	if int(val) == maxIterations-1 {
		return color.Black
	}

	t := val / float64(maxIterations-1)
	return colorful.Hsv(t*p.Sweep+p.Offset, 1.0, 1.0)
}

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

func RenderEscapeTime(results ResultsReader, palette ColorSampler) bitmap {
	rows, cols := results.Dimensions()
	bitmap := newBitmap(rows, cols)
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			_, _, iterations := results.At(row, col)
			bitmap[row][col] = palette.SampleColor(float64(iterations))
		}
	}

	return bitmap
}

func RenderSmoothedEscapeTime(results ResultsReader, palette ColorSampler) bitmap {
	rows, cols := results.Dimensions()
	bitmap := newBitmap(rows, cols)
	log2 := math.Log(2.0)
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			z, _, iterations := results.At(row, col)
			var colorVal float64
			if iterations < maxIterations-1 {
				lz := math.Log(cmplx.Abs(z))
				nu := math.Log(lz/log2) / log2
				colorVal = float64(iterations+1) - nu
			} else {
				colorVal = float64(iterations)
			}
			bitmap[row][col] = palette.SampleColor(colorVal)
		}
	}
	return bitmap
}
