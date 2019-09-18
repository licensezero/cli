package inventory

import "licensezero.com/cli/api"
import "licensezero.com/cli/data"
import "os"
import "path"

// Project describes a License Zero project in inventory.
type Project struct {
	Type     string                  `json:"type"`
	Path     string                  `json:"path"`
	Scope    string                  `json:"scope,omitempty"`
	Name     string                  `json:"name"`
	Version  string                  `json:"version"`
	Envelope ProjectManifestEnvelope `json:"envelope"`
}

// ProjectManifestEnvelope describes a signed project manifest.
type ProjectManifestEnvelope struct {
	LicensorSignature string          `json:"licensorSignature"`
	AgentSignature    string          `json:"agentSignature"`
	Manifest          ProjectManifest `json:"license"`
}

// ProjectManifest describes project data from licensezero.json.
type ProjectManifest struct {
	// Note: These declaration must appear in the order so as to
	// serialize in the correct order for signature verification.
	Homepage     string `json:"homepage"`
	Jurisdiction string `json:"jurisdiction"`
	Name         string `json:"name"`
	ProjectID    string `json:"projectID"`
	PublicKey    string `json:"publicKey"`
	Terms        string `json:"terms"`
	Version      string `json:"version"`
}

// Projects describes the categorization of projects in inventory.
type Projects struct {
	Licensable []Project `json:"licensable"`
	Licensed   []Project `json:"licensed"`
	Waived     []Project `json:"waived"`
	Unlicensed []Project `json:"unlicensed"`
	Ignored    []Project `json:"ignored"`
	Invalid    []Project `json:"invalid"`
}

// Inventory finds License Zero projects included or referenced in a working tree.
func Inventory(home string, cwd string, ignoreNC bool, ignoreR bool) (*Projects, error) {
	identity, _ := data.ReadIdentity(home)
	var licenses []data.LicenseEnvelope
	var waivers []data.WaiverEnvelope
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
	projects, err := readProjects(cwd)
	if err != nil {
		return nil, err
	}
	agentPublicKey, err := api.FetchAgentPublicKey()
	if err != nil {
		return nil, err
	}
	var returned Projects
	for _, result := range projects {
		licensor, err := api.Project(result.Envelope.Manifest.ProjectID)
		if err != nil {
			returned.Invalid = append(returned.Invalid, result)
			continue
		}
		err = CheckMetadata(&result.Envelope, licensor.PublicKey, agentPublicKey)
		if err != nil {
			returned.Invalid = append(returned.Invalid, result)
			continue
		} else {
			returned.Licensable = append(returned.Licensable, result)
		}
		if haveLicense(&result, licenses, identity) {
			returned.Licensed = append(returned.Licensed, result)
			continue
		}
		if haveWaiver(&result, waivers, identity) {
			returned.Waived = append(returned.Waived, result)
			continue
		}
		if identity != nil && ownProject(&result, identity) {
			continue
		}
		terms := result.Envelope.Manifest.Terms
		if (terms == "noncommercial" || terms == "prosperity") && ignoreNC {
			returned.Ignored = append(returned.Ignored, result)
			continue
		}
		if (terms == "reciprocal" || terms == "parity") && ignoreR {
			returned.Ignored = append(returned.Ignored, result)
			continue
		}
		returned.Unlicensed = append(returned.Unlicensed, result)
	}
	return &returned, nil
}

func haveLicense(project *Project, licenses []data.LicenseEnvelope, identity *data.Identity) bool {
	for _, license := range licenses {
		if license.ProjectID == project.Envelope.Manifest.ProjectID &&
			license.Manifest.Licensee.Name == identity.Name &&
			license.Manifest.Licensee.Jurisdiction == identity.Jurisdiction &&
			license.Manifest.Licensee.EMail == identity.EMail {
			return true
		}
	}
	return false
}

func haveWaiver(project *Project, waivers []data.WaiverEnvelope, identity *data.Identity) bool {
	for _, waiver := range waivers {
		if waiver.ProjectID == project.Envelope.Manifest.ProjectID &&
			waiver.Manifest.Beneficiary.Name == identity.Name &&
			waiver.Manifest.Beneficiary.Jurisdiction == identity.Jurisdiction {
			return true
		}
	}
	return false
}

func ownProject(project *Project, identity *data.Identity) bool {
	return project.Envelope.Manifest.Name == identity.Name &&
		project.Envelope.Manifest.Jurisdiction == identity.Jurisdiction
}

func readProjects(cwd string) ([]Project, error) {
	descenders := []func(string) ([]Project, error){
		readNPMProjects,
		readRubyGemsProjects,
		readGoDeps,
		readCargoCrates,
		recurseLicenseZeroFiles,
	}
	returned := []Project{}
	for _, descender := range descenders {
		projects, err := descender(cwd)
		if err == nil {
			for _, project := range projects {
				projectID := project.Envelope.Manifest.ProjectID
				if alreadyHaveProject(returned, projectID) {
					continue
				}
				returned = append(returned, project)
			}
		}
	}
	return returned, nil
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
