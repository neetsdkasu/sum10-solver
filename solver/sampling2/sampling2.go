// SUM10-SOLVER
// author: Leonardone @ NEETSDKASU

package sampling2

import (
	"context"
	"fmt"
	"github.com/neetsdkasu/sum10-solver/game"
	"github.com/neetsdkasu/sum10-solver/problem"
	"github.com/neetsdkasu/sum10-solver/search"
	"github.com/neetsdkasu/sum10-solver/solver"
	"log"
	"math/rand"
	"sort"
	"time"
)

type Sampling2 struct {
	NumOfTarget     int
	NumOfPlaying    int
	LimitNumOfState int
}

func init() {
	solver.Register(&Sampling2{10, 10, 100})
	solver.Register(&Sampling2{20, 20, 100})
	solver.Register(&Sampling2{30, 30, 100})
	solver.Register(&Sampling2{200, 30, 600})
	solver.Register(&Sampling2{60, 60, 600})

	// 微妙、多様性(？)を確保できてない･･･
}

func (self *Sampling2) Name() string {
	return fmt.Sprint(
		"Sampling2-",
		self.NumOfTarget,
		"-",
		self.NumOfPlaying,
		"-",
		self.LimitNumOfState,
	)
}

func (self *Sampling2) Description() string {
	return fmt.Sprintf(
		"選択肢から最大%d個選び、それぞれ%d個のランダム解を生成し、平均スコアを計算し、全ての候補のなかから平均スコアの高かったものを１つ選び、他は候補として残す(平均スコア上位から最大%d個)。このプロセスの中のランダム解で出た最大スコアのものを解とする",
		self.NumOfTarget,
		self.NumOfPlaying,
		self.LimitNumOfState,
	)
}

func (self *Sampling2) Search(startTime time.Time, runningSeconds int, prob *problem.Problem) (solver.Solution, error) {
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

type State struct {
	Game      *game.Game
	Sum       int
	FirstTime bool
}

func (self *Sampling2) run(ctx context.Context, prob *problem.Problem) <-chan solver.Solution {
	ch := make(chan solver.Solution)

	numOfTarget := self.NumOfTarget
	numOfPlaying := self.NumOfPlaying
	limitNumOfState := self.LimitNumOfState

	go func() {
		defer close(ch)
		rng := rand.New(rand.NewSource(time.Now().Unix()))

		cur := game.New(prob)
		best := cur

		idleList := make([]*State, 0, limitNumOfState+200)
		idleList = append(idleList, &State{cur, 0, false})

		additionalList := make([]*State, 0, numOfTarget)

		for len(idleList) > 0 {
			select {
			case <-ctx.Done():
				return
			default:
			}

			size := len(idleList) - 1
			curState := idleList[size]
			idleList = idleList[:size]

			cur = curState.Game
			list := search.Search(cur)
			if len(list) == 0 {
				continue
			}

			if len(list) > numOfTarget {
				rng.Shuffle(len(list), func(i, j int) {
					list[i], list[j] = list[j], list[i]
				})
				list = list[:numOfTarget]
			}

			additionalList = additionalList[:0]

			for _, sel := range list {
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

				// 多様性(？)確保のため(？)
				// 深い階層での評価を雑に下げる（深い層では平均スコアが高くなりそうだから？）
				// 下げ方が無根拠なのでよろしくない
				sum -= cur.Steps * numOfPlaying

				state := &State{tmp, sum, true}
				additionalList = append(additionalList, state)
			}

			sort.Slice(additionalList, func(a, b int) bool {
				return additionalList[a].Sum < additionalList[b].Sum
			})

			if additionalList[0].Sum == additionalList[len(additionalList)-1].Sum {
				// 選択肢に幅がないと同スコアが並びそうという勝手な妄想から…
				// 浅い階層でここに来ることは少なそう（願望）
				// 中間の階層で運悪くここに突入することはありえるかも？（分からん…）
				// 変化する可能性の低い高スコアStateが並ぶとidleList内の多様性を破壊するので…（？）
				if curState.FirstTime {
					curState.FirstTime = false
					idleList = append(idleList, curState)
				}
				continue
			}

			idleList = append(idleList, additionalList...)
			sort.Slice(idleList, func(a, b int) bool {
				diff := idleList[a].Sum - idleList[b].Sum
				if diff == 0 {
					// 同じスコアが発生するかはわからんが（分析してないので）
					// 多様性（？）確保のため、浅い階層を優先する
					return idleList[a].Game.Steps < idleList[b].Game.Steps
				}
				return diff < 0
			})

			if len(idleList) > limitNumOfState {
				pos := len(idleList) - limitNumOfState
				copy(idleList[:limitNumOfState], idleList[pos:])
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
