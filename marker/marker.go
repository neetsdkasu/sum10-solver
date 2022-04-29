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
			count := marker.dfs(row, col)
			return count == markedCount
		}
	}
	return false
}

func Copy(dst, src *Marker) {
	util.CopyField(dst.Field, src.Field)
}

func (marker *Marker) makeCopy() *Marker {
	dst := New()
	Copy(dst, marker)
	return dst
}

func (marker *Marker) dfs(row0, col0 int) int {
	field := marker.makeCopy().Field
	stack := append([]int{}, row0, col0)
	count := 0
	for 0 < len(stack) {
		length := len(stack)
		row, col := stack[length-2], stack[length-1]
		stack = stack[:length-2]
		if row < 0 || col < 0 || len(field) <= row || len(field[row]) <= col {
			continue
		}
		if !field[row][col] {
			continue
		}
		count++
		field[row][col] = false
		stack = append(stack, row-1, col)
		stack = append(stack, row+1, col)
		stack = append(stack, row, col-1)
		stack = append(stack, row, col+1)
	}
	return count
}
