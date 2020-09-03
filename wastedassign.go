package wastedassign

import (
	"fmt"

	"github.com/sanposhiho/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ssa"
)

const doc = "wastedassign finds wasted assignment statements."

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "wastedassign",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		buildssa.Analyzer,
	},
}

func run(pass *analysis.Pass) (interface{}, error) {
	s := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA)

	for _, sf := range s.SrcFuncs {
		for _, bl := range sf.Blocks {
			for _, ist := range bl.Instrs {
				switch ist.(type) {
				case *ssa.Store:
					var buf [10]*ssa.Value
					for _, op := range ist.Operands(buf[:0]) {
						if (*op) != nil && opInLocals(sf.Locals, op) {
							if isNextOperationToOpIsStore([]*ssa.BasicBlock{bl}, op, 0) {
								pass.Reportf(ist.Pos(), "wasted assignment")
							}
						}
					}
				}
			}
		}
	}
	return nil, nil
}

func opInLocals(locals []*ssa.Alloc, op *ssa.Value) bool {
	for _, l := range locals {
		if *op == ssa.Value(l) {
			return true
		}
	}
	return false
}

// 次のblockまでみて、storeが連続であるかを調べる
func isNextOperationToOpIsStore(bls []*ssa.BasicBlock, currentOp *ssa.Value, depth int) bool {

	// depth == 0の時は少なくとも一つstoreが見つかるので一回めは飛ばす
	skipStore := true

	// SuccsのSuccsが全てnilだった場合はtrueを返したい
	flag := true
	for _, bl := range bls {
		for _, ist := range bl.Instrs {
			fmt.Print("\n")
			fmt.Print("\n")
			fmt.Print(ist)
			fmt.Print("\n")

			switch w := ist.(type) {
			case *ssa.Store:
				fmt.Print("store\n")
				fmt.Print(w.Addr)
				fmt.Print("store\n")
				fmt.Print(w.Val)
				var buf [10]*ssa.Value
				fmt.Print(ist.Operands(buf[:0]))
				for _, op := range ist.Operands(buf[:0]) {
					if op == currentOp {
						if !skipStore {
							// 連続storeなのでtrue
							return true
						}
						skipStore = false
					}
				}
			case *ssa.MakeClosure:
				fmt.Print("makeuuuuuuuuuuuuuuuuuuuu\n")
			default:
				fmt.Print("default\n")
				var buf [10]*ssa.Value
				fmt.Print(ist.Operands(buf[:0]))
				for _, op := range ist.Operands(buf[:0]) {
					if op == currentOp {
						// 連続storeではなかった
						return false
					}
				}
			}
		}
		// 次のBlockにcurrentOpに対する操作がなかった

		if bl.Succs != nil {
			flag = false
		}

		if bl.Succs != nil && isNextOperationToOpIsStore(bl.Succs, currentOp, depth+1) {
			// その次にはあった&true
			return true
		}
	}
	return flag
}
