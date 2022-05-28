package main

import (
	"fmt"
	"github.com/neetsdkasu/sum10-solver/solver"
	"log"
	"os"
)

func runSolver(seed int, name string, limitSeconds, tryCount int) (err error) {
	var file *os.File

	fileName := fmt.Sprintf("result%05d.txt", seed)
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

	fmt.Printf("ソルバ『%s』でSEED=%dの解を探します\n", name, seed)

	sol, _ := solver.FindSolver(name)
	fmt.Println("ソルバの詳細： ", sol.Description())

	log.Println("開始します")

	if err = solver.Solve(file, seed, name, limitSeconds, tryCount); err != nil {
		return
	}

	log.Println("結果を", fileName, "に保存しました")

	return
}
