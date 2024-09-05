// analyzers check a call of os.Exit inside main function.
package analyzers

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

var OsExitAnalyzer = &analysis.Analyzer{
	Name: "exitcheck",
	Doc:  "check for os.Exit",
	Run:  run,
}

type Pass struct {
	Fset  *token.FileSet
	Pkg   *types.Package
	Files []*ast.File
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		if pass.Pkg.Name() != "main" {
			continue
		}

		var ParentFunc string
		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.FuncDecl:
				ParentFunc = x.Name.Name

			case *ast.ExprStmt:
				if ParentFunc != "main" {
					return true
				}

				c, ok := x.X.(*ast.CallExpr)
				if !ok {
					return true
				}

				s, ok := c.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}

				i, ok := s.X.(*ast.Ident)
				if !ok {
					return true
				}

				if i.Name == "os" && s.Sel.Name == "Exit" {
					// fmt.Printf("exitcheck: os.Exit() in %v\n", pass.Fset.Position())
					pass.Reportf(x.Pos(), "exit function in main file")
				}
			}
			return true
		})
	}
	return nil, nil
}
