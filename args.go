package main

import (
	"flag"
	"fmt"
	"os"
	"sum10-solver/solver"
	"sum10-solver/util"
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
		fmt.Println("SUM10 Puzzle:", url)
		fmt.Println()
	}
}

type Args struct {
	seed              int
	withStatistics    bool
	compMode          bool
	compLimitSeconds  int
	compNumOfTestcase int
}

const (
	argSeed          = "seed"
	argStatistics    = "statistics"
	argCompMode      = "comp"
	argLimitSeconds  = "comptime"
	argNumOfTestcase = "compsize"
)

const (
	MinLimitSeconds     = 1
	MaxLimitSeconds     = 300
	DefaultLimitSeconds = 5

	MinNumOfTestcase     = 1
	MaxNumOfTestcase     = 100
	DefaultNumOfTestcase = 10
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

	flag.Parse()

	if args.compMode {
		args.validateCompMode()
	} else {
		args.validateNormalMode()
	}

	return args
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
