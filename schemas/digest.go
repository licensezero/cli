package schemas

// Digest is a JSON schema.
const Digest = `{
  "$schema": "http://json-schema.org/schema#",
  "$id": "https://protocol.licensezero.com/1.0.0-pre/digest.json",
  "title": "hex-encoded SHA256 digest",
  "type": "string",
  "pattern": "^[0-9a-f]{64}$"
}`
