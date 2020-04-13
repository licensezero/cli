package schemas

// PrivateKey is a JSON schema.
const PrivateKey = `{
  "$schema": "http://json-schema.org/schema#",
  "$id": "https://protocol.licensezero.com/1.0.0-pre/private-key.json",
  "title": "hex-encoded ed25519 private key",
  "type": "string",
  "pattern": "^[0-9a-f]{128}$"
}`
