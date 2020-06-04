// Copyright 2020 Andrew Quinn. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gofrac

import (
	"errors"
	"math/cmplx"
	"runtime"
	"sync"
)

//var debug = log.New(os.Stdout, "DEBUG: ", log.LstdFlags)

// :-\
var glob = struct {
	maxIterations    int
	invMaxIterations float64
}{100, .01}

// MaxIterations gets the maximum number of iterations that should be performed
// before considering a point to be inside a set.
func MaxIterations() int {
	return glob.maxIterations
}

func setMaxIterations(iterations int) error {
	if iterations < 1 {
		return errors.New("gofrac: the maximum iteration count must be greater than zero")
	}
	glob.maxIterations = iterations
	glob.invMaxIterations = 1 / float64(iterations)
	return nil
}

// Fraccer maps a point in the complex plane to the result of a fractal calculation
type Fraccer interface {
	// Frac performs iterations of a fractal equation for a complex number
	// given by loc.
	Frac(loc complex128) *Result
}

// FracIt applies the fractal calculation given by f to every sample in the
// domain d. The maximum number of iterations to be performed is given by
// iterations.
func FracIt(d DomainReader, f Fraccer, iterations int) (*Results, error) {
	err := setMaxIterations(iterations)
	if err != nil {
		return nil, err
	}

	rows, cols := d.Dimensions()
	if cols < 1 || rows < 1 {
		return nil, errors.New("gofrac: the domain must be sampled at least once along each axis")
	}

	results := NewResults(rows, cols)
	defer results.Done()

	rowJobs := make(chan int, rows)

	numWorkers := runtime.NumCPU()
	wg := sync.WaitGroup{}
	for worker := 0; worker < numWorkers; worker++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for row := range rowJobs {
				for col := 0; col < cols; col++ {
					loc, err := d.At(col, row)
					if err != nil {
						panic(err)
					}
					r := f.Frac(loc)
					results.SetResult(row, col, r.Z, r.C, r.Iterations)
				}
			}
		}()
	}

	for row := 0; row < rows; row++ {
		rowJobs <- row
	}

	close(rowJobs)
	wg.Wait()

	return &results, nil
}

// Quadratic stores the information needed by a quadratic fractal.
type Quadratic struct {
	// Radius is the bailout radius of a fractal calculation.
	Radius float64
}

func (q Quadratic) q(z complex128, c complex128) *Result {
	count := 0
	for mod := cmplx.Abs(z); mod <= q.Radius; mod, count = cmplx.Abs(z), count+1 {
		z = z*z + c
		if count == glob.maxIterations-1 {
			break
		}
	}
	return &Result{
		Z:          z,
		C:          c,
		Iterations: count,
	}
}

// The Mandelbrot set, which results from iterating the function
// f_c(z) = z^2 + c, for all complex numbers c and z_0 = 0.
type Mandelbrot struct {
	Quadratic
}

// NewMandelbrot constructs a Mandelbrot struct with a given bailout radius.
func NewMandelbrot(radius float64) *Mandelbrot {
	return &Mandelbrot{Quadratic{radius}}
}

func (m Mandelbrot) Frac(loc complex128) *Result {
	return m.q(0, loc)
}

// JuliaQ is the quadratic Julia set, which results from iterating the function
// f_C(z) = z^2 + C for a all complex numbers z and a given complex number C.
type JuliaQ struct {
	Quadratic
	C complex128
}

// NewJuliaQ constructs a new JuliaQ struct with a given bailout radius and a
// complex parameter c corresponding to the C in f_C(z) = z^2 + C.
func NewJuliaQ(radius float64, c complex128) *JuliaQ {
	return &JuliaQ{Quadratic{radius}, c}
}

func (j JuliaQ) Frac(loc complex128) *Result {
	return j.q(loc, j.C)
}
