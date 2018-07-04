package readme

import "errors"
import "io/ioutil"
import "os"
import "path"
import "regexp"

// READMEPattern matches README file names.
var READMEPattern = regexp.MustCompile(`(?i)^readme(\.(txt|md|markdown|mdown|mkdn|textile|rdoc|org|creole|mediawiki|wiki|rst|asciidoc|adoc|asc|pod))?$`)

// FindREADME returns the file name of a README in the directory.
func FindREADME(directoryPath string) (string, error) {
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
		if READMEPattern.MatchString(name) {
			matches = append(matches, name)
		}
	}
	if len(matches) == 0 {
		return "", errors.New("could not find README file")
	}
	if len(matches) != 1 {
		return "", errors.New("found multiple README files")
	}
	return matches[0], nil
}

// ReadREADME returns the file name and contents of the README in a directory.
func ReadREADME(directoryPath string) (string, []byte, error) {
	fileName, err := FindREADME(directoryPath)
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
