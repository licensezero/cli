package subcommands

import "encoding/json"
import "io/ioutil"
import "os"
import "path"

type Result struct {
	Type     string
	Path     string
	Metadata Metadata
}

type Metadata struct {
	Name             string            `json:"name"`
	Version          string            `json:"version"`
	ProjectManifests []ProjectManifest `json:"licensezero"`
}

type ProjectManifest struct {
	LicensorSignature string            `json:"agentSignature"`
	AgentSignature    string            `json:"licensorSignature"`
	LicenseManifests  []LicenseManifest `json:"license"`
}

type LicenseManifest struct {
	Jurisdiction string `json:"jurisdiction"`
	Name         string `json:"name"`
	ProjectID    string `json:"projectID"`
	PublicKey    string `json:"publicKey"`
	Terms        string `json:"terms"`
	Version      string `json:"version"`
	Homepage     string `json:"homepage"`
}

type Results struct {
	Licensable []Result
	Licensed   []Result
	Waived     []Result
	Unlicensed []Result
	Ignored    []Result
	Invalid    []string
}

func Inventory(paths Paths, ignoreNC bool, ignoreR bool) (*Results, error) {
	// identity, _ := readIdentity(paths.Home)
	/*
		if identity {
			licenses, _ := readLicenses(paths.Home, identity)
			waivers, _ := ReadWaivers(paths.Home, identity)
		}
	*/
	tree, err := ReadTrees(paths.CWD)
	if err != nil {
		return nil, err
	}
	return &Results{
		Ignored: tree,
	}, nil
}

func ReadTrees(cwd string) ([]Result, error) {
	return ReadNPMTree(cwd)
}

/*
import "errors"
import "os/exec"

type npmOutput struct {
	Version      string
	Dependencies map[string]npmOutput
}

func ListNPMTree(cwd string) ([]Result, error) {
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

func FlattenNPMTree(output npmOutput, above []string) []Result {
	var returned []Result
	for name, manifest := range output.Dependencies {
		returned = append(returned, Result{
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

func ReadNPMTree(cwd string) ([]Result, error) {
	var returned []Result
	node_modules := path.Join(cwd, "node_modules")
	files, directoryReadError := ioutil.ReadDir(node_modules)
	if directoryReadError != nil {
		if os.IsNotExist(directoryReadError) {
			return []Result{}, nil
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
			var parsed Metadata
			json.Unmarshal(data, &parsed)
			result := Result{
				Type:     "npm",
				Path:     directory,
				Metadata: parsed,
			}
			returned = append(returned, result)
		}
		below, recursionError := ReadNPMTree(directory)
		if recursionError != nil {
			return nil, recursionError
		}
		returned = append(returned, below...)
	}
	return returned, nil
}
