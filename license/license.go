package license

import "errors"
import "io/ioutil"
import "os"
import "path"
import "regexp"

// LicensePatterns lists patterns that match generic license file names.
var LicensePatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)^license(\.(txt|md|markdown|mdown|mkdn|textile|rdoc|org|creole|mediawiki|wiki|rst|asciidoc|adoc|asc|pod))?$`),
	regexp.MustCompile(`(?i)^copying(\.(txt|md|markdown|mdown|mkdn|textile|rdoc|org|creole|mediawiki|wiki|rst|asciidoc|adoc|asc|pod))?$`),
}

const notFound = "coult not find any generic license file"

// FindLicense returns the file name of a generic license file in a directory.
func FindLicense(directoryPath string) (string, error) {
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
		for _, pattern := range LicensePatterns {
			if pattern.MatchString(name) {
				matches = append(matches, name)
			}
		}
	}
	if len(matches) == 0 {
		return "", errors.New(notFound)
	}
	if len(matches) != 1 {
		return "", errors.New("found multiple generic license files")
	}
	return matches[0], nil
}

// IsNotFound returns true for errors about not findy any license file.
func IsNotFound(err error) bool {
	return err.Error() == notFound
}

// ReadLicense returns the file name and contents of the generic license file in a directory.
func ReadLicense(directoryPath string) (string, []byte, error) {
	fileName, err := FindLicense(directoryPath)
	if err != nil {
		return "", nil, err
	}
	filePath := path.Join(directoryPath, fileName)
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", nil, err
	}
	return filePath, data, nil
}
