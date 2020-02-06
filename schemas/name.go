package schemas

// Name contains the JSON subschema for licensor, licenses,
// and vendor names.
const Name = `{
  "$schema": "http://json-schema.org/schema#",
  "$id": "https://schemas.licensezero.com/1.0.0-pre/name.json",
  "title": "personal or organization name",
  "type": "string",
  "minLength": 3,
  "examples": [
    "John Doe",
    "Artless Devices LLC"
  ]
}`
