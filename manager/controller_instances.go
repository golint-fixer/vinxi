package manager

import (
	"gopkg.in/vinxi/vinxi.v0"
	"gopkg.in/vinxi/vinxi.v0/config"
)

// JSONInstance represents the Instance entity for JSON serialization.
type JSONInstance struct {
	ID          string          `json:"id"`
	Name        string          `json:"name,omitempty"`
	Description string          `json:"description,omitempty"`
	Metadata    []config.Config `json:"metadata,omitempty"`
	Scopes      []JSONScope     `json:"scopes"`
}

func createInstance(instance *Instance) *vinxi.Metadata {
	return instance.Metadata()
}

func createInstances(instances []*Instance) []*vinxi.Metadata {
	list := []*vinxi.Metadata{}
	for _, instance := range instances {
		list = append(list, createInstance(instance))
	}
	return list
}

// InstancesController represents the rules entity HTTP controller.
type InstancesController struct{}

func (InstancesController) List(ctx *Context) {
	ctx.Send(createInstances(ctx.Manager.Instances()))
}

func (InstancesController) Get(ctx *Context) {
	ctx.Send(createInstance(ctx.Instance))
}

func (InstancesController) Delete(ctx *Context) {
	if ctx.Manager.RemoveInstance(ctx.Instance.ID()) {
		ctx.SendNoContent()
	} else {
		ctx.SendError(500, "Cannot remove instance")
	}
}
