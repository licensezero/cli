package data

import "os"
import "path"

func configPath(home string) string {
	return path.Join(home, ".config", "licensezero")
}

func makeConfigDirectory(home string) error {
	path := configPath(home)
	return os.MkdirAll(path, 0744)
}
