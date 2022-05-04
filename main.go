package main

import (
	"fmt"
	"math/rand"
	"sum10-solver/game"
	"sum10-solver/marker"
	"sum10-solver/problem"
	"sum10-solver/search"
	"time"
)

func main() {

	problem := problem.New(5531)

	showField(problem.Field)
	fmt.Println("-----------------------------")

	game0 := game.New(problem)

	time0 := time.Now()

	rand.Seed(time0.Unix())

	scores := make([]int, 300)

	for i := 0; i < 1000; i++ {
		game := game0
		for {
			markerList := search.Search(game.Field)
			if len(markerList) == 0 {
				break
			}
			sel := rand.Intn(len(markerList))
			marker := markerList[sel]
			var err error
			if game, err = game.Take(marker); err != nil {
				fmt.Println(err)
				return
			}
			// showGameWithMark(game.Prev, marker)
			// fmt.Println("-----------------------------")
		}
		scores[game.Score]++
		// showGame(game)
	}
	time1 := time.Now()

	for sc, cnt := range scores {
		if cnt == 0 {
			continue
		}
		fmt.Println("SCORE", sc, "COUNT", cnt)
	}

	fmt.Println("time", time1.Sub(time0))
}

func showField(field [][]int) {
	for _, line := range field {
		for _, value := range line {
			fmt.Printf(" %2d", value)
		}
		fmt.Println()
	}
}

func showFieldWithMark(field [][]int, marker *marker.Marker) {
	for row, line := range field {
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

func showGameWithMark(game *game.Game, marker *marker.Marker) {
	fmt.Println("Steps:", game.Steps)
	fmt.Println("Score:", game.Score)
	showFieldWithMark(game.Field, marker)
}

func showGame(game *game.Game) {
	fmt.Println("Steps:", game.Steps)
	fmt.Println("Score:", game.Score)
	showField(game.Field)
}
