package schemas

// Artifact is a JSON schema.
const Artifact = `{
  "$schema": "http://json-schema.org/schema#",
  "$id": "https://protocol.licensezero.com/1.0.0-pre/artifact.json",
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
            "$ref": "url.json",
            "examples": [
              "https://broker.licensezero.com"
            ]
          },
          "offerID": {
            "$ref": "id.json"
          },
          "public": {
            "title": "public license class",
            "type": "string",
            "enum": [
              "noncommercial",
              "share alike"
            ]
          }
        }
      }
    }
  }
}`
