package gofrac

import (
	"errors"
	"fmt"
)

// DomainReader reads values from a discretization of a bounded 2D space
type DomainReader interface {
	At(i int, j int) (loc complex128, err error)
	Dimensions() (rows int, cols int)
}

type Domain struct {
	x0, y0       float64
	xDist, yDist float64
	xs, ys       int

	wInv float64
	hInv float64
}

func (r *Domain) At(i int, j int) (loc complex128, err error) {
	if i < 0 || i >= r.xs || j < 0 || j >= r.ys {
		return 0, errors.New("gofrac: sample is out of bounds")
	}

	ti := float64(i) * r.wInv
	re := ti*r.xDist + r.x0

	tj := 1.0 - float64(j)*r.hInv
	im := tj*r.yDist + r.y0

	return complex(re, im), nil
}

func (r *Domain) Dimensions() (rows int, cols int) {
	return r.ys, r.xs
}

func NewDomain(x0, y0, x1, y1 float64, xSamples, ySamples int) *Domain {
	if xSamples <= 0 || ySamples <= 0 {
		fmt.Println("gofrac: a positive number of samples must be taken along both the x and y axes")
		return &Domain{xs: 1, ys: 1, wInv: 1.0, hInv: 1.0}
	}

	return &Domain{
		x0:    x0,
		y0:    y0,
		xs:    xSamples,
		ys:    ySamples,
		xDist: x1 - x0,
		yDist: y1 - y0,
		wInv:  1.0 / float64(xSamples),
		hInv:  1.0 / float64(ySamples),
	}
}
