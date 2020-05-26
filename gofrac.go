package gofrac

import (
	"errors"
	"math/cmplx"
)

//var debug = log.New(os.Stdout, "DEBUG: ", log.LstdFlags)

const maxIterations = 35

// Mandelbrot generates the Mandelbrot set and saves the results to a 2D slice of Results. The parameters
// w and h are the number of samples to be taken along the horizontal and vertical axes of the domain, respectively.
func Mandelbrot(w int, h int) (*Results, error) {
	if w == 0 || h == 0 {
		return nil, errors.New("gofrac: w and h must both be greater than zero")
	}

	domain, err := NewRectangularDomain(-2.5, -1.0, 1.0, 1.0, w, h)
	if err != nil {
		return nil, errors.New("gofrac: Could not initialize a domain")
	}

	results := NewResults(h, w)

	for row := 0; row < domain.RowCount(); row++ {
		for col := 0; col < domain.ColCount(row); col++ {
			var z complex128 = 0
			x, y, _ := domain.At(col, row)
			var c = complex(x, y)

			count := 0
			for mod := cmplx.Abs(z); mod <= 16.0; mod, count = cmplx.Abs(z), count+1 {
				z = z*z + c

				if count == maxIterations-1 {
					break
				}
			}
			results.SetResult(row, col, z, c, count)
		}
	}

	return results, nil
}
