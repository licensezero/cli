package schemas

// Register is a JSON schema.
const Register = `{
  "$schema": "http://json-schema.org/schema#",
  "$id": "https://protocol.licensezero.com/1.0.0-pre/register.json",
  "type": "object",
  "required": [
    "updated",
    "keys"
  ],
  "properties": {
    "updated": {
      "$ref": "time.json"
    },
    "keys": {
      "type": "object",
      "patternProperties": {
        "^[0-9a-f]{64}$": {
          "type": "object",
          "required": [
            "from"
          ],
          "properties": {
            "from": {
              "$ref": "time.json"
            },
            "through": {
              "$ref": "time.json"
            }
          }
        }
      }
    }
  }
}`
