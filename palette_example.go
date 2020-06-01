package gofrac

import "github.com/lucasb-eyer/go-colorful"

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

var PrettyBlends = BlendedBandedPalette{PrettyBands}

var PrettyBlends2 = BlendedBandedPalette{PrettyBands2}

var PrettyPeriodic = PeriodicPalette{
	period:        1,
	BandedPalette: PrettyBands,
}

var PrettyPeriodic2 = PeriodicPalette{
	period:        10,
	BandedPalette: PrettyBands2,
}
