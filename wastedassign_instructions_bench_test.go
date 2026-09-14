package wastedassign_test

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"strings"
	"testing"

	"github.com/sanposhiho/wastedassign/v2"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

func BenchmarkAnalyzerInstructions(b *testing.B) {
	for _, count := range []int{16, 128, 1024} {
		b.Run(fmt.Sprintf("stores=%d", count), func(b *testing.B) {
			var source strings.Builder
			source.WriteString("package p\nfunc F() int {\nvalue := 0\n")
			for i := range count {
				fmt.Fprintf(&source, "value = %d\n", i)
			}
			source.WriteString("return value\n}\n")

			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "p.go", source.String(), 0)
			require.NoError(b, err)

			info := &types.Info{
				Types:      make(map[ast.Expr]types.TypeAndValue),
				Defs:       make(map[*ast.Ident]types.Object),
				Uses:       make(map[*ast.Ident]types.Object),
				Implicits:  make(map[ast.Node]types.Object),
				Selections: make(map[*ast.SelectorExpr]*types.Selection),
				Scopes:     make(map[ast.Node]*types.Scope),
			}
			var config types.Config
			files := []*ast.File{file}
			pkg, err := config.Check("p", fset, files, info)
			require.NoError(b, err)

			pass := &analysis.Pass{
				Analyzer:  wastedassign.Analyzer,
				Fset:      fset,
				Files:     files,
				Pkg:       pkg,
				TypesInfo: info,
				Report:    func(analysis.Diagnostic) {},
				ResultOf:  make(map[*analysis.Analyzer]any),
			}
			inspection, err := inspect.Analyzer.Run(pass)
			require.NoError(b, err)
			pass.ResultOf[inspect.Analyzer] = inspection
			b.ReportAllocs()

			for b.Loop() {
				_, err := wastedassign.Analyzer.Run(pass)
				require.NoError(b, err)
			}
		})
	}
}
