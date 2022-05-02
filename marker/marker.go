package marker

import (
	"sum10-solver/util"
)

type Marker struct {
	Field [][]bool
}

func New() *Marker {
	field := util.MakeEmptyField[bool]()
	return &Marker{Field: field}
}

func (marker *Marker) Has(row, col int) bool {
	return marker.Field[row][col]
}

func (marker *Marker) Set(row, col int) {
	marker.Field[row][col] = true
}

func (marker *Marker) Unset(row, col int) {
	marker.Field[row][col] = false
}

func (marker *Marker) Clear() {
	for _, line := range marker.Field {
		for col := range line {
			line[col] = false
		}
	}
}

func (marker *Marker) Count() int {
	count := 0
	for _, line := range marker.Field {
		for _, mark := range line {
			if mark {
				count++
			}
		}
	}
	return count
}

func (marker *Marker) IsValid() bool {
	markedCount := marker.Count()
	if markedCount == 0 {
		return false
	}
	for row, line := range marker.Field {
		for col, mark := range line {
			if !mark {
				continue
			}
			count := marker.GetCopy().dfsDelete(row, col)
			return count == markedCount
		}
	}
	return false
}

func Copy(dst, src *Marker) {
	util.CopyField(dst.Field, src.Field)
}

func (marker *Marker) GetCopy() *Marker {
	dst := New()
	Copy(dst, marker)
	return dst
}

func (marker *Marker) dfsDelete(row, col int) int {
	if !util.FieldContains(row, col) {
		return 0
	}
	if !marker.Field[row][col] {
		return 0
	}
	marker.Field[row][col] = false
	return 1 +
		marker.dfsDelete(row+1, col) +
		marker.dfsDelete(row-1, col) +
		marker.dfsDelete(row, col+1) +
		marker.dfsDelete(row, col-1)
}

func (marker *Marker) Sum(field [][]int) int {
	sum := 0
	for row, line := range marker.Field {
		for col, mark := range line {
			if mark {
				sum += field[row][col]
			}
		}
	}
	return sum
}

func (marker *Marker) Fill(field [][]int, value int) {
	for row, line := range marker.Field {
		for col, mark := range line {
			if mark {
				field[row][col] = value
			}
		}
	}
}
