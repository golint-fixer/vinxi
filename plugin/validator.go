package plugin

import (
	"errors"

	"gopkg.in/vinxi/vinxi.v0/config"
)

// getField is used to find and retrieve the param Field configuration.
func getField(name string, params Params) (Field, bool) {
	for _, field := range params {
		if field.Name == name {
			return field, true
		}
	}
	return Field{}, false
}

// Validate is used to validates plugin params
// and report the proper param specific error.
func Validate(params Params, opts config.Config) error {
	// Set defaults params
	for _, param := range params {
		// Set defaults
		if param.Default != nil && !opts.Exists(param.Name) {
			opts.Set(param.Name, param.Default)
		}
		// Validate mandatory params
		if param.Mandatory && !opts.Exists(param.Name) {
			return errors.New("Missing required param: " + param.Name)
		}
	}

	// Validate params
	for name, value := range opts {
		field, exists := getField(name, params)
		if !exists {
			continue
		}

		// Cast type to verify type contract
		var kind string
		switch value.(type) {
		case string:
			kind = "string"
		case int:
			kind = "int"
		case bool:
			kind = "bool"
		}

		if kind != field.Type {
			return errors.New("Invalid type for param: " + name)
		}

		if field.Validator != nil {
			err := field.Validator(value, opts)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
