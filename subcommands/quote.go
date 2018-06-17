package subcommands

import "flag"
import "fmt"
import "os"

var Quote = Subcommand{
	Description: "Quote missing private licenses.",
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("quote", flag.ExitOnError)
		noNoncommercial := flagSet.Bool("no-noncommercial", false, "Ignore L0-NC dependencies.")
		noReciprocal := flagSet.Bool("no-reciprocal", false, "Ignore L0-R dependencies.")
		flagSet.Parse(args)
		projects, err := Inventory(paths, *noNoncommercial, *noReciprocal)
		if err != nil {
			os.Stderr.WriteString("Could not read dependeny tree.")
			os.Exit(1)
		} else {
			for _, project := range projects.Ignored {
				fmt.Println(project.Type + " " + project.Name + "@" + project.Version + " " + project.Path)
				manifest := project.Manifest
				fmt.Println("#" + manifest.ProjectID + " (" + manifest.Terms + ")")
			}
			os.Exit(0)
		}
	},
}
