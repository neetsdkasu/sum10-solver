package first

import (
	"github.com/neetsdkasu/sum10-solver/game"
	"github.com/neetsdkasu/sum10-solver/problem"
	"github.com/neetsdkasu/sum10-solver/search"
	"github.com/neetsdkasu/sum10-solver/solver"
	"time"
)

type First struct{}

var instance First

func init() {
	solver.Register(instance)
}

func (First) Name() string {
	return "First"
}

func (First) Description() string {
	return "合計10になる選択のしかたで最初に見つかったものを選んでいくだけ"
}

func (First) Search(startTime time.Time, runningSeconds int, prob *problem.Problem) (solver.Solution, error) {
	var err error
	cur := game.New(prob)
	for {
		list := search.Search(cur)
		if len(list) == 0 {
			break
		}
		cur, err = cur.Take(list[0])
		if err != nil {
			return nil, err
		}
	}
	sol := solver.ToSolution(cur)
	return sol, nil
}
