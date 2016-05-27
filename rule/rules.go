package rule

import (
	"errors"

	"gopkg.in/vinxi/vinxi.v0/config"
)

// ErrRuleNotFound is used when a rule does not exists.
var ErrRuleNotFound = errors.New("vinxi: rule does not exists")

// Rules is used to store the existent rules globally.
var Rules = make(map[string]Info)

// Factory represents the rule factory function interface.
type Factory func(config.Config) (Matcher, error)

// NewFunc represents the Rule constructor factory function interface.
type NewFunc func(config.Config) (Rule, error)

// Validator represents the rule config field validator function interface.
type Validator func(interface{}, config.Config) error

// Info represents the rule entity fields
// storing the name, description and factory function
// used to initialize the fields.
type Info struct {
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Params      Params  `json:"params,omitempty"`
	Factory     Factory `json:"-"`
}

// Params represents the list of supported config fields by rules.
type Params []Field

// Field is used to declare specific config fields supported by rules.
type Field struct {
	Name        string      `json:"name,omitempty"`
	Type        string      `json:"type,omitempty"`
	Description string      `json:"description,omitempty"`
	Mandatory   bool        `json:"mandatory,omitempty"`
	Examples    []string    `json:"examples,omitempty"`
	Default     interface{} `json:"default,omitempty"`
	Validator   Validator   `json:"-"`
}

// Register registers the given rule in the current store.
func Register(rule Info) {
	Rules[rule.Name] = rule
}

// Init is used to initialize a new rule by name identifier
// based on the given config options.
func Init(name string, opts config.Config) (Rule, error) {
	if !Exists(name) {
		return nil, ErrRuleNotFound
	}
	return NewWithConfig(Rules[name], opts)
}

// Get is used to find and retrieve a rule factory function.
func Get(name string) NewFunc {
	rule, ok := Rules[name]
	if ok {
		return New(rule)
	}
	return nil
}

// GetFactory is used to find and retrieve a rule factory function.
func GetFactory(name string) Factory {
	rule, ok := Rules[name]
	if ok {
		return rule.Factory
	}
	return nil
}

// GetInfo is used to find and retrieve a rule info struct, if exists.
func GetInfo(name string) Info {
	return Rules[name]
}

// Exists is used to check if a given rule name exists.
func Exists(name string) bool {
	_, ok := Rules[name]
	return ok
}
