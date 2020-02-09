package subcommands

import (
	"os"
)

// Subcommand describes a CLI subcommand.
type Subcommand struct {
	Tag         string
	Description string
	Handler     func(args []string, stdin, stdout, stderr *os.File) int
}
