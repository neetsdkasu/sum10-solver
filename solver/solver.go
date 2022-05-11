package solver

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"sort"
	"sum10-solver/game"
	"sum10-solver/problem"
	"sum10-solver/util"
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

func (res *Result) IsError() bool {
	return res.Duration == 0
}

func (res *Result) IsTimeout() bool {
	return res.Duration < 0
}

func (res *Result) IsValid() bool {
	return res.Game != nil && res.Duration > 0
}

func Comp(file io.Writer, runningSeconds, numOfTestcase, seed int) (err error) {
	log.Println("Running Comp-Mode")
	log.Println("  running limit:", runningSeconds, "sec.")
	log.Println("  number of testcase:", numOfTestcase)

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
	if _, err = fmt.Fprintln(file, "----------------------------------------------"); err != nil {
		return
	}

	results := make([][]*Result, len(solvers))
	for i := range results {
		results[i] = make([]*Result, numOfTestcase)
	}

	for i, solver := range solvers {
		if _, err = fmt.Fprintf(file, "Entry No. %3d\n", i); err != nil {
			return
		}
		if _, err = fmt.Fprintln(file, " ", solver.Name()); err != nil {
			return
		}
		if _, err = fmt.Fprintln(file, " ", solver.Description()); err != nil {
			return
		}

		log.Printf("process: No. %3d %s", i, solver.Name())

		for k, prob := range problemList {
			log.Printf("  [%3d/%3d] Seed: %5d", k+1, numOfTestcase, prob.Seed())
			sol, dur := process(runningSeconds, prob, solver)
			results[i][k] = &Result{
				Solution: sol,
				Game:     sol.Replay(prob),
				Duration: dur,
			}
		}
	}

	if _, err = fmt.Fprintln(file, "----------------------------------------------"); err != nil {
		return
	}

	/* * * * * * * * * * * * * * * * * * * * */

	if _, err = fmt.Fprint(file, "                 ENTRY NO:"); err != nil {
		return
	}
	for i := range solvers {
		if _, err = fmt.Fprintf(file, "  %3d", i); err != nil {
			return
		}
	}
	if _, err = fmt.Fprintln(file); err != nil {
		return
	}

	/* * * * * * * * * * * * * * * * * * * * */

	if _, err = fmt.Fprintln(file, "==============================================================================="); err != nil {
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

	if _, err = fmt.Fprintln(file, "==============================================================================="); err != nil {
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
		select {
		case ch <- sol:
		default:
		}
	}()
	return startTime, ch
}
