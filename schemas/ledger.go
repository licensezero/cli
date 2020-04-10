package schemas

// Ledger is a JSON schema.
const Ledger = `{
  "$schema": "http://json-schema.org/schema#",
  "$id": "https://protocol.licensezero.com/1.0.0-pre/ledger.json",
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
