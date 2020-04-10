package schemas

// Broker is a JSON schema.
const Broker = `{
  "$schema": "http://json-schema.org/schema#",
  "$id": "https://protocol.licensezero.com/1.0.0-pre/broker.json",
  "title": "license broker",
  "comment": "information on the party that sold the license, such as an agent or reseller, if the seller did not sell the license themself",
  "type": "object",
  "required": [
    "email",
    "jurisdiction",
    "name",
    "website"
  ],
  "additionalProperties": false,
  "properties": {
    "email": {
      "type": "string",
      "format": "email"
    },
    "jurisdiction": {
      "$ref": "jurisdiction.json"
    },
    "name": {
      "$ref": "name.json",
      "example": [
        "Artless Devices LLC"
      ]
    },
    "website": {
      "$ref": "url.json"
    }
  }
}`
