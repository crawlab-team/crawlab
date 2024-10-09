package vcs

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"strings"
)

type GitOption func(c *GitClient)

func WithPath(path string) GitOption {
	return func(c *GitClient) {
		c.path = path
	}
}

func WithRemoteUrl(url string) GitOption {
	return func(c *GitClient) {
		c.remoteUrl = url
	}
}

func WithIsMem() GitOption {
	return func(c *GitClient) {
		c.isMem = true
	}
}

func WithAuthType(authType GitAuthType) GitOption {
	return func(c *GitClient) {
		c.authType = authType
	}
}

func WithUsername(username string) GitOption {
	return func(c *GitClient) {
		c.username = username
	}
}

func WithPassword(password string) GitOption {
	return func(c *GitClient) {
		c.password = password
	}
}

func WithPrivateKey(key string) GitOption {
	return func(c *GitClient) {
		c.privateKey = key
	}
}

func WithDefaultInit(init bool) GitOption {
	return func(c *GitClient) {
		c.defaultInit = init
	}
}

func WithPrivateKeyPath(path string) GitOption {
	return func(c *GitClient) {
		c.privateKeyPath = path
	}
}

type GitCloneOption func(o *git.CloneOptions)

func WithURL(url string) GitCloneOption {
	return func(o *git.CloneOptions) {
		o.URL = url
	}
}

type GitCheckoutOption func(o *git.CheckoutOptions)

func WithBranch(branch string) GitCheckoutOption {
	return func(o *git.CheckoutOptions) {
		if strings.HasPrefix(branch, "refs/heads") {
			o.Branch = plumbing.ReferenceName(branch)
		} else {
			o.Branch = plumbing.NewBranchReferenceName(branch)
		}
	}
}

func WithHash(hash string) GitCheckoutOption {
	return func(o *git.CheckoutOptions) {
		h := plumbing.NewHash(hash)
		if h.IsZero() {
			return
		}
		o.Hash = h
	}
}

type GitCommitOption func(o *git.CommitOptions)

type GitPullOption func(o *git.PullOptions)

func WithRemoteNamePull(name string) GitPullOption {
	return func(o *git.PullOptions) {
		o.RemoteName = name
	}
}

func WithBranchNamePull(branch string) GitPullOption {
	return func(o *git.PullOptions) {
		o.ReferenceName = plumbing.NewBranchReferenceName(branch)
	}
}

func WithAuthPull(auth transport.AuthMethod) GitPullOption {
	return func(o *git.PullOptions) {
		if auth != nil {
			o.Auth = auth
		}
	}
}

type GitPushOption func(o *git.PushOptions)

func WithAuthPush(auth transport.AuthMethod) GitPushOption {
	return func(o *git.PushOptions) {
		o.Auth = auth
	}
}

type GitResetOption func(o *git.ResetOptions)

func WithMode(mode git.ResetMode) GitResetOption {
	return func(o *git.ResetOptions) {
		o.Mode = mode
	}
}
