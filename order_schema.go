package cli

const orderSchema = `{
  "$schema": "http://json-schema.org/schema#",
  "$id": "https://schemas.licensezero.com/1.0.0-pre/order.json",
  "type": "object",
  "required": [
    "licensee",
    "offerIDs"
  ],
  "properties": {
    "licensee": {
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
          "$ref": "iso31662.json"
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
