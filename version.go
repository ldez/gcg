package main

import (
	"fmt"
	"runtime"
)

var (
	version = "dev"
	commit  = "I don't remember exactly"
	date    = "I don't remember exactly"
)

// DisplayVersion Displays version.
func DisplayVersion() {
	fmt.Printf(`gcg:
 version     : %s
 commit      : %s
 build date  : %s
 go version  : %s
 go compiler : %s
 platform    : %s/%s
`, version, commit, date, runtime.Version(), runtime.Compiler, runtime.GOOS, runtime.GOARCH)
}
