package cli

const bundleSchema = `{
  "$schema": "http://json-schema.org/schema#",
  "$id": "https://schemas.licensezero.com/1.0.0-pre/bundle.json",
  "type": "array",
  "items": {
    "$ref": "receipt.json"
  }
}`
