// SUM10-SOLVER
// author: Leonardone @ NEETSDKASU

package main

func main() {
	args := parseArgs()

	var err error

	if args.isSolverMode() {
		err = runSolver(args.seed, args.solverName, args.solverLimitSeconds, args.solverTryCount)
	} else if args.compMode {
		err = runComp(args.compLimitSeconds, args.compNumOfTestcase, args.seed)
	} else {
		err = run(args.seed, args.withStatistics)
	}

	if err != nil {
		panic(err)
	}
}
