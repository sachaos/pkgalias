package pkgalias

import (
	"fmt"
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
)

const doc = "pkgalias is ..."

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "pkgalias",
	Doc:  doc,
	Run:  runWithConfig(loadConfig()),
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

type fullpath = string
type alias = string

type config struct {
	settings []*expectedPackage `yaml:"settings"`
}

type expectedPackage struct {
	alias    alias    `yaml:"alias"`
	fullpath fullpath `yaml:"fullpath"`
}

func loadConfig() *config {
	var cnf config
	confFile, err := os.Open(".pkgalias.yaml")
	if err != nil {
		return &cnf
	}

	if err := yaml.NewDecoder(confFile).Decode(&cnf); err != nil {
		return &cnf
	}

	return &cnf
}

func runWithConfig(c *config) func (pass *analysis.Pass) (interface{}, error) {
	return func (pass *analysis.Pass) (interface{}, error) {
		inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

		var aliasDict = map[fullpath]alias{}
		for _, expectedPackage := range c.settings {
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
}
