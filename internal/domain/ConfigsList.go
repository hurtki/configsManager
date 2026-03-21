package domain

type ConfigsList struct {
	Configs map[string]string
}

func NewConfigsList(cfgs map[string]string) *ConfigsList {
	return &ConfigsList{
		Configs: cfgs,
	}
}

func GetDefaultConfigsList(pathToAppConfig string) *ConfigsList {
	configs := make(map[string]string)
	configs["cm_config"] = pathToAppConfig
	return &ConfigsList{
		Configs: configs,
	}
}

func (cl *ConfigsList) SetConfig(name, path string) {
	cl.Configs[name] = path
}

func (cl *ConfigsList) GetAllKeys() []string {
	configs := cl.Configs
	keys := make([]string, 0, len(configs))
	for k := range configs {
		keys = append(keys, k)
	}
	return keys
}

func (cl *ConfigsList) GetPath(name string) (string, bool) {
	val, ok := cl.Configs[name]
	if ok {
		return val, true
	} else {
		return "", false
	}
}

func (cl *ConfigsList) RemoveConfig(name string) {
	delete(cl.Configs, name)
}

func (cl *ConfigsList) HasKey(name string) bool {
	_, ok := cl.Configs[name]
	return ok
}
