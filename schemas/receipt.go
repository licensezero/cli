package schemas

// Receipt is a JSON schema.
const Receipt = `{
  "$schema": "http://json-schema.org/schema#",
  "$id": "https://protocol.licensezero.com/1.0.0-pre/receipt.json",
  "title": "license receipt",
  "comment": "A receipt represents confirmation of the sale of a software license.",
  "type": "object",
  "required": [
    "key",
    "license",
    "signature"
  ],
  "additionalProperties": false,
  "properties": {
    "key": {
      "title": "public signing key of the broker server",
      "$ref": "public-key.json"
    },
    "license": {
      "title": "license manifest",
      "type": "object",
      "required": [
        "form",
        "values"
      ],
      "properties": {
        "form": {
          "title": "license form",
          "type": "string",
          "minLength": 1
        },
        "values": {
          "type": "object",
          "required": [
            "server",
            "effective",
            "buyer",
            "seller",
            "sellerID",
            "offerID",
            "orderID"
          ],
          "additionalProperties": false,
          "properties": {
            "server": {
              "title": "license server",
              "$ref": "url.json"
            },
            "effective": {
              "title": "effective date",
              "$ref": "time.json"
            },
            "expires": {
              "title": "expiration date of the license",
              "$ref": "time.json"
            },
            "offerID": {
              "$ref": "id.json"
            },
            "orderID": {
              "$ref": "id.json"
            },
            "price": {
              "title": "purchase price",
              "$ref": "price.json"
            },
            "buyer": {
              "$ref": "buyer.json"
            },
            "recurring": {
              "const": true
            },
            "seller": {
              "$ref": "seller.json"
            },
            "sellerID": {
              "$ref": "id.json"
            },
            "broker": {
              "$ref": "broker.json"
            }
          }
        }
      }
    },
    "signature": {
      "title": "signature of the license broker server",
      "$ref": "signature.json"
    }
  }
}`
