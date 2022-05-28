// SUM10-SOLVER
// author: Leonardone @ NEETSDKASU

package solver

import (
	"context"
	"fmt"
	"github.com/neetsdkasu/sum10-solver/game"
	"github.com/neetsdkasu/sum10-solver/problem"
	"github.com/neetsdkasu/sum10-solver/show"
	"github.com/neetsdkasu/sum10-solver/util"
	"io"
	"log"
	"math/rand"
	"sort"
	"time"
)

type Solver interface {
	Name() string
	Description() string
	Search(startTime time.Time, runningSeconds int, prob *problem.Problem) (Solution, error)
}

var (
	solvers          = []Solver{}
	uniqueSolverName = make(map[string]bool)
)

func Count() int {
	return len(solvers)
}

func SolverList() []string {
	list := make([]string, len(solvers))
	for i, solver := range solvers {
		list[i] = solver.Name()
	}
	return list
}

func FindSolver(name string) (Solver, bool) {
	if _, ok := uniqueSolverName[name]; !ok {
		return nil, false
	}
	for _, solver := range solvers {
		if name == solver.Name() {
			return solver, true
		}
	}
	return nil, false
}

func Register(solver Solver) {
	name := solver.Name()
	if _, ok := uniqueSolverName[name]; ok {
		log.Println("Duplicate solver name: ", name)
		return
	}
	solvers = append(solvers, solver)
	uniqueSolverName[name] = true
}

type Result struct {
	Solution Solution
	Game     *game.Game
	Duration time.Duration
}

func newResult(prob *problem.Problem, sol Solution, dur time.Duration) *Result {
	return &Result{
		Solution: sol,
		Game:     sol.Replay(prob),
		Duration: dur,
	}
}

func (res *Result) IsError() bool {
	return res.Duration == 0
}

func (res *Result) IsTimeout() bool {
	return res.Duration < 0
}

func (res *Result) IsValid() bool {
	return res.Game != nil && res.Duration > 0
}

func Solve(file io.Writer, seed int, name string, runningSeconds, tryCount int) (err error) {
	log.Println("Running Solver-Mode")
	log.Println("  seed          :", seed)
	log.Println("  solver name   :", name)
	log.Println("  running limit :", runningSeconds, "sec.")
	log.Println("  try count     :", tryCount)

	prob := problem.New(uint32(seed))
	solver, _ := FindSolver(name)

	var best *Result = nil

	for i := 1; i <= tryCount; i++ {
		log.Println("process [", i, "/", tryCount, "]")
		sol, dur := process(runningSeconds, prob, solver)
		res := newResult(prob, sol, dur)
		if res.IsValid() {
			score := res.Game.Score
			if best == nil || best.Game.Score < score {
				best = res
				if tryCount > 1 {
					log.Println("  SUCCESS: score", score, "(UPDATE BEST)")
				} else {
					log.Println("  SUCCESS: score", score)
				}
			} else {
				log.Println("  SUCCESS: score", score)
			}
		} else if res.IsTimeout() {
			log.Println("  FAILURE: time out")
		} else {
			log.Println("  FAILURE: error")
		}
	}

	if best != nil {
		err = showBestSolution(file, best)
		return
	}

	log.Println("NO SOLUTION FOUND")

	if _, err = fmt.Fprintln(file, "SEED:", prob.Seed()); err != nil {
		return
	}

	if _, err = fmt.Fprintln(file, "----------------------------------------------"); err != nil {
		return
	}

	if err = show.ShowField(file, prob); err != nil {
		return
	}

	if _, err = fmt.Fprintln(file, "----------------------------------------------"); err != nil {
		return
	}

	_, err = fmt.Fprintln(file, "解を見つけることができませんでした")

	return
}

func Comp(file io.Writer, runningSeconds, numOfTestcase, seed int) (err error) {
	log.Println("Running Comp-Mode")
	log.Println("  running limit      :", runningSeconds, "sec.")
	log.Println("  number of testcase :", numOfTestcase)
	log.Println("  number of solver   :", len(solvers))

	problemList := make([]*problem.Problem, numOfTestcase)
	if util.IsValidSeed(seed) {
		log.Println("  testcases use SINGLE seed (", seed, ")")
		prob := problem.New(uint32(seed))
		for i := range problemList {
			problemList[i] = prob
		}
	} else {
		log.Println("  testcases use RANDOM seed")
		rand.Seed(time.Now().Unix())
		set := make(map[int]bool)
		for i := range problemList {
			for {
				seed = rand.Intn(util.MaxSeed - util.MinSeed + 1)
				if _, ok := set[seed]; !ok {
					set[seed] = true
					break
				}
			}
			problemList[i] = problem.New(uint32(seed))
		}
		sort.Slice(problemList, func(i, j int) bool {
			return problemList[i].Seed() < problemList[j].Seed()
		})
	}

	if _, err = fmt.Fprintln(file, "RUNNING LIMIT:", runningSeconds, "sec."); err != nil {
		return
	}
	if _, err = fmt.Fprintln(file, "--------------------------------------------------------------------------------------------"); err != nil {
		return
	}

	results := make([][]*Result, len(solvers))
	for i := range results {
		results[i] = make([]*Result, numOfTestcase)
	}

	for i, solver := range solvers {
		if _, err = fmt.Fprintf(file, "Entry No. %3d\n", i+1); err != nil {
			return
		}
		if _, err = fmt.Fprintln(file, " ", solver.Name()); err != nil {
			return
		}
		if _, err = fmt.Fprintln(file, " ", solver.Description()); err != nil {
			return
		}

		log.Printf("process: No. %3d %s", i+1, solver.Name())

		for k, prob := range problemList {
			log.Printf("  [%3d/%3d] Seed: %5d", k+1, numOfTestcase, prob.Seed())
			sol, dur := process(runningSeconds, prob, solver)
			results[i][k] = newResult(prob, sol, dur)
		}
	}

	/* * * * * * * * * * * * * * * * * * * * */

	if _, err = fmt.Fprint(file, "--------------------------"); err != nil {
		return
	}
	for _ = range solvers {
		if _, err = fmt.Fprint(file, "-----"); err != nil {
			return
		}
	}
	if _, err = fmt.Fprintln(file); err != nil {
		return
	}

	/* * * * * * * * * * * * * * * * * * * * */

	if _, err = fmt.Fprint(file, "                 ENTRY NO:"); err != nil {
		return
	}
	for i := range solvers {
		if _, err = fmt.Fprintf(file, "  %3d", i+1); err != nil {
			return
		}
	}
	if _, err = fmt.Fprintln(file); err != nil {
		return
	}

	/* * * * * * * * * * * * * * * * * * * * */

	if _, err = fmt.Fprint(file, "=========================="); err != nil {
		return
	}
	for _ = range solvers {
		if _, err = fmt.Fprint(file, "====="); err != nil {
			return
		}
	}
	if _, err = fmt.Fprintln(file); err != nil {
		return
	}

	/* * * * * * * * * * * * * * * * * * * * */

	for k, prob := range problemList {
		best := 0
		for _, result := range results {
			if result[k].IsValid() {
				best = util.Max(best, result[k].Game.Score)
			}
		}
		if _, err = fmt.Fprintf(file, "[%3d] SEED %5d, MAX %3d:", k+1, prob.Seed(), best); err != nil {
			return
		}
		for _, result := range results {
			if result[k].IsError() {
				if _, err = fmt.Fprint(file, "  ERR"); err != nil {
					return
				}
			} else if result[k].IsTimeout() {
				if _, err = fmt.Fprint(file, "  ---"); err != nil {
					return
				}
			} else if result[k].IsValid() {
				score := result[k].Game.Score
				if _, err = fmt.Fprintf(file, "  %3d", score); err != nil {
					return
				}
			} else {
				// 解が正しくない （合計が10にならない、連結でない、など）
				if _, err = fmt.Fprint(file, "  BAD"); err != nil {
					return
				}
			}
		}
		if _, err = fmt.Fprintln(file); err != nil {
			return
		}
	}

	/* * * * * * * * * * * * * * * * * * * * */

	if _, err = fmt.Fprint(file, "=========================="); err != nil {
		return
	}
	for _ = range solvers {
		if _, err = fmt.Fprint(file, "====="); err != nil {
			return
		}
	}
	if _, err = fmt.Fprintln(file); err != nil {
		return
	}

	/* * * * * * * * * * * * * * * * * * * * */

	if _, err = fmt.Fprint(file, "                MIN SCORE:"); err != nil {
		return
	}
	for i := range solvers {
		min := 999
		for _, result := range results[i] {
			if result.IsValid() {
				min = util.Min(min, result.Game.Score)
			}
		}
		if min == 999 {
			if _, err = fmt.Fprint(file, "  ---"); err != nil {
				return
			}
		} else {
			if _, err = fmt.Fprintf(file, "  %3d", min); err != nil {
				return
			}
		}
	}
	if _, err = fmt.Fprintln(file); err != nil {
		return
	}

	/* * * * * * * * * * * * * * * * * * * * */

	if _, err = fmt.Fprint(file, "                MAX SCORE:"); err != nil {
		return
	}
	for i := range solvers {
		max := -1
		for _, result := range results[i] {
			if result.IsValid() {
				max = util.Max(max, result.Game.Score)
			}
		}
		if max == -1 {
			if _, err = fmt.Fprint(file, "  ---"); err != nil {
				return
			}
		} else {
			if _, err = fmt.Fprintf(file, "  %3d", max); err != nil {
				return
			}
		}
	}
	if _, err = fmt.Fprintln(file); err != nil {
		return
	}

	/* * * * * * * * * * * * * * * * * * * * */

	if _, err = fmt.Fprint(file, "            AVERAGE SCORE:"); err != nil {
		return
	}
	for i := range solvers {
		total := 0
		count := 0
		for _, result := range results[i] {
			if result.IsValid() {
				total += result.Game.Score
				count++
			}
		}
		if count == 0 {
			if _, err = fmt.Fprint(file, "  ---"); err != nil {
				return
			}
		} else {
			average := total / count
			if _, err = fmt.Fprintf(file, "  %3d", average); err != nil {
				return
			}
		}
	}
	if _, err = fmt.Fprintln(file); err != nil {
		return
	}

	/* * * * * * * * * * * * * * * * * * * * */

	if _, err = fmt.Fprint(file, "           *AVERAGE SCORE:"); err != nil {
		return
	}
	for i := range solvers {
		total := 0
		for _, result := range results[i] {
			if result.IsValid() {
				total += result.Game.Score
			}
		}
		average := total / numOfTestcase
		if _, err = fmt.Fprintf(file, "  %3d", average); err != nil {
			return
		}
	}
	if _, err = fmt.Fprintln(file); err != nil {
		return
	}

	/* * * * * * * * * * * * * * * * * * * * */

	if _, err = fmt.Fprint(file, "--------------------------"); err != nil {
		return
	}
	for _ = range solvers {
		if _, err = fmt.Fprint(file, "-----"); err != nil {
			return
		}
	}
	if _, err = fmt.Fprintln(file); err != nil {
		return
	}

	/* * * * * * * * * * * * * * * * * * * * */

	if problemList[0].Seed() == problemList[numOfTestcase-1].Seed() {
		var best *Result = nil
		for _, line := range results {
			for _, res := range line {
				if res == nil {
					continue
				}
				if res.Game == nil {
					continue
				}
				if best == nil || res.Game.Score > best.Game.Score {
					best = res
				}
			}
		}
		if best != nil {
			if err = showBestSolution(file, best); err != nil {
				return
			}
		}
	}

	return
}

func process(runningSeconds int, prob *problem.Problem, solver Solver) (Solution, time.Duration) {
	const Margin = 50 * time.Millisecond
	startTime, ch := run(runningSeconds, prob, solver)
	deadline := startTime.Add(Margin + time.Duration(int64(runningSeconds))*time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()
	select {
	case <-ctx.Done():
		log.Println("  timeout")
		go func() {
			_, _ = <-ch
		}()
		return nil, -1
	case sol, ok := <-ch:
		if ok {
			dur := time.Now().Sub(startTime)
			return sol, dur
		} else {
			return nil, 0
		}
	}
}

func run(runningSeconds int, prob *problem.Problem, solver Solver) (time.Time, <-chan Solution) {
	ch := make(chan Solution)
	startTime := time.Now()
	go func() {
		defer close(ch)
		sol, err := solver.Search(startTime, runningSeconds, prob)
		if err != nil {
			log.Println(err)
			return
		}
		ch <- sol
	}()
	return startTime, ch
}

func showBestSolution(file io.Writer, best *Result) (err error) {

	prob := best.Game.Problem

	if _, err = fmt.Fprintln(file, "SEED:", prob.Seed()); err != nil {
		return
	}

	if _, err = fmt.Fprintln(file, "----------------------------------------------"); err != nil {
		return
	}

	if err = show.ShowField(file, prob); err != nil {
		return
	}

	if _, err = fmt.Fprintln(file, "----------------------------------------------"); err != nil {
		return
	}

	if _, err = fmt.Fprintln(file, "BEST SOLUTION (SCORE", best.Game.Score, ")"); err != nil {
		return
	}

	if _, err = fmt.Fprintln(file, "----------------------------------------------"); err != nil {
		return
	}

	steps := make([]*game.Game, best.Game.Steps)

	cur := best.Game
	for cur.Steps > 0 {
		steps[cur.Prev.Steps] = cur
		cur = cur.Prev
	}

	for _, step := range steps {
		if err = show.ShowGameWithMark(file, step.Prev, step.Taked); err != nil {
			return
		}
		if _, err = fmt.Fprintln(file, "----------------------------------------------"); err != nil {
			return
		}
	}

	if err = show.ShowGame(file, best.Game); err != nil {
		return
	}

	if _, err = fmt.Fprintln(file, "----------------------------------------------"); err != nil {
		return
	}

	return
}
