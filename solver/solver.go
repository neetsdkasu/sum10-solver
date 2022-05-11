package solver

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"sort"
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

func Register(solver Solver) {
	name := solver.Name()
	if _, ok := uniqueSolverName[name]; ok {
		log.Println("Duplicate solver name: ", name)
		return
	}
	solvers = append(solvers, solver)
	uniqueSolverName[name] = true
}

func Comp(file io.Writer, runningSeconds, numOfTestcase, seed int) (err error) {
	log.Println("Running Comp-Mode")
	log.Println("  limit running seconds:", runningSeconds, "sec.")
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
			sol := process(runningSeconds, prob, solver)
			result := sol.Replay(prob)
			if result != nil {
				log.Println("  ok.")
			}
		}
	}

	if _, err = fmt.Fprintln(file, "----------------------------------------------"); err != nil {
		return
	}

	return
}

func process(runningSeconds int, prob *problem.Problem, solver Solver) Solution {
	const Margin = 10 * time.Millisecond
	startTime, ch := run(runningSeconds, prob, solver)
	deadline := startTime.Add(Margin + time.Duration(int64(runningSeconds))*time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()
	select {
	case <-ctx.Done():
		return nil
	case sol, ok := <-ch:
		if ok {
			return sol
		} else {
			return nil
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
