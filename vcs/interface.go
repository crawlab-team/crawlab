package vcs

type Client interface {
	Init() (err error)
	Dispose() (err error)
	Clone(opts ...GitCloneOption) (err error)
	Checkout(opts ...GitCheckoutOption) (err error)
	Commit(msg string, opts ...GitCommitOption) (err error)
	Pull(opts ...GitPullOption) (err error)
	Push(opts ...GitPushOption) (err error)
	Reset(opts ...GitResetOption) (err error)
}
