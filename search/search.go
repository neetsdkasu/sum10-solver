package search

import (
	"sum10-solver/marker"
	"sum10-solver/util"
)

// いずれも8x8の64マスが前提の処理

type state struct {
	Cells     uint64
	Neigbours uint64
	Sum       int
}

func shift(row, col int) int {
	return (row << 3) | col
}

func newState(row, col, value int) *state {
	s := &state{}
	s.Cells = 1 << shift(row, col)
	if util.FieldContains(row+1, col) {
		s.Neigbours |= 1 << shift(row+1, col)
	}
	if util.FieldContains(row-1, col) {
		s.Neigbours |= 1 << shift(row-1, col)
	}
	if util.FieldContains(row, col+1) {
		s.Neigbours |= 1 << shift(row, col+1)
	}
	if util.FieldContains(row, col-1) {
		s.Neigbours |= 1 << shift(row, col-1)
	}
	s.Sum = value
	return s
}

func (s *state) HasCell(row, col int) bool {
	return (s.Cells & (1 << shift(row, col))) != 0
}

func (s *state) ToMarker() *marker.Marker {
	marker := marker.New()
	for row := 0; row < util.RowCount; row++ {
		for col := 0; col < util.ColCount; col++ {
			if s.HasCell(row, col) {
				marker.Set(row, col)
			}
		}
	}
	return marker
}

func (s *state) Merge(o *state) *state {
	if (s.Neigbours & o.Cells) == 0 {
		return nil
	}
	if s.Sum+o.Sum > util.Sum {
		return nil
	}
	r := &state{}
	r.Cells = s.Cells | o.Cells
	r.Neigbours = (s.Neigbours | o.Neigbours) &^ r.Cells
	r.Sum = s.Sum + o.Sum
	return r
}

func Search(field [][]int) []*marker.Marker {

	stateList := []*state{}

	// 8x8の64マスに0～9の数字が均等に散らばっていることが前提の処理
	// (仮に前提が崩れて全マスが0だったら組み合わせ爆発が起こる)
	for row, line := range field {
		for col, value := range line {
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

	markerList := []*marker.Marker{}
	for _, s := range stateList {
		if s.Sum == util.Sum {
			markerList = append(markerList, s.ToMarker())
		}
	}
	return markerList
}
