package schemas

// Offer contains the JSON schema for API offer responses.
const Offer = `{
  "$schema": "http://json-schema.org/schema#",
  "$id": "https://schemas.licensezero.com/1.0.0-pre/offer.json",
  "type": "object",
  "required": [
    "licensorID",
    "pricing",
    "url"
  ],
  "additionalProperties": true,
  "properties": {
    "licensorID": {
      "type": "string",
      "format": "uuid"
    },
    "pricing": {
      "type": "object",
      "properties": {
        "single": {
          "$ref": "price.json"
        },
        "site": {
          "$ref": "price.json"
        }
      },
      "patternProperties": {
        "^d+$": {
          "$ref": "price.json"
        }
      }
    },
    "url": {
      "type": "string",
      "format": "uri"
    }
  }
}`
