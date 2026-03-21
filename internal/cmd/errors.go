package cmd

import "errors"

var (
	ErrUserAborted = errors.New("user chose to quit")
)
