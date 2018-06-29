package inventory

import "encoding/xml"
import "io/ioutil"
import "path"

type POM struct {
	GroupID    string `xml:"groupId"`
	ArtifactID string `xml:"artifactId"`
	Version    string `xml:"version"`
}

func findMavenPackageInfo(directoryPath string) *Project {
	pom_xml := path.Join(directoryPath, "pom.xml")
	data, err := ioutil.ReadFile(pom_xml)
	if err != nil {
		return nil
	}
	var parsed POM
	xml.Unmarshal(data, &parsed)
	if err != nil {
		return nil
	}
	if parsed.ArtifactID == "" {
		return nil
	}
	return &Project{
		Type:    "maven",
		Name:    parsed.ArtifactID,
		Scope:   parsed.GroupID,
		Version: parsed.Version,
	}
}
