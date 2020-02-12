package schemas

// Artifact is a JSON schema.
const Artifact = `{
  "$schema": "http://json-schema.org/schema#",
  "$id": "https://schemas.licensezero.com/1.0.0-pre/artifact.json",
  "type": "object",
  "required": [
    "offers"
  ],
  "additionalProperties": false,
  "properties": {
    "offers": {
      "type": "array",
      "items": {
        "type": "object",
        "required": [
          "server",
          "offerID"
        ],
        "additionalProperties": false,
        "properties": {
          "server": {
            "title": "licensing server",
            "type": "string",
            "format": "uri",
            "pattern": "^https://",
            "examples": [
              "https://broker.licensezero.com"
            ]
          },
          "offerID": {
            "title": "UUIDv4 offer identifier",
            "type": "string",
            "format": "uuid"
          },
          "public": {
            "title": "public license identifier",
            "type": "string",
            "pattern": "^[A-Za-z0-9-.]+",
            "examples": [
              "Parity-7.0.0"
            ]
          }
        }
      }
    }
  }
}`
