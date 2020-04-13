package schemas

// Buyer is a JSON schema.
const Buyer = `{
  "$schema": "http://json-schema.org/schema#",
  "$id": "https://protocol.licensezero.com/1.0.0-pre/buyer.json",
  "title": "buyer",
  "comment": "The buyer is the one receiving the license.",
  "type": "object",
  "required": [
    "email",
    "jurisdiction",
    "name"
  ],
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
      "examples": [
        "Joe Buyer"
      ]
    }
  }
}`
