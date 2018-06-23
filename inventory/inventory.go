package inventory

import "bytes"
import "encoding/json"
import "errors"
import "github.com/licensezero/cli/data"
import "io/ioutil"
import "net/http"
import "os"
import "path"

type Project struct {
	Type     string
	Path     string
	Scope    string
	Name     string
	Version  string
	Envelope ProjectManifestEnvelope
}

type ProjectManifestEnvelope struct {
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
		agentPublicKey, err := fetchAgentPublicKey()
		err = CheckMetadata(&result, agentPublicKey)
		if err != nil {
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
		if result.Envelope.Manifest.Terms == "noncommercial" && ignoreNC {
			returned.Ignored = append(returned.Ignored, result)
			continue
		}
		if result.Envelope.Manifest.Terms == "reciprocal" && ignoreR {
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
		if license.ProjectID == project.Envelope.Manifest.ProjectID {
			return true
		}
	}
	return false
}

func HaveWaiver(project *Project, waivers []data.WaiverManifest) bool {
	// TODO: Check waiver identity.
	for _, waiver := range waivers {
		if waiver.ProjectID == project.Envelope.Manifest.ProjectID {
			return true
		}
	}
	return false
}

func OwnProject(project *Project, identity *data.Identity) bool {
	return project.Envelope.Manifest.Name == identity.Name &&
		project.Envelope.Manifest.Jurisdiction == identity.Jurisdiction
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

type KeyRequest struct {
	Action string `json:"action"`
}

type KeyResponse struct {
	Key string `json:"key"`
}

func fetchAgentPublicKey() (string, error) {
	bodyData := KeyRequest{Action: "key"}
	body, err := json.Marshal(bodyData)
	if err != nil {
		return "", errors.New("error encoding agent key request body")
	}
	response, err := http.Post("https://licensezero.com/api/v0", "application/json", bytes.NewBuffer(body))
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", errors.New("error reading agent key response body")
	}
	var parsed KeyResponse
	err = json.Unmarshal(responseBody, &parsed)
	if err != nil {
		return "", errors.New("error parsing agent key response body")
	}
	return parsed.Key, nil
}
