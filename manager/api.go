package manager

// Initialize HTTP controllers
var index indexController
var rules scopesController
var scopes scopesController
var plugins pluginsController
var instances instancesController

// routes stores the registered routes.
var routes = []*Route{}

// route registers a new HTTP route based on the
// given verb and path expression.
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

	// Manager admin level plugins routes
	route("GET", "/manager", index.Manager)
	route("GET", "/manager/plugins", plugins.List)
	route("POST", "/manager/plugins", plugins.Create)
	route("GET", "/manager/plugins/:plugin", plugins.Get)
	route("DELETE", "/manager/plugins/:plugin", plugins.Delete)

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
	route("POST", "/scopes/:scope/plugins", plugins.Create)
	route("GET", "/scopes/:scope/plugins/:plugin", plugins.Get)
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
	route("POST", "/instances/:instance/scopes", scopes.Create)
	route("GET", "/instances/:instance/scopes/:scope", scopes.Get)
	route("DELETE", "/instances/:instance/scopes/:scope", scopes.Delete)

	// Instance-specific, scope-specific plugins
	route("GET", "/instances/:instance/scopes/:scope/plugins", plugins.List)
	route("POST", "/instances/:instance/scopes/:scope/plugins", plugins.Create)
	route("GET", "/instances/:instance/scopes/:scope/plugins/:plugin", plugins.Get)
	route("DELETE", "/instances/:instance/scopes/:scope/plugins/:plugin", plugins.Delete)

	// Instance-specific, scope-specific rules
	route("GET", "/instances/:instance/scopes/:scope/rules", rules.List)
	route("POST", "/instances/:instance/scopes/:scope/rules", rules.Create)
	route("GET", "/instances/:instance/scopes/:scope/rules/:rule", rules.Get)
	route("DELETE", "/instances/:instance/scopes/:scope/rules/:rule", rules.Delete)
}
