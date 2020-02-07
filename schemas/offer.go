package schemas

// Offer is a JSON schema.
const Offer = `{
  "$schema": "http://json-schema.org/schema#",
  "$id": "https://schemas.licensezero.com/1.0.0-pre/offer.json",
  "type": "object",
  "required": [
    "sellerID",
    "pricing",
    "url"
  ],
  "additionalProperties": true,
  "properties": {
    "sellerID": {
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
