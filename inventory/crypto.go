package inventory

import "bytes"
import "encoding/hex"
import "encoding/json"
import "errors"
import "golang.org/x/crypto/ed25519"

func CheckMetadata(project *Project, agentKeyHex string) error {
	serialized, err := json.Marshal(project.Envelope.Manifest)
	compacted := bytes.NewBuffer([]byte{})
	err = json.Compact(compacted, serialized)
	if err != nil {
		return errors.New("could not serialize Manifest")
	}
	err = checkManifestSignature(
		project.Envelope.Manifest.PublicKey,
		project.Envelope.LicensorSignature,
		compacted.Bytes(),
		"licensor",
	)
	if err != nil {
		return err
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
	signatureBytes := make([]byte, 64)
	_, err := hex.Decode(signatureBytes, []byte(signature))
	if err != nil {
		return errors.New("invalid " + source + "signature")
	}
	publicKeyBytes := make([]byte, 32)
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
