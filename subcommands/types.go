package subcommands

type Subcommand struct {
	Description string
	Handler     func([]string)
}
