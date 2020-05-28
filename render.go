package gofrac

import (
	"github.com/lucasb-eyer/go-colorful"
	"image/color"
	"math"
	"math/cmplx"
	"runtime"
	"sync"
)

type ColorSampler interface {
	SampleColor(val float64) color.Color
}

type SpectralPalette struct {
	Sweep  float64
	Offset float64
}

func (p *SpectralPalette) SampleColor(val float64) color.Color {
	if int(val) == maxIterations-1 {
		return color.Black
	}

	t := val / float64(maxIterations-1)
	return colorful.Hsv(t*p.Sweep+p.Offset, 1.0, 1.0)
}

type BandedPalette struct {
	bands []color.Color
}

func NewUniformBandedPalette(colors ...color.Color) BandedPalette {
	return BandedPalette{bands: colors}
}

func (p *BandedPalette) SampleColor(val float64) color.Color {
	if int(val) == maxIterations-1 {
		return color.Black
	}
	t := val / float64(maxIterations-1)
	scaledVal := t * float64(len(p.bands)-1)
	sv := int(scaledVal)
	c1, _ := colorful.MakeColor(p.bands[sv])
	c2, _ := colorful.MakeColor(p.bands[sv+1])
	return c1.BlendHcl(c2, scaledVal-math.Floor(scaledVal))
}

var PrettyBands = NewUniformBandedPalette(
	colorful.Hsv(24.0, 0.38, 0.33),
	colorful.Hsv(158.0, 0.48, 0.73),
	colorful.Hsv(58.0, 0.72, 0.83),
	colorful.Hsv(58.0, 0.32, 0.95),
	colorful.Hsv(24.0, 0.86, 0.97),
)

var PrettyBands2 = NewUniformBandedPalette(
	colorful.Hsv(27.0, 0.75, 0.25),
	colorful.Hsv(188.0, 0.35, 0.82),
	colorful.Hsv(175.0, 0.13, 0.91),
	colorful.Hsv(35.0, 0.17, 0.85),
	colorful.Hsv(52.0, 0.06, 1.00),
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
