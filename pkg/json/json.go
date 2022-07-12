package json

import (
	"github.com/xeipuuv/gojsonschema"
)

func Validate(schema, jsonStr string) (*gojsonschema.Result, error) {
	schemaLoader := gojsonschema.NewStringLoader(schema)
	jsonStrLoader := gojsonschema.NewStringLoader(jsonStr)
	return gojsonschema.Validate(schemaLoader, jsonStrLoader)
}
