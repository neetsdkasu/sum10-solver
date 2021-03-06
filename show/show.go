// SUM10-SOLVER
// author: Leonardone @ NEETSDKASU

package show

import (
	"fmt"
	"github.com/neetsdkasu/sum10-solver/game"
	"github.com/neetsdkasu/sum10-solver/marker"
	"github.com/neetsdkasu/sum10-solver/util"
	"io"
)

func ShowField(file io.Writer, field util.FieldViewer) (err error) {
	for row := 0; row < util.RowCount; row++ {
		for col := 0; col < util.ColCount; col++ {
			value := field.Get(row, col)
			_, err = fmt.Fprintf(file, " %2d", value)
			if err != nil {
				return
			}
		}
		_, err = fmt.Fprintln(file)
		if err != nil {
			return
		}
	}
	return
}

func ShowFieldWithMark(file io.Writer, field util.FieldViewer, marker *marker.Marker) (err error) {
	for row := 0; row < util.RowCount; row++ {
		for col := 0; col < util.ColCount; col++ {
			value := field.Get(row, col)
			if marker.Has(row, col) {
				_, err = fmt.Fprintf(file, " *%d", value)
			} else {
				_, err = fmt.Fprintf(file, " %2d", value)
			}
			if err != nil {
				return
			}
		}
		_, err = fmt.Fprintln(file)
		if err != nil {
			return
		}
	}
	return
}

func ShowGameWithMark(file io.Writer, game *game.Game, marker *marker.Marker) (err error) {
	if _, err = fmt.Fprintln(file, "Steps:", game.Steps); err != nil {
		return
	}
	if _, err = fmt.Fprintln(file, "Score:", game.Score); err != nil {
		return
	}
	return ShowFieldWithMark(file, game, marker)
}

func ShowGame(file io.Writer, game *game.Game) (err error) {
	if _, err = fmt.Fprintln(file, "Steps:", game.Steps); err != nil {
		return
	}
	if _, err = fmt.Fprintln(file, "Score:", game.Score); err != nil {
		return
	}
	return ShowField(file, game)
}
