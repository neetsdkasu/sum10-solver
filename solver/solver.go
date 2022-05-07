package solver

import (
	"context"
	"io"
	"log"
	"sum10-solver/marker"
	"sum10-solver/problem"
)

type Solver interface {
	Name() string
	Description() string
	Search(ctx context.Context, problem *problem.Problem) []*marker.Marker
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

func Comp(file io.Writer, runningLimitTime, numOfTestcase int) error {

	return nil
}
