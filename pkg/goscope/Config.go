package goscope

// Config is the global instance of the application's configuration.
var Config Environment //nolint:gochecknoglobals

// Initialize the configuration instance to the values provided by the user.
func configSetup(config *Environment) {
	if config == nil {
		panic("Please provide a pointer to a valid and instantiated GoScopeApplicationEnvironment.")
	}

	Config = *config
}
