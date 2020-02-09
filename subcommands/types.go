package subcommands

// Subcommand describes a CLI subcommand.
type Subcommand struct {
	Tag         string
	Description string
	Handler     func([]string, Paths)
}
