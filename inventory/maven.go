package inventory

import "encoding/xml"
import "io/ioutil"
import "path"

type pom struct {
	GroupID    string `xml:"groupId"`
	ArtifactID string `xml:"artifactId"`
	Version    string `xml:"version"`
}

func findMavenPackageInfo(directoryPath string) *Offer {
	pomFile := path.Join(directoryPath, "pom.xml")
	data, err := ioutil.ReadFile(pomFile)
	if err != nil {
		return nil
	}
	var parsed pom
	xml.Unmarshal(data, &parsed)
	if err != nil {
		return nil
	}
	if parsed.ArtifactID == "" {
		return nil
	}
	return &Offer{
		Artifact: ArtifactData{
			Type:    "maven",
			Name:    parsed.ArtifactID,
			Scope:   parsed.GroupID,
			Version: parsed.Version,
		},
	}
}
