package inventory

import "bytes"
import "os"
import "os/exec"
import "path"
import "strings"

func findPythonPackageInfo(directoryPath string) *DescenderResult {
	setup := path.Join(directoryPath, "setup.py")
	_, err := os.Stat(setup)
	if err != nil {
		return nil
	}
	command := exec.Command("python", "setup.py", "--name", "--version")
	var stdout bytes.Buffer
	command.Stdout = &stdout
	err = command.Run()
	if err != nil {
		return nil
	}
	output := string(stdout.Bytes())
	lines := strings.Split(output, "\n")
	name := strings.TrimSpace(lines[0])
	version := strings.TrimSpace(lines[1])
	return &DescenderResult{
		Type:    "python",
		Name:    name,
		Version: version,
	}
}
