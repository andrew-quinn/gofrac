package gofrac

import (
	"errors"
	"github.com/lucasb-eyer/go-colorful"
	"image/color"
	"math/cmplx"
)

//var debug = log.New(os.Stdout, "DEBUG: ", log.LstdFlags)

const maxIterations = 50

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
	return newSpectralPalette(360.0, 0.0)
}

// bitmap stores a 2D field of color.Color that can be used to generate images
type bitmap [][]color.Color

// newBitmap initializes a bitmap
func newBitmap(r int, c int) *bitmap {
	b := make(bitmap, r)
	for r := range b {
		b[r] = make([]color.Color, c)
	}
	return &b
}

// Mandelbrot generates the Mandelbrot set and encodes the escape time as an element of a color palette. The parameters
// w and h are the number of samples to be taken along the horizontal and vertical axes of the domain, respectively.
func Mandelbrot(w int, h int) (*bitmap, error) {
	if w == 0 || h == 0 {
		return nil, errors.New("gofrac: w and h must both be greater than zero")
	}

	palette := newRainbowPalette()
	bitmap := newBitmap(h, w)

	hInv := 1.0 / float64(h)
	wInv := 1.0 / float64(w)

	for row := 0; row < h; row++ {
		ty := float64(row) * hInv
		y0 := 2.0*ty - 1.0

		for col := 0; col < w; col++ {
			var z complex128 = 0
			tx := float64(col) * wInv
			x0 := 3.5*tx - 2.5
			var c = complex(x0, y0)

			count := 0
			for mod := cmplx.Abs(z); mod <= 4.0; mod, count = cmplx.Abs(z), count+1 {
				z = z*z + c

				if count == maxIterations-1 {
					break
				}
			}
			(*bitmap)[row][col] = palette[count]
		}
	}
	return bitmap, nil
}
