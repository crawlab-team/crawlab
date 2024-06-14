package interfaces

// Git interface
type Git interface {
	Model
	GetUrl() (url string)
	SetUrl(url string)
	GetAuthType() (authType string)
	SetAuthType(authType string)
	GetUsername() (username string)
	SetUsername(username string)
	GetPassword() (password string)
	SetPassword(password string)
	GetCurrentBranch() (currentBranch string)
	SetCurrentBranch(currentBranch string)
	GetAutoPull() (autoPull bool)
	SetAutoPull(autoPull bool)
}
