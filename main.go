package main

import "flag"
import "fmt"
import "os"

func main() {
	subcommands := map[string]func([]string){
		"buy":   buy,
		"quote": quote,
		// "buy": "[--do-not-open]",
		// "identify": "<name> <jurisdiction> <email>",
		// "import-license": "<file>",
		// "import-waiver": "<file>",
		// "license": "<project id> (--noncommercial | --reciprocal)",
		// "lock": "<project id> <date>",
		// "offer": "<PRICE> [--relicense CENTS]",
		// "reprice": "<PRICE> [--relicense CENTS]",
		// "purchased": "<URL>",
		// "quote": "[--no-noncommercial] [--no-reciprocal]",
		// "readme": "",
		// "register": "",
		// "reset-token": "",
		// "retract": "<project id>",
		// "set-licensor-id": "<licensor ID>",
		// "sponsor": "<project id> [--do-not-open]",
		// "waive": "<project id> -b NAME -j CODE -d DAYS",
		// "whoami": "",
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

func quote(args []string) {
	flagSet := flag.NewFlagSet("quote", flag.ExitOnError)
	noNoncommercial := flagSet.Bool("no-noncommercial", false, "Ignore L0-NC dependencies.")
	noReciprocal := flagSet.Bool("no-reciprocal", false, "Ignore L0-R dependencies.")
	flagSet.Parse(args)
	if *noNoncommercial {
		fmt.Println("No L0-NC")
	}
	if *noReciprocal {
		fmt.Println("No L0-R")
	}
	os.Exit(0)
}

func buy(args []string) {
	flagSet := flag.NewFlagSet("buy", flag.ExitOnError)
	doNotOpen := flagSet.Bool("do-not-open", false, "Do not open checkout page.")
	flagSet.Parse(args)
	if *doNotOpen {
		fmt.Println("not opening")
	}
	os.Exit(0)
}
