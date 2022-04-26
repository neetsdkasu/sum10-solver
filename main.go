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
			fmt.Print(" ", value)
		}
		fmt.Println()
	}

	game := game.New(problem)

	_ = game

	marker := marker.New()

	fmt.Println(marker.IsValid())

	marker.Field[4][4] = true
	fmt.Println(marker.IsValid())

	marker.Field[5][4] = true
	marker.Field[5][3] = true
	marker.Field[6][2] = true
	marker.Field[6][1] = true
	fmt.Println(marker.IsValid())

	marker.Field[5][2] = true
	fmt.Println(marker.IsValid())

}
