package schemas

// Key contains the JSON subschema for signing keys.
const Key = `{
  "$schema": "http://json-schema.org/schema#",
  "$id": "https://schemas.licensezero.com/1.0.0-pre/key.json",
  "title": "hex-encoded ed25519 public key",
  "type": "string",
  "pattern": "^[0-9a-f]{64}$"
}`
