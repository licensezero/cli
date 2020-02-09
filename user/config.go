package user

import (
	"github.com/mitchellh/go-homedir"
	"os"
	"path"
)

// ConfigPath computes the path of the CLI's configuration directory.
func ConfigPath() (configPath string, err error) {
	configPath = os.Getenv("LICENSEZERO_CONFIG")
	if configPath != "" {
		return
	}
	home, err := homedir.Dir()
	if err != nil {
		return
	}
	return path.Join(home, ".config", "licensezero"), nil
}

func makeConfigDirectory() error {
	configPath, err := ConfigPath()
	if err != nil {
		return err
	}
	return os.MkdirAll(configPath, 0755)
}
