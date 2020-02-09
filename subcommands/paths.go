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
func NewPaths() (paths Paths, err error) {
	home, err := homedir.Dir()
	if err != nil {
		return
	}
	cwd, err := os.Getwd()
	if err != nil {
		return
	}
	return Paths{Home: home, CWD: cwd}, nil
}
