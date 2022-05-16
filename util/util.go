// SUM10-SOLVER
// author: Leonardone @ NEETSDKASU

package util

const (
	RowCount = 8
	ColCount = 8

	Sum      = 10
	Hole     = -1
	Obstacle = 11
)

const (
	MinSeed = 0
	MaxSeed = 99999
	NoSeed  = -1
)

type FieldViewer interface {
	Get(row, col int) int
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func IsValidSeed(seed int) bool {
	return MinSeed <= seed && seed <= MaxSeed
}

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

func FillField[T int | bool](field [][]T, value T) {
	for _, line := range field {
		for col := range line {
			line[col] = value
		}
	}
}

func Swap(field [][]int, row1, col1, row2, col2 int) {
	field[row1][col1], field[row2][col2] = field[row2][col2], field[row1][col1]
}
