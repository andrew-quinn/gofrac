package gofrac

type Result struct {
	Z          complex128
	C          complex128
	Iterations int
	NFactor    float64
}

type Results [][]Result

func NewResults(rows int, cols int) Results {
	r := make([][]Result, rows)
	for row := range r {
		r[row] = make([]Result, cols)
	}
	return r
}

func (r Results) SetResult(row int, col int, z complex128, c complex128, iterations int) {
	r[row][col].Z = z
	r[row][col].C = c
	r[row][col].Iterations = iterations
}

func (r Results) At(row int, col int) *Result {
	return &r[row][col]
}

func (r Results) Dimensions() (rows int, cols int) {
	rows = len(r)
	if rows > 0 {
		cols = len(r[0])
	} else {
		cols = 0
	}
	return rows, cols
}

func calculateAccumulatedHistogram(r Results) (hist []int) {
	hist = make([]int, glob.maxIterations)

	// regular histogram
	for row := range r {
		for col := range r[row] {
			n := r[row][col].Iterations
			if n < glob.maxIterations-1 {
				hist[r[row][col].Iterations]++
			}
		}
	}

	// accumulate it
	for i, n := range hist {
		if i == 0 {
			continue
		}
		hist[i] = n + hist[i-1]
	}

	return hist
}

func setNFactors(r Results, hist []int) {
	invTotal := 1.0

	if glob.maxIterations > 1 {
		lastDivergent := hist[len(hist)-2]
		// non-degenerate case
		if lastDivergent > 0 {
			// only diverging locations need to be normalized
			invTotal = 1.0 / float64(lastDivergent)
		}
	}

	for row := range r {
		for col := range r[row] {
			i := r[row][col].Iterations
			if i < glob.maxIterations-1 {
				r[row][col].NFactor = float64(hist[i]) * invTotal
			}
		}
	}
}

func (r Results) Done() {
	hist := calculateAccumulatedHistogram(r)
	setNFactors(r, hist)
}
