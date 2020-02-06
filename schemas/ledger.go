package schemas

// Ledger contains the JSON schema for API ledger responses.
const Ledger = `{
  "$schema": "http://json-schema.org/schema#",
  "$id": "https://schemas.licensezero.com/1.0.0-pre/ledger.json",
  "type": "array",
  "items": {
    "type": "object",
    "required": [
      "digest",
      "signature",
      "time"
    ],
    "properties": {
      "digest": {
        "$ref": "digest.json"
      },
      "signature": {
        "$ref": "signature.json"
      },
      "time": {
        "$ref": "time.json"
      }
    }
  }
}`
