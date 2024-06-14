package vcs

import "errors"

var (
	ErrInvalidArgsLength               = errors.New("invalid arguments length")
	ErrUnsupportedType                 = errors.New("unsupported type")
	ErrInvalidAuthType                 = errors.New("invalid auth type")
	ErrInvalidOptions                  = errors.New("invalid options")
	ErrRepoAlreadyExists               = errors.New("repo already exists")
	ErrInvalidRepoPath                 = errors.New("invalid repo path")
	ErrUnableToGetCurrentBranch        = errors.New("unable to get current branch")
	ErrUnableToCloneWithEmptyRemoteUrl = errors.New("unable to clone with empty remote url")
	ErrInvalidHeadRef                  = errors.New("invalid head ref")
	ErrNoMatchedRemoteBranch           = errors.New("no matched remote branch")
)
