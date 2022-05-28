// SUM10-SOLVER
// author: Leonardone @ NEETSDKASU

package near_edge

import (
	"github.com/neetsdkasu/sum10-solver/game"
	"github.com/neetsdkasu/sum10-solver/marker"
	"github.com/neetsdkasu/sum10-solver/problem"
	"github.com/neetsdkasu/sum10-solver/search"
	"github.com/neetsdkasu/sum10-solver/solver"
	"github.com/neetsdkasu/sum10-solver/util"
	"time"
)

type NearEdge struct{}

func init() {
	solver.Register(NearEdge{})
}

func (NearEdge) Name() string {
	return "NearEdge"
}

func (NearEdge) Description() string {
	return "詰める側の辺に近そうなものを選ぶ"
}

func (NearEdge) Search(startTime time.Time, runningSeconds int, prob *problem.Problem) (solver.Solution, error) {

	cur := game.New(prob)

	for {
		list := search.Search(cur)
		if len(list) == 0 {
			break
		}
		var sel *marker.Marker
		switch cur.Steps % 4 {
		case 0:
			sel = findNearTop(list)
		case 1:
			sel = findNearRight(list)
		case 2:
			sel = findNearBottom(list)
		case 3:
			sel = findNearLeft(list)
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

func findNearTop(list []*marker.Marker) *marker.Marker {
	return findNear(list, func(row, col int) int {
		return row
	})
}

func findNearRight(list []*marker.Marker) *marker.Marker {
	return findNear(list, func(row, col int) int {
		return util.ColCount - col
	})
}

func findNearBottom(list []*marker.Marker) *marker.Marker {
	return findNear(list, func(row, col int) int {
		return util.RowCount - row
	})
}

func findNearLeft(list []*marker.Marker) *marker.Marker {
	return findNear(list, func(row, col int) int {
		return col
	})
}

func findNear(list []*marker.Marker, calcDist func(row, col int) int) *marker.Marker {
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
