package services

type ConfigsList struct {
	configs map[string]string
}

func GetDefaultConfigsList(pathToAppConfig string) *ConfigsList {
	configs := make(map[string]string)
	configs["cm_config"] = pathToAppConfig
	return &ConfigsList{
		configs: configs,
	}
}

func (cl *ConfigsList) SetConfig(name, path string) {
	cl.configs[name] = path
}

func (cl *ConfigsList) GetAllKeys() []string {
	configs := cl.configs
	keys := make([]string, 0, len(configs))
	for k := range configs {
		keys = append(keys, k)
	}
	return keys
}

func (cl *ConfigsList) GetPath(name string) (string, bool) {
	val, ok := cl.configs[name]
	if ok {
		return val, true
	} else {
		return "", false
	}
}

func (cl *ConfigsList) RemoveConfig(name string) {
	delete(cl.configs, name)
}

func (cl *ConfigsList) HasKey(name string) bool {
	_, ok := cl.configs[name]
	return ok
}
