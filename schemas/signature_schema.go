package schemas

const Signature = `{
  "$schema": "http://json-schema.org/schema#",
  "$id": "https://schemas.licensezero.com/1.0.0-pre/signature.json",
  "title": "hex-encoded ed25519 detached signature",
  "type": "string",
  "pattern": "^[0-9a-f]{128}$"
}`
