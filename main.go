package main

import "github.com/licensezero/cli/subcommands"
import "fmt"
import "github.com/mitchellh/go-homedir"
import "os"

var commands = map[string]subcommands.Subcommand{
	"buy":             subcommands.Buy,
	"identify":        subcommands.Identify,
	"import":          subcommands.Import,
	"license":         subcommands.License,
	"lock":            subcommands.Lock,
	"offer":           subcommands.Offer,
	"purchased":       subcommands.Purchased,
	"quote":           subcommands.Quote,
	"readme":          subcommands.README,
	"register":        subcommands.Register,
	"reprice":         subcommands.Reprice,
	"reset":           subcommands.Reset,
	"retract":         subcommands.Retract,
	"set-licensor-id": subcommands.SetLicensorID,
	"sponsor":         subcommands.Sponsor,
	"waive":           subcommands.Waive,
	"whoami":          subcommands.WhoAmI,
}

func main() {
	home, homeError := homedir.Dir()
	if homeError != nil {
		os.Stderr.WriteString("Could not find home directory.\n")
		os.Exit(1)
	}
	cwd, cwdError := os.Getwd()
	if cwdError != nil {
		os.Stderr.WriteString("Could not find working directory.\n")
		os.Exit(1)
	}
	paths := subcommands.Paths{Home: home, CWD: cwd}
	arguments := os.Args
	if len(arguments) > 1 {
		subcommand := os.Args[1]
		if value, ok := commands[subcommand]; ok {
			value.Handler(os.Args[2:], paths)
		} else {
			showUsage()
			os.Exit(1)
		}
	} else {
		showUsage()
		os.Exit(0)
	}
}

func showUsage() {
	fmt.Println("Manage License Zero dependencies.")
	fmt.Println("")
	fmt.Println("Subcommands:")
	for name, info := range commands {
		fmt.Println("\t" + name + ": " + info.Description)
	}
}
