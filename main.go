package main

import "go.hasteful.org/ffrelayctl/cmd"

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	versionInfo := cmd.VersionInfo{
		Version: version,
		Commit:  commit,
		Date:    date,
	}
	cmd.Execute(versionInfo)
}
