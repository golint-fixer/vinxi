package manager

// Initialize HTTP controllers
var index IndexController
var rules ScopesController
var scopes ScopesController
var plugins PluginsController
var instances InstancesController

// routes stores the registered routes.
var routes = []*Route{}

func route(method, path string, fn RouteHandler) {
	routes = append(routes, &Route{
		Path:    path,
		Method:  method,
		Handler: fn,
	})
}

func init() {
	// General routes
	route("GET", "/", index.Get)
	route("GET", "/catalog", index.Catalog)

	// Global plugins routes
	route("GET", "/plugins", plugins.List)
	route("POST", "/plugins", plugins.Create)
	route("GET", "/plugins/:plugin", plugins.Get)
	route("DELETE", "/plugins/:plugin", plugins.Delete)

	// Global scopes routes
	route("GET", "/scopes", scopes.List)
	route("POST", "/scopes", scopes.Create)
	route("GET", "/scopes/:scope", scopes.Get)
	route("DELETE", "/scopes/:scope", scopes.Delete)

	// Scope-specific plugins routes
	route("GET", "/scopes/:scope/plugins", plugins.List)
	route("GET", "/scopes/:scope/plugins/:plugin", plugins.Get)
	route("POST", "/scopes/:scope/plugins/:plugin", plugins.Create)
	route("DELETE", "/scopes/:scope/plugins/:plugin", plugins.Delete)

	// Scope-specific rules routes
	route("GET", "/scopes/:scope/rules", rules.List)
	route("GET", "/scopes/:scope/rules/:rule", rules.Get)
	route("POST", "/scopes/:scope/rules/:rule", rules.Create)
	route("DELETE", "/scopes/:scope/rules/:rule", rules.Delete)

	// Instances routes
	route("GET", "/instances", instances.List)
	route("GET", "/instances/:instance", instances.Get)
	route("DELETE", "/instances/:instance", instances.Delete)

	// Instance-specific scopes
	route("GET", "/instances/:instance/scopes", scopes.List)
	route("GET", "/instances/:instance/scopes/:scope", scopes.Get)
	route("POST", "/instances/:instance/scopes/:scope", scopes.Create)
	route("DELETE", "/instances/:instance/scopes/:scope", scopes.Delete)

	// Instance-specific, scope-specific plugins
	route("GET", "/instances/:instance/scopes/:scope/plugins", plugins.List)
	route("GET", "/instances/:instance/scopes/:scope/plugins/:plugin", plugins.Get)
	route("POST", "/instances/:instance/scopes/:scope/plugins/:plugin", plugins.Create)
	route("DELETE", "/instances/:instance/scopes/:scope/plugins/:plugin", plugins.Delete)

	// Instance-specific, scope-specific rules
	route("GET", "/instances/:instance/scopes/:scope/rules", rules.List)
	route("GET", "/instances/:instance/scopes/:scope/rules/:rule", rules.Get)
	route("POST", "/instances/:instance/scopes/:scope/rules/:rule", rules.Create)
	route("DELETE", "/instances/:instance/scopes/:scope/rules/:rule", rules.Delete)
}
