package main

import (
	"fmt"
	"io"
	"sum10-solver/game"
	"sum10-solver/marker"
)

func showField(file io.Writer, field [][]int) (err error) {
	for _, line := range field {
		for _, value := range line {
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

func showFieldWithMark(file io.Writer, field [][]int, marker *marker.Marker) (err error) {
	for row, line := range field {
		for col, value := range line {
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

func showGameWithMark(file io.Writer, game *game.Game, marker *marker.Marker) (err error) {
	if _, err = fmt.Fprintln(file, "Steps:", game.Steps); err != nil {
		return
	}
	if _, err = fmt.Fprintln(file, "Score:", game.Score); err != nil {
		return
	}
	return showFieldWithMark(file, game.Field, marker)
}

func showGame(file io.Writer, game *game.Game) (err error) {
	if _, err = fmt.Fprintln(file, "Steps:", game.Steps); err != nil {
		return
	}
	if _, err = fmt.Fprintln(file, "Score:", game.Score); err != nil {
		return
	}
	return showField(file, game.Field)
}
