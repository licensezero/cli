package subcommands

import "os"
import "strings"

// Fail prints an error message and calls os.Exit(1).
func Fail(message string) {
	if strings.HasSuffix(message, "\n") {
		os.Stderr.WriteString(message)
	} else {
		os.Stderr.WriteString(message + "\n")
	}
	os.Exit(1)
}
