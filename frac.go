package gofrac

import (
	"errors"
	"fmt"
	"math/cmplx"
	"runtime"
	"sync"
)

//var debug = log.New(os.Stdout, "DEBUG: ", log.LstdFlags)

var glob = struct {
	maxIterations    int
	invMaxIterations float64
}{100, .01}

func setMaxIterations(iterations int) {
	if iterations < 1 {
		fmt.Println("gofrac: the maximum iteration count must be greater than zero")
	}
	glob.maxIterations = iterations
	glob.invMaxIterations = 1 / float64(iterations)
}

type Frac interface {
	frac(loc complex128) *Result
}

func FracIt(d DomainReader, f Frac, iterations int) (*Results, error) {
	setMaxIterations(iterations)
	rows, cols := d.Dimensions()
	if cols < 1 || rows < 1 {
		return nil, errors.New("gofrac: the domain must be sampled at least once along each axis")
	}

	results := NewResults(rows, cols)

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
					r := f.frac(loc)
					results.SetResult(row, col, r.z, r.c, r.iterations)
				}
			}
		}()
	}

	for row := 0; row < rows; row++ {
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

		if count == glob.maxIterations-1 {
			break
		}
	}
	return &Result{
		z:          z,
		c:          c,
		iterations: count,
	}
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
		if count == glob.maxIterations-1 {
			break
		}
	}
	return &Result{
		z:          z,
		c:          j.c,
		iterations: count,
	}
}
