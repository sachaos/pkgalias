package pkgalias

import (
	"fmt"
	"go/ast"
	"go/token"
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


type config struct {
	Settings []*expectedPackage `yaml:"settings"`
}

type expectedPackage struct {
	Alias    string `yaml:"alias"`
	Fullpath string `yaml:"fullpath"`
}

func loadConfig() *config {
	var cnf config
	confFile, err := os.Open(".pkgalias.yaml")
	if err != nil {
		fmt.Println("failed to load .pkgalias.yaml")
		return &cnf
	}

	if err := yaml.NewDecoder(confFile).Decode(&cnf); err != nil {
		fmt.Println("failed to decode config")
		return &cnf
	}

	return &cnf
}

func runWithConfig(c *config) func (pass *analysis.Pass) (interface{}, error) {
	return func (pass *analysis.Pass) (interface{}, error) {
		inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

		var aliasDict = map[string]string{}
		for _, expectedPackage := range c.Settings {
			aliasDict[expectedPackage.Fullpath] = expectedPackage.Alias
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
								msg := fmt.Sprintf(`invalid alias: use "%s" instead of "%s"`, expAlias, actName)
								pass.Report(analysis.Diagnostic{
									Pos:            n.Pos(),
									Message:        msg,
									SuggestedFixes: []analysis.SuggestedFix{
										replacementSuggest(ident, expAlias, msg),
									},
								})
							}
						}
					}
				}
			case *ast.ImportSpec:
				path := strings.Trim(n.Path.Value, "\"")
				if expAlias, ok := aliasDict[path]; ok {
					if expAlias != "" {
						if n.Name != nil {
							if n.Name.Name != expAlias {
								msg := fmt.Sprintf(`invalid alias: package name should be "%s", replace this`, expAlias)
								pass.Report(analysis.Diagnostic{
									Pos:            n.Name.Pos(),
									Message:        msg,
									SuggestedFixes: []analysis.SuggestedFix{
										replacementSuggest(n.Name, expAlias, msg),
									},
								})
							}
						} else {
							msg := fmt.Sprintf(`invalid alias: package name should be "%s", insert this.`, expAlias)
							pass.Report(analysis.Diagnostic{
								Pos:            n.Pos(),
								Message:        msg,
								SuggestedFixes: []analysis.SuggestedFix{
									insertSuggest(n.Pos(), expAlias, msg),
								},
							})
						}
					} else {
						if n.Name != nil {
							msg := `invalid alias: package name should not be specified, delete this.`
							pass.Report(analysis.Diagnostic{
								Pos:            n.Name.Pos(),
								Message:        msg,
								SuggestedFixes: []analysis.SuggestedFix{
									deleteSuggest(n.Name, msg),
								},
							})
						}
					}
				}
			}
		})

		return nil, nil
	}
}

func replacementSuggest(node ast.Node, newText string, message string) analysis.SuggestedFix {
	return analysis.SuggestedFix{
		Message:   message,
		TextEdits: []analysis.TextEdit{
			{
				Pos:     node.Pos(),
				End:     node.End(),
				NewText: []byte(newText),
			},
		},
	}
}

func insertSuggest(pos token.Pos, newText string, message string) analysis.SuggestedFix {
	return analysis.SuggestedFix{
		Message:   message,
		TextEdits: []analysis.TextEdit{
			{
				Pos:     pos,
				End:     pos,
				NewText: []byte(newText),
			},
		},
	}
}

func deleteSuggest(node ast.Node, message string) analysis.SuggestedFix {
	return analysis.SuggestedFix{
		Message:   message,
		TextEdits: []analysis.TextEdit{
			{
				Pos:     node.Pos(),
				End:     node.End(),
				NewText: []byte(""),
			},
		},
	}
}
