package sync_cmd

import "errors"

var (
	ErrPullBothFlagsRequired            = errors.New("both --all and --sp flags must be set when no arguments are provided")
	ErrPullAllFlagNotSupported          = errors.New("--all flag is not supported when one argument is provided")
	ErrPullAllAndSpFlagsNotSupported    = errors.New("--all and --sp flags are not supported with two arguments provided")
	ErrPullMoreThanTwoArgumentsProvided = errors.New("'cm sync pull' doesn't support more than 2 arguments")
)
