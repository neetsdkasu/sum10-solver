// SUM10-SOLVER
// author: Leonardone @ NEETSDKASU

package main

import (
	"flag"
	"fmt"
	"github.com/neetsdkasu/sum10-solver/solver"
	"github.com/neetsdkasu/sum10-solver/util"
	"os"
)

const InvalidArgumentExitCode = 2

func init() {
	const url = "https://neetsdkasu.github.io/game/sum10/"

	oldUsage := flag.Usage

	flag.Usage = func() {
		fmt.Println()
		fmt.Println("通常モード")
		fmt.Println("  SUM10 Puzzleのスコアの良さそうな解を探すプログラム")
		fmt.Println("  求めた解を resutl?????.txt ファイルに出力します (?????はSEEDの値)")
		fmt.Println()

		oldUsage()

		fmt.Println()
		fmt.Println("ソルバ一覧")
		fmt.Println(solver.SolverList())
		fmt.Println()

		fmt.Println()
		fmt.Println("Example:")
		fmt.Printf("  %s -%s 5531\n", flag.CommandLine.Name(), argSeed)
		fmt.Println("        通常モードで実行。SEED=5531の時の解を探し、その解をファイル(result05531.txt)に出力する。")
		fmt.Printf("  %s -%s 5531 -%s RandomWalk\n", flag.CommandLine.Name(), argSeed, argSolverName)
		fmt.Println("        solverモードで実行。SEED=5531の時の解をRandomWalkソルバに探させ、その解をファイル(result05531.txt)に出力する。")
		fmt.Printf("  %s -%s 5531 -%s\n", flag.CommandLine.Name(), argSeed, argCompMode)
		fmt.Println("        compモードで実行。SEED=5531の時の解を複数のソルバに探させ、スコアを比較した結果をファイル(comp05531.txt)に出力する。")
		fmt.Printf("  %s -%s\n", flag.CommandLine.Name(), argCompMode)
		fmt.Println("        compモードで実行。ランダムに選ばれた", DefaultNumOfTestcase, "種類のSEEDのそれぞれの時の解を複数のソルバに探させ、スコアを比較した結果をファイル(comp.txt)に出力する。")
		fmt.Println()

		fmt.Println()
		fmt.Println("SUM10 Puzzle:", url)
		fmt.Println()
	}
}

type Args struct {
	seed               int
	withStatistics     bool
	compMode           bool
	compLimitSeconds   int
	compNumOfTestcase  int
	solverName         string
	solverLimitSeconds int
	solverTryCount     int
}

const (
	argSeed               = "seed"
	argStatistics         = "statistics"
	argCompMode           = "comp"
	argLimitSeconds       = "comptime"
	argNumOfTestcase      = "compsize"
	argSolverName         = "solver"
	argSolverLimitSeconds = "soltime"
	argSolverTryCount     = "solsize"
)

const (
	MinLimitSeconds     = 1
	MaxLimitSeconds     = 600
	DefaultLimitSeconds = 5

	MinNumOfTestcase     = 1
	MaxNumOfTestcase     = 100
	DefaultNumOfTestcase = 10

	MinSolverLimitSeconds     = 1
	MaxSolverLimitSeconds     = 3600
	DefaultSolverLimitSeconds = 5

	MinSolverTryCount     = 1
	MaxSolverTryCount     = 100
	DefaultSolverTryCount = 1
)

func parseArgs() *Args {
	args := &Args{}

	flag.IntVar(&args.seed, argSeed, util.NoSeed,
		fmt.Sprint("SUM10パズルのSEEDの指定 (", util.MinSeed, " ～ ", util.MaxSeed, ")"))

	flag.BoolVar(&args.withStatistics, argStatistics, false,
		"出力ファイルに初手の選択と最終スコアの関係の統計情報を含めます [通常モード]")

	flag.BoolVar(&args.compMode, argCompMode, false,
		fmt.Sprint("compモードで実行 (ソルバ数: ", solver.Count(), ")"))

	flag.IntVar(&args.compLimitSeconds, argLimitSeconds, DefaultLimitSeconds,
		fmt.Sprint("各ソルバがテストケースごとに使ってよい最大時間(秒) (", MinLimitSeconds, " ～ ", MaxLimitSeconds, ") [compモード]"))

	flag.IntVar(&args.compNumOfTestcase, argNumOfTestcase, DefaultNumOfTestcase,
		fmt.Sprint("使用するテストケースの数 (", MinNumOfTestcase, " ～ ", MaxNumOfTestcase, ") [compモード]"))

	flag.StringVar(&args.solverName, argSolverName, "", "指定した名前のソルバに解かせてみるsolverモードで実行")

	flag.IntVar(&args.solverLimitSeconds, argSolverLimitSeconds, DefaultSolverLimitSeconds,
		fmt.Sprint("ソルバがテストケースごとに使ってよい最大時間(秒) (", MinSolverLimitSeconds, " ～ ", MaxSolverLimitSeconds, ") [solverモード]"))

	flag.IntVar(&args.solverTryCount, argSolverTryCount, DefaultSolverTryCount,
		fmt.Sprint("ソルバの試行回数 (", MinSolverTryCount, " ～ ", MaxSolverTryCount, ") [solverモード]"))

	flag.Parse()

	if args.isSolverMode() {
		args.validateSolverMode()
	} else if args.compMode {
		args.validateCompMode()
	} else {
		args.validateNormalMode()
	}

	return args
}

func (args *Args) isSolverMode() bool {
	return args.solverName != ""
}

func (args *Args) validateSolverMode() {
	if !util.IsValidSeed(args.seed) {
		fmt.Println("solverモードでは", argSeed, "の指定が必須です")
		fmt.Println(argSeed, "は", util.MinSeed, "から", util.MaxSeed, "の範囲内の数字で指定してください")
		flag.Usage()
		os.Exit(InvalidArgumentExitCode)
	}

	if args.compMode {
		fmt.Println("solverモードとcompモードを同時に指定できません")
		flag.Usage()
		os.Exit(InvalidArgumentExitCode)
	}

	if !isValidSolverLimitSeconds(args.solverLimitSeconds) {
		fmt.Println(argSolverLimitSeconds, "は", MinSolverLimitSeconds, "から", MaxSolverLimitSeconds, "の範囲内の数字で指定してください")
		flag.Usage()
		os.Exit(InvalidArgumentExitCode)
		return
	}

	if !isValidSolverTryCount(args.solverTryCount) {
		fmt.Println(argSolverTryCount, "は", MinSolverTryCount, "から", MaxSolverTryCount, "の範囲内の数字で指定してください")
		flag.Usage()
		os.Exit(InvalidArgumentExitCode)
		return
	}

	if args.withStatistics {
		fmt.Println(argStatistics, "はsolverモードでは指定できません")
		flag.Usage()
		os.Exit(InvalidArgumentExitCode)
		return
	}

	if _, ok := solver.FindSolver(args.solverName); !ok {
		fmt.Println("ソルバ名が不正です")
		flag.Usage()
		os.Exit(InvalidArgumentExitCode)
		return
	}
}

func (args *Args) validateNormalMode() {
	if !util.IsValidSeed(args.seed) {
		fmt.Println("通常モードでは", argSeed, "の指定が必須です")
		fmt.Println(argSeed, "は", util.MinSeed, "から", util.MaxSeed, "の範囲内の数字で指定してください")
		flag.Usage()
		os.Exit(InvalidArgumentExitCode)
	}
}

func (args *Args) validateCompMode() {
	if !isValidLimitSeconds(args.compLimitSeconds) {
		fmt.Println(argLimitSeconds, "は", MinLimitSeconds, "から", MaxLimitSeconds, "の範囲内の数字で指定してください")
		flag.Usage()
		os.Exit(InvalidArgumentExitCode)
		return
	}

	if !isValidNumOfTestcase(args.compNumOfTestcase) {
		fmt.Println(argNumOfTestcase, "は", MinNumOfTestcase, "から", MaxNumOfTestcase, "の範囲内の数字で指定してください")
		flag.Usage()
		os.Exit(InvalidArgumentExitCode)
		return
	}

	if args.withStatistics {
		fmt.Println(argStatistics, "はcompモードでは指定できません")
		flag.Usage()
		os.Exit(InvalidArgumentExitCode)
		return
	}

	if !util.IsValidSeed(args.seed) {
		args.seed = util.NoSeed
	}
}

func isValidLimitSeconds(limitSeconds int) bool {
	return MinLimitSeconds <= limitSeconds &&
		limitSeconds <= MaxLimitSeconds
}

func isValidNumOfTestcase(numOfTestcase int) bool {
	return MinNumOfTestcase <= numOfTestcase &&
		numOfTestcase <= MaxNumOfTestcase
}

func isValidSolverLimitSeconds(limitSeconds int) bool {
	return MinSolverLimitSeconds <= limitSeconds &&
		limitSeconds <= MaxSolverLimitSeconds
}

func isValidSolverTryCount(tryCount int) bool {
	return MinSolverTryCount <= tryCount &&
		tryCount <= MaxSolverTryCount
}
