package main

import (
	"fmt"
	"log"
	"os"
)

func run(seed int, withStatistics bool) (err error) {
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

	println(fmt.Sprintf("SUM10パズルのSEED=%dに対するランダム解を大量に生成し、その中からスコアの一番高い解を見つけます", seed))
	println("この作業には数十分以上の時間がかかります")
	log.Println("開始します")

	if err = findGoodSolution(file, uint32(seed), withStatistics); err != nil {
		return
	}

	log.Println("結果を", fileName, "に保存しました")

	return
}
