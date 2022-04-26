package game

import (
	"sum10-solver/problem"
)

type Game struct {
	Problem *problem.Problem
	Step    int
	Score   int
	Field   [][]int
}

func New(problem *problem.Problem) *Game {
	if problem == nil {
		return nil
	}
	field := make([][]int, len(problem.Field))
	for row := range field {
		field[row] = make([]int, len(problem.Field[row]))
		copy(field[row], problem.Field[row])
	}
	return &Game{
		Problem: problem,
		Step:    0,
		Score:   0,
		Field:   field,
	}
}
