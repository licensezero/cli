package schemas

// Digest contains the JSON subschema for cryptoraphic digests.
const Digest = `{
  "$schema": "http://json-schema.org/schema#",
  "$id": "https://schemas.licensezero.com/1.0.0-pre/digest.json",
  "title": "hex-encoded SHA256 digest",
  "type": "string",
  "pattern": "^[0-9a-f]{64}$"
}`
