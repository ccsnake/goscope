package web

import "embed"

//go:embed templates/goscope-components/*.gohtml templates/goscope-views/*.gohtml
var TemplateFiles embed.FS //nolint:gochecknoglobals
