package schemas

// Name is a JSON schema.
const Name = `{
  "$schema": "http://json-schema.org/schema#",
  "$id": "https://protocol.licensezero.com/1.0.0-pre/name.json",
  "title": "personal or organization name",
  "type": "string",
  "minLength": 3,
  "examples": [
    "John Doe",
    "Artless Devices LLC"
  ]
}`
