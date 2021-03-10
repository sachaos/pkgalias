package pkgalias

import (
	"fmt"
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"strings"
)

const doc = "pkgalias is ..."

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "pkgalias",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

type fullpath = string
type alias = string

type expectedPackage struct {
	alias alias
	fullpath fullpath
}

var testSets = []*expectedPackage{
	{
		alias: "m",
		fullpath: "math",
	},
	{
		alias: "format",
		fullpath: "fmt",
	},
	{
		alias: "opesys",
		fullpath: "os",
	},
	{
		alias: "",
		fullpath: "net",
	},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	var aliasDict = map[fullpath]alias{}
	for _, expectedPackage := range testSets {
		aliasDict[expectedPackage.fullpath] = expectedPackage.alias
	}

	nodeFilter := []ast.Node{
		(*ast.ImportSpec)(nil),
		(*ast.SelectorExpr)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.SelectorExpr:
			var ident *ast.Ident
			var ok bool

			if ident, ok = n.X.(*ast.Ident); !ok {
				return
			}

			if obj := pass.TypesInfo.ObjectOf(ident); obj != nil {
				if pkg, ok := obj.(*types.PkgName); ok {
					if expAlias, ok := aliasDict[pkg.Imported().Path()]; ok {
						if expAlias == "" {
							expAlias = pkg.Imported().Name()
						}

						actName := ident.Name
						if expAlias != actName {
							pass.Reportf(n.Pos(), fmt.Sprintf(`invalid alias: use "%s" instead of "%s"`, expAlias, actName))
						}
					}
				}
			}
		case *ast.ImportSpec:
			path := strings.Trim(n.Path.Value, "\"")
			if expAlias, ok := aliasDict[path]; ok {
				if n.Name != nil {
					if expAlias != n.Name.Name {
						pass.Reportf(n.Pos(), fmt.Sprintf(`invalid alias: package name should be "%s"`, expAlias))
					}
				} else {
					if expAlias != "" {
						pass.Reportf(n.Pos(), `invalid alias: package name should not be specified`)
					}
				}
			}
		}
	})

	return nil, nil
}
