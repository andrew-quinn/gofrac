package gofrac

type Result struct {
	z          complex128
	c          complex128
	iterations int
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
	r[row][col].z = z
	r[row][col].c = c
	r[row][col].iterations = iterations
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



}

	}
}
