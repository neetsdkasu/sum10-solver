// SUM10-SOLVER
// author: Leonardone @ NEETSDKASU

package far_edge

import (
	"github.com/neetsdkasu/sum10-solver/game"
	"github.com/neetsdkasu/sum10-solver/marker"
	"github.com/neetsdkasu/sum10-solver/problem"
	"github.com/neetsdkasu/sum10-solver/search"
	"github.com/neetsdkasu/sum10-solver/solver"
	"github.com/neetsdkasu/sum10-solver/util"
	"time"
)

type FarEdge struct{}

func init() {
	solver.Register(FarEdge{})
}

func (FarEdge) Name() string {
	return "FarEdge"
}

func (FarEdge) Description() string {
	return "詰める側の辺から遠そうなものを選ぶ"
}

func (FarEdge) Search(startTime time.Time, runningSeconds int, prob *problem.Problem) (solver.Solution, error) {

	cur := game.New(prob)

	for {
		list := search.Search(cur)
		if len(list) == 0 {
			break
		}
		var sel *marker.Marker
		switch cur.Steps % 4 {
		case 0:
			sel = findFarTop(list)
		case 1:
			sel = findFarRight(list)
		case 2:
			sel = findFarBottom(list)
		case 3:
			sel = findFarLeft(list)
		}
		if tmp, err := cur.Take(sel); err == nil {
			cur = tmp
		} else {
			return nil, err
		}
	}

	sol := solver.ToSolution(cur)
	return sol, nil
}

func findFarTop(list []*marker.Marker) *marker.Marker {
	return findFar(list, func(row, col int) int {
		return util.RowCount - row
	})
}

func findFarRight(list []*marker.Marker) *marker.Marker {
	return findFar(list, func(row, col int) int {
		return col
	})
}

func findFarBottom(list []*marker.Marker) *marker.Marker {
	return findFar(list, func(row, col int) int {
		return row
	})
}

func findFarLeft(list []*marker.Marker) *marker.Marker {
	return findFar(list, func(row, col int) int {
		return util.ColCount - col
	})
}

func findFar(list []*marker.Marker, calcDist func(row, col int) int) *marker.Marker {
	var best *marker.Marker
	var bestValue = 99999.0
	for _, sel := range list {
		count := 0
		sum := 0
		for row := 0; row < util.RowCount; row++ {
			for col := 0; col < util.ColCount; col++ {
				if sel.Has(row, col) {
					count++
					sum += calcDist(row, col)
				}
			}
		}
		if count > 0 {
			tmp := float64(sum) / float64(count)
			if tmp < bestValue {
				best = sel
				bestValue = tmp
			}
		}
	}
	return best
}
