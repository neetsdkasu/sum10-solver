// SUM10-SOLVER
// author: Leonardone @ NEETSDKASU

package sampling

import (
	"context"
	"fmt"
	"github.com/neetsdkasu/sum10-solver/game"
	"github.com/neetsdkasu/sum10-solver/problem"
	"github.com/neetsdkasu/sum10-solver/search"
	"github.com/neetsdkasu/sum10-solver/solver"
	"log"
	"math/rand"
	"time"
)

type Sampling struct {
	NumOfTarget  int
	NumOfPlaying int
}

func init() {
	solver.Register(&Sampling{3, 5})
	solver.Register(&Sampling{5, 5})
	solver.Register(&Sampling{5, 10})
	solver.Register(&Sampling{10, 10})
	solver.Register(&Sampling{15, 10})
	solver.Register(&Sampling{10, 15})
	solver.Register(&Sampling{15, 15})
	solver.Register(&Sampling{20, 20})
	solver.Register(&Sampling{30, 10})
	solver.Register(&Sampling{10, 30})
	solver.Register(&Sampling{30, 30})
	solver.Register(&Sampling{40, 40})
	solver.Register(&Sampling{200, 30})

	// 数を増やしていっても、途中の一度選んでしまった箇所に最終スコアが左右されるので
	// あまり、よくないソルバー
	// ただのランダムよりは平均的に高い得点にはなる
}

func (self *Sampling) Name() string {
	return fmt.Sprint("Sampling-", self.NumOfTarget, "-", self.NumOfPlaying)
}

func (self *Sampling) Description() string {
	return fmt.Sprintf(
		"選択肢から最大%d個選び、それぞれ%d個のランダム解を生成し、平均スコアの高かったものを選ぶ。このプロセスの中のランダム解で出た最大スコアのものを解とする",
		self.NumOfTarget,
		self.NumOfPlaying,
	)
}

func (self *Sampling) Search(startTime time.Time, runningSeconds int, prob *problem.Problem) (solver.Solution, error) {
	deadline := startTime.Add(time.Duration(int64(runningSeconds)) * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()
	var best solver.Solution
	ch := self.run(ctx, prob)
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

func (self *Sampling) run(ctx context.Context, prob *problem.Problem) <-chan solver.Solution {
	ch := make(chan solver.Solution)

	numOfTarget := self.NumOfTarget
	numOfPlaying := self.NumOfPlaying

	go func() {
		defer close(ch)
		rng := rand.New(rand.NewSource(time.Now().Unix()))

		cur := game.New(prob)
		best := cur

		for {
			list := search.Search(cur)
			if len(list) == 0 {
				break
			}

			if len(list) > numOfTarget {
				rng.Shuffle(len(list), func(i, j int) {
					list[i], list[j] = list[j], list[i]
				})
				list = list[:numOfTarget]
			}

			bestIndex := 0
			bestSum := -1

			for index, sel := range list {
				sum := 0
				tmp, err := cur.Take(sel)
				if err != nil {
					log.Println(err)
					return
				}
				for p := 0; p < numOfPlaying; p++ {
					res, err := play(rng, tmp)
					if err != nil {
						log.Println(err)
						return
					}
					if res.Score > best.Score {
						best = res
						sol := solver.ToSolution(best)
						select {
						case <-ctx.Done():
							return
						case ch <- sol:
						}
					} else {
						select {
						case <-ctx.Done():
							return
						default:
						}
					}
					sum += res.Score
				}
				if sum > bestSum {
					bestSum = sum
					bestIndex = index
				}
			}

			var err error
			cur, err = cur.Take(list[bestIndex])
			if err != nil {
				log.Println(err)
				return
			}
		}
	}()

	return ch
}

func play(rng *rand.Rand, cur *game.Game) (*game.Game, error) {
	var err error
	for {
		list := search.Search(cur)
		if len(list) == 0 {
			return cur, nil
		}
		sel := rng.Intn(len(list))
		cur, err = cur.Take(list[sel])
		if err != nil {
			return nil, err
		}
	}
}
