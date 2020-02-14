package subcommands

const versionDescription = "Print version."

// Version prints the CLI version.
var Version = &Subcommand{
	Tag:         "misc",
	Description: versionDescription,
	Handler: func(env Environment) int {
		if env.Rev == "" {
			env.Stdout.WriteString("Development Build\n")
		} else {
			env.Stdout.WriteString(env.Rev)
		}
		return 0
	},
}
