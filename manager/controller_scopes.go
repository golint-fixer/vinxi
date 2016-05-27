package manager

// JSONScope represents the scope entity for JSON serialization.
type JSONScope struct {
	ID      string       `json:"id"`
	Name    string       `json:"name,omitempty"`
	Rules   []JSONRule   `json:"rules"`
	Plugins []JSONPlugin `json:"plugins"`
}

func createScope(scope *Scope) JSONScope {
	return JSONScope{
		ID:      scope.ID,
		Name:    scope.Name,
		Rules:   createRules(scope),
		Plugins: createPlugins(scope.Plugins.All()),
	}
}

func createScopes(scopes []*Scope) []JSONScope {
	list := []JSONScope{}
	for _, scope := range scopes {
		list = append(list, createScope(scope))
	}
	return list
}

// ScopesController represents the scopes entity HTTP controller.
type ScopesController struct{}

func (ScopesController) List(ctx *Context) {
	var scopes []*Scope
	if ctx.Instance != nil {
		scopes = ctx.Instance.Scopes()
	} else {
		scopes = ctx.Manager.Scopes()
	}
	ctx.Send(createScopes(scopes))
}

func (ScopesController) Get(ctx *Context) {
	ctx.Send(createScope(ctx.Scope))
}

func (ScopesController) Delete(ctx *Context) {
	if ctx.Manager.RemoveScope(ctx.Scope.ID) {
		ctx.SendNoContent()
	} else {
		ctx.SendError(500, "Cannot remove scope")
	}
}

func (ScopesController) Create(ctx *Context) {
	type data struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	var scope data
	err := ctx.ParseBody(&scope)
	if err != nil {
		return
	}
	if scope.Name == "" {
		ctx.SendError(400, "Missing required param: name")
		return
	}

	instance := NewScope(scope.Name, scope.Description)
	if ctx.Instance != nil {
		ctx.Instance.UseScope(instance)
	} else {
		ctx.Manager.UseScope(instance)
	}

	ctx.Send(createScope(instance))
}
