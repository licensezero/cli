package schemas

// Order is a JSON schema.
const Order = `{
  "$schema": "http://json-schema.org/schema#",
  "$id": "https://protocol.licensezero.com/1.0.0-pre/order.json",
  "type": "object",
  "required": [
    "email",
    "jurisdiction",
    "name",
    "offerIDs[]"
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
      "type": "string",
      "minLength": 3
    },
    "offerIDs[]": {
      "type": "array",
      "items": {
        "$ref": "id.json"
      }
    }
  }
}`
