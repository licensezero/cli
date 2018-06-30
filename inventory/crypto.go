package inventory

import "bytes"
import "encoding/hex"
import "encoding/json"
import "errors"
import "golang.org/x/crypto/ed25519"

type agentSignaturePackage struct {
	Manifest          ProjectManifest `json:"license"`
	LicensorSignature string          `json:"licensorSignature"`
}

// CheckMetadata verifies signatures to package metadata.
func CheckMetadata(project *Project, licensorKeyHex string, agentKeyHex string) error {
	serialized, err := json.Marshal(project.Envelope.Manifest)
	compacted := bytes.NewBuffer([]byte{})
	err = json.Compact(compacted, serialized)
	if err != nil {
		return errors.New("could not serialize Manifest")
	}
	err = checkManifestSignature(
		licensorKeyHex,
		project.Envelope.LicensorSignature,
		compacted.Bytes(),
		"licensor",
	)
	if err != nil {
		return err
	}
	serialized, err = json.Marshal(agentSignaturePackage{
		Manifest:          project.Envelope.Manifest,
		LicensorSignature: project.Envelope.LicensorSignature,
	})
	compacted = bytes.NewBuffer([]byte{})
	err = json.Compact(compacted, serialized)
	if err != nil {
		return errors.New("could not serialize agent signature packet")
	}
	err = checkManifestSignature(
		agentKeyHex,
		project.Envelope.AgentSignature,
		compacted.Bytes(),
		"agent",
	)
	if err != nil {
		return err
	}
	return nil
}

func checkManifestSignature(publicKey string, signature string, json []byte, source string) error {
	signatureBytes := make([]byte, hex.DecodedLen(len(signature)))
	_, err := hex.Decode(signatureBytes, []byte(signature))
	if err != nil {
		return errors.New("invalid " + source + "signature")
	}
	publicKeyBytes := make([]byte, hex.DecodedLen(len(publicKey)))
	_, err = hex.Decode(publicKeyBytes, []byte(publicKey))
	if err != nil {
		return errors.New("invalid " + source + " public key")
	}
	signatureValid := ed25519.Verify(
		publicKeyBytes,
		json,
		signatureBytes,
	)
	if !signatureValid {
		return errors.New("invalid " + source + " signature")
	}
	return nil
}
