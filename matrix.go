package mat

import "errors"

type Matrix struct {
	RowCount    int
	ColumnCount int
	Data        []float64
}

func NewMatrix(rows, columns int, data []float64) (Matrix, error) {
	if rows < 1 || columns < 1 {
		return Matrix{}, errors.New("mat.NewMatrix: matrix dimensions must be > 0")
	}
	if rows*columns != len(data) {
		return Matrix{}, errors.New("mat.NewMatrix: len(data) must be rows * columns")
	}
	return Matrix{
		RowCount:    rows,
		ColumnCount: columns,
		Data:        data,
	}, nil
}

func (m *Matrix) At(row, column int) float64 {
	return m.Data[row*m.ColumnCount+column]
}

func (m *Matrix) Set(row, column int, value float64) {
	m.Data[row*m.ColumnCount+column] = value
}

func (m Matrix) Row(y int) Matrix {
	r := Matrix{
		RowCount:    1,
		ColumnCount: m.ColumnCount,
		Data:        make([]float64, m.ColumnCount),
	}
	copy(r.Data, m.Data[y*m.ColumnCount:(y+1)*m.ColumnCount])
	return r
}

func (m Matrix) Column(x int) Matrix {
	c := Matrix{
		RowCount:    m.RowCount,
		ColumnCount: 1,
		Data:        make([]float64, m.RowCount),
	}
	for i := range c.Data {
		c.Data[i] = m.Data[x+i*m.ColumnCount]
	}
	return c
}

func (m Matrix) Transposed() Matrix {
	t := Matrix{
		RowCount:    m.ColumnCount,
		ColumnCount: m.RowCount,
		Data:        make([]float64, len(m.Data)),
	}
	for y := 0; y < t.RowCount; y++ {
		for x := 0; x < t.ColumnCount; x++ {
			t.Set(y, x, m.At(x, y))
		}
	}
	return t
}

func Multiply(m1, m2 Matrix, ms ...Matrix) (Matrix, error) {
	if m1.ColumnCount != m2.RowCount {
		return Matrix{}, errors.New("mat.Multiply: first matrix column count must match second matrix row count")
	}
	p := Matrix{
		RowCount:    m1.RowCount,
		ColumnCount: m2.ColumnCount,
		Data:        make([]float64, m1.RowCount*m2.ColumnCount),
	}

	i := 0
	for y := 0; y < p.RowCount; y++ {
		for x := 0; x < p.ColumnCount; x++ {
			for j := 0; j < m1.ColumnCount; j++ {
				p.Data[i] += m1.At(y, j) * m2.At(j, x)
			}
			i++
		}
	}

	if len(ms) == 0 {
		return p, nil
	}
	return Multiply(p, ms[0], ms[1:]...)
}
