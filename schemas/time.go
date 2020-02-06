package schemas

// Time contains the JSON subschema for times.
const Time = `{
  "$schema": "http://json-schema.org/schema#",
  "$id": "https://schemas.licensezero.com/1.0.0-pre/time.json",
  "title": "ISO 8601 UTC date and time",
  "type": "string",
  "format": "date-time",
  "pattern": "Z$"
}`
