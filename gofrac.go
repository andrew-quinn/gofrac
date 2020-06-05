// Copyright 2020 Andrew Quinn. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package gofrac implements a fractal generation library
package gofrac

import (
	"errors"
	"image"
)

// GetImage performs an iterated fractal calculation within a domain and
// generates an image.RGBA of the output.
//
// f is the fractal generator (i.e., Mandelbrot, Julia, etc.).
// d is the input domain.
// plotter maps outputs of f to floating point values.
// palette maps floating point values to colors.
// maxIterations gives the number of iterations to be performed before
// considering a point to have converged.
func GetImage(f Fraccer, d DomainReader, plotter Plotter, palette ColorSampler, maxIterations int) (*image.RGBA, error) {
	if maxIterations < 1 {
		return nil, errors.New("gofrac: maximum iteration count must be greater than zero")
	}

	plotter.SetMaxIterations(maxIterations)

	results, err := FracIt(d, f, maxIterations)
	if err != nil {
		return nil, err
	}

	bitmap := Render(results, plotter, palette)
	h, w := d.Dimensions()
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y, row := range bitmap {
		for x, clr := range row {
			img.Set(x, y, clr)
		}
	}
	return img, nil
}
