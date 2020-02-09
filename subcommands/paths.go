package subcommands

import (
	"github.com/mitchellh/go-homedir"
	"os"
)

// Paths describes the paths in which the CLI is run.
type Paths struct {
	Home string
	CWD  string
}

// NewPaths fetches home and working directories.
func NewPaths() Paths {
	home, homeError := homedir.Dir()
	if homeError != nil {
		Fail("Could not find home directory.")
	}
	cwd, cwdError := os.Getwd()
	if cwdError != nil {
		Fail("Could not find working directory.")
	}
	return Paths{Home: home, CWD: cwd}
}
