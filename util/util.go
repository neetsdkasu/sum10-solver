package util

const (
	RowCount = 8
	ColCount = 8
)

func FieldContains(row, col int) bool {
	return 0 <= row && row < RowCount &&
		0 <= col && col < ColCount
}

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

func Swap(field [][]int, row1, col1, row2, col2 int) {
	field[row1][col1], field[row2][col2] = field[row2][col2], field[row1][col1]
}
