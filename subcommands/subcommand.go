package subcommands

import (
	"io"
	"licensezero.com/licensezero/api"
)

// Subcommand describes a CLI subcommand.
type Subcommand struct {
	Tag         string
	Description string
	Handler     func(
		args []string, stdin InputDevice,
		stdout io.StringWriter,
		stderr io.StringWriter,
		client api.Client,
	) int
}
