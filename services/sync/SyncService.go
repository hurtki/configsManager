package sync_services

var cloudManagerFileName = "cloud_manger.json"

type SyncService interface {
	// Authorization
	Auth(provider string, token string) error
	Logout(provider string) error // пустая строка = всех

	// Pulling
	PullAll(samePlace bool) ([]*ConfigObj, error)
	PullOne(key string) (*ConfigObj, error)

	// Pushing
	Push(configs []*ConfigObj, force bool) map[*ConfigObj]error
}

type SyncServiceImpl struct {
	CloudManager CloudManager
}
