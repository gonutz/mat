package mat_test

import (
	"testing"

	"github.com/gonutz/check"
	"github.com/gonutz/mat"
)

func TestMatrix(t *testing.T) {
	m, err := mat.NewMatrix(2, 3, []float64{
		1, 2, 3,
		4, 5, 6,
	})
	check.Eq(t, err, nil)
	check.Eq(t, m.RowCount, 2)
	check.Eq(t, m.ColumnCount, 3)
	check.Eq(t, m.Data, []float64{1, 2, 3, 4, 5, 6})
}

func TestMatrixAt(t *testing.T) {
	m, _ := mat.NewMatrix(2, 3, []float64{
		1, 2, 3,
		4, 5, 6,
	})
	check.Eq(t, m.At(0, 0), 1)
	check.Eq(t, m.At(0, 1), 2)
	check.Eq(t, m.At(0, 2), 3)
	check.Eq(t, m.At(1, 0), 4)
	check.Eq(t, m.At(1, 1), 5)
	check.Eq(t, m.At(1, 2), 6)
}

func TestMatrixSet(t *testing.T) {
	m, _ := mat.NewMatrix(2, 3, []float64{
		0, 0, 0,
		0, 0, 0,
	})
	m.Set(0, 0, 1)
	m.Set(0, 1, 2)
	m.Set(0, 2, 3)
	m.Set(1, 0, 4)
	m.Set(1, 1, 5)
	m.Set(1, 2, 6)
	check.Eq(t, m.Data, []float64{1, 2, 3, 4, 5, 6})
}

func TestMatrixTransposed(t *testing.T) {
	m, _ := mat.NewMatrix(2, 3, []float64{
		1, 2, 3,
		4, 5, 6,
	})
	trans := m.Transposed()
	check.Eq(t, trans.RowCount, 3)
	check.Eq(t, trans.ColumnCount, 2)
	check.Eq(t, trans.Data, []float64{
		1, 4,
		2, 5,
		3, 6,
	})
}

func TestMultiply(t *testing.T) {
	a, _ := mat.NewMatrix(2, 3, []float64{
		1, 2, 3,
		4, 5, 6,
	})
	b, _ := mat.NewMatrix(3, 4, []float64{
		3, 4, 5, 6,
		4, 5, 6, 7,
		5, 6, 7, 8,
	})
	c, err := mat.Multiply(a, b)
	check.Eq(t, err, nil)
	check.Eq(t, c.RowCount, 2)
	check.Eq(t, c.ColumnCount, 4)
	check.Eq(t, c.Data, []float64{
		26, 32, 38, 44,
		62, 77, 92, 107,
	})

	d, _ := mat.NewMatrix(4, 1, []float64{
		2,
		3,
		4,
		5,
	})
	e, err := mat.Multiply(a, b, d)
	check.Eq(t, err, nil)
	check.Eq(t, e.RowCount, 2)
	check.Eq(t, e.ColumnCount, 1)
	check.Eq(t, e.Data, []float64{520, 1258})

	f, _ := mat.NewMatrix(2, 4, []float64{
		1, 2, 3, 4,
		5, 6, 7, 8,
	})
	_, err = mat.Multiply(a, f)
	check.Neq(t, err, nil)
}
