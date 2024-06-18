package utils

import (
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/interfaces"
	vcs "github.com/crawlab-team/crawlab/vcs"
)

func InitGitClientAuth(g interfaces.Git, gitClient *vcs.GitClient) {
	// set auth
	switch g.GetAuthType() {
	case constants.GitAuthTypeHttp:
		gitClient.SetAuthType(vcs.GitAuthTypeHTTP)
		gitClient.SetUsername(g.GetUsername())
		gitClient.SetPassword(g.GetPassword())
	case constants.GitAuthTypeSsh:
		gitClient.SetAuthType(vcs.GitAuthTypeSSH)
		gitClient.SetUsername(g.GetUsername())
		gitClient.SetPrivateKey(g.GetPassword())
	}
}
