package gofrac

import (
	"errors"
	"math/cmplx"
)

//var debug = log.New(os.Stdout, "DEBUG: ", log.LstdFlags)

const maxIterations = 35

type Frac interface {
	frac(loc complex128) *Result
}

func fracIt(d DiscreteDomain, f Frac) (*Results, error) {
	h := d.RowCount()
	w := d.ColCount()
	if w == 0 || h == 0 {
		return nil, errors.New("gofrac: w and h must both be greater than zero")
	}

	results := NewResults(h, w)

	for row := 0; row < h; row++ {
		for col := 0; col < w; col++ {
			loc, _ := d.At(col, row)
			r := f.frac(loc)

			results.SetResult(row, col, r.z, r.c, r.iterations)
		}
	}

	return results, nil
}

type mandelbrot struct{}

func (_ mandelbrot) frac(loc complex128) *Result {
	var z complex128 = 0
	var c = loc

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
	d, err := NewRectangularDomain(-2.5, -1.0, 1.0, 1.0, w, h)
	if err != nil {
		o, _ := NewRectangularDomain(0.0, 0.0, 0.0, 0.0, 1, 1)
		d = o
	}

	r, err := fracIt(d, mandelbrot{})
	if err != nil {
		// oops
	}
	return r
}

type julia struct {
	c complex128
}

func (j julia) frac(loc complex128) *Result {
	z := loc
	radius := 2.0
	count := 0
	for mod := cmplx.Abs(z); mod <= radius; mod, count = cmplx.Abs(z), count+1 {
		z = z*z + j.c
		if count == maxIterations-1 {
			break
		}
	}
	return &Result{
		z:          z,
		c:          j.c,
		iterations: count,
	}
}

func Julia(w int, h int, c complex128) *Results {
	x := 1.6
	y := 1.0
	d, err := NewRectangularDomain(-x, -y, x, y, w, h)
	if err != nil {
		o, _ := NewRectangularDomain(0.0, 0.0, 0.0, 0.0, 1, 1)
		d = o
	}

	j := julia{c}
	r, err := fracIt(d, j)
	if err != nil {
		// oops
	}
	return r
}
