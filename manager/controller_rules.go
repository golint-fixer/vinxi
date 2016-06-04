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

// rulesController represents the rules entity HTTP controller.
type rulesController struct{}

func (rulesController) List(ctx *Context) {
	ctx.Send(createRules(ctx.Scope))
}

func (rulesController) Get(ctx *Context) {
	ctx.Send(createRule(ctx.Rule))
}

func (rulesController) Delete(ctx *Context) {
	if ctx.Scope.RemoveRule(ctx.Rule.ID()) {
		ctx.SendNoContent()
	} else {
		ctx.SendError(500, "Cannot remove rule")
	}
}

func (rulesController) Create(ctx *Context) {
	type data struct {
		Name   string        `json:"name"`
		Params config.Config `json:"config"`
	}

	var input data
	err := ctx.ParseBody(&input)
	if err != nil {
		return
	}

	if input.Name == "" {
		ctx.SendError(400, "Missing required param: name")
		return
	}

	factory := rule.Get(input.Name)
	if factory == nil {
		ctx.SendNotFound("Rule not found")
		return
	}

	instance, err := factory(input.Params)
	if err != nil {
		ctx.SendError(400, "Cannot create rule: "+err.Error())
		return
	}

	ctx.Scope.UseRule(instance)
	ctx.Send(createRule(instance))
}
