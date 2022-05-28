// SUM10-SOLVER
// author: Leonardone @ NEETSDKASU

package main

import (
	"fmt"
	"github.com/neetsdkasu/sum10-solver/game"
	"github.com/neetsdkasu/sum10-solver/marker"
	"github.com/neetsdkasu/sum10-solver/problem"
	"github.com/neetsdkasu/sum10-solver/search"
	"github.com/neetsdkasu/sum10-solver/show"
	"github.com/neetsdkasu/sum10-solver/util"
	"io"
	"log"
	"math"
	"math/rand"
	"time"
)

const NumOfSearching = 500000
const Progress = NumOfSearching / 50

func findGoodSolution(file io.Writer, seed uint32, withStatistics bool) (err error) {
	sw := NewSolWriter(file, seed)
	sw.Seed()
	sw.Bar()
	sw.Problem()
	sw.Bar()
	sw.Search()
	sw.Time()
	sw.Bar()
	sw.Solution()
	sw.Bar()
	sw.Best()
	sw.Bar()
	if withStatistics {
		sw.Statistics()
	}
	return sw.Error()
}

type SolWriter struct {
	err        error
	file       io.Writer
	prob       *problem.Problem
	firstSteps []*marker.Marker
	best       *game.Game
	maxSel     int
	statistics [][]int
	scores     []int
	time       time.Duration
}

func NewSolWriter(file io.Writer, seed uint32) *SolWriter {
	writer := &SolWriter{}
	writer.file = file
	writer.prob = problem.New(seed)
	return writer
}

func (this *SolWriter) Error() error {
	return this.err
}

func (this *SolWriter) Bar() {
	if this.err != nil {
		return
	}
	const Bar = "---------------------------------------------------------------------------------------"
	_, this.err = fmt.Fprintln(this.file, Bar)
}

func (this *SolWriter) Seed() {
	if this.err != nil {
		return
	}
	_, this.err = fmt.Fprintln(this.file, "SEED:", this.prob.Seed())
}

func (this *SolWriter) Problem() {
	if this.err != nil {
		return
	}
	this.err = show.ShowField(this.file, this.prob)
}

func (this *SolWriter) Time() {
	if this.err != nil {
		return
	}
	_, this.err = fmt.Fprintln(this.file, "TIME:", this.time)
}

func (this *SolWriter) Search() {
	if this.err != nil {
		return
	}

	game0 := game.New(this.prob)

	firstSteps := search.Search(this.prob)
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

		if i%Progress == 0 {
			dur := time.Now().Sub(time0).String()
			log.Printf("%6d / %6d (%s)\n", i, NumOfSearching, dur)
		}

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
			if cur, this.err = cur.Take(target); this.err != nil {
				return
			}
		}
		scores[cur.Score]++
		statistics[cur.Score][firstSel]++

		if cur.Score > best.Score {
			best = cur
		}
	}
	time1 := time.Now()

	this.maxSel = maxSel
	this.best = best
	this.time = time1.Sub(time0)
	this.statistics = statistics
	this.scores = scores
	this.firstSteps = firstSteps
}

func (this *SolWriter) Solution() {
	if this.err != nil {
		return
	}
	best := this.best

	_, this.err = fmt.Fprintln(this.file, "BEST SOLUTION (SCORE", best.Score, ")")
	if this.err != nil {
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
			this.Bar()
			if this.err != nil {
				return
			}
			this.err = show.ShowGameWithMark(this.file, step.Prev, step.Taked)
			if this.err != nil {
				return
			}
		}
	}
}

func (this *SolWriter) Best() {
	if this.err != nil {
		return
	}
	this.err = show.ShowGame(this.file, this.best)
}

func (this *SolWriter) Statistics() {
	if this.err != nil {
		return
	}

	_, this.err = fmt.Fprintln(this.file, "STATISTICS OF FIRST STEP")

	this.Bar()
	this.Indexes()
	this.TableBar()
	this.CountTable()
	this.TableBar()
	this.Total()
	this.TableBar()
	this.MinScore()
	this.MaxScore()
	this.AverageScore()
	this.ModeScore()
	this.StandardDeviation()
	this.Bar()
	this.FirstSteps()
	this.Bar()
	this.MinScoreMap()
	this.Bar()
	this.MaxScoreMap()
	this.Bar()
	this.AverageScoreMap()
	this.Bar()
	this.ModeScoreMap()
	this.Bar()
	this.StandardDeviationMap()
	this.Bar()
}

func (this *SolWriter) Indexes() {
	if this.err != nil {
		return
	}
	maxSel := this.maxSel

	_, this.err = fmt.Fprint(this.file, "       FIRST STEP INDEX: ")
	if this.err != nil {
		return
	}

	for i := 0; i <= maxSel; i++ {
		_, this.err = fmt.Fprintf(this.file, "%5d", i)
		if this.err != nil {
			return
		}
	}

	_, this.err = fmt.Fprintln(this.file)
}

func (this *SolWriter) TableBar() {
	if this.err != nil {
		return
	}
	file := this.file
	maxSel := this.maxSel

	_, this.err = fmt.Fprint(file, "=========================")
	if this.err != nil {
		return
	}

	for i := 0; i <= maxSel; i++ {
		_, this.err = fmt.Fprint(file, "=====")
		if this.err != nil {
			return
		}
	}

	_, this.err = fmt.Fprintln(file)
}

func (this *SolWriter) CountTable() {
	if this.err != nil {
		return
	}
	file := this.file
	maxSel := this.maxSel
	scores := this.scores
	statistics := this.statistics

	for sc, cnt := range scores {
		if cnt == 0 {
			continue
		}
		_, this.err = fmt.Fprintf(file, "SCORE %3d, COUNT %6d: ", sc, cnt)
		if this.err != nil {
			return
		}
		for i := 0; i <= maxSel; i++ {
			_, this.err = fmt.Fprintf(file, "%5d", statistics[sc][i])
			if this.err != nil {
				return
			}
		}
		_, this.err = fmt.Fprintln(file)
		if this.err != nil {
			return
		}
	}
}

func (this *SolWriter) Total() {
	if this.err != nil {
		return
	}
	file := this.file
	maxSel := this.maxSel
	statistics := this.statistics

	_, this.err = fmt.Fprint(file, "                  TOTAL: ")
	if this.err != nil {
		return
	}

	for i := 0; i <= maxSel; i++ {
		total := 0
		for _, line := range statistics {
			total += line[i]
		}
		_, this.err = fmt.Fprintf(file, "%5d", total)
		if this.err != nil {
			return
		}
	}

	_, this.err = fmt.Fprintln(file)
}

func (this *SolWriter) ModeScore() {
	if this.err != nil {
		return
	}
	file := this.file
	maxSel := this.maxSel
	statistics := this.statistics

	_, this.err = fmt.Fprint(file, "             MODE SCORE: ")
	if this.err != nil {
		return
	}

	for i := 0; i <= maxSel; i++ {
		modeScore := -1
		count := 0
		for sc, line := range statistics {
			if line[i] > count {
				modeScore = sc
				count = line[i]
			}
		}
		if modeScore == -1 {
			_, this.err = fmt.Fprint(file, "  ---")
		} else {
			_, this.err = fmt.Fprintf(file, "%5d", modeScore)
		}
		if this.err != nil {
			return
		}
	}

	_, this.err = fmt.Fprintln(file)
}

func (this *SolWriter) StandardDeviation() {
	if this.err != nil {
		return
	}
	file := this.file
	maxSel := this.maxSel
	statistics := this.statistics

	_, this.err = fmt.Fprint(file, "     STANDARD DEVIATION: ")
	if this.err != nil {
		return
	}

	for i := 0; i <= maxSel; i++ {
		total := uint64(0)
		total2 := uint64(0)
		count := 0

		for sc, line := range statistics {
			if line[i] > 0 {
				total += uint64(sc) * uint64(line[i])
				total2 += uint64(sc*sc) * uint64(line[i])
				count += line[i]
			}
		}
		if count == 0 {
			_, this.err = fmt.Fprint(file, "  ---")
		} else {
			v := float64(total2)/float64(count) -
				math.Pow(float64(total)/float64(count), 2)
			sd := int(math.Floor(math.Sqrt(v)))
			_, this.err = fmt.Fprintf(file, "%5d", sd)
		}
		if this.err != nil {
			return
		}
	}

	_, this.err = fmt.Fprintln(file)
}

func (this *SolWriter) MinScore() {
	if this.err != nil {
		return
	}
	file := this.file
	maxSel := this.maxSel
	statistics := this.statistics

	_, this.err = fmt.Fprint(file, "              MIN SCORE: ")
	if this.err != nil {
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
			_, this.err = fmt.Fprint(file, "  ---")
		} else {
			_, this.err = fmt.Fprintf(file, "%5d", minScore)
		}
		if this.err != nil {
			return
		}
	}

	_, this.err = fmt.Fprintln(file)
}

func (this *SolWriter) MaxScore() {
	if this.err != nil {
		return
	}
	file := this.file
	maxSel := this.maxSel
	statistics := this.statistics

	_, this.err = fmt.Fprint(file, "              MAX SCORE: ")
	if this.err != nil {
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
			_, this.err = fmt.Fprint(file, "  ---")
		} else {
			_, this.err = fmt.Fprintf(file, "%5d", maxScore)
		}
		if this.err != nil {
			return
		}
	}

	_, this.err = fmt.Fprintln(file)
}

func (this *SolWriter) AverageScore() {
	if this.err != nil {
		return
	}
	file := this.file
	maxSel := this.maxSel
	statistics := this.statistics

	_, this.err = fmt.Fprint(file, "          AVERAGE SCORE: ")
	if this.err != nil {
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

			_, this.err = fmt.Fprintf(file, "%5d", average)
		} else {
			_, this.err = fmt.Fprint(file, "  ---")
		}
		if this.err != nil {
			return
		}
	}

	_, this.err = fmt.Fprintln(file)
}

func (this *SolWriter) FirstSteps() {
	if this.err != nil {
		return
	}
	prob := this.prob
	file := this.file
	firstSteps := this.firstSteps

	_, this.err = fmt.Fprintln(file, "FIRST STEP DETAILS")
	if this.err != nil {
		return
	}

	for i, step := range firstSteps {
		_, this.err = fmt.Fprintf(file, "---------- %3d ----------\n", i)
		if this.err != nil {
			return
		}
		this.err = show.ShowFieldWithMark(file, prob, step)
		if this.err != nil {
			return
		}
	}
}

func (this *SolWriter) MinScoreMap() {
	if this.err != nil {
		return
	}
	statistics := this.statistics
	file := this.file
	firstSteps := this.firstSteps

	_, this.err = fmt.Fprintln(file, "FIRST STEP MIN SCORE MAP")
	if this.err != nil {
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
				_, this.err = fmt.Fprint(file, " ---")
			} else {
				_, this.err = fmt.Fprintf(file, "%4d", value)
			}
			if this.err != nil {
				return
			}

		}
		_, this.err = fmt.Fprintln(file)
		if this.err != nil {
			return
		}
	}
}

func (this *SolWriter) MaxScoreMap() {
	if this.err != nil {
		return
	}
	statistics := this.statistics
	file := this.file
	firstSteps := this.firstSteps

	_, this.err = fmt.Fprintln(file, "FIRST STEP MAX SCORE MAP")
	if this.err != nil {
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
				_, this.err = fmt.Fprint(file, " ---")
			} else {
				_, this.err = fmt.Fprintf(file, "%4d", value)
			}
			if this.err != nil {
				return
			}
		}
		_, this.err = fmt.Fprintln(file)
		if this.err != nil {
			return
		}
	}
}

func (this *SolWriter) AverageScoreMap() {
	if this.err != nil {
		return
	}
	statistics := this.statistics
	file := this.file
	firstSteps := this.firstSteps

	_, this.err = fmt.Fprintln(file, "FIRST STEP AVERAGE SCORE MAP")
	if this.err != nil {
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
				_, this.err = fmt.Fprintf(file, " %5.1f", value)
			} else {
				_, this.err = fmt.Fprint(file, " ---.-")
			}
			if this.err != nil {
				return
			}
		}
		_, this.err = fmt.Fprintln(file)
		if this.err != nil {
			return
		}
	}
}

func (this *SolWriter) ModeScoreMap() {
	if this.err != nil {
		return
	}
	statistics := this.statistics
	file := this.file
	firstSteps := this.firstSteps

	_, this.err = fmt.Fprintln(file, "FIRST STEP MODE SCORE MAP")
	if this.err != nil {
		return
	}

	scoreField := make([][][]int, 400)
	for sc := range scoreField {
		scoreField[sc] = util.MakeEmptyField[int]()
		util.FillField(scoreField[sc], -1)
	}

	for sc, cnts := range statistics {
		for i, step := range firstSteps {
			if cnts[i] == 0 {
				continue
			}
			for row := 0; row < util.RowCount; row++ {
				for col := 0; col < util.ColCount; col++ {
					if step.Has(row, col) {
						scoreField[sc][row][col] += cnts[i]
					}
				}
			}
		}
	}

	for row := 0; row < util.RowCount; row++ {
		for col := 0; col < util.ColCount; col++ {
			score := 0
			count := 0
			for sc, field := range scoreField {
				if field[row][col] > count {
					score = sc
					count = field[row][col]
				}
			}
			if count == 0 {
				_, this.err = fmt.Fprint(file, " ---")
			} else {
				_, this.err = fmt.Fprintf(file, "%4d", score)
			}
			if this.err != nil {
				return
			}
		}
		_, this.err = fmt.Fprintln(file)
		if this.err != nil {
			return
		}
	}
}

func (this *SolWriter) StandardDeviationMap() {
	if this.err != nil {
		return
	}
	statistics := this.statistics
	file := this.file
	firstSteps := this.firstSteps

	_, this.err = fmt.Fprintln(file, "FIRST STEP STANDARD DEVIATION MAP")
	if this.err != nil {
		return
	}

	scoreField := make([][][]int, 400)
	for sc := range scoreField {
		scoreField[sc] = util.MakeEmptyField[int]()
		util.FillField(scoreField[sc], -1)
	}

	for sc, cnts := range statistics {
		for i, step := range firstSteps {
			if cnts[i] == 0 {
				continue
			}
			for row := 0; row < util.RowCount; row++ {
				for col := 0; col < util.ColCount; col++ {
					if step.Has(row, col) {
						scoreField[sc][row][col] += cnts[i]
					}
				}
			}
		}
	}

	for row := 0; row < util.RowCount; row++ {
		for col := 0; col < util.ColCount; col++ {
			total := uint64(0)
			total2 := uint64(0)
			count := 0
			for sc, field := range scoreField {
				if cnt := field[row][col]; cnt > 0 {
					total += uint64(sc) * uint64(cnt)
					total2 += uint64(sc) * uint64(sc) * uint64(cnt)
					count += cnt
				}
			}
			if count == 0 {
				_, this.err = fmt.Fprint(file, " --.-")
			} else {
				v := float64(total2)/float64(count) -
					math.Pow(float64(total)/float64(count), 2)
				sd := math.Sqrt(v)
				_, this.err = fmt.Fprintf(file, "%5.1f", sd)
			}
			if this.err != nil {
				return
			}
		}
		_, this.err = fmt.Fprintln(file)
		if this.err != nil {
			return
		}
	}
}
