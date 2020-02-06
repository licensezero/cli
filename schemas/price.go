package schemas

// Price contains the JSON subschema for prices.
const Price = `{
  "$schema": "http://json-schema.org/schema#",
  "$id": "https://schemas.licensezero.com/1.0.0-pre/price.json",
  "title": "price",
  "type": "object",
  "required": [
    "amount",
    "currency"
  ],
  "additionalProperties": false,
  "properties": {
    "amount": {
      "title": "purchase price in minor units of currency",
      "type": "integer",
      "minimum": 1,
      "examples": [
        0,
        100
      ]
    },
    "currency": {
      "title": "purchase price currency code",
      "$ref": "currency.json",
      "examples": [
        "USD"
      ]
    }
  }
}`
