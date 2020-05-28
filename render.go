package gofrac

import (
	"image/color"
	"runtime"
	"sync"
)

// bitmap stores a 2D field of color.Color that can be used to generate images
type bitmap [][]color.Color

// NewBitmap initializes a bitmap
func NewBitmap(r int, c int) bitmap {
	b := make(bitmap, r)
	for r := range b {
		b[r] = make([]color.Color, c)
	}
	return b
}

func Render(plotter Plotter, results ResultsReader, palette ColorSampler) bitmap {
	rows, cols := results.Dimensions()
	bitmap := NewBitmap(rows, cols)

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
