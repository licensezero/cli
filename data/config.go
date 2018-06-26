package data

import "os"
import "path"

func ConfigPath(home string) string {
	fromEnvironment := os.Getenv("LICENSEZERO_CONFIG")
	if fromEnvironment != "" {
		return fromEnvironment
	} else {
		return path.Join(home, ".config", "licensezero")
	}
}

func makeConfigDirectory(home string) error {
	path := ConfigPath(home)
	return os.MkdirAll(path, 0755)
}
