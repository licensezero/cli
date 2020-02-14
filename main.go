package main

import (
	"fmt"
	"io"
	"licensezero.com/licensezero/subcommands"
	"net/http"
	"os"
	"sort"
)

// Rev represents the current build revision.  Set via ldflags.
var Rev string

var commands = map[string]*subcommands.Subcommand{
	"backup":   subcommands.Backup,
	"bugs":     subcommands.Bugs,
	"buy":      subcommands.Buy,
	"identify": subcommands.Identify,
	"import":   subcommands.Import,
	"latest":   subcommands.Latest,
	"quote":    subcommands.Quote,
	"register": subcommands.Register,
	"reset":    subcommands.Reset,
	"token":    subcommands.Token,
	"version":  subcommands.Version,
	"verify":   subcommands.Verify,
	"whoami":   subcommands.WhoAmI,
}

func main() {
	code := run(
		os.Args[1:],
		&subcommands.StandardInputDevice{File: os.Stdin},
		os.Stdout,
		os.Stderr,
		&http.Client{Transport: http.DefaultTransport},
	)
	os.Exit(code)
}

func run(
	arguments []string,
	input subcommands.InputDevice,
	stdout, stderr io.StringWriter,
	client *http.Client,
) int {
	if len(arguments) > 0 {
		subcommand := arguments[0]
		if value, ok := commands[subcommand]; ok {
			return value.Handler(subcommands.Environment{
				Rev:       Rev,
				Arguments: arguments[1:],
				Stdin:     input,
				Stdout:    stdout,
				Stderr:    stderr,
				Client:    client,
			})
		}
		showUsage(stdout)
		return 1
	}
	showUsage(stdout)
	return 0
}

func showUsage(stdout io.StringWriter) {
	stdout.WriteString("Manage License Zero projects and dependencies.\n\nSubcommands:\n")
	buyer := map[string]*subcommands.Subcommand{}
	seller := map[string]*subcommands.Subcommand{}
	misc := map[string]*subcommands.Subcommand{}
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
	listSubcommands(stdout, "For Buyers", buyer)
	listSubcommands(stdout, "For Sellers", seller)
	listSubcommands(stdout, "Miscellaneous", misc)
}

func listSubcommands(stdout io.StringWriter, header string, list map[string]*subcommands.Subcommand) {
	stdout.WriteString("\n  " + header + ":\n\n")
	longestSubcommand := 0
	var names []string
	for name := range list {
		length := len(name) + 1
		if length > longestSubcommand {
			longestSubcommand = length
		}
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		info := list[name]
		stdout.WriteString(
			fmt.Sprintf("  %-"+fmt.Sprintf("%d", longestSubcommand)+"s %s\n", name, info.Description),
		)
	}
}
