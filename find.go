package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"sum10-solver/game"
	"sum10-solver/problem"
	"sum10-solver/search"
	"time"
)

func findGoodSolution(file io.Writer, seed uint32, withStatistics bool) (err error) {
	const Bar = "-----------------------------"

	log.Printf("searching solutions of seed %d puzzle by random walk\n", seed)
	log.Println("wait for tens of minutes ...")

	problem := problem.New(seed)

	if _, err = fmt.Fprintln(file, "SEED:", seed); err != nil {
		return
	}
	if _, err = fmt.Fprintln(file, Bar); err != nil {
		return
	}

	if err = showField(file, problem.Field); err != nil {
		return
	}
	if _, err = fmt.Fprintln(file, Bar); err != nil {
		return
	}

	game0 := game.New(problem)

	numOfFirstStep := len(search.Search(problem.Field))

	scores := make([]int, 400)
	statistics := make([][]int, 400)
	for i := range statistics {
		statistics[i] = make([]int, numOfFirstStep)
	}

	maxSel := 0

	best := game0

	const NumOfSearching = 500000
	const Progress = NumOfSearching / 50

	time0 := time.Now()
	rand.Seed(time0.Unix())

	for i := 0; i < NumOfSearching; i++ {
		if i%Progress == 0 {
			dur := time.Now().Sub(time0).String()
			log.Printf("%6d / %6d (%s)\n", i, NumOfSearching, dur)
		}
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
			if game, err = game.Take(marker); err != nil {
				return
			}
		}
		scores[game.Score]++
		statistics[game.Score][firstSel]++

		if game.Score > best.Score {
			best = game
		}
	}
	time1 := time.Now()

	if _, err = fmt.Fprintln(file, "TIME:", time1.Sub(time0)); err != nil {
		return
	}

	if _, err = fmt.Fprintln(file, Bar); err != nil {
		return
	}

	if _, err = fmt.Fprintln(file, "BEST SOLUTION (SCORE", best.Score, ")"); err != nil {
		return
	}

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
			if _, err = fmt.Fprintln(file, Bar); err != nil {
				return
			}
			if err = showGameWithMark(file, step.Prev, step.Taked); err != nil {
				return
			}
		}
	}

	if _, err = fmt.Fprintln(file, Bar); err != nil {
		return
	}

	if err = showGame(file, best); err != nil {
		return
	}

	if _, err = fmt.Fprintln(file, Bar); err != nil {
		return
	}

	if withStatistics {

		if _, err = fmt.Fprintln(file, "STATISTICS OF FIRST STEP"); err != nil {
			return
		}

		if _, err = fmt.Fprintln(file, Bar); err != nil {
			return
		}

		/*  *  *  *  *  *  *  *  *  *  *  */

		if _, err = fmt.Fprint(file, "       FIRST STEP INDEX: "); err != nil {
			return
		}

		for i := 0; i <= maxSel; i++ {
			if _, err = fmt.Fprintf(file, "%5d", i); err != nil {
				return
			}
		}
		if _, err = fmt.Fprintln(file); err != nil {
			return
		}

		/*  *  *  *  *  *  *  *  *  *  *  */

		if _, err = fmt.Fprint(file, "========================="); err != nil {
			return
		}

		for i := 0; i <= maxSel; i++ {
			if _, err = fmt.Fprint(file, "====="); err != nil {
				return
			}
		}
		if _, err = fmt.Fprintln(file); err != nil {
			return
		}

		/*  *  *  *  *  *  *  *  *  *  *  */

		for sc, cnt := range scores {
			if cnt == 0 {
				continue
			}
			if _, err = fmt.Fprintf(file, "SCORE %3d, COUNT %6d: ", sc, cnt); err != nil {
				return
			}
			for i := 0; i <= maxSel; i++ {
				if _, err = fmt.Fprintf(file, "%5d", statistics[sc][i]); err != nil {
					return
				}
			}
			if _, err = fmt.Fprintln(file); err != nil {
				return
			}
		}

		/*  *  *  *  *  *  *  *  *  *  *  */

		if _, err = fmt.Fprint(file, "========================="); err != nil {
			return
		}

		for i := 0; i <= maxSel; i++ {
			if _, err = fmt.Fprint(file, "====="); err != nil {
				return
			}
		}
		if _, err = fmt.Fprintln(file); err != nil {
			return
		}

		/*  *  *  *  *  *  *  *  *  *  *  */

		if _, err = fmt.Fprint(file, "                  TOTAL: "); err != nil {
			return
		}

		for i := 0; i <= maxSel; i++ {
			total := 0
			for _, line := range statistics {
				total += line[i]
			}
			if _, err = fmt.Fprintf(file, "%5d", total); err != nil {
				return
			}
		}
		if _, err = fmt.Fprintln(file); err != nil {
			return
		}

		/*  *  *  *  *  *  *  *  *  *  *  */

		if _, err = fmt.Fprint(file, "========================="); err != nil {
			return
		}

		for i := 0; i <= maxSel; i++ {
			if _, err = fmt.Fprint(file, "====="); err != nil {
				return
			}
		}
		if _, err = fmt.Fprintln(file); err != nil {
			return
		}

		/*  *  *  *  *  *  *  *  *  *  *  */

		if _, err = fmt.Fprint(file, "              MIN SCORE: "); err != nil {
			return
		}

		for i := 0; i <= maxSel; i++ {
			minScore := 9999
			for sc, line := range statistics {
				if line[i] > 0 {
					minScore = sc
					break
				}
			}
			if _, err = fmt.Fprintf(file, "%5d", minScore); err != nil {
				return
			}
		}
		if _, err = fmt.Fprintln(file); err != nil {
			return
		}

		/*  *  *  *  *  *  *  *  *  *  *  */

		if _, err = fmt.Fprint(file, "              MAX SCORE: "); err != nil {
			return
		}

		for i := 0; i <= maxSel; i++ {
			maxScore := 0
			for sc, line := range statistics {
				if line[i] > 0 {
					maxScore = sc
				}
			}
			if _, err = fmt.Fprintf(file, "%5d", maxScore); err != nil {
				return
			}
		}
		if _, err = fmt.Fprintln(file); err != nil {
			return
		}

		/*  *  *  *  *  *  *  *  *  *  *  */

		if _, err = fmt.Fprint(file, "          AVERAGE SCORE: "); err != nil {
			return
		}

		for i := 0; i <= maxSel; i++ {
			total := 0
			score := uint64(0)
			for sc, line := range statistics {
				total += line[i]
				score += uint64(sc) * uint64(line[i])
			}
			average := score / uint64(total)
			if _, err = fmt.Fprintf(file, "%5d", average); err != nil {
				return
			}
		}
		if _, err = fmt.Fprintln(file); err != nil {
			return
		}

		/*  *  *  *  *  *  *  *  *  *  *  */

		if _, err = fmt.Fprintln(file, Bar); err != nil {
			return
		}

		if _, err = fmt.Fprintln(file, "FIRST STEP DETAILS"); err != nil {
			return
		}

		for i, marker := range search.Search(problem.Field) {
			if _, err = fmt.Fprintf(file, "---------- %3d ----------\n", i); err != nil {
				return
			}
			if err = showFieldWithMark(file, problem.Field, marker); err != nil {
				return
			}
		}

		if _, err = fmt.Fprintln(file, Bar); err != nil {
			return
		}

	} // withStatistics

	return
}
