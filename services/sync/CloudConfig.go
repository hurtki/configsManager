package sync_services

type CloudConfigInfo struct {
	ConfigKey string
	Checksum  [32]byte
}

func DefaultCloudManagerConfigFile() []CloudConfigInfo {
	return []CloudConfigInfo{
		{
			ConfigKey: "test",
			Checksum:  [32]byte{},
		},
	}
}
