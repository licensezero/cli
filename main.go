package main

import "fmt"
import "os"
import "./subcommands"

func main() {
	subcommands := map[string]func([]string){
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

	arguments := os.Args
	if len(arguments) > 1 {
		subcommand := os.Args[1]
		if value, ok := subcommands[subcommand]; ok {
			value(os.Args[2:])
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
	fmt.Println(`Manage License Zero dependences.

Usage:
  licensezero (--help | --version)
  licensezero buy [--do-not-open]
  licensezero identify <name> <jurisdiction> <email>
  licensezero import-license <file>
  licensezero import-waiver <file>
  licensezero license <project id> (--noncommercial | --reciprocal)
  licensezero lock <project id> <date>
  licensezero offer <PRICE> [--relicense CENTS]
  licensezero reprice <PRICE> [--relicense CENTS]
  licensezero purchased <URL>
  licensezero quote [--no-noncommercial] [--no-reciprocal]
  licensezero readme
  licensezero register
  licensezero reset-token
  licensezero retract <project id>
  licensezero set-licensor-id <licensor ID>
  licensezero sponsor <project id> [--do-not-open]
  licensezero waive <project id> -b NAME -j CODE -d DAYS
  licensezero whoami`)
}
