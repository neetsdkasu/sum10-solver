// SUM10-SOLVER
// author: Leonardone @ NEETSDKASU

package search

import (
	"github.com/neetsdkasu/sum10-solver/marker"
	"github.com/neetsdkasu/sum10-solver/util"
	"sort"
)

// いずれも8x8の64マスが前提の処理

func Search(field util.FieldViewer) []*marker.Marker {

	stateList := []*state{}

	// 8x8の64マスに0～9の数字が均等に散らばっていることが前提の処理
	// (仮に前提が崩れて全マスが0だったら組み合わせ爆発が起こる)
	for row := 0; row < util.RowCount; row++ {
		for col := 0; col < util.ColCount; col++ {
			value := field.Get(row, col)
			if value == util.Hole || value == util.Obstacle {
				continue
			}
			cur := newState(row, col, value)
			list := stateList[:]
			stateList = append(stateList, cur)
			for _, s := range list {
				if m := s.Merge(cur); m != nil {
					stateList = append(stateList, m)
				}
			}
		}
	}

	sort.Slice(stateList, func(a, b int) bool {
		nzA := stateList[a].NonZero
		nzB := stateList[b].NonZero
		if nzA < nzA {
			return true
		} else if nzA > nzB {
			return false
		} else {
			return stateList[a].Count < stateList[b].Count
		}
	})

	markerList := []*marker.Marker{}
	for _, s := range stateList {
		if s.Sum == util.Sum {
			markerList = append(markerList, s.ToMarker())
		}
	}
	return markerList
}
