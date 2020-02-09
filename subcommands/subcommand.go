package subcommands

import (
	"io"
	"net/http"
)

// Subcommand describes a CLI subcommand.
type Subcommand struct {
	Tag         string
	Description string
	Handler     func(
		args []string, stdin InputDevice,
		stdout io.StringWriter,
		stderr io.StringWriter,
		client *http.Client,
	) int
}
