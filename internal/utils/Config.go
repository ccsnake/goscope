package utils

import (
	"github.com/averageflow/goscope/v3/pkg/goscope"
)

// Config is the global instance of the application's configuration.
var Config goscope.Environment //nolint:gochecknoglobals

// Initialize the configuration instance to the values provided by the user.
func ConfigSetup(config *goscope.Environment) {
	if config == nil {
		panic("Please provide a pointer to a valid and instantiated GoScopeApplicationEnvironment.")
	}

	Config = *config
}
