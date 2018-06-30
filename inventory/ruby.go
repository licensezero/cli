package inventory

import "bytes"
import "encoding/json"
import "github.com/yookoala/realpath"
import "io/ioutil"
import "os/exec"
import "path"
import "regexp"
import "strings"

func readRubyGemsProjects(packagePath string) ([]Project, error) {
	// Run `bundle show` in the working directory to list dependencies.
	showNamesAndVersions := exec.Command("bundle", "show")
	var first bytes.Buffer
	showNamesAndVersions.Stdout = &first
	err := showNamesAndVersions.Run()
	if err != nil {
		return nil, err
	}
	namesAndVersions := strings.Split(string(first.Bytes()), "\n")
	// Run `bundle show --paths` to list dependencies' paths.
	showPaths := exec.Command("bundle", "show", "--paths")
	var second bytes.Buffer
	showPaths.Stdout = &second
	err = showPaths.Run()
	if err != nil {
		return nil, err
	}
	paths := strings.Split(string(second.Bytes()), "\n")
	var returned []Project
	// Parse each line of output.
	re, _ := regexp.Compile(`^\s+\*\s+([^(]+) \((.+)\)$`)
	for i, line := range namesAndVersions[1:] {
		result := re.FindStringSubmatch(line)
		if len(result) == 0 {
			continue
		}
		name := result[1]
		version := result[2]
		gemPath := paths[i]
		// Try to read a licensezero.json file there.
		jsonFile := path.Join(gemPath, "licensezero.json")
		data, err := ioutil.ReadFile(jsonFile)
		if err != nil {
			continue
		}
		var parsed LicenseZeroJSONFile
		json.Unmarshal(data, &parsed)
		for _, envelope := range parsed.Envelopes {
			if alreadyHaveProject(returned, envelope.Manifest.ProjectID) {
				continue
			}
			project := Project{
				Path:     gemPath,
				Envelope: envelope,
				Type:     "rubygem",
				Name:     name,
				Version:  version,
			}
			realDirectory, err := realpath.Realpath(gemPath)
			if err != nil {
				project.Path = realDirectory
			} else {
				project.Path = gemPath
			}
			returned = append(returned, project)
		}
	}
	return returned, nil
}
