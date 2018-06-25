package subcommands

import "strings"
import "sort"

func flagsList(mapping map[string]string) string {
	returned := ""
	var flags []string
	var longest int
	for key, _ := range mapping {
		flags = append(flags, key)
		length := len(key) + 2
		if length > longest {
			longest = length
		}
	}
	sort.Strings(flags)
	for _, key := range flags {
		returned = returned +
			"  --" + key + strings.Repeat(" ", longest-len(key)) +
			mapping[key] + "\n"
	}
	return returned
}
