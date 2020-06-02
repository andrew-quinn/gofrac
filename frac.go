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

func MaxIterations() int {
	return glob.maxIterations
}

func setMaxIterations(iterations int) {
	if iterations < 1 {
		fmt.Println("gofrac: the maximum iteration count must be greater than zero")
	}
	glob.maxIterations = iterations
	glob.invMaxIterations = 1 / float64(iterations)
}

type Frac interface {
	Frac(loc complex128) *Result
}

func FracIt(d DomainReader, f Frac, iterations int) (*Results, error) {
	setMaxIterations(iterations)
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

type Quadratic struct {
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

type Mandelbrot struct {
	Quadratic
}

func NewMandelbrot(radius float64) *Mandelbrot {
	return &Mandelbrot{Quadratic{radius}}
}

func (m Mandelbrot) Frac(loc complex128) *Result {
	return m.q(0, loc)
}

type JuliaQ struct {
	Quadratic
	C complex128
}

func NewJuliaQ(radius float64, c complex128) *JuliaQ {
	return &JuliaQ{Quadratic{radius}, c}
}

func (j JuliaQ) Frac(loc complex128) *Result {
	return j.q(loc, j.C)
}
