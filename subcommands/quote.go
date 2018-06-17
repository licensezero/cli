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
		results, err := Inventory(paths, *noNoncommercial, *noReciprocal)
		if err != nil {
			os.Stderr.WriteString("Could not read dependeny tree.")
			os.Exit(1)
		} else {
			for _, result := range results.Ignored {
				metadata := result.Metadata
				fmt.Println(result.Type + " " + metadata.Name + "@" + metadata.Version + " " + result.Path)
				for _, projectManifest := range metadata.ProjectManifests {
					for _, licenseManifest := range projectManifest.LicenseManifests {
						fmt.Println(licenseManifest.Terms)
					}
				}
			}
			os.Exit(0)
		}
	},
}
