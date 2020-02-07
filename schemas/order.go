package schemas

// Order contains the JSON schema for API order responses.
const Order = `{
  "$schema": "http://json-schema.org/schema#",
  "$id": "https://schemas.licensezero.com/1.0.0-pre/order.json",
  "type": "object",
  "required": [
    "buyer",
    "offerIDs"
  ],
  "properties": {
    "buyer": {
      "type": "object",
      "required": [
        "email",
        "jurisdiction",
        "name"
      ],
      "additionalProperties": true,
      "properties": {
        "email": {
          "type": "string",
          "format": "email"
        },
        "jurisdiction": {
          "$ref": "jurisdiction.json"
        },
        "name": {
          "type": "string",
          "minLength": 3
        }
      }
    },
    "offerIDs": {
      "type": "array",
      "items": {
        "type": "string",
        "format": "uuid"
      }
    }
  }
}`
