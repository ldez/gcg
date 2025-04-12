package main

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

const (
	defaultVersion = "dev"
	unknownInfo    = "unknown"
)

var (
	version = defaultVersion
	commit  = unknownInfo
	date    = unknownInfo
)

// DisplayVersion Displays version.
func DisplayVersion() {
	info, ok := debug.ReadBuildInfo()
	if ok {
		if version == defaultVersion && info.Main.Version != "" {
			version = info.Main.Version
		}

		for _, setting := range info.Settings {
			switch {
			case setting.Key == "vcs.time" && date == unknownInfo:
				date = setting.Value
			case setting.Key == "vcs.revision" && commit == unknownInfo:
				commit = setting.Value
			}
		}
	}

	fmt.Printf(`gcg:
 version     : %s
 commit      : %s
 build date  : %s
 go version  : %s
 go compiler : %s
 platform    : %s/%s
`, version, commit, date, runtime.Version(), runtime.Compiler, runtime.GOOS, runtime.GOARCH)
}
