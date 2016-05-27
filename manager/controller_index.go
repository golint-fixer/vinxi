package manager

import (
	"os"
	"runtime"

	"gopkg.in/vinxi/vinxi.v0"
	"gopkg.in/vinxi/vinxi.v0/plugin"
	"gopkg.in/vinxi/vinxi.v0/rule"
)

type info struct {
	Hostname      string            `json:"hostname"`
	Version       string            `json:"version"`
	Runtime       string            `json:"runtime"`
	Platform      string            `json:"platform"`
	NumCPU        int               `json:"cpus"`
	NumGoroutines int               `json:"gorutines"`
	Links         map[string]string `json:"links"`
}

// IndexController represents the base routes HTTP controller.
type IndexController struct{}

func (IndexController) Get(ctx *Context) {
	hostname, _ := os.Hostname()
	links := map[string]string{
		"catalog":   "/catalog",
		"plugins":   "/plugins",
		"scopes":    "/scopes",
		"instances": "/instances",
		"manager":   "/manager",
	}

	ctx.Send(info{
		Hostname:      hostname,
		Version:       vinxi.Version,
		Platform:      runtime.GOOS,
		Runtime:       runtime.Version(),
		NumCPU:        runtime.NumCPU(),
		NumGoroutines: runtime.NumGoroutine(),
		Links:         links,
	})
}

func (IndexController) Catalog(ctx *Context) {
	rules := []rule.Info{}
	for _, rule := range rule.Rules {
		rules = append(rules, rule)
	}

	plugins := []plugin.Info{}
	for _, plugin := range plugin.Plugins {
		plugins = append(plugins, plugin)
	}

	catalog := struct {
		Rules   []rule.Info   `json:"rules"`
		Plugins []plugin.Info `json:"plugins"`
	}{
		Rules:   rules,
		Plugins: plugins,
	}

	ctx.Send(catalog)
}

func (IndexController) Manager(ctx *Context) {
	info := struct {
		Links map[string]string `json:"links"`
	}{
		Links: map[string]string{"plugins": "/manager/plugins"},
	}

	ctx.Send(info)
}
