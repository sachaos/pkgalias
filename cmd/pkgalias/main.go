package main

import (
	"golang.org/x/tools/go/analysis/unitchecker"
	"pkgalias"
)

func main() { unitchecker.Main(pkgalias.Analyzer) }
