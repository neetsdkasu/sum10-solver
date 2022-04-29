package game

import (
	"errors"
	m "sum10-solver/marker"
	p "sum10-solver/problem"
	"sum10-solver/util"
)

const Sum = 10

type Game struct {
	Problem *p.Problem
	Steps   int
	Score   int
	Field   [][]int
}

func New(problem *p.Problem) *Game {
	if problem == nil {
		return nil
	}
	field := util.MakeEmptyField[int]()
	util.CopyField(field, problem.Field)
	return &Game{
		Problem: problem,
		Steps:   0,
		Score:   0,
		Field:   field,
	}
}

func Copy(dst, src *Game) {
	field := dst.Field
	*dst = *src
	util.CopyField(field, src.Field)
	dst.Field = field
}

var (
	InvalidMarker  = errors.New("Invalid Marker")
	UnsatisfiedSum = errors.New("Unsatisfied Sum")
)

func (game *Game) Take(marker m.Marker) (next *Game, err error) {
	if !marker.IsValid() {
		return nil, InvalidMarker
	}
	sum := 0
	for row, line := range marker.Field {
		for col, mark := range line {
			if mark {
				sum += game.Field[row][col]
			}
		}
	}
	if sum != Sum {
		return nil, UnsatisfiedSum
	}
	panic("TODO")
}
