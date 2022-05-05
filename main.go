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
	statistcs := make([][]int, 300)
	for i := range statistcs {
		statistcs[i] = make([]int, 100)
	}
	maxSel := 0

	best := game0

	for i := 0; i < 500000; i++ {
		game := game0
		firstSel := -1
		for {
			markerList := search.Search(game.Field)
			if len(markerList) == 0 {
				break
			}
			sel := rand.Intn(len(markerList))
			if firstSel < 0 {
				firstSel = sel
				if sel > maxSel {
					maxSel = sel
				}
			}
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
		statistcs[game.Score][firstSel]++

		if game.Score > best.Score {
			best = game
		}
	}
	time1 := time.Now()

	for i, marker := range search.Search(problem.Field) {
		fmt.Println("--------", i, "-------")
		showFieldWithMark(problem.Field, marker)
	}

	fmt.Print("SCORE ---, COUNT ------: ")
	for i := 0; i <= maxSel; i++ {
		fmt.Printf("%5d", i)
	}
	fmt.Println()

	for sc, cnt := range scores {
		if cnt == 0 {
			continue
		}
		fmt.Printf("SCORE %3d, COUNT %6d: ", sc, cnt)
		for i := 0; i <= maxSel; i++ {
			fmt.Printf("%5d", statistcs[sc][i])
		}
		fmt.Println()
	}

	fmt.Println("time", time1.Sub(time0))

	fmt.Println("best solution")
	steps := []*game.Game{}
	for best != nil {
		steps = append(steps, best)
		best = best.Prev
	}
	best = steps[0]
	for len(steps) > 0 {
		pos := len(steps) - 1
		step := steps[pos]
		steps = steps[:pos]
		if step.Taked != nil {
			fmt.Println("----------------------------")
			showGameWithMark(step.Prev, step.Taked)
		}
	}
	fmt.Println("----------------------------")
	showGame(best)
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
