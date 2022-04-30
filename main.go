package main

import (
	"fmt"
	"sum10-solver/game"
	"sum10-solver/marker"
	"sum10-solver/problem"
)

func main() {

	problem := problem.New(5531)

	for _, line := range problem.Field {
		for _, value := range line {
			fmt.Printf(" %2d", value)
		}
		fmt.Println()
	}

	game := game.New(problem)

	marker := marker.New()

	fmt.Println(marker.IsValid())

	marker.Set(4, 4)
	fmt.Println(marker.IsValid())

	marker.Set(5, 4)
	marker.Set(5, 3)
	marker.Set(6, 2)
	marker.Set(6, 1)
	fmt.Println(marker.IsValid())

	marker.Set(5, 2)
	fmt.Println(marker.IsValid())

	next, err := game.Take(marker)
	if err != nil {
		fmt.Println(err)
		return
	}
	showGameWithMark(game, marker)
	showGame(next)

	steps := [][]int{
		[]int{2, 0, 3, 0, 3, 1, 4, 1},
		[]int{0, 6, 0, 7},
		[]int{2, 5, 3, 2, 3, 4, 3, 5, 4, 2, 4, 3, 4, 4},
		[]int{2, 1, 2, 2, 2, 3, 2, 4},
	}

	game = next

	for _, step := range steps {
		marker.Clear()
		for i := 0; i < len(step); i += 2 {
			marker.Set(step[i], step[i+1])
		}
		fmt.Println(marker.IsValid())

		game, err = game.Take(marker)
		if err != nil {
			fmt.Println(err)
			return
		}
		showGameWithMark(game.Prev, game.Taked)
		showGame(game)
	}
}

func showGameWithMark(game *game.Game, marker *marker.Marker) {
	fmt.Println("Steps:", game.Steps)
	fmt.Println("Score:", game.Score)
	for row, line := range game.Field {
		for col, value := range line {
			if marker.Has(row, col) {
				fmt.Printf(" *%d", value)
			} else {
				fmt.Printf(" %2d", value)
			}
		}
		fmt.Println()
	}
}
func showGame(game *game.Game) {
	fmt.Println("Steps:", game.Steps)
	fmt.Println("Score:", game.Score)
	for _, line := range game.Field {
		for _, value := range line {
			fmt.Printf(" %2d", value)
		}
		fmt.Println()
	}
}
