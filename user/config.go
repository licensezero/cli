package user

import (
	"github.com/mitchellh/go-homedir"
	"os"
	"path"
)

// ConfigPath computes the path of the CLI's configuration directory.
func ConfigPath() string {
	fromEnvironment := os.Getenv("LICENSEZERO_CONFIG")
	if fromEnvironment != "" {
		return fromEnvironment
	}
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	return path.Join(home, ".config", "licensezero")
}

func makeConfigDirectory() error {
	path := ConfigPath()
	return os.MkdirAll(path, 0755)
}
