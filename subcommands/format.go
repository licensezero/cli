package subcommands

import "fmt"
import "strconv"

func commission(percent uint) string {
	return strconv.Itoa(int(percent)) + "%"
}

func currency(cents uint) string {
	if cents < 100 {
		if cents < 10 {
			return "$0.0" + strconv.Itoa(int(cents))
		}
		return "$0." + strconv.Itoa(int(cents))
	}
	asString := fmt.Sprintf("%d", cents)
	return "$" + asString[:len(asString)-2] + "." + asString[len(asString)-2:]
}
