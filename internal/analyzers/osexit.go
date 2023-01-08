package analyzers

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// Analyzer is an analyzer that checks if os.Exit is called from main.
var Analyzer = &analysis.Analyzer{
	Name: "osExitCheck",
	Doc:  "check if os.Exit is called from main",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, f := range pass.Files {
		if f.Name.String() != "main" {
			continue
		}
		ast.Inspect(f, func(n ast.Node) bool {
			if call, ok := n.(*ast.CallExpr); ok {
				var fun *ast.SelectorExpr
				if fun, ok = call.Fun.(*ast.SelectorExpr); !ok || fun.Sel.Name != "Exit" {
					return true
				}
				var id *ast.Ident
				if id, ok = fun.X.(*ast.Ident); !ok || id.Name != "os" {
					return true
				}
				pass.Report(analysis.Diagnostic{
					Pos:     call.Pos(),
					Message: fmt.Sprintf("os.Exit(...) call in main"),
				})
			}
			return true
		})
	}
	return nil, nil
}
