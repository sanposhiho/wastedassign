package main

import (
	"github.com/sanposhiho/wastedassign"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(wastedassign.Analyzer) }
