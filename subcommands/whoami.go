package subcommands

import (
	"licensezero.com/licensezero/user"
)

const whoAmIDescription = "Show your identity information."

// WhoAmI prints identity information.
var WhoAmI = &Subcommand{
	Tag:         "misc",
	Description: whoAmIDescription,
	Handler: func(env Environment) int {
		identity, err := user.ReadIdentity()
		if err != nil {
			env.Stderr.WriteString("Could not read identity file.\n")
			return 1
		}
		env.Stdout.WriteString("Name: " + identity.Name + "\n")
		env.Stdout.WriteString("Jurisdiction: " + identity.Jurisdiction + "\n")
		env.Stdout.WriteString("E-Mail: " + identity.EMail + "\n")
		return 0
	},
}
