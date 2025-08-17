package services

type AppConfig struct {
	// Editor command that command "open" uses
	Editor *string `json:"editor"`
	// what cm will do when it sees an already exisging key
	// "default" - cm will ask you what to do + will notice to change this setting
	// "o" - cm will overwrite the existing key
	// "n" - cm will automatically create a new name
	// "ask" - cm will always ask you what to do
	IfKeyExists *string `json:"if_key_exists"`
}

// default pointers for default config
func ptrString(s string) *string { return &s }

// default config
var defaultConfig = AppConfig{
	Editor:      ptrString("vim"),
	IfKeyExists: ptrString("default"),
}

func NewDefaultAppConfig() *AppConfig {
	return &defaultConfig
}

// validate_IfKeyExists() returns True if everythink write and false if somethink is wrong with IfKeyExists field
func (cfg *AppConfig) validate_IfKeyExists() bool {
	if cfg.IfKeyExists == nil {
		return true
	}
	switch *cfg.IfKeyExists {
	case "o", "default", "n", "ask":
		return false
	default:
		return true
	}
}

// validateAppConfig() can validate and insert default values
// if validateAppConfig() changed at least one field it returns true
func (cfg *AppConfig) validateAppConfig() bool {
	changed := false

	if cfg.Editor == nil {
		editor := *defaultConfig.Editor
		cfg.Editor = &editor
		changed = true
	}

	if cfg.validate_IfKeyExists() {
		def := *defaultConfig.IfKeyExists
		cfg.IfKeyExists = &def
		changed = true
	}

	return changed
}
