package rules

import (
	"gopkg.in/vinxi/vinxi.v0/manager"
)

// Rules is used to store the existent rules globally.
var Rules = make(Store)

// RuleFactory represents the rule factory function interface.
type RuleFactory func(manager.Config) manager.Rule

// Rule represents the rule entity fields
// storing the name, description and factory function
// used to initialize the fields.
type Rule struct {
	Name        string
	Description string
	Factory     RuleFactory
}

// Store represents the rules store used
// to register and fetch rules.
type Store map[string]Rule

// Register registers the given rule in the current store.
func (s Store) Register(rule Rule) {
	s[rule.Name] = rule
}
