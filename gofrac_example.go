// Copyright 2020 Andrew Quinn. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gofrac

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"strings"
)

var UHDRes = struct {
	w int
	h int
}{w: 3840, h: 2160}

func writeExample(filename string, img *image.RGBA) {
	if !strings.HasSuffix(filename, ".png") {
		filename += ".png"
	}
	outFile, err := os.Create(filename)
	if err != nil {
		panic("Could not write example to file: " + err.Error())
	}
	defer outFile.Close()

	png.Encode(outFile, img)
}

func handleExampleError(err error) {
	fmt.Printf("gofrac: Example image generation failed with the following error: %v", err)
}

// MandelbrotExample generates the classic Mandelbrot image and stores it as
// "mandelbrot.png".
func MandelbrotExample() {
	domain, err := NewDomain(-2.5, -1.0, 1.0, 1.0, UHDRes.w, UHDRes.h)
	if err != nil {
		handleExampleError(err)
		return
	}

	img, err := GetImage(
		NewMandelbrot(8000.0),
		domain,
		SmoothedEscapeTimePlotter{},
		&PrettyPeriodic,
		2500,
	)
	if err != nil {
		handleExampleError(err)
		return
	}

	writeExample("mandelbrot.png", img)
}

// JuliaQExample generates an interesting quadratic Julia set and stores it as
// "julia.png".
func JuliaQExample() {
	domain, err := NewDomain(-1.6, -1.0, 1.6, 1.0, UHDRes.w, UHDRes.h)
	if err != nil {
		handleExampleError(err)
		return
	}

	img, err := GetImage(
		NewJuliaQ(1024.0, complex(-0.8, 0.156)),
		domain,
		&SmoothedEscapeTimePlotter{},
		&SpectralPalette{Sweep: 360.0},
		200,
	)
	if err != nil {
		handleExampleError(err)
		return
	}

	writeExample("julia.png", img)
}
