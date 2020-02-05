package cli

const keySchema = `{
  "$schema": "http://json-schema.org/schema#",
  "$id": "https://schemas.licensezero.com/1.0.0-pre/key.json",
  "title": "hex-encoded ed25519 public key",
  "type": "string",
  "pattern": "^[0-9a-f]{64}$"
}`
