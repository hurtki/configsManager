package sync_services

import "github.com/hurtki/configsManager/services"

// Dependencies for sync_... commands
type Deps struct {
	AppConfigService   services.AppConfigService
	ConfigsListService services.ConfigsListService
	OsService          services.OsService
}
