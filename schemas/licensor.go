package schemas

const Licensor = `{
  "$schema": "http://json-schema.org/schema#",
  "$id": "https://schemas.licensezero.com/1.0.0-pre/licensor.json",
  "type": "object",
  "required": [
    "jurisdiction",
    "name"
  ],
  "properties": {
    "jurisdiction": {
      "$ref": "jurisdiction.json"
    },
    "name": {
      "type": "string",
      "minLength": 3
    }
  }
}`
