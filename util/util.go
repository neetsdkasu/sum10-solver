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

	MinLimitSeconds     = 1
	MaxLimitSeconds     = 60
	DefaultLimitSeconds = 10

	MinNumOfTestcase     = 1
	MaxNumOfTestcase     = 100
	DefaultNumOfTestcase = 10
)

func IsValidSeed(seed int) bool {
	return MinSeed <= seed && seed <= MaxSeed
}

func IsValidLimitSeconds(limitSeconds int) bool {
	return MinLimitSeconds <= limitSeconds &&
		limitSeconds <= MaxLimitSeconds
}

func IsValidNumOfTestcase(numOfTestcase int) bool {
	return MinNumOfTestcase <= numOfTestcase &&
		numOfTestcase <= MaxNumOfTestcase
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

func Swap(field [][]int, row1, col1, row2, col2 int) {
	field[row1][col1], field[row2][col2] = field[row2][col2], field[row1][col1]
}
