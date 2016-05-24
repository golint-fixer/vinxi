package plugins

import (
	// Ugly but unique way to autoload subpackages
	_ "gopkg.in/vinxi/vinxi.v0/plugins/forward"
	_ "gopkg.in/vinxi/vinxi.v0/plugins/static"
)
