package schemas

// Bundle is a JSON schema.
const Bundle = `{
  "$schema": "http://json-schema.org/schema#",
  "$id": "https://schemas.licensezero.com/1.0.0-pre/bundle.json",
  "type": "object",
  "required": [
    "created",
    "receipts"
  ],
  "properties": {
    "created": {
      "$ref": "time.json"
    },
    "receipts": {
      "type": "array",
      "items": {
        "$ref": "receipt.json"
      }
    }
  }
}`
