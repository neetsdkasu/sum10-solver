package util

const (
	RowCount = 8
	ColCount = 8
)

func MakeEmptyField[T int | bool]() [][]T {
	field := make([][]T, RowCount)
	for row := range field {
		field[row] = make([]T, ColCount)
	}
	return field
}

func CopyField[T int | bool](dst, src [][]T) {
	for row, line := range src {
		copy(dst[row], line)
	}
}
