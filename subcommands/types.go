package subcommands

type Paths struct {
	Home string
	CWD  string
}

type Subcommand struct {
	Description string
	Handler     func([]string, Paths)
}
