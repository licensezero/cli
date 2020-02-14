package subcommands

import (
	"io"
	"net/http"
)

// Environment wraps arguments to subcommands.
type Environment struct {
	Rev       string
	Arguments []string
	Stdin     InputDevice
	Stdout    io.StringWriter
	Stderr    io.StringWriter
	Client    *http.Client
}

// Subcommand describes a CLI subcommand.
type Subcommand struct {
	Tag         string
	Description string
	Handler     func(Environment) int
}
