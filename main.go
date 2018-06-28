package main

import "fmt"
import "github.com/licensezero/cli/subcommands"
import "github.com/mitchellh/go-homedir"
import "os"
import "sort"

var Rev string // Set via ldflags.

var commands = map[string]subcommands.Subcommand{
	"backup":   subcommands.Backup,
	"buy":      subcommands.Buy,
	"bugs":     subcommands.Bugs,
	"identify": subcommands.Identify,
	"import":   subcommands.Import,
	"license":  subcommands.License,
	"lock":     subcommands.Lock,
	"offer":    subcommands.Offer,
	"quote":    subcommands.Quote,
	"readme":   subcommands.README,
	"register": subcommands.Register,
	"reprice":  subcommands.Reprice,
	"reset":    subcommands.Reset,
	"retract":  subcommands.Retract,
	"sponsor":  subcommands.Sponsor,
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
	os.Stdout.WriteString("Manage License Zero projects and dependencies.\n\nSubcommands:\n")
	buyer := map[string]subcommands.Subcommand{}
	seller := map[string]subcommands.Subcommand{}
	misc := map[string]subcommands.Subcommand{}
	for key, value := range commands {
		switch value.Tag {
		case "buyer":
			buyer[key] = value
		case "seller":
			seller[key] = value
		default:
			misc[key] = value
		}
	}
	listSubcommands("For Buyers", buyer)
	listSubcommands("For Sellers", seller)
	listSubcommands("Miscellaneous", misc)
}

func listSubcommands(header string, list map[string]subcommands.Subcommand) {
	os.Stdout.WriteString("\n  " + header + ":\n\n")
	longestSubcommand := 0
	var names []string
	for name, _ := range list {
		length := len(name) + 1
		if length > longestSubcommand {
			longestSubcommand = length
		}
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		info := list[name]
		fmt.Printf("  %-"+fmt.Sprintf("%d", longestSubcommand)+"s %s\n", name, info.Description)
	}
}
