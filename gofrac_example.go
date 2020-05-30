package gofrac

import (
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

func MandelbrotExample() {
	img := GetImage(
		Mandelbrot{radius: 8.0},
		NewDomain(-2.5, -1.0, 1.0, 1.0, UHDRes.w, UHDRes.h),
		SmoothedEscapeTimePlotter{},
		&PrettyBands,
		25,
	)
	writeExample("mandelbrot.png", img)
}

func JuliaExample() {
	img := GetImage(
		Julia{c: complex(-0.8, 0.156), radius: 1024.0},
		NewDomain(-1.6, -1.0, 1.6, 1.0, UHDRes.w, UHDRes.h),
		&SmoothedEscapeTimePlotter{},
		&SpectralPalette{Sweep: 360.0},
		200,
	)
	writeExample("julia.png", img)
}
