package game

import (
	"errors"
	m "sum10-solver/marker"
	p "sum10-solver/problem"
	"sum10-solver/util"
)

const (
	Sum      = 10
	Hole     = -1
	Obstacle = 11
)

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

func (game *Game) GetCopy() *Game {
	dst := &Game{}
	dst.Field = util.MakeEmptyField[int]()
	Copy(dst, game)
	return dst
}

var (
	InvalidMarker  = errors.New("Invalid Marker")
	UnsatisfiedSum = errors.New("Unsatisfied Sum")
)

func (game *Game) Take(marker *m.Marker) (next *Game, err error) {
	if !marker.IsValid() {
		return nil, InvalidMarker
	}
	sum := marker.Sum(game.Field)
	if sum != Sum {
		return nil, UnsatisfiedSum
	}
	next = game.GetCopy()
	marker.Fill(next.Field, Hole)
	next.fill()
	next.Steps++
	next.Score += Sum
	return next, nil
}

func (game *Game) fill() {
	switch game.Steps & 3 {
	case 0:
		game.moveToBottom()
	case 1:
		game.moveToLeft()
	case 2:
		game.moveToTop()
	case 3:
		game.moveToRight()
	}
	for _, line := range game.Field {
		for col, value := range line {
			if value == Hole {
				line[col] = Obstacle
			}
		}
	}
}

func (game *Game) moveToBottom() {
	field := game.Field
	for col := 0; col < util.ColCount; col++ {
		holeRow := util.RowCount - 1
		for curRow := util.RowCount - 1; 0 <= curRow; curRow-- {
			if field[curRow][col] == Hole {
				continue
			}
			for ; curRow < holeRow; holeRow-- {
				if field[holeRow][col] == Hole {
					util.Swap(field, holeRow, col, curRow, col)
					break
				}
			}
		}
	}
}

func (game *Game) moveToLeft() {
	for _, line := range game.Field {
		holeCol := 0
		for curCol := 0; curCol < util.ColCount; curCol++ {
			if line[curCol] == Hole {
				continue
			}
			for ; holeCol < curCol; holeCol++ {
				if line[holeCol] == Hole {
					line[holeCol], line[curCol] = line[curCol], line[holeCol]
					break
				}
			}
		}
	}
}

func (game *Game) moveToTop() {
	field := game.Field
	for col := 0; col < util.ColCount; col++ {
		holeRow := 0
		for curRow := 0; curRow < util.RowCount; curRow++ {
			if field[curRow][col] == Hole {
				continue
			}
			for ; holeRow < curRow; holeRow++ {
				if field[holeRow][col] == Hole {
					util.Swap(field, curRow, col, holeRow, col)
					break
				}
			}
		}
	}
}

func (game *Game) moveToRight() {
	for _, line := range game.Field {
		holeCol := util.ColCount - 1
		for curCol := util.ColCount - 1; 0 <= curCol; curCol-- {
			if line[curCol] == Hole {
				continue
			}
			for ; curCol < holeCol; holeCol-- {
				if line[holeCol] == Hole {
					line[curCol], line[holeCol] = line[holeCol], line[curCol]
					break
				}
			}
		}
	}
}
