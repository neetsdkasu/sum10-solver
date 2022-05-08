package solver

import (
	"io"
	"log"
	"sum10-solver/marker"
	"sum10-solver/problem"
	"time"
)

type Solver interface {
	Name() string
	Description() string
	Search(startTime time.Time, runningSeconds int, problem *problem.Problem) []*marker.Marker
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

func Comp(file io.Writer, runningSeconds, numOfTestcase int) error {
	problem := problem.New(5531)
	for _, solver := range solvers {
		log.Println("Solver:", solver.Name())
		log.Println("Description:", solver.Description())
		_ = solver.Search(time.Now(), runningSeconds, problem)
	}
	return nil
}
