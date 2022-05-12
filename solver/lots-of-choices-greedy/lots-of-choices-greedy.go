package lots_of_choices_greedy

import (
	"context"
	"fmt"
	"github.com/neetsdkasu/sum10-solver/game"
	"github.com/neetsdkasu/sum10-solver/marker"
	"github.com/neetsdkasu/sum10-solver/problem"
	"github.com/neetsdkasu/sum10-solver/search"
	"github.com/neetsdkasu/sum10-solver/solver"
	"time"
)

type LotsOfChoicesGreedy struct {
	depth int
}

func init() {
	solver.Register(&LotsOfChoicesGreedy{1})
	solver.Register(&LotsOfChoicesGreedy{2})
}

func (this *LotsOfChoicesGreedy) Name() string {
	return fmt.Sprint("LotsOfChoicesGreedy-", this.depth)
}

func (this *LotsOfChoicesGreedy) Description() string {
	if this.depth == 1 {
		return "次の手で選択肢が一番多くなる手を選ぶ (1テストケースあたり数秒以上かけないと解が求まらないことがある)"
	} else {
		return "次の次の手で選択肢が一番多くなる手を選ぶ (1テストケースあたり2分以上かけないと解が求まらない…)"
	}
}

func (this *LotsOfChoicesGreedy) Search(startTime time.Time, runningSeconds int, prob *problem.Problem) (solver.Solution, error) {
	depth1 := this.depth == 1

	deadline := startTime.Add(time.Duration(int64(runningSeconds))*time.Second - 20*time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	cur := game.New(prob)

Loop:
	for {
		list := search.Search(cur)
		if len(list) == 0 {
			break
		}
		var ch <-chan *game.Game
		if depth1 {
			ch = depth1Select(cur, list)
		} else {
			ch = depth2Select(cur, list)
		}
		select {
		case <-ctx.Done():
			break Loop
		case tmp, ok := <-ch:
			if ok && tmp != nil {
				cur = tmp
			} else {
				break Loop
			}
		}
	}

	sol := solver.ToSolution(cur)
	return sol, nil
}

func depth1Select(cur *game.Game, list []*marker.Marker) <-chan *game.Game {
	ch := make(chan *game.Game)

	go func() {
		defer close(ch)

		var bestNext *game.Game = nil
		bestCount := -1

		for _, sel := range list {
			tmpNext, err := cur.Take(sel)
			if err != nil {
				continue
			}
			choicesCount := len(search.Search(tmpNext))
			if choicesCount > bestCount {
				bestCount = choicesCount
				bestNext = tmpNext
			}
		}

		select {
		case ch <- bestNext:
		default:
		}
	}()

	return ch
}

func depth2Select(cur *game.Game, list []*marker.Marker) <-chan *game.Game {
	ch := make(chan *game.Game)

	go func() {
		defer close(ch)

		var bestNext1 *game.Game = nil
		var bestNext2 *game.Game = nil
		bestCount1 := -1
		bestCount2 := -1

		for _, sel := range list {
			tmpNext1, err := cur.Take(sel)
			if err != nil {
				continue
			}
			list1 := search.Search(tmpNext1)
			choicesCount1 := len(list1)
			if choicesCount1 > bestCount1 {
				bestCount1 = choicesCount1
				bestNext1 = tmpNext1
			}
			if choicesCount1 == 0 {
				continue
			}
			for _, sel1 := range list1 {
				tmpNext2, err := tmpNext1.Take(sel1)
				if err != nil {
					continue
				}
				choicesCount2 := len(search.Search(tmpNext2))
				if choicesCount2 > bestCount2 {
					bestCount2 = choicesCount2
					bestNext2 = tmpNext1
				}
			}
		}

		if bestNext2 != nil {
			select {
			case ch <- bestNext2:
			default:
			}
		} else if bestNext1 != nil {
			select {
			case ch <- bestNext1:
			default:
			}
		}
	}()

	return ch
}
