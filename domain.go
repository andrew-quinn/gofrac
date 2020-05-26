package gofrac

import (
	"errors"
)

// DiscreteDomain is a discretization of a bounded 2D space
type DiscreteDomain interface {
	At(i int, j int) (x float64, y float64, err error)
	RowCount() (rows int)
	ColCount(rowIdx int) (cols int)
}

type RectangularDomain struct {
	x0, y0       float64
	xDist, yDist float64
	xs, ys       int

	wInv float64
	hInv float64
}

func (r *RectangularDomain) At(i int, j int) (x float64, y float64, err error) {
	if i < 0 || i >= r.xs || j < 0 || j >= r.ys {
		return x, y, errors.New("gofrac: sample is out of bounds")
	}

	ti := float64(i) * r.wInv
	x = ti*r.xDist + r.x0

	tj := 1.0 - float64(j)*r.hInv
	y = tj*r.yDist + r.y0

	return x, y, nil
}

func (r *RectangularDomain) RowCount() (rows int) {
	return r.ys
}

func (r *RectangularDomain) ColCount(_ int) (colCount int) {
	return r.xs
}

func NewRectangularDomain(x0, y0, x1, y1 float64, xSamples, ySamples int) (*RectangularDomain, error) {
	if xSamples <= 0 || ySamples <= 0 {
		return &RectangularDomain{}, errors.New("gofrac: a positive number of samples must be taken along both the x and y axes")
	}
	return &RectangularDomain{
		x0:    x0,
		y0:    y0,
		xs:    xSamples,
		ys:    ySamples,
		xDist: x1 - x0,
		yDist: y1 - y0,
		wInv:  1.0 / float64(xSamples),
		hInv:  1.0 / float64(ySamples),
	}, nil
}
