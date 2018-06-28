package subcommands

import "os"
import "strings"

func Fail(message string) {
	if strings.HasSuffix(message, "\n") {
		os.Stderr.WriteString(message)
	} else {
		os.Stderr.WriteString(message + "\n")
	}
	os.Exit(1)
}
