package vcs

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
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

func WithDefaultBranch(branch string) GitOption {
	return func(c *GitClient) {
		c.defaultBranch = branch
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

func WithAuthClone(auth transport.AuthMethod) GitCloneOption {
	return func(o *git.CloneOptions) {
		o.Auth = auth
	}
}

func WithRemoteName(name string) GitCloneOption {
	return func(o *git.CloneOptions) {
		o.RemoteName = name
	}
}

func WithSingleBranch(singleBranch bool) GitCloneOption {
	return func(o *git.CloneOptions) {
		o.SingleBranch = singleBranch
	}
}

func WithNoCheckout(noCheckout bool) GitCloneOption {
	return func(o *git.CloneOptions) {
		o.NoCheckout = noCheckout
	}
}

func WithDepthClone(depth int) GitCloneOption {
	return func(o *git.CloneOptions) {
		o.Depth = depth
	}
}

func WithRecurseSubmodules(recurseSubmodules git.SubmoduleRescursivity) GitCloneOption {
	return func(o *git.CloneOptions) {
		o.RecurseSubmodules = recurseSubmodules
	}
}

func WithTags(tags git.TagMode) GitCloneOption {
	return func(o *git.CloneOptions) {
		o.Tags = tags
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

func WithAll(all bool) GitCommitOption {
	return func(o *git.CommitOptions) {
		o.All = all
	}
}

func WithAuthor(author *object.Signature) GitCommitOption {
	return func(o *git.CommitOptions) {
		o.Author = author
	}
}

func WithCommitter(committer *object.Signature) GitCommitOption {
	return func(o *git.CommitOptions) {
		o.Committer = committer
	}
}

func WithParents(parents []plumbing.Hash) GitCommitOption {
	return func(o *git.CommitOptions) {
		o.Parents = parents
	}
}

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

func WithDepthPull(depth int) GitPullOption {
	return func(o *git.PullOptions) {
		o.Depth = depth
	}
}

func WithAuthPull(auth transport.AuthMethod) GitPullOption {
	return func(o *git.PullOptions) {
		if auth != nil {
			o.Auth = auth
		}
	}
}

func WithRecurseSubmodulesPull(recurseSubmodules git.SubmoduleRescursivity) GitPullOption {
	return func(o *git.PullOptions) {
		o.RecurseSubmodules = recurseSubmodules
	}
}

func WithForcePull(force bool) GitPullOption {
	return func(o *git.PullOptions) {
		o.Force = force
	}
}

type GitPushOption func(o *git.PushOptions)

func WithRemoteNamePush(name string) GitPushOption {
	return func(o *git.PushOptions) {
		o.RemoteName = name
	}
}

func WithRefSpecs(specs []config.RefSpec) GitPushOption {
	return func(o *git.PushOptions) {
		o.RefSpecs = specs
	}
}

func WithAuthPush(auth transport.AuthMethod) GitPushOption {
	return func(o *git.PushOptions) {
		o.Auth = auth
	}
}

func WithPrune(prune bool) GitPushOption {
	return func(o *git.PushOptions) {
		o.Prune = prune
	}
}

func WithForcePush(force bool) GitPushOption {
	return func(o *git.PushOptions) {
		o.Force = force
	}
}

type GitResetOption func(o *git.ResetOptions)

func WithCommit(commit plumbing.Hash) GitResetOption {
	return func(o *git.ResetOptions) {
		o.Commit = commit
	}
}

func WithMode(mode git.ResetMode) GitResetOption {
	return func(o *git.ResetOptions) {
		o.Mode = mode
	}
}
