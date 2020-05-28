package gofrac

import (
	"image/color"
	"math"
	"math/cmplx"
	"runtime"
	"sync"
)

// bitmap stores a 2D field of color.Color that can be used to generate images
type bitmap [][]color.Color

// newBitmap initializes a bitmap
func newBitmap(r int, c int) bitmap {
	b := make(bitmap, r)
	for r := range b {
		b[r] = make([]color.Color, c)
	}
	return b
}

type Plotter interface {
	Plot(r *Result) float64
}

type EscapeTimePlotter struct{}

func (p *EscapeTimePlotter) Plot(r *Result) float64 {
	return float64(r.iterations)
}

var invLog2 = 1.0 / math.Log(2.0)

type SmoothedEscapeTimePlotter struct{}

func (p *SmoothedEscapeTimePlotter) Plot(r *Result) float64 {
	if r.iterations < maxIterations-1 {
		lz := math.Log(cmplx.Abs(r.z))
		nu := math.Log(lz*invLog2) * invLog2
		return float64(r.iterations+1) - nu
	} else {
		return float64(r.iterations)
	}

}

func RenderEscapeTime(results ResultsReader, palette ColorSampler) bitmap {
	p := EscapeTimePlotter{}
	return render(&p, results, palette)
}

func RenderSmoothedEscapeTime(results ResultsReader, palette ColorSampler) bitmap {
	p := SmoothedEscapeTimePlotter{}
	return render(&p, results, palette)
}

func render(plotter Plotter, results ResultsReader, palette ColorSampler) bitmap {
	rows, cols := results.Dimensions()
	bitmap := newBitmap(rows, cols)

	rowJobs := make(chan int, rows)

	numWorkers := runtime.NumCPU()
	wg := sync.WaitGroup{}
	for worker := 0; worker < numWorkers; worker++ {
		wg.Add(1)
		go func() {
			for row := range rowJobs {
				for col := 0; col < cols; col++ {
					r := results.At(row, col)
					val := plotter.Plot(r)
					bitmap[row][col] = palette.SampleColor(val)
				}
			}
			wg.Done()
		}()
	}

	for row := 0; row < rows; row++ {
		rowJobs <- row
	}

	close(rowJobs)
	wg.Wait()

	return bitmap
}
