package wastedassign

import (
	"go/token"

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
		for _, local := range sf.Locals {
			var isAfterStore bool
			var storePos token.Pos
			for _, rf := range *(local.Referrers()) {
				switch rf.(type) {
				case *ssa.Store:
					if isAfterStore {
						pass.Reportf(storePos, "Inefficient assignment")
					}
					storePos = rf.Pos()
					isAfterStore = true
				default:
					isAfterStore = false
				}
			}
			// find the value, reassigned but never used afterwards
			if isAfterStore {
				pass.Reportf(storePos, "reassigned, but never used afterwards")
			}
		}
	}
	return nil, nil
}
