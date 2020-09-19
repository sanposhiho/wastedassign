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
			blCopy := *bl
			for _, ist := range bl.Instrs {
				blCopy.Instrs = rmInstrFromInstrs(blCopy.Instrs, ist)
				switch ist.(type) {
				case *ssa.Store:
					var buf [10]*ssa.Value
					for _, op := range ist.Operands(buf[:0]) {
						if (*op) != nil && opInLocals(sf.Locals, op) {
							if reason := isNextOperationToOpIsStore([]*ssa.BasicBlock{&blCopy}, op, nil); reason != notWasted {
								pass.Reportf(ist.Pos(), reason.String())
							}
						}
					}
				}
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
func isNextOperationToOpIsStore(bls []*ssa.BasicBlock, currentOp *ssa.Value, haveCheckedMap *map[int]bool) wastedReason {
	wastedReasons := []wastedReason{}
	wastedReasonsCurrentBls := []wastedReason{}

	if haveCheckedMap == nil {
		haveCheckedMap = &map[int]bool{}
	}

	for _, bl := range bls {
		if (*haveCheckedMap)[bl.Index] == true {
			continue
		}
		(*haveCheckedMap)[bl.Index] = true
		breakFlag := false
		for _, ist := range bl.Instrs {
			if breakFlag {
				break
			}
			switch w := ist.(type) {
			case *ssa.Store:
				var buf [10]*ssa.Value
				for _, op := range ist.Operands(buf[:0]) {
					if *op == *currentOp {
						if w.Addr.Name() == (*currentOp).Name() {
							wastedReasonsCurrentBls = append(wastedReasonsCurrentBls, reassignedSoon)
							breakFlag = true
							break
						} else {
							return notWasted
						}
					}
				}
			default:
				var buf [10]*ssa.Value
				for _, op := range ist.Operands(buf[:0]) {
					if *op == *currentOp {
						// 連続storeではなかった
						return notWasted
					}
				}
			}
		}
		if len(bl.Succs) != 0 && !breakFlag {
			wastedReason := isNextOperationToOpIsStore(rmSameBlock(bl.Succs, bl), currentOp, haveCheckedMap)
			if wastedReason == notWasted {
				return notWasted
			}
			wastedReasons = append(wastedReasons, wastedReason)
		}
	}

	wastedReasons = append(wastedReasons, wastedReasonsCurrentBls...)

	if len(wastedReasons) != 0 {
		if containReassignedSoon(wastedReasons) {
			return reassignedSoon
		}
		return noUseUntilReturn
	}
	return noUseUntilReturn
}

func rmSameBlock(bls []*ssa.BasicBlock, currentBl *ssa.BasicBlock) []*ssa.BasicBlock {
	rto := []*ssa.BasicBlock{}

	for _, bl := range bls {
		if bl != currentBl {
			rto = append(rto, bl)
		}
	}
	return rto
}

func containNotWasted(ws []wastedReason) bool {
	for _, w := range ws {
		if w == notWasted {
			return true
		}
	}
	return false
}

func containReassignedSoon(ws []wastedReason) bool {
	for _, w := range ws {
		if w == reassignedSoon {
			return true
		}
	}
	return false
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
