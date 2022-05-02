package game

import (
	"errors"
	"sum10-solver/marker"
	"sum10-solver/problem"
	"sum10-solver/util"
)

const (
	Sum      = 10
	Hole     = -1
	Obstacle = 11
)

type Game struct {
	Problem *problem.Problem
	Steps   int
	Score   int
	Field   [][]int
	Prev    *Game
	Taked   *marker.Marker
}

func New(problem *problem.Problem) *Game {
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
		Prev:    nil,
		Taked:   nil,
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

func (game *Game) Take(marker *marker.Marker) (next *Game, err error) {
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
	next.Prev = game
	next.Taked = marker.GetCopy()
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

func (game *Game) IsGameOver() bool {
	return !game.findSum10OnHorizontalLine() &&
		!game.findSum10OnVerticalLine() &&
		!game.findSum10()
}

func (game *Game) findSum10OnHorizontalLine() bool {
	for _, line := range game.Field {
		head := 0
		tail := 0
		sum := 0
		for {
			for sum < Sum && head < util.ColCount {
				sum += line[head]
				head++
			}
			if sum < Sum {
				break
			}
			if sum == Sum {
				return true
			}
			for sum >= Sum && tail < head {
				sum -= line[tail]
				tail++
			}
		}
	}
	return false
}

func (game *Game) findSum10OnVerticalLine() bool {
	field := game.Field
	for col := 0; col < util.ColCount; col++ {
		head := 0
		tail := 0
		sum := 0
		for {
			for sum < Sum && head < util.RowCount {
				sum += field[head][col]
				head++
			}
			if sum < Sum {
				break
			}
			if sum == Sum {
				return true
			}
			for sum >= Sum && tail < head {
				sum -= field[tail][col]
				tail++
			}
		}
	}
	return false
}

func (game *Game) findSum10() bool {
	const HalfValue = Sum / 2
	marker := marker.New()
	for row, line := range game.Field {
		for col, value := range line {
			if value <= HalfValue {
				if game.findSum10ByDfs(marker, row, col, 0) {
					return true
				}
			}
		}
	}
	return false
}

func (game *Game) findSum10ByDfs(marker *marker.Marker, row, col, sum int) bool {
	if !util.FieldContains(row, col) {
		return false
	}
	if marker.Has(row, col) {
		return false
	}
	sum += game.Field[row][col]
	if sum >= Sum {
		return sum == Sum
	}
	marker.Set(row, col)
	defer marker.Unset(row, col)
	return game.findSum10ByDfs(marker, row+1, col, sum) ||
		game.findSum10ByDfs(marker, row-1, col, sum) ||
		game.findSum10ByDfs(marker, row, col+1, sum) ||
		game.findSum10ByDfs(marker, row, col-1, sum)
}
