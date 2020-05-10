package main

import "fmt"
import "licensezero.com/cli/subcommands"
import "github.com/mitchellh/go-homedir"
import "os"
import "sort"

// Rev represents the current build revision.  Set via ldflags.
var Rev string

var commands = map[string]*subcommands.Subcommand{
	"backup":   subcommands.Backup,
	"bugs":     subcommands.Bugs,
	"identify": subcommands.Identify,
	"latest":   subcommands.Latest,
	"lock":     subcommands.Lock,
	"offer":    subcommands.Offer,
	"projects": subcommands.Projects,
	"register": subcommands.Register,
	"reprice":  subcommands.Reprice,
	"reset":    subcommands.Reset,
	"retract":  subcommands.Retract,
	"token":    subcommands.Token,
	"version":  subcommands.Version,
	"waive":    subcommands.Waive,
	"whoami":   subcommands.WhoAmI,
}

func main() {
	home, homeError := homedir.Dir()
	if homeError != nil {
		subcommands.Fail("Could not find home directory.")
	}
	cwd, cwdError := os.Getwd()
	if cwdError != nil {
		subcommands.Fail("Could not find working directory.")
	}
	paths := subcommands.Paths{Home: home, CWD: cwd}
	arguments := os.Args
	if len(arguments) > 1 {
		subcommand := os.Args[1]
		if value, ok := commands[subcommand]; ok {
			if subcommand == "version" || subcommand == "latest" {
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
	os.Stdout.WriteString("Manage License Zero projects and dependencies.\n\nSubcommands:\n")
	longestSubcommand := 0
	var names []string
	for name := range commands {
		length := len(name) + 1
		if length > longestSubcommand {
			longestSubcommand = length
		}
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		info := commands[name]
		fmt.Printf("  %-"+fmt.Sprintf("%d", longestSubcommand)+"s %s\n", name, info.Description)
	}
}
