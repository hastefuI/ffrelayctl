// Command ffrelayctl is a command line client for the Firefox Relay API.
//
// Deprecated: ffrelayctl has moved to go.hasteful.org/ffrelayctl. Install it
// with "go install go.hasteful.org/ffrelayctl@latest".
package main

import "github.com/hastefuI/ffrelayctl/cmd"

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
