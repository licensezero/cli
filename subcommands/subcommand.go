package subcommands

import (
	"io"
)

// Subcommand describes a CLI subcommand.
type Subcommand struct {
	Tag         string
	Description string
	Handler     func(args []string, stdin InputDevice, stdout, stderr io.StringWriter) int
}
