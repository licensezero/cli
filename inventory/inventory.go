package inventory

import "github.com/licensezero/cli/data"
import "os"
import "path"

type Project struct {
	Type     string
	Path     string
	Scope    string
	Name     string
	Version  string
	Manifest ProjectManifest
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

func Inventory(home string, cwd string, ignoreNC bool, ignoreR bool) (*Projects, error) {
	identity, _ := data.ReadIdentity(home)
	var licenses []data.LicenseManifest
	var waivers []data.WaiverManifest
	if identity != nil {
		readLicenses, err := data.ReadLicenses(home)
		if err != nil {
			licenses = readLicenses
		}
		readWaivers, err := data.ReadWaivers(home)
		if err != nil {
			waivers = readWaivers
		}
	}
	projects, err := ReadProjects(cwd)
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

func HaveLicense(project *Project, licenses []data.LicenseManifest) bool {
	// TODO: Check license identity.
	for _, license := range licenses {
		if license.ProjectID == project.Manifest.ProjectID {
			return true
		}
	}
	return false
}

func HaveWaiver(project *Project, waivers []data.WaiverManifest) bool {
	// TODO: Check waiver identity.
	for _, waiver := range waivers {
		if waiver.ProjectID == project.Manifest.ProjectID {
			return true
		}
	}
	return false
}

func OwnProject(project *Project, identity *data.Identity) bool {
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