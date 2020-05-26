package gofrac

type ResultsSetter interface {
	SetResult(row int, col int, z complex128, c complex128, iterations int)
}

type ResultsReader interface {
	At(row int, col int) (z complex128, c complex128, iterations int)
	Dimensions() (rows int, cols int)
}

type Result struct {
	z          complex128
	c          complex128
	iterations int
}

type Results struct {
	rows    int
	cols    int
	results [][]Result
}

func NewResults(rows int, cols int) *Results {
	return &Results{
		rows: rows,
		cols: cols,
		results: func() [][]Result {
			results := make([][]Result, rows)
			for row := range results {
				results[row] = make([]Result, cols)
			}
			return results
		}(),
	}
}

func (r *Results) SetResult(row int, col int, z complex128, c complex128, iterations int) {
	r.results[row][col].z = z
	r.results[row][col].c = c
	r.results[row][col].iterations = iterations
}

func (r *Results) At(row int, col int) (z complex128, c complex128, iterations int) {
	res := r.results[row][col]
	return res.z, res.c, res.iterations
}

func (r *Results) Dimensions() (rows int, cols int) {
	return r.rows, r.cols
}

// TODO: Measure whether this makes much of a difference for escape time visualization
/*
type EscapeTimeResults struct {
	counts [][]int
}

func (e *EscapeTimeResults) SetResult(row int, col int, _ complex128, _ complex128, iterations int) {
	e.counts[row][col] = iterations
}

func (e *EscapeTimeResults) At(row int, col int) (_ complex128, _ complex128, iterations int) {
	return 0, 0, e.counts[row][col]
}

func NewEscapeTimeCounter(rows int, cols int) *EscapeTimeResults {
	counts := make([][]int, rows)
	for row := range counts {
		counts[row] = make([]int, cols)
	}
	return &EscapeTimeResults{counts}
}
*/
