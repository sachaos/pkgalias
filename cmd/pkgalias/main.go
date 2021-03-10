package main

import (
	"golang.org/x/tools/go/analysis/multichecker"
	"github.com/sachaos/pkgalias"
)

func main() { multichecker.Main(pkgalias.Analyzer) }
