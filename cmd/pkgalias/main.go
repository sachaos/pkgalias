package main

import (
	"github.com/sachaos/pkgalias"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(pkgalias.Analyzer) }
