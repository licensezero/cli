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
	showAll := exec.Command("bundle", "show")
	var stdout bytes.Buffer
	showAll.Stdout = &stdout
	err := showAll.Run()
	if err != nil {
		return nil, err
	}
	showAllOutput := string(stdout.Bytes())
	lines := strings.Split(showAllOutput, "\n")
	var returned []Project
	// Parse each line of output for Gem name and version.
	re, _ := regexp.Compile(`^\s+\*\s+([^(]+) \((.+)\)$`)
	for _, line := range lines {
		result := re.FindStringSubmatch(line)
		if len(result) == 0 {
			continue
		}
		name := result[1]
		version := result[2]
		// Run `bundle show $name` to find the Gem's local path.
		gemPath, err := getRubyGemPathFromBundler(name)
		if err != nil {
			continue
		}
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

func getRubyGemPathFromBundler(name string) (string, error) {
	showGem := exec.Command("bundle", "show", name)
	var stdout bytes.Buffer
	showGem.Stdout = &stdout
	err := showGem.Run()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(stdout.Bytes())), nil
}
