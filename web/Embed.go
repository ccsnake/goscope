package web

import "embed"

//go:embed templates/components/*.gohtml templates/views/*.gohtml
var TemplateFiles embed.FS
