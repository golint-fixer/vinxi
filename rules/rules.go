package rules

import (
	"gopkg.in/vinxi/vinxi.v0/config"
)

// Rules is used to store the existent rules globally.
var Rules = make(Store)

// RuleFactory represents the rule factory function interface.
type RuleFactory func(config.Config) Rule

// Rule represents the rule entity fields
// storing the name, description and factory function
// used to initialize the fields.
type Rule struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Config      Config `json:"config,omitempty"`
	Factory     RuleFactory
}

// Config represents the list of supported config fields by rules.
type Config []ConfigField

// ConfigField is used to declare specific config fields supported by rules.
type ConfigField struct {
	Name        string `json:"name,omitempty"`
	Type        string `json:"type,omitempty"`
	Description string `json:"description,omitempty"`
	Mandatory   bool   `json:"mandatory,omitempty"`
	Example     string `json:"example,omitempty"`
}

// Store represents the rules store used
// to register and fetch rules.
type Store map[string]Rule

// Register registers the given rule in the current store.
func (s Store) Register(rule Rule) {
	s[rule.Name] = rule
}
