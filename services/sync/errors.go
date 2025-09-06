package sync_services

import "errors"

var (
	ErrFileDoesntExist          = errors.New("file doesn't exists")
	ErrAuthProviderDoesntExist  = errors.New("given provider doesn't exist")
	ErrTokenNotFoundInSecrets   = errors.New("token not found in secrets, try cm sync auth")
	ErrKeyNotFoundInCloud       = errors.New("key not found in cloud")
	ErrFileChanged              = errors.New("file already exists and has changed")
	ErrNotAuthenticated         = errors.New("token not found in storage")
	ErrNothingToPush            = errors.New("all configs are synced, nothing to push")
	ErrUnauthorizedRequest      = errors.New("not valid access token in request to dropbox api using sdk")
	ErrRetrieveTokenFromStorage = errors.New("retrieving tokens from storage went wrong (you likely typed passphrase wrong!)")
)
