package gofrac

import (
	"fmt"
	"image"
)

func GetImage(f Frac, d DomainReader, plotter Plotter, palette ColorSampler, maxIterations int) *image.RGBA {
	results, err := FracIt(d, f, maxIterations)
	if err != nil {
		fmt.Println("gofrac: An error occurred while generating the fractal: ", err.Error())
		return image.NewRGBA(image.Rect(0, 0, 0, 0))
	}

	bitmap := Render(results, plotter, palette)
	h, w := d.Dimensions()
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y, row := range bitmap {
		for x, clr := range row {
			img.Set(x, y, clr)
		}
	}
	return img
}
