package schemas

// URL is a JSON schema.
const URL = `{
  "$schema": "http://json-schema.org/schema#",
  "$id": "https://schemas.licensezero.com/1.0.0-pre/url.json",
  "title": "HTTPS URL",
  "type": "string",
  "format": "url",
  "pattern": "^https://"
}`
