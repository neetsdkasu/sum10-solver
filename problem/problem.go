// SUM10-SOLVER
// author: Leonardone @ NEETSDKASU

package problem

import (
	mt "bitbucket.org/neetsdkasu/mersenne_twister_go"
	"github.com/neetsdkasu/sum10-solver/util"
	"math"
)

type Problem struct {
	seed  uint32
	field [][]int
}

func (prob *Problem) Seed() uint32 {
	return prob.seed
}

func (prob *Problem) Get(row, col int) int {
	return prob.field[row][col]
}

func New(seed uint32) *Problem {
	mt := mt.NewMersenneTwister().Init(seed)
	numbers := make([]int, util.RowCount*util.ColCount)
	for i := range numbers {
		numbers[i] = i % 10
	}
	for i := len(numbers) - 1; 0 < i; i-- {
		k := int(math.Floor(mt.Real2() * float64(i+1)))
		numbers[i], numbers[k] = numbers[k], numbers[i]
	}
	field := util.MakeEmptyField[int]()
	for row := range field {
		for col := range field[row] {
			field[row][col] = numbers[row*util.ColCount+col]
		}
	}
	return &Problem{
		seed:  seed,
		field: field,
	}
}
