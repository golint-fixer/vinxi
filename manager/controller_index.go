package manager

import (
	"os"
	"runtime"

	"gopkg.in/vinxi/vinxi.v0"
	"gopkg.in/vinxi/vinxi.v0/plugin"
	"gopkg.in/vinxi/vinxi.v0/rule"
)

// IndexController represents the base routes HTTP controller.
type IndexController struct{}

func (IndexController) Get(ctx *Context) {
	hostname, _ := os.Hostname()

	info := struct {
		Hostname string            `json:"hostname"`
		Version  string            `json:"version"`
		Platform string            `json:"platform"`
		Links    map[string]string `json:"links"`
	}{
		Hostname: hostname,
		Version:  vinxi.Version,
		Platform: runtime.GOOS,
		Links: map[string]string{
			"catalog":   "/catalog",
			"plugins":   "/plugins",
			"scopes":    "/scopes",
			"instances": "/instances",
		},
	}

	ctx.Send(info)
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
