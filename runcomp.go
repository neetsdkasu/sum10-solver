package main

import (
	"fmt"
	"github.com/neetsdkasu/sum10-solver/solver"
	"github.com/neetsdkasu/sum10-solver/util"
	"log"
	"os"
)

// list up solvers
import (
	_ "github.com/neetsdkasu/sum10-solver/solver/random-walk"

	_ "github.com/neetsdkasu/sum10-solver/solver/lots-of-choices"

	_ "github.com/neetsdkasu/sum10-solver/solver/fewer-choices"

	_ "github.com/neetsdkasu/sum10-solver/solver/middle-choices"

	_ "github.com/neetsdkasu/sum10-solver/solver/first"
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
