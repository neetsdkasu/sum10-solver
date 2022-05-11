package main

import (
	"fmt"
	"log"
	"os"
	"sum10-solver/solver"
	"sum10-solver/util"
)

// list up solvers
import (
	_ "sum10-solver/solver/random-walk"

	_ "sum10-solver/solver/lots-of-choices-greedy"

	_ "sum10-solver/solver/fewer-choices-greedy"

	_ "sum10-solver/solver/middle-choices-greedy"
)

func runComp(limitSeconds, numOfTestcase, seed int) (err error) {
	var file *os.File
	var fileName string

	if util.IsValidSeed(seed) {
		fileName = fmt.Sprintf("comp%05d.txt", seed)
	} else {
		fileName = "comp.txt"
	}
	file, err = os.Create(fileName)
	if err != nil {
		return
	}
	defer func() {
		err2 := file.Close()
		if err2 != nil {
			if err == nil {
				err = err2
			} else {
				log.Println(err2)
			}
		}
	}()

	if err = solver.Comp(file, limitSeconds, numOfTestcase, seed); err != nil {
		return
	}

	log.Println("save result to " + fileName)

	return
}
