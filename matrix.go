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

func (m Matrix) Copy() Matrix {
	c := Matrix{
		RowCount:    m.RowCount,
		ColumnCount: m.ColumnCount,
		Data:        make([]float64, len(m.Data)),
	}
	copy(c.Data, m.Data)
	return c
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

func (m *Matrix) Reshape(rows, columns int) error {
	if m.RowCount*m.ColumnCount != rows*columns {
		return errors.New("mat.Matrix.Reshape: new shape must have the same number of elements as old shape")
	}
	m.RowCount = rows
	m.ColumnCount = columns
	return nil
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

func (m Matrix) Inverse() (Matrix, error) {
	if m.RowCount != m.ColumnCount {
		return Matrix{}, errors.New("mat.Matrix.Inverse: matrix is not invertible, it is not square")
	}

	if m.RowCount == 1 {
		return m, nil
	}

	if m.RowCount == 2 {
		det := (m.Data[0]*m.Data[3] - m.Data[1]*m.Data[2])
		if det == 0 {
			return Matrix{}, errors.New("mat.Matrix.Inverse: matrix is not invertible, determinant is zero")
		}
		scale := 1.0 / det
		return Matrix{
			RowCount:    2,
			ColumnCount: 2,
			Data: []float64{
				scale * m.Data[3],
				-scale * m.Data[1],
				-scale * m.Data[2],
				scale * m.Data[0],
			},
		}, nil
	}

	n := m.RowCount
	left := m.Copy()
	right, _ := Identity(n)
	tempRow := make([]float64, n)
	swapRows := func(y1, y2 int) {
		if y1 == y2 {
			return
		}

		i1 := y1 * n
		i2 := y2 * n

		copy(tempRow, left.Data[i1:i1+n])
		copy(left.Data[i1:i1+n], left.Data[i2:i2+n])
		copy(left.Data[i2:i2+n], tempRow)

		copy(tempRow, right.Data[i1:i1+n])
		copy(right.Data[i1:i1+n], right.Data[i2:i2+n])
		copy(right.Data[i2:i2+n], tempRow)
	}
	scaleRow := func(row int, f float64) {
		start := row * n
		end := start + n
		for i := start; i < end; i++ {
			left.Data[i] *= f
			right.Data[i] *= f
		}
	}
	firstNonZero := func(x int) int {
		for y := x; y < n; y++ {
			if left.At(y, x) != 0 {
				return y
			}
		}
		return x
	}
	addScaledRow := func(rowToChange int, scale float64, sourceRow int) {
		dest := rowToChange * n
		src := sourceRow * n
		for x := 0; x < n; x++ {
			left.Data[dest+x] += scale * left.Data[src+x]
			right.Data[dest+x] += scale * right.Data[src+x]
		}
	}
	for i := 0; i < n; i++ {
		swapRows(i, firstNonZero(i))

		pivot := left.At(i, i)
		if pivot == 0 {
			return Matrix{}, errors.New("mat.Matrix.Inverse: matrix is not invertible")
		}
		scaleRow(i, 1.0/pivot)

		for y := i + 1; y < n; y++ {
			addScaledRow(y, -left.At(y, i), i)
		}
	}
	for i := n - 1; i >= 0; i-- {
		for y := 0; y < i; y++ {
			addScaledRow(y, -left.At(y, i), i)
		}
	}

	return right, nil
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

func Identity(dimension int) (Matrix, error) {
	if dimension <= 0 {
		return Matrix{}, errors.New("mat.Identity: dimension must be greater than zero")
	}
	m := Matrix{
		RowCount:    dimension,
		ColumnCount: dimension,
		Data:        make([]float64, dimension*dimension),
	}
	for i := 0; i < dimension; i++ {
		m.Data[i*(dimension+1)] = 1
	}
	return m, nil
}
