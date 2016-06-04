package manager

import (
	"gopkg.in/vinxi/vinxi.v0/config"
	"gopkg.in/vinxi/vinxi.v0/plugin"
)

// JSONPlugin represents the Plugin entity for JSON serialization.
type JSONPlugin struct {
	ID          string        `json:"id"`
	Name        string        `json:"name,omitempty"`
	Description string        `json:"description,omitempty"`
	Enabled     bool          `json:"enabled,omitempty"`
	Config      config.Config `json:"config,omitempty"`
	Metadata    config.Config `json:"metadata,omitempty"`
}

func createPlugins(plugins []plugin.Plugin) []JSONPlugin {
	list := []JSONPlugin{}
	for _, plugin := range plugins {
		list = append(list, createPlugin(plugin))
	}
	return list
}

func createPlugin(p plugin.Plugin) JSONPlugin {
	return JSONPlugin{
		ID:          p.ID(),
		Name:        p.Name(),
		Description: p.Description(),
		Config:      p.Config(),
		Metadata:    p.Metadata(),
	}
}

// pluginsController represents the plugins entity HTTP controller.
type pluginsController struct{}

func (pluginsController) List(ctx *Context) {
	var layer *plugin.Layer
	if ctx.AdminPlugins != nil {
		layer = ctx.AdminPlugins
	} else if ctx.Scope != nil {
		layer = ctx.Scope.Plugins
	} else {
		layer = ctx.Manager.Plugins
	}
	ctx.Send(createPlugins(layer.All()))
}

func (pluginsController) Get(ctx *Context) {
	ctx.Send(createPlugin(ctx.Plugin))
}

func (pluginsController) Delete(ctx *Context) {
	if ctx.Manager.RemovePlugin(ctx.Plugin.ID()) {
		ctx.SendNoContent()
	} else {
		ctx.SendError(500, "Cannot remove plugin")
	}
}

func (p pluginsController) Create(ctx *Context) {
	type data struct {
		Name   string        `json:"name"`
		Params config.Config `json:"config"`
	}

	var plu data
	err := ctx.ParseBody(&plu)
	if err != nil {
		return
	}

	if plu.Name == "" {
		ctx.SendError(400, "Missing required param: name")
		return
	}

	factory := plugin.Get(plu.Name)
	if factory == nil {
		ctx.SendNotFound("Plugin not found")
		return
	}

	instance, err := factory(plu.Params)
	if err != nil {
		ctx.SendError(400, "Cannot create plugin: "+err.Error())
		return
	}

	p.registerPlugin(ctx, instance)
	ctx.Send(createPlugin(instance))
}

func (pluginsController) registerPlugin(ctx *Context, instance plugin.Plugin) {
	if ctx.Scope != nil {
		ctx.Scope.UsePlugin(instance)
	} else {
		ctx.Manager.UsePlugin(instance)
	}
}
