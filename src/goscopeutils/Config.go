package goscopeutils

import "github.com/averageflow/goscope/v2/src/goscopetypes"

// Config is the global instance of the application's configuration.
var Config goscopetypes.GoScopeApplicationEnvironment //nolint:gochecknoglobals

// Initialize the configuration instance to the values provided by the user.
func ConfigSetup(config *goscopetypes.GoScopeApplicationEnvironment) {
	if config == nil {
		panic("Please provide a pointer to a valid and instantiated GoScopeApplicationEnvironment.")
	}

	Config = *config
}
