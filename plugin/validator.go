package plugin

import (
	"errors"

	"gopkg.in/vinxi/vinxi.v0/config"
)

func getField(name string, params Params) (Field, bool) {
	for _, field := range params {
		if field.Name == name {
			return field, true
		}
	}
	return Field{}, false
}

// Validate is used to validate plugin params
// and report the proper param specific error.
func Validate(params Params, opts config.Config) error {
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
