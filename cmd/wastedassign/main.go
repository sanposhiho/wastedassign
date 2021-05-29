package main

import (
	"github.com/sanposhiho/wastedassign/v2"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(wastedassign.Analyzer) }
