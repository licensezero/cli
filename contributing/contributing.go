package contributing

import "errors"
import "io/ioutil"
import "os"
import "path"
import "regexp"

// CONTRIBUTINGPattern matches CONTRIBUTING file names.
var CONTRIBUTINGPattern = regexp.MustCompile(`(?i)^contributing(\.(txt|md|markdown|mdown|mkdn|textile|rdoc|org|creole|mediawiki|wiki|rst|asciidoc|adoc|asc|pod))?$`)

// FindCONTRIBUTING returns the file name of a CONTRIBUTING in the directory.
func FindCONTRIBUTING(directoryPath string) (string, error) {
	// Read the entries in the directory.
	directory, err := os.Open(directoryPath)
	if err != nil {
		return "", err
	}
	defer directory.Close()
	entries, err := directory.Readdir(-1)
	if err != nil {
		return "", err
	}
	// Compile a list of names that match our pattern.
	var matches []string
	for _, entry := range entries {
		if entry.Mode()&os.ModeType != 0 {
			// Not a regular file.
			continue
		}
		name := entry.Name()
		if CONTRIBUTINGPattern.MatchString(name) {
			matches = append(matches, name)
		}
	}
	if len(matches) == 0 {
		return "", errors.New("could not find CONTRIBUTING file")
	}
	if len(matches) != 1 {
		return "", errors.New("found multiple CONTRIBUTING files")
	}
	return matches[0], nil
}

// ReadCONTRIBUTING returns the file name and contents of the CONTRIBUTING in a directory.
func ReadCONTRIBUTING(directoryPath string) (string, []byte, error) {
	fileName, err := FindCONTRIBUTING(directoryPath)
	if err != nil {
		return "", nil, err
	}
	path := path.Join(directoryPath, fileName)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", nil, err
	}
	return fileName, data, nil
}
