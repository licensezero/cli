package inventory

import (
	"bytes"
	"os/exec"
	"regexp"
	"strings"
)

func findRubyGems(packagePath string) (findings []*Finding, err error) {
	// Run `bundle show` in the working directory to list dependencies.
	showNamesAndVersions := exec.Command("bundle", "show")
	showNamesAndVersions.Dir = packagePath
	var first bytes.Buffer
	showNamesAndVersions.Stdout = &first
	err = showNamesAndVersions.Run()
	if err != nil {
		return
	}
	namesAndVersions := strings.Split(string(first.Bytes()), "\n")
	// Run `bundle list --paths` to list dependencies' paths.
	showPaths := exec.Command("bundle", "list", "--paths")
	showPaths.Dir = packagePath
	var second bytes.Buffer
	showPaths.Stdout = &second
	err = showPaths.Run()
	if err != nil {
		return
	}
	paths := strings.Split(string(second.Bytes()), "\n")
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
		fromJSON, err := readLicenseZeroJSON(gemPath)
		if err != nil {
			continue
		}
		for _, finding := range fromJSON {
			finding.Type = "rubygem"
			finding.Name = name
			finding.Version = version
			findings = append(findings, finding)
		}
	}
	return
}
