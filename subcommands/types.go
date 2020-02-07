package subcommands

// Paths describes the paths in which the CLI is run.
type Paths struct {
	Home string
	CWD  string
}

// Subcommand describes a CLI subcommand.
type Subcommand struct {
	Tag         string
	Description string
	Handler     func([]string, Paths)
}
