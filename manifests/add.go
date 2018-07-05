package manifests

import "bytes"
import "encoding/json"
import "errors"
import "io/ioutil"
import "os"
import "path"
import "strings"

type manifestType struct {
	File    string
	Handler func(string, string) error
}

var plans = []*manifestType{
	&manifestType{
		File:    "package.json",
		Handler: AddToPackageJSONFilesArray,
	},
	&manifestType{
		File:    "MANIFEST.in",
		Handler: IncludeInMANIFESTIN,
	},
	&manifestType{
		File:    "Manifest.txt",
		Handler: IncludeInMANIFESTTXT,
	},
}

// AddToPackageJSONFilesArray adds a file to a directory's MANIFEST.in file, if it has one.
func AddToPackageJSONFilesArray(directoryPath, filePath string) error {
	packageJSONPath := path.Join(directoryPath, "package.json")
	data, err := ioutil.ReadFile(packageJSONPath)
	if err != nil {
		return err
	}
	var parsed map[string]interface{}
	json.Unmarshal(data, &parsed)
	files, ok := parsed["files"].([]interface{})
	if !ok {
		return errors.New("no files array in package.json")
	}
	for _, existing := range files {
		if existingString, ok := existing.(string); ok {
			if existingString == filePath {
				return errors.New("already in package.json files array")
			}
		}
	}
	files = append(files, filePath)
	parsed["files"] = files
	serialized := new(bytes.Buffer)
	encoder := json.NewEncoder(serialized)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(parsed)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(packageJSONPath, serialized.Bytes(), 0644)
	if err != nil {
		return err
	}
	return nil
}

// IncludeInMANIFESTIN adds a file to the files array in a directory's package.json file, if it has one.
func IncludeInMANIFESTIN(directoryPath, filePath string) error {
	manifestPath := path.Join(directoryPath, "Manifest.txt")
	data, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		return err
	}
	text := string(data)
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		if line == filePath {
			return errors.New("already in Manifest.txt")
		}
	}
	text = text + "\n" + filePath + "\n"
	err = ioutil.WriteFile(manifestPath, []byte(text), 0644)
	if err != nil {
		return err
	}
	return nil
}

// IncludeInMANIFESTTXT adds a file to a directory's Manifest.txt file, if it has one.
func IncludeInMANIFESTTXT(directoryPath, filePath string) error {
	manifestPath := path.Join(directoryPath, "MANIFEST.in")
	data, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		return err
	}
	text := string(data)
	lines := strings.Split(text, "\n")
	include := "include " + filePath
	for _, line := range lines {
		if line == include {
			return errors.New("already in MANIFEST.in")
		}
	}
	text = text + "\n" + include + "\n"
	err = ioutil.WriteFile(manifestPath, []byte(text), 0644)
	if err != nil {
		return err
	}
	return nil
}

// AddToManifests adds a file to package manifests for various languages.
func AddToManifests(directoryPath, filePath string) ([]string, error) {
	// Read the entries in the directory.
	directory, err := os.Open(directoryPath)
	if err != nil {
		return nil, err
	}
	defer directory.Close()
	entries, err := directory.Readdir(-1)
	if err != nil {
		return nil, err
	}
	// Iterate entries.
	var modified []string
	for _, entry := range entries {
		if entry.Mode()&os.ModeType != 0 {
			// Not a regular file.
			continue
		}
		name := entry.Name()
		for _, plan := range plans {
			if plan.File == name {
				err = plan.Handler(directoryPath, filePath)
				if err == nil {
					modified = append(modified, name)
				}
			}
		}
	}
	return modified, nil
}
