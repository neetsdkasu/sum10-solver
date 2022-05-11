package solver

import (
	"sum10-solver/game"
	"sum10-solver/marker"
	"sum10-solver/problem"
)

type Solution []*marker.Marker

func ToSolution(g *game.Game) (sol Solution) {
	if g == nil {
		return
	}
	list := make([]*marker.Marker, g.Steps)
	for g.Steps > 0 {
		list[g.Prev.Steps] = g.Taked
		g = g.Prev
	}
	sol = Solution(list)
	return
}

func (sol Solution) Replay(prob *problem.Problem) *game.Game {
	cur := game.New(prob)
	if cur == nil {
		return nil
	}
	for _, sel := range []*marker.Marker(sol) {
		var err error
		cur, err = cur.Take(sel)
		if err != nil {
			return nil
		}
	}
	return cur
}
