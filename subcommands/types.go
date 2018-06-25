package subcommands

type Paths struct {
	Home string
	CWD  string
}

type Subcommand struct {
	Tag         string
	Description string
	Handler     func([]string, Paths)
}
