package subcommands

import "flag"
import "licensezero.com/cli/api"
import "licensezero.com/cli/data"
import "io/ioutil"
import "os"

const lockDescription = "Lock pricing and availability."

// Lock fixes pricing and availability.
var Lock = &Subcommand{
	Tag:         "seller",
	Description: lockDescription,
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("lock", flag.ExitOnError)
		offerID := offerIDFlag(flagSet)
		id := idFlag(flagSet)
		unlock := flagSet.String("unlock", "", "")
		silent := silentFlag(flagSet)
		flagSet.SetOutput(ioutil.Discard)
		flagSet.Usage = lockUsage
		flagSet.Parse(args)
		if *unlock == "" || (*offerID == "" && *id == "") {
			lockUsage()
		}
		if *offerID == "" && *id == "" {
			lockUsage()
		}
		if *offerID != "" && *id != "" {
			lockUsage()
		}
		if *offerID != "" {
			*id = *offerID
		}
		if !validID(*id) {
			invalidID()
		}
		licensor, err := data.ReadLicensor(paths.Home)
		if err != nil {
			Fail(licensorHint)
		}
		err = api.Lock(licensor, *id, *unlock)
		if err != nil {
			Fail("Error sending lock request: " + err.Error())
		}
		if !*silent {
			os.Stdout.WriteString("Locked pricing.\n")
		}
		os.Exit(0)
	},
}

func lockUsage() {
	usage := lockDescription + "\n\n" +
		"Usage:\n" +
		"  licensezero lock --id ID --unlock DATE\n\n" +
		"Options:\n" +
		flagsList(map[string]string{
			"id ID":           idLine,
			"silent":          silentLine,
			"unlock DATETIME": "Unlock date and time, RFC 3339 format.",
		})
	os.Stderr.WriteString(usage)
	os.Exit(1)
}
