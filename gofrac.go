package gofrac

import (
	"fmt"
	"image"
)

func GetImage(f Frac, d *Domain, plotter Plotter, palette ColorSampler, maxIterations int) *image.RGBA {
	results, err := fracIt(d, f, maxIterations)
	if err != nil {
		fmt.Println("gofrac: An error occurred while generating the fractal: ", err.Error())
		return image.NewRGBA(image.Rect(0, 0, 0, 0))
	}

	bitmap := Render(plotter, results, palette)
	img := image.NewRGBA(image.Rect(0, 0, d.ColCount(), d.RowCount()))
	for y, row := range bitmap {
		for x, clr := range row {
			img.Set(x, y, clr)
		}
	}
	return img
}
