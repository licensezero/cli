package main

import "fmt"
import "github.com/licensezero/cli/subcommands"
import "github.com/mitchellh/go-homedir"
import "os"
import "sort"

var Rev string // Set via ldflags.

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
	"version":   subcommands.Version,
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
			if subcommand == "version" {
				value.Handler([]string{Rev}, paths)
			} else {
				value.Handler(os.Args[2:], paths)
			}
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
	usage := "Manage License Zero projects and dependencies.\n\nSubcommands:\n"
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
		usage += fmt.Sprintf("  %-"+fmt.Sprintf("%d", longestSubcommand)+"s %s\n", name, info.Description)
	}
	os.Stdout.WriteString(usage)
}
