package middle_choices_greedy

import (
	"github.com/neetsdkasu/sum10-solver/game"
	"github.com/neetsdkasu/sum10-solver/problem"
	"github.com/neetsdkasu/sum10-solver/search"
	"github.com/neetsdkasu/sum10-solver/solver"
	"sort"
	"time"
)

type MiddleChoicesGreedy struct{}

var instance MiddleChoicesGreedy

func init() {
	solver.Register(instance)
}

func (MiddleChoicesGreedy) Name() string {
	return "MiddleChoicesGreedy"
}

func (MiddleChoicesGreedy) Description() string {
	return "次の手で選択肢が１つ以上ある手の中からをその数の大きさ順で中央の位置になる手を選ぶ"
}

type Item struct {
	Game         *game.Game
	ChoicesCount int
}

func (MiddleChoicesGreedy) Search(startTime time.Time, runningSeconds int, prob *problem.Problem) (solver.Solution, error) {

	cur := game.New(prob)

	for {
		list := search.Search(cur)
		if len(list) == 0 {
			break
		}
		var zero *game.Game = nil
		items := []*Item{}
		for _, sel := range list {
			tmpNext, err := cur.Take(sel)
			if err != nil {
				continue
			}
			choicesCount := len(search.Search(tmpNext))
			if 0 < choicesCount {
				items = append(items, &Item{tmpNext, choicesCount})
			} else {
				zero = tmpNext
			}
		}
		if 0 < len(items) {
			sort.Slice(items, func(i, j int) bool {
				return items[i].ChoicesCount < items[j].ChoicesCount
			})
			mid := len(items) / 2
			cur = items[mid].Game
		} else if zero != nil {
			cur = zero
		} else {
			break
		}
	}

	sol := solver.ToSolution(cur)
	return sol, nil
}
