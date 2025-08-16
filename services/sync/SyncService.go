package sync_services

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
	AuthManager  AuthManager
}

func NewSyncServiceImpl() (*SyncServiceImpl, error) {
	authManager := NewAuthManagerImpl()
	authToken, err := authManager.GetToken("dropbox")
	if err != nil {
		return nil, err
	}
	return &SyncServiceImpl{
		AuthManager:  authManager,
		CloudManager: NewCloudManagerImpl(authToken),
	}, nil
}
