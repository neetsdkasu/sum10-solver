// SUM10-SOLVER
// author: Leonardone @ NEETSDKASU

package fewer_choices

import (
	"github.com/neetsdkasu/sum10-solver/game"
	"github.com/neetsdkasu/sum10-solver/problem"
	"github.com/neetsdkasu/sum10-solver/search"
	"github.com/neetsdkasu/sum10-solver/solver"
	"time"
)

type FewerChoices struct{}

var instance FewerChoices

func init() {
	solver.Register(instance)
}

func (FewerChoices) Name() string {
	return "FewerChoices"
}

func (FewerChoices) Description() string {
	return "次の手で選択肢が１つ以上ある手の中から一番少ない選択肢になる手を選ぶ"
}

func (FewerChoices) Search(startTime time.Time, runningSeconds int, prob *problem.Problem) (solver.Solution, error) {

	cur := game.New(prob)

	for {
		list := search.Search(cur)
		if len(list) == 0 {
			break
		}
		var bestNext *game.Game = nil
		bestCount := -1
		for _, sel := range list {
			tmpNext, err := cur.Take(sel)
			if err != nil {
				continue
			}
			choicesCount := len(search.Search(tmpNext))
			if 0 < choicesCount {
				if bestCount <= 0 || choicesCount < bestCount {
					bestCount = choicesCount
					bestNext = tmpNext
				}
			} else if bestCount < 0 {
				// 最終手の選択
				bestCount = choicesCount
				bestNext = tmpNext
			}
		}
		if bestNext == nil {
			break
		}
		cur = bestNext
	}

	sol := solver.ToSolution(cur)
	return sol, nil
}
