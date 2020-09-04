package wastedassign

import (
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
			blCopy := bl
			for _, ist := range bl.Instrs {
				switch ist.(type) {
				case *ssa.Store:
					var buf [10]*ssa.Value
					for _, op := range ist.Operands(buf[:0]) {
						if (*op) != nil && opInLocals(sf.Locals, op) {
							if reason := isNextOperationToOpIsStore([]*ssa.BasicBlock{blCopy}, op, 0, sf.Locals); reason != notWasted {
								pass.Reportf(ist.Pos(), reason.String())
							}
						}
					}
				}
				blCopy.Instrs = rmInstrFromInstrs(blCopy.Instrs, ist)
			}
		}
	}
	return nil, nil
}

type wastedReason string

const (
	noUseUntilReturn wastedReason = "reassigned, but never used afterwards"
	reassignedSoon   wastedReason = "wasted assignment"
	notWasted        wastedReason = ""
)

func (wr wastedReason) String() string {
	switch wr {
	case noUseUntilReturn:
		return "reassigned, but never used afterwards"
	case reassignedSoon:
		return "wasted assignment"
	case notWasted:
		return ""
	}
	return ""
}

// 次のblockまでみて、storeが連続であるかを調べる
func isNextOperationToOpIsStore(bls []*ssa.BasicBlock, currentOp *ssa.Value, depth int, locals []*ssa.Alloc) wastedReason {

	// depth == 0の時は少なくとも一つstoreが見つかるので一回めは飛ばすためのflag
	skipStore := depth == 0

	// blsが全てSuccsを持っていなかった場合を判別するためのflag
	noNextSuccs := true

	for _, bl := range bls {
		for _, ist := range bl.Instrs {
			switch w := ist.(type) {
			case *ssa.Store:
				var buf [10]*ssa.Value
				for _, op := range ist.Operands(buf[:0]) {
					if *op == *currentOp {
						if w.Addr.Name() == (*currentOp).Name() {
							if !skipStore {
								// 連続storeなのでtrue
								return reassignedSoon
							}
							skipStore = false
						} else if w.Val.Name() == (*currentOp).Name() {
							return ""
						}
					}
				}
			default:
				var buf [10]*ssa.Value
				for _, op := range ist.Operands(buf[:0]) {
					if *op == *currentOp {
						// 連続storeではなかった
						return ""
					}
				}
			}
		}

		if len(bl.Succs) != 0 {
			noNextSuccs = false
			wastedReason := isNextOperationToOpIsStore(bl.Succs, currentOp, depth+1, locals)
			if wastedReason != "" {
				return wastedReason
			}
			// SuccsにcurrentOpに対する操作がなかった
		}
	}

	if noNextSuccs {
		return noUseUntilReturn
	}
	return notWasted
}

func rmInstrFromInstrs(instrs []ssa.Instruction, instrToRm ssa.Instruction) []ssa.Instruction {
	var rto []ssa.Instruction
	for _, i := range instrs {
		if i != instrToRm {
			rto = append(rto, i)
		}
	}
	return rto
}

func opInLocals(locals []*ssa.Alloc, op *ssa.Value) bool {
	for _, l := range locals {
		if *op == ssa.Value(l) {
			return true
		}
	}
	return false
}
