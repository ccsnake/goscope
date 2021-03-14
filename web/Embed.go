package web

import "embed"

//go:embed templates/goscope-components/*.gohtml templates/goscope-views/*.gohtml
var TemplateFiles embed.FS //nolint:gochecknoglobals

//go:embed styles/*.css
var StyleFiles embed.FS

//go:embed scripts/*.js
var ScriptFiles embed.FS
