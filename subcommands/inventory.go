package subcommands

import "encoding/json"
import "io/ioutil"
import "os"
import "path"

type Project struct {
	Type     string
	Path     string
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
		if OwnProject(&result, identity) {
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

/*
import "errors"
import "os/exec"

type npmOutput struct {
	Version      string
	Dependencies map[string]npmOutput
}

func ListNPMTree(cwd string) ([]Project, error) {
	npmPath, pathError := exec.LookPath("npm")
	if pathError != nil {
		return nil, errors.New("Could not find npm.")
	}
	fmt.Println(npmPath)
	command := exec.Command(npmPath, "ls", "--json")
	command.Dir = cwd
	output, err := command.Output()
	if err != nil {
		panic(err)
		return nil, err
	}
	var parsed npmOutput
	jsonError := json.Unmarshal(output, &parsed)
	if jsonError != nil {
		return nil, errors.New("Could not parse npm ls --json output.")
	}
	flattened := FlattenNPMTree(parsed, []string{})
	return flattened, nil
}

func FlattenNPMTree(output npmOutput, above []string) []Project {
	var returned []Project
	for name, manifest := range output.Dependencies {
		returned = append(returned, Project{
			Type:    "npm",
			Name:    name,
			Version: manifest.Version,
			Path:    append(above, name),
		})
		below := FlattenNPMTree(manifest, append(above, name))
		returned = append(returned, below...)
	}
	return returned
}
*/

func ReadNPMProjects(cwd string) ([]Project, error) {
	var returned []Project
	node_modules := path.Join(cwd, "node_modules")
	files, directoryReadError := ioutil.ReadDir(node_modules)
	if directoryReadError != nil {
		if os.IsNotExist(directoryReadError) {
			return []Project{}, nil
		} else {
			return nil, directoryReadError
		}
	}
	for _, entry := range files {
		name := entry.Name()
		directory := path.Join(node_modules, name)
		package_json := path.Join(directory, "package.json")
		data, fileReadError := ioutil.ReadFile(package_json)
		if fileReadError == nil {
			var parsed PackageJSONFile
			json.Unmarshal(data, &parsed)
			for _, envelope := range parsed.Envelopes {
				project := Project{
					Type:     "npm",
					Path:     directory,
					Name:     parsed.Name,
					Version:  parsed.Version,
					Manifest: envelope.Manifest,
				}
				returned = append(returned, project)
			}
		}
		below, recursionError := ReadNPMProjects(directory)
		if recursionError != nil {
			return nil, recursionError
		}
		returned = append(returned, below...)
	}
	return returned, nil
}
