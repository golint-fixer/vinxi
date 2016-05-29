package rules

import (
	// Ugly but unique way to autoload subpackages
	_ "gopkg.in/vinxi/vinxi.v0/rules/path"
	_ "gopkg.in/vinxi/vinxi.v0/rules/vhost"
)
