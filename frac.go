package gofrac

import (
	"errors"
	"math/cmplx"
	"runtime"
	"sync"
)

//var debug = log.New(os.Stdout, "DEBUG: ", log.LstdFlags)

const maxIterations = 1000

type Frac interface {
	frac(loc complex128) *Result
}

func fracIt(d DiscreteDomain, f Frac, iterations int) (*Results, error) {
	h := d.RowCount()
	w := d.ColCount()
	if w < 1 || h < 1 {
		return nil, errors.New("gofrac: w and h must both be greater than zero")
	}

	results := NewResults(h, w)

	rowJobs := make(chan int, h)

	numWorkers := runtime.NumCPU()
	wg := sync.WaitGroup{}
	for worker := 0; worker < numWorkers; worker++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for row := range rowJobs {
				for col := 0; col < w; col++ {
					loc, err := d.At(col, row)
					if err != nil {
						panic(err)
					}
					r := f.frac(loc)
					results.SetResult(row, col, r.z, r.c, r.iterations)
				}
			}
		}()
	}

	for row := 0; row < h; row++ {
		rowJobs <- row
	}

	close(rowJobs)
	wg.Wait()

	return results, nil
}

type Mandelbrot struct{}

func (_ Mandelbrot) frac(loc complex128) *Result {
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
func DefaultMandelbrot(w int, h int) *Results {
	d, err := NewDomain(-2.5, -1.0, 1.0, 1.0, w, h)
	if err != nil {
		o, _ := NewDomain(0.0, 0.0, 0.0, 0.0, 1, 1)
		d = o
	}

	r, err := fracIt(d, Mandelbrot{})
	if err != nil {
		// oops
	}
	return r
}

type Julia struct {
	c complex128
}

func (j Julia) frac(loc complex128) *Result {
	z := loc
	radius := 1024.0
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

func DefaultJulia(w int, h int, c complex128) *Results {
	x := 1.6
	y := 1.0
	d, err := NewDomain(-x, -y, x, y, w, h)
	if err != nil {
		o, _ := NewDomain(0.0, 0.0, 0.0, 0.0, 1, 1)
		d = o
	}

	j := Julia{c}
	r, err := fracIt(d, j)
	if err != nil {
		// oops
	}
	return r
}
