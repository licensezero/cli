package subcommands

import "encoding/json"
import "log"
import "io/ioutil"
import "os"
import "path"
import "strings"
import realpath "github.com/yookoala/realpath"

type Project struct {
	Type     string
	Path     string
	Scope    string
	Name     string
	Version  string
	Manifest ProjectManifest
}

type PackageJSONFile struct {
	Name      string                   `json:"name"`
	Version   string                   `json:"version"`
	Envelopes []ProjectManifestEnvlope `json:"licensezero"`
}

type ProjectManifestEnvlope struct {
	LicensorSignature string          `json:"agentSignature"`
	AgentSignature    string          `json:"licensorSignature"`
	Manifest          ProjectManifest `json:"license"`
}

type ProjectManifest struct {
	Jurisdiction string `json:"jurisdiction"`
	Name         string `json:"name"`
	ProjectID    string `json:"projectID"`
	PublicKey    string `json:"publicKey"`
	Terms        string `json:"terms"`
	Version      string `json:"version"`
	Homepage     string `json:"homepage"`
}

type Projects struct {
	Licensable []Project
	Licensed   []Project
	Waived     []Project
	Unlicensed []Project
	Ignored    []Project
	Invalid    []Project
}

func Inventory(paths Paths, ignoreNC bool, ignoreR bool) (*Projects, error) {
	identity, _ := readIdentity(paths.Home)
	var licenses []LicenseManifest
	var waivers []WaiverManifest
	if identity != nil {
		readLicenses, err := ReadLicenses(paths.Home)
		if err != nil {
			licenses = readLicenses
		}
		readWaivers, err := ReadWaivers(paths.Home)
		if err != nil {
			waivers = readWaivers
		}
	}
	projects, err := ReadProjects(paths.CWD)
	if err != nil {
		return nil, err
	}
	var returned Projects
	for _, result := range projects {
		if !ValidMetadata(&result) {
			returned.Invalid = append(returned.Invalid, result)
			continue
		} else {
			returned.Licensable = append(returned.Licensable, result)
		}
		if HaveLicense(&result, licenses) {
			returned.Licensed = append(returned.Licensed, result)
			continue
		}
		if HaveWaiver(&result, waivers) {
			returned.Waived = append(returned.Waived, result)
			continue
		}
		if identity != nil && OwnProject(&result, identity) {
			continue
		}
		if result.Manifest.Terms == "noncommercial" && ignoreNC {
			returned.Ignored = append(returned.Ignored, result)
			continue
		}
		if result.Manifest.Terms == "reciprocal" && ignoreR {
			returned.Ignored = append(returned.Ignored, result)
			continue
		}
		returned.Unlicensed = append(returned.Unlicensed, result)
	}
	return &returned, nil
}

func HaveLicense(project *Project, licenses []LicenseManifest) bool {
	// TODO: Check license identity.
	for _, license := range licenses {
		if license.ProjectID == project.Manifest.ProjectID {
			return true
		}
	}
	return false
}

func HaveWaiver(project *Project, waivers []WaiverManifest) bool {
	// TODO: Check waiver identity.
	for _, waiver := range waivers {
		if waiver.ProjectID == project.Manifest.ProjectID {
			return true
		}
	}
	return false
}

func OwnProject(project *Project, identity *Identity) bool {
	return project.Manifest.Name == identity.Name &&
		project.Manifest.Jurisdiction == identity.Jurisdiction
}

func ValidMetadata(result *Project) bool {
	// TODO
	return true
}

func ReadProjects(cwd string) ([]Project, error) {
	return ReadNPMProjects(cwd)
}

func ReadNPMProjects(packagePath string) ([]Project, error) {
	var returned []Project
	node_modules := path.Join(packagePath, "node_modules")
	entries, err := readAndStatDir(node_modules)
	if err != nil {
		if os.IsNotExist(err) {
			return []Project{}, nil
		} else {
			return nil, err
		}
	}
	processProject := func(directory string, scope *string) error {
		package_json := path.Join(directory, "package.json")
		data, err := ioutil.ReadFile(package_json)
		if err != nil {
			return err
		}
		var parsed PackageJSONFile
		json.Unmarshal(data, &parsed)
		anyNewProjects := false
		for _, envelope := range parsed.Envelopes {
			if alreadyHaveProject(returned, envelope.Manifest.ProjectID) {
				continue
			}
			anyNewProjects = true
			project := Project{
				Type:     "npm",
				Path:     directory,
				Name:     parsed.Name,
				Version:  parsed.Version,
				Manifest: envelope.Manifest,
			}
			realDirectory, err := realpath.Realpath(directory)
			if err != nil {
				project.Path = realDirectory
			} else {
				project.Path = directory
			}
			if scope != nil {
				project.Scope = *scope
			}
			returned = append(returned, project)
		}
		if anyNewProjects {
			below, recursionError := ReadNPMProjects(directory)
			if recursionError != nil {
				return recursionError
			}
			returned = append(returned, below...)
		}
		return nil
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		log.Println(name)
		if strings.HasPrefix(name, "@") { // ./node_modules/@scope/package
			scope := name[1:]
			scopePath := path.Join(node_modules, name)
			scopeEntries, err := readAndStatDir(scopePath)
			if err != nil {
				if os.IsNotExist(err) {
					continue
				} else {
					return nil, err
				}
			}
			for _, scopeEntry := range scopeEntries {
				if !scopeEntry.IsDir() {
					continue
				}
				directory := path.Join(node_modules, name, scopeEntry.Name())
				err := processProject(directory, &scope)
				if err != nil {
					return nil, err
				}
			}
		} else { // ./node_modules/package/
			directory := path.Join(node_modules, name)
			err := processProject(directory, nil)
			if err != nil {
				return nil, err
			}
		}
	}
	return returned, nil
}

func alreadyHaveProject(projects []Project, projectID string) bool {
	for _, other := range projects {
		if other.Manifest.ProjectID == projectID {
			return true
		}
	}
	return false
}

func isDirSymlink(linkPath string) (bool, error) {
	resolvedPath, err := os.Readlink(linkPath)
	if err != nil {
		return false, err
	}
	info, err := os.Stat(resolvedPath)
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

func isSymlink(entry os.FileInfo) bool {
	return entry.Mode()&os.ModeSymlink != 0
}

// Like ioutil.ReadDir, but don't sort, and read all symlinks.
func readAndStatDir(directoryPath string) ([]os.FileInfo, error) {
	directory, err := os.Open(directoryPath)
	if err != nil {
		return nil, err
	}
	entries, err := directory.Readdir(-1)
	directory.Close()
	if err != nil {
		return nil, err
	}
	returned := make([]os.FileInfo, len(entries))
	for i, entry := range entries {
		if isSymlink(entry) {
			linkPath := path.Join(directoryPath, entry.Name())
			targetPath, err := os.Readlink(linkPath)
			if err != nil {
				return nil, err
			}
			if !path.IsAbs(targetPath) {
				targetPath = path.Join(path.Dir(directoryPath), targetPath)
			}
			newEntry, err := os.Stat(targetPath)
			if err != nil {
				return nil, err
			}
			returned[i] = newEntry
		} else {
			returned[i] = entry
		}
	}
	return returned, nil
}
