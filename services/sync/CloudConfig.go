package sync_services

type CloudConfigRegistry struct {
	Configs map[string][32]byte
}

func (r *CloudConfigRegistry) GetAllKeys() []string {
	keys := []string{}
	for key := range r.Configs {
		keys = append(keys, key)
	}
	return keys
}

func (r *CloudConfigRegistry) SetChecksum(key string, checksum [32]byte) {
	r.Configs[key] = checksum
}
