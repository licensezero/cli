package inventory

import "os/exec"
import "bytes"
import "strings"

func readCabalDeps(packagePath string) ([]Project, error) {
	var returned []Project
	// Read the names of all Cabal dependencies.
	deps, err := cabalListDeps(packagePath)
	if err != nil {
		return nil, err
	}
	// Iterate the package names.
	for _, dep := range deps {
		// Get package information.
		info, err := cabalListPackageInfo(dep)
		if err != nil {
			continue
		}
		// Try to read licensezero.json in the package's path.
		projects, err := ReadLicenseZeroJSON(info.Dir)
		if err != nil {
			continue
		}
		for _, project := range projects {
			projectID := project.Envelope.Manifest.ProjectID
			if alreadyHaveProject(returned, projectID) {
				continue
			}
			project.Type = "cabal"
			project.Name = info.Name
		}
	}
	return returned, nil
}

func cabalListDeps(packagePath string) ([]string, error) {
	list := exec.Command("ghc-pkg", "field", "-f", "{{ join .Deps \"\\n\" }}")
	list.Dir = packagePath
	var stdout bytes.Buffer
	list.Stdout = &stdout
	err := list.Run()
	if err != nil {
		return nil, err
	}
	deps := strings.Split(string(stdout.Bytes()), "\n")
	// Remove empty string after final newline.
	if len(deps) != 0 {
		deps = deps[0 : len(deps)-1]
	}
	return deps, nil
}

type cabalPackageInfo struct {
	Name       string
	Dir        string
	ImportPath string
}

func cabalListPackageInfo(name string) (*cabalPackageInfo, error) {
	list := exec.Command("ghc-pkg", "field", name)
	var stdout bytes.Buffer
	list.Stdout = &stdout
	err := list.Run()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(stdout.Bytes()), "\n")
	return &cabalPackageInfo{
		Name:       lines[0],
		Dir:        lines[1],
		ImportPath: lines[2],
	}, nil
}
