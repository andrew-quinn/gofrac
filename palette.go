package gofrac

import (
	"github.com/lucasb-eyer/go-colorful"
	"image/color"
	"math"
)

type ColorSampler interface {
	SampleColor(val float64) color.Color
}

type SpectralPalette struct {
	Sweep  float64
	Offset float64
}

func (p *SpectralPalette) SampleColor(val float64) color.Color {
	if int(val) == glob.maxIterations-1 {
		return color.Black
	}

	t := val / float64(glob.maxIterations-1)
	return colorful.Hsv(t*p.Sweep+p.Offset, 1.0, 1.0)
}

type BandedPalette struct {
	bands []color.Color
}

func NewUniformBandedPalette(colors ...color.Color) BandedPalette {
	return BandedPalette{bands: colors}
}

func (p *BandedPalette) SampleColor(val float64) color.Color {
	if int(val) == glob.maxIterations-1 {
		return color.Black
	}

	i := int(float64(len(p.bands)) * val / float64(glob.maxIterations-1))
	return p.bands[i]
}

type BlendedBandedPalette struct {
	BandedPalette
}

func (p *BlendedBandedPalette) SampleColor(val float64) color.Color {
	if int(val) == glob.maxIterations-1 {
		return color.Black
	}
	t := val / float64(glob.maxIterations-1)
	scaledVal := t * float64(len(p.bands)-1)
	sv := int(scaledVal)
	c1, _ := colorful.MakeColor(p.bands[sv])
	c2, _ := colorful.MakeColor(p.bands[sv+1])
	return c1.BlendLab(c2, scaledVal-math.Floor(scaledVal))
}

type PeriodicPalette struct {
	BandedPalette
	period int
}

func (p *PeriodicPalette) SampleColor(val float64) color.Color {
	if int(val) == glob.maxIterations-1 {
		return color.Black
	}

	idx := (int(val) / p.period) % len(p.bands)
	return p.bands[idx]
}
