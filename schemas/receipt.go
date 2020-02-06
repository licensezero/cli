package schemas

// Receipt contains the JSON schema for API receipt responses.
const Receipt = `{
  "$schema": "http://json-schema.org/schema#",
  "$id": "https://schemas.licensezero.com/1.0.0-pre/receipt.json",
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
      "title": "public signing key of the license vendor",
      "$ref": "key.json"
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
            "api",
            "effective",
            "licensee",
            "licensor",
            "offerID",
            "orderID"
          ],
          "additionalProperties": false,
          "properties": {
            "api": {
              "title": "license API",
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
              "title": "offer identifier",
              "type": "string",
              "format": "uuid"
            },
            "orderID": {
              "title": "order identifier",
              "type": "string",
              "format": "uuid"
            },
            "price": {
              "title": "purchase price",
              "$ref": "price.json"
            },
            "licensee": {
              "title": "licensee",
              "comment": "The licensee is the one receiving the license.",
              "type": "object",
              "required": [
                "email",
                "jurisdiction",
                "name"
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
                  "$ref": "name.json",
                  "examples": [
                    "Joe Licensee"
                  ]
                }
              }
            },
            "licensor": {
              "title": "licensor",
              "comment": "The licensor is the one giving the license.",
              "type": "object",
              "required": [
                "email",
                "jurisdiction",
                "licensorID",
                "name"
              ],
              "properties": {
                "email": {
                  "type": "string",
                  "format": "email"
                },
                "jurisdiction": {
                  "$ref": "jurisdiction.json"
                },
                "licensorID": {
                  "title": "licensor identifier",
                  "type": "string",
                  "format": "uuid"
                },
                "name": {
                  "$ref": "name.json",
                  "examples": [
                    "Joe Licensor"
                  ]
                }
              }
            },
            "vendor": {
              "title": "licesne vendor",
              "comment": "information on the party that sold the license, such as an agent or reseller, if the licensor did not sell the license themself",
              "type": "object",
              "required": [
                "email",
                "jurisdiction",
                "name",
                "website"
              ],
              "additionalProperties": false,
              "properties": {
                "email": {
                  "type": "string",
                  "format": "email"
                },
                "jurisdiction": {
                  "$ref": "jurisdiction.json"
                },
                "name": {
                  "$ref": "name.json",
                  "example": [
                    "Artless Devices LLC"
                  ]
                },
                "website": {
                  "$ref": "url.json"
                }
              }
            }
          }
        }
      }
    },
    "signature": {
      "title": "signature of the license vendor",
      "$ref": "signature.json"
    }
  }
}`
