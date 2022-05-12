package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"sum10-solver/game"
	"sum10-solver/problem"
	"sum10-solver/search"
	"sum10-solver/util"
	"time"
)

const NumOfSearching = 500000
const Progress = NumOfSearching / 50

func findGoodSolution(file io.Writer, seed uint32, withStatistics bool) (err error) {
	const Bar = "---------------------------------------------------------------------------------------"

	println(fmt.Sprintf("SUM10パズルのSEED=%dに対するランダム解を大量に生成し、その中からスコアの一番高い解を見つけます", seed))
	println("この作業には数十分以上の時間がかかります")
	log.Println("開始します")

	prob := problem.New(seed)

	if _, err = fmt.Fprintln(file, "SEED:", seed); err != nil {
		return
	}
	if _, err = fmt.Fprintln(file, Bar); err != nil {
		return
	}

	if err = showField(file, prob); err != nil {
		return
	}
	if _, err = fmt.Fprintln(file, Bar); err != nil {
		return
	}

	game0 := game.New(prob)

	firstSteps := search.Search(prob)
	numOfFirstStep := len(firstSteps)

	scores := make([]int, 400)
	statistics := make([][]int, 400)
	for i := range statistics {
		statistics[i] = make([]int, numOfFirstStep)
	}

	maxSel := 0

	best := game0

	time0 := time.Now()
	rand.Seed(time0.Unix())

	for i := 0; i < NumOfSearching; i++ {
		cur := game0
		firstSel := -1
		for {
			markerList := search.Search(cur)
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
			target := markerList[sel]
			if cur, err = cur.Take(target); err != nil {
				return
			}
		}
		scores[cur.Score]++
		statistics[cur.Score][firstSel]++

		if cur.Score > best.Score {
			best = cur
		}

		if i%Progress == 0 {
			dur := time.Now().Sub(time0).String()
			log.Printf("%6d / %6d (%s)\n", i+1, NumOfSearching, dur)
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
			if minScore == 9999 {
				if _, err = fmt.Fprint(file, "  ---"); err != nil {
					return
				}
			} else {
				if _, err = fmt.Fprintf(file, "%5d", minScore); err != nil {
					return
				}
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
			maxScore := -1
			for sc, line := range statistics {
				if line[i] > 0 {
					maxScore = sc
				}
			}
			if maxScore < 0 {
				if _, err = fmt.Fprint(file, "  ---"); err != nil {
					return
				}
			} else {
				if _, err = fmt.Fprintf(file, "%5d", maxScore); err != nil {
					return
				}
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
			if total > 0 {
				average := score / uint64(total)

				if _, err = fmt.Fprintf(file, "%5d", average); err != nil {
					return
				}
			} else {
				if _, err = fmt.Fprint(file, "  ---"); err != nil {
					return
				}
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

		for i, marker := range search.Search(prob) {
			if _, err = fmt.Fprintf(file, "---------- %3d ----------\n", i); err != nil {
				return
			}
			if err = showFieldWithMark(file, prob, marker); err != nil {
				return
			}
		}

		/*  *  *  *  *  *  *  *  *  *  *  */

		if _, err = fmt.Fprintln(file, Bar); err != nil {
			return
		}

		if _, err = fmt.Fprintln(file, "FIRST STEP MIN SCORE MAP"); err != nil {
			return
		}

		minScoreField := util.MakeEmptyField[int]()
		util.FillField(minScoreField, 9999)

		for sc, cnts := range statistics {
			for i, step := range firstSteps {
				if cnts[i] == 0 {
					continue
				}
				for row := 0; row < util.RowCount; row++ {
					for col := 0; col < util.ColCount; col++ {
						if step.Has(row, col) {
							minScoreField[row][col] = util.Min(minScoreField[row][col], sc)
						}
					}
				}
			}
		}

		for row := 0; row < util.RowCount; row++ {
			for col := 0; col < util.ColCount; col++ {
				value := minScoreField[row][col]
				if value == 9999 {
					if _, err = fmt.Fprint(file, " ---"); err != nil {
						return
					}
				} else {
					if _, err = fmt.Fprintf(file, "%4d", value); err != nil {
						return
					}
				}

			}
			if _, err = fmt.Fprintln(file); err != nil {
				return
			}
		}

		/*  *  *  *  *  *  *  *  *  *  *  */

		if _, err = fmt.Fprintln(file, Bar); err != nil {
			return
		}

		if _, err = fmt.Fprintln(file, "FIRST STEP MAX SCORE MAP"); err != nil {
			return
		}

		maxScoreField := util.MakeEmptyField[int]()
		util.FillField(maxScoreField, -1)
		for sc, cnts := range statistics {
			for i, step := range firstSteps {
				if cnts[i] == 0 {
					continue
				}
				for row := 0; row < util.RowCount; row++ {
					for col := 0; col < util.ColCount; col++ {
						if step.Has(row, col) {
							maxScoreField[row][col] = util.Max(maxScoreField[row][col], sc)
						}
					}
				}
			}
		}

		for row := 0; row < util.RowCount; row++ {
			for col := 0; col < util.ColCount; col++ {
				value := maxScoreField[row][col]
				if value < 0 {
					if _, err = fmt.Fprint(file, " ---"); err != nil {
						return
					}
				} else {
					if _, err = fmt.Fprintf(file, "%4d", value); err != nil {
						return
					}
				}

			}
			if _, err = fmt.Fprintln(file); err != nil {
				return
			}
		}

		/*  *  *  *  *  *  *  *  *  *  *  */

		if _, err = fmt.Fprintln(file, Bar); err != nil {
			return
		}

		if _, err = fmt.Fprintln(file, "FIRST STEP AVERAGE SCORE MAP"); err != nil {
			return
		}

		scoreField := util.MakeEmptyField[int]()
		countField := util.MakeEmptyField[int]()
		for sc, cnts := range statistics {
			for i, step := range firstSteps {
				cnt := cnts[i]
				if cnt == 0 {
					continue
				}
				sum := sc * cnt
				for row := 0; row < util.RowCount; row++ {
					for col := 0; col < util.ColCount; col++ {
						if step.Has(row, col) {
							scoreField[row][col] += sum
							countField[row][col] += cnt
						}
					}
				}
			}
		}

		for row := 0; row < util.RowCount; row++ {
			for col := 0; col < util.ColCount; col++ {
				if countField[row][col] > 0 {
					value := float64(scoreField[row][col]) / float64(countField[row][col])

					if _, err = fmt.Fprintf(file, " %5.1f", value); err != nil {
						return
					}
				} else {
					if _, err = fmt.Fprint(file, " ---.-"); err != nil {
						return
					}
				}

			}
			if _, err = fmt.Fprintln(file); err != nil {
				return
			}
		}

		/*  *  *  *  *  *  *  *  *  *  *  */

		if _, err = fmt.Fprintln(file, Bar); err != nil {
			return
		}
	} // withStatistics

	return
}
