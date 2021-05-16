package main

import (
	"flag"
	"fmt"

	"github.com/kjk/u"
)

var must = u.Must
var logf = u.Logf

func main() {
	u.CdUpDir("notion")
	//logf("currDirAbs: '%s'\n", u.CurrDirAbsMust())
	var (
		flgWc bool
	)
	{
		flag.BoolVar(&flgWc, "wc", false, "wc -l on source files")
		flag.Parse()
	}
}

var srcFiles = u.MakeAllowedFileFilterForExts(".go", ".js", ".html", ".css")
var excludeDirs = u.MakeExcludeDirsFilter("node_modules", "tmpdata")
var allFiles = u.MakeFilterAnd(srcFiles, excludeDirs)

func doLineCount() int {
	stats := u.NewLineStats()
	err := stats.CalcInDir(".", allFiles, true)
	if err != nil {
		fmt.Printf("doLineCount: stats.wcInDir() failed with '%s'\n", err)
		return 1
	}
	u.PrintLineStats(stats)
	return 0
}
