package main

import "fmt"
import "os"
import "./subcommands"

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
	arguments := os.Args
	if len(arguments) > 1 {
		subcommand := os.Args[1]
		if value, ok := commands[subcommand]; ok {
			value.Handler(os.Args[2:])
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
	fmt.Println("Manage License Zero dependences.")
	fmt.Println("")
	fmt.Println("Subcommands:")
	for name, info := range commands {
		fmt.Println("\t" + name + ": " + info.Description)
	}
}
