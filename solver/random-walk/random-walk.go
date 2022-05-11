package randomWalk

import (
	"context"
	"math/rand"
	"sum10-solver/game"
	"sum10-solver/problem"
	"sum10-solver/search"
	"sum10-solver/solver"
	"time"
)

type RandomWalk struct{}

var randomWalk = &RandomWalk{}

func init() {
	solver.Register(randomWalk)
}

func (*RandomWalk) Name() string {
	return "RandomWalk"
}

func (*RandomWalk) Description() string {
	return "ランダムな解を時間いっぱい生成して一番スコアがよいものを選ぶ"
}

func (*RandomWalk) Search(startTime time.Time, runningSeconds int, problem *problem.Problem) (solver.Solution, error) {
	deadline := startTime.Add(time.Duration(int64(runningSeconds)) * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()
	var best = solver.Solution{}
	ch := run(ctx, problem)
	for {
		select {
		case <-ctx.Done():
			return best, nil
		case sol, ok := <-ch:
			if ok {
				best = sol
			} else {
				return best, nil
			}
		}
	}
}

func run(ctx context.Context, prob *problem.Problem) <-chan solver.Solution {
	ch := make(chan solver.Solution, 100)

	game0 := game.New(prob)

	go func() {
		defer close(ch)

		rng := rand.New(rand.NewSource(time.Now().Unix()))

		best := game0

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			var err error = nil
			cur := game0
			for err == nil {
				list := search.Search(cur)
				if len(list) == 0 {
					break
				}
				sel := rng.Intn(len(list))
				cur, err = cur.Take(list[sel])
			}
			if err != nil {
				continue
			}
			if cur.Score <= best.Score {
				continue
			}
			best = cur
			sol := solver.ToSolution(best)
			select {
			case <-ctx.Done():
				return
			case ch <- sol:
			}
		}
	}()

	return ch
}
