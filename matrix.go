package godifft

type matrix struct {
	m    int
	n    int
	data []int
}

func newMatrix(m, n int) *matrix {
	return &matrix{m, n, make([]int, m*n)}
}

func (m *matrix) get(i, j int) int {
	return m.data[m.index(i, j)]
}

func (m *matrix) set(i, j int, v int) {
	m.data[m.index(i, j)] = v
}

func (m *matrix) index(i, j int) int {
	return m.n*i + j
}
