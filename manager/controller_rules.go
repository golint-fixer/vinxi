package manager

import (
	"gopkg.in/vinxi/vinxi.v0/config"
	"gopkg.in/vinxi/vinxi.v0/rule"
)

// JSONRule represents the Rule entity for JSON serialization.
type JSONRule struct {
	ID          string        `json:"id"`
	Name        string        `json:"name,omitempty"`
	Description string        `json:"description,omitempty"`
	Config      config.Config `json:"config,omitempty"`
	Metadata    config.Config `json:"metadata,omitempty"`
}

func createRules(scope *Scope) []JSONRule {
	rules := []JSONRule{}
	for _, rule := range scope.Rules.All() {
		rules = append(rules, createRule(rule))
	}
	return rules
}

func createRule(rule rule.Rule) JSONRule {
	return JSONRule{
		ID:          rule.ID(),
		Name:        rule.Name(),
		Description: rule.Description(),
		Config:      rule.Config(),
	}
}

// RulesController represents the rules entity HTTP controller.
type RulesController struct{}

func (RulesController) List(ctx *Context) {
	ctx.Send(createRules(ctx.Scope))
}

func (RulesController) Get(ctx *Context) {
	ctx.Send(createRule(ctx.Rule))
}

func (RulesController) Create(ctx *Context) {
	ctx.Send(createRule(ctx.Rule))
}

func (RulesController) Delete(ctx *Context) {
	if ctx.Scope.RemoveRule(ctx.Rule.ID()) {
		ctx.SendNoContent()
	} else {
		ctx.SendError(500, "Cannot remove rule")
	}
}
