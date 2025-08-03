package services

type AppConfig struct {
	// Editor command that command "open" uses
	Editor *string `json:"editor"`
	// Is cm going to overwrite an existing key in configs list
	// if true cm won't ask you
	// if false cm will ask you "If you want to overwrite"
	ForceOverwrite *bool `json:"overwrite_if_exists"`
	// is cm going to add the path if it doesn't exist
	// if true cm won't ask you
	// if false cm will ask you "If you want to add the non existing path"
	ForceAddPath *bool `json:"force_add_path"`
}

// default pointers for default config
func ptrBool(b bool) *bool       { return &b }
func ptrString(s string) *string { return &s }

// default config
var defaultConfig = AppConfig{
	Editor:         ptrString("vim"),
	ForceOverwrite: ptrBool(false),
	ForceAddPath:   ptrBool(false),
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

	if cfg.ForceOverwrite == nil {
		def := *defaultConfig.ForceOverwrite
		cfg.ForceOverwrite = &def
		changed = true
	}

	if cfg.ForceAddPath == nil {
		def := *defaultConfig.ForceAddPath
		cfg.ForceAddPath = &def
		changed = true
	}

	return changed
}
