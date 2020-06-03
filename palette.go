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

func (p SpectralPalette) SampleColor(val float64) color.Color {
	if int(val) == glob.maxIterations-1 {
		return color.Black
	}

	t := val / float64(glob.maxIterations-1)
	return colorful.Hsv(t*p.Sweep+p.Offset, 1.0, 1.0)
}

type BandedPalette color.Palette

func NewUniformBandedPalette(colors ...color.Color) BandedPalette {
	return colors
}

func (p BandedPalette) SampleColor(val float64) color.Color {
	if int(val) == glob.maxIterations-1 {
		return color.Black
	}

	i := int(float64(len(p)) * val / float64(glob.maxIterations-1))
	return p[i]
}

type BlendedBandedPalette color.Palette

func NewUniformBlendedBandedPalette(colors ...color.Color) BlendedBandedPalette {
	return BlendedBandedPalette(NewUniformBandedPalette(colors...))
}

func (p BlendedBandedPalette) SampleColor(val float64) color.Color {
	if int(val) == glob.maxIterations-1 {
		return color.Black
	}
	t := val / float64(glob.maxIterations-1)
	scaledVal := t * float64(len(p)-1)
	sv := int(scaledVal)
	c1, _ := colorful.MakeColor(p[sv])
	c2, _ := colorful.MakeColor(p[sv+1])
	return c1.BlendLab(c2, scaledVal-math.Floor(scaledVal))
}

type PeriodicPalette struct {
	BandedPalette
	Period int
}

func (p PeriodicPalette) SampleColor(val float64) color.Color {
	if int(val) == glob.maxIterations-1 {
		return color.Black
	}

	idx := (int(val) / p.Period) % len(p.BandedPalette)
	return p.BandedPalette[idx]
}
