package solver

import (
	"io"
	"log"
	"sum10-solver/marker"
	"sum10-solver/problem"
	"sum10-solver/util"
	"time"
)

type Solver interface {
	Name() string
	Description() string
	Search(startTime time.Time, runningSeconds int, prob *problem.Problem) []*marker.Marker
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

func Comp(file io.Writer, runningSeconds, numOfTestcase, seed int) error {
	if util.IsValidSeed(seed) {
		return compSpecialSeed(file, runningSeconds, numOfTestcase, seed)
	} else {
		return compRandomTestcases(file, runningSeconds, runningSeconds)
	}
}

/*

    // 仮
	prob := problem.New(5531)
	for _, solver := range solvers {
		log.Println("Solver:", solver.Name())
		log.Println("Description:", solver.Description())
		list := solver.Search(time.Now(), runningSeconds, prob)
		cur := game.New(prob)
		for _, step := range list {
			var err error
			cur, err = cur.Take(step)
			if err != nil {
				break
			}
		}
		if cur != nil {
			log.Println("Get Score:", cur.Score)
		}
		<-time.After(time.Second)
	}
	return nil

*/
