package schemas

// Offer is a JSON schema.
const Offer = `{
  "$schema": "http://json-schema.org/schema#",
  "$id": "https://protocol.licensezero.com/1.0.0-pre/offer.json",
  "type": "object",
  "required": [
    "sellerID",
    "pricing",
    "url"
  ],
  "additionalProperties": true,
  "properties": {
    "sellerID": {
      "$ref": "id.json"
    },
    "pricing": {
      "type": "object",
      "properties": {
        "single": {
          "$ref": "price.json"
        }
      },
      "patternProperties": {
        "^\\d+$": {
          "$ref": "price.json"
        }
      }
    },
    "url": {
      "$ref": "url.json"
    }
  }
}`
