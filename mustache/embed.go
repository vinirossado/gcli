package mustache

import "embed"

//go:embed create/*.mustache
var CreateTemplateFS embed.FS
