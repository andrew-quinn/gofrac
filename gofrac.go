package gofrac

import (
	"errors"
	"math/cmplx"
)

//var debug = log.New(os.Stdout, "DEBUG: ", log.LstdFlags)

const maxIterations = 35

type Frac interface {
	frac(re float64, im float64) *Result
}

func fracIt(w int, h int, f Frac) (*Results, error) {
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
			re, im, _ := domain.At(col, row)
			r := f.frac(re, im)

			results.SetResult(row, col, r.z, r.c, r.iterations)
		}
	}

	return results, nil
}

type mandelbrot struct{}

func (_ mandelbrot) frac(re float64, im float64) *Result {
	var z complex128 = 0
	var c = complex(re, im)

	count := 0
	for mod := cmplx.Abs(z); mod <= 6.0; mod, count = cmplx.Abs(z), count+1 {
		z = z*z + c

		if count == maxIterations-1 {
			break
		}
	}
	return &Result{
		z:          z,
		c:          c,
		iterations: count,
	}
}

// Mandelbrot generates the Mandelbrot set and saves the results to a 2D slice of Results. The parameters
// w and h are the number of samples to be taken along the horizontal and vertical axes of the domain, respectively.
func Mandelbrot(w int, h int) *Results {
	r, err := fracIt(w, h, mandelbrot{})
	if err != nil {
		// oops
	}
	return r
}
