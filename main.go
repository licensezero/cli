package main

import "fmt"
import "github.com/licensezero/cli/subcommands"
import "github.com/mitchellh/go-homedir"
import "os"
import "sort"

var commands = map[string]subcommands.Subcommand{
	"buy":       subcommands.Buy,
	"identify":  subcommands.Identify,
	"import":    subcommands.Import,
	"license":   subcommands.License,
	"lock":      subcommands.Lock,
	"offer":     subcommands.Offer,
	"purchased": subcommands.Purchased,
	"quote":     subcommands.Quote,
	"readme":    subcommands.README,
	"register":  subcommands.Register,
	"reprice":   subcommands.Reprice,
	"reset":     subcommands.Reset,
	"retract":   subcommands.Retract,
	"sponsor":   subcommands.Sponsor,
	"token":     subcommands.Token,
	"waive":     subcommands.Waive,
	"whoami":    subcommands.WhoAmI,
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
	longestSubcommand := 0
	var names []string
	for name, _ := range commands {
		if len(name) > longestSubcommand {
			longestSubcommand = len(name) + 1
		}
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		info := commands[name]
		fmt.Printf("  %-"+fmt.Sprintf("%d", longestSubcommand)+"s %s\n", name, info.Description)
	}
}
