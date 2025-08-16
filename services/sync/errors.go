package sync_services

import "errors"

var (
	ErrFileDoesntExist         = errors.New("file doesn't exists")
	ErrAuthProviderDoesntExist = errors.New("given provider doesn't exist")
	ErrTokenNotFoundInSecrets  = errors.New("token not found in secrets, try cm sync auth")
	ErrKeyNotFoundInCloud      = errors.New("token not found in cloud")
	ErrFileChanged             = errors.New("file already exists and has changed")
)
