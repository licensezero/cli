package inventory

import (
	"bytes"
	"os"
	"os/exec"
	"path"
	"strings"
)

func readPythonPackageMetadata(directoryPath string) *Finding {
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
	return &Finding{
		Type:    "python",
		Name:    name,
		Version: version,
	}
}
