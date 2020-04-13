package schemas

// PublicKey is a JSON schema.
const PublicKey = `{
  "$schema": "http://json-schema.org/schema#",
  "$id": "https://protocol.licensezero.com/1.0.0-pre/public-key.json",
  "title": "hex-encoded ed25519 public key",
  "type": "string",
  "pattern": "^[0-9a-f]{64}$"
}`
