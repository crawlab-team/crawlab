package vcs

import (
	"errors"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab/trace"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	gitssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/memory"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"
)

var headRefRegexp, _ = regexp.Compile("^ref: (.*)")

type GitClient struct {
	// settings
	path           string
	remoteUrl      string
	isMem          bool
	authType       GitAuthType
	username       string
	password       string
	privateKey     string
	privateKeyPath string
	defaultBranch  string
	defaultInit    bool

	// internals
	r *git.Repository
}

func (c *GitClient) Init() (err error) {
	initType := c.getInitType()
	switch initType {
	case GitInitTypeFs:
		err = c.initFs()
	case GitInitTypeMem:
		err = c.initMem()
	}
	if err != nil {
		return err
	}
	return nil
}

func (c *GitClient) Dispose() (err error) {
	switch c.getInitType() {
	case GitInitTypeFs:
		if err := os.RemoveAll(c.path); err != nil {
			return trace.TraceError(err)
		}
	case GitInitTypeMem:
		GitMemStorages.Delete(c.path)
		GitMemFileSystem.Delete(c.path)
	}
	return nil
}

func (c *GitClient) Clone() (err error) {
	// validate
	if c.remoteUrl == "" {
		return ErrUnableToCloneWithEmptyRemoteUrl
	}

	// auth
	auth, err := c.getGitAuth()
	if err != nil {
		return err
	}

	// options
	o := &git.CloneOptions{
		URL:  c.remoteUrl,
		Auth: auth,
	}

	// clone
	_, err = git.PlainClone(c.path, false, o)
	if err != nil {
		return err
	}

	return nil
}

func (c *GitClient) Checkout(opts ...GitCheckoutOption) (err error) {
	// worktree
	wt, err := c.r.Worktree()
	if err != nil {
		return trace.TraceError(err)
	}

	// apply options
	o := &git.CheckoutOptions{}
	for _, opt := range opts {
		opt(o)
	}

	// checkout to the branch
	if err := wt.Checkout(o); err != nil {
		return trace.TraceError(err)
	}

	return nil
}

func (c *GitClient) Commit(msg string, opts ...GitCommitOption) (err error) {
	// worktree
	wt, err := c.r.Worktree()
	if err != nil {
		return trace.TraceError(err)
	}

	// apply options
	o := &git.CommitOptions{}
	for _, opt := range opts {
		opt(o)
	}

	// commit
	if _, err := wt.Commit(msg, o); err != nil {
		return trace.TraceError(err)
	}

	return nil
}

func (c *GitClient) Pull(opts ...GitPullOption) (err error) {
	// worktree
	wt, err := c.r.Worktree()
	if err != nil {
		return trace.TraceError(err)
	}

	// auth
	auth, err := c.getGitAuth()
	if err != nil {
		return err
	}
	if auth != nil {
		opts = append(opts, WithAuthPull(auth))
	}

	// apply options
	o := &git.PullOptions{}
	for _, opt := range opts {
		opt(o)
	}

	// pull
	if err := wt.Pull(o); err != nil {
		if errors.Is(err, transport.ErrEmptyRemoteRepository) {
			return nil
		}
		if errors.Is(err, transport.ErrEmptyUploadPackRequest) {
			return nil
		}
		if errors.Is(err, git.NoErrAlreadyUpToDate) {
			return nil
		}
		if errors.Is(err, git.ErrNonFastForwardUpdate) {
			return nil
		}
		return trace.TraceError(err)
	}

	return nil
}

func (c *GitClient) Push(opts ...GitPushOption) (err error) {
	// auth
	auth, err := c.getGitAuth()
	if err != nil {
		return err
	}
	if auth != nil {
		opts = append(opts, WithAuthPush(auth))
	}

	// apply options
	o := &git.PushOptions{}
	for _, opt := range opts {
		opt(o)
	}

	// push
	if err := c.r.Push(o); err != nil {
		return trace.TraceError(err)
	}

	return nil
}

func (c *GitClient) Reset(opts ...GitResetOption) (err error) {
	// apply options
	o := &git.ResetOptions{
		Mode: git.HardReset,
	}
	for _, opt := range opts {
		opt(o)
	}

	// worktree
	wt, err := c.r.Worktree()
	if err != nil {
		return err
	}

	// reset
	if err := wt.Reset(o); err != nil {
		return err
	}

	// clean
	if err := wt.Clean(&git.CleanOptions{Dir: true}); err != nil {
		return err
	}

	return nil
}

func (c *GitClient) CreateBranch(branch, remote string, ref *plumbing.Reference) (err error) {
	return c.createBranch(branch, remote, ref)
}

func (c *GitClient) CheckoutBranchFromRef(branch string, ref *plumbing.Reference, opts ...GitCheckoutOption) (err error) {
	return c.CheckoutBranchWithRemote(branch, "", ref, opts...)
}

func (c *GitClient) CheckoutBranchWithRemoteFromRef(branch, remote string, ref *plumbing.Reference, opts ...GitCheckoutOption) (err error) {
	return c.CheckoutBranchWithRemote(branch, remote, ref, opts...)
}

func (c *GitClient) CheckoutBranch(branch string, opts ...GitCheckoutOption) (err error) {
	return c.CheckoutBranchWithRemote(branch, "", nil, opts...)
}

func (c *GitClient) CheckoutBranchWithRemote(branch, remote string, ref *plumbing.Reference, opts ...GitCheckoutOption) (err error) {
	if remote == "" {
		remote = GitRemoteNameOrigin
	}

	// remote
	if _, err := c.r.Remote(remote); err != nil {
		return trace.TraceError(err)
	}

	// check if the branch exists
	_, err = c.r.Branch(branch)
	if err != nil {
		if errors.Is(err, git.ErrBranchNotFound) {
			// create a new branch if it does not exist
			if err := c.createBranch(branch, remote, ref); err != nil {
				return err
			}
			_, err = c.r.Branch(branch)
			if err != nil {
				return trace.TraceError(err)
			}
		} else {
			return trace.TraceError(err)
		}
	}

	// add to options
	opts = append(opts, WithBranch(branch))

	return c.Checkout(opts...)
}

func (c *GitClient) CheckoutHash(hash string, opts ...GitCheckoutOption) (err error) {
	// add to options
	opts = append(opts, WithHash(hash))

	return c.Checkout(opts...)
}

func (c *GitClient) MoveBranch(from, to string) (err error) {
	wt, err := c.r.Worktree()
	if err != nil {
		return trace.TraceError(err)
	}
	if err := wt.Checkout(&git.CheckoutOptions{
		Create: true,
		Branch: plumbing.NewBranchReferenceName(to),
	}); err != nil {
		return trace.TraceError(err)
	}
	fromRef, err := c.r.Reference(plumbing.NewBranchReferenceName(from), false)
	if err != nil {
		return trace.TraceError(err)
	}
	if err := c.r.Storer.RemoveReference(fromRef.Name()); err != nil {
		return trace.TraceError(err)
	}
	return nil
}

func (c *GitClient) CommitAll(msg string, opts ...GitCommitOption) (err error) {
	// worktree
	wt, err := c.r.Worktree()
	if err != nil {
		return trace.TraceError(err)
	}

	// add all files
	if _, err := wt.Add("."); err != nil {
		return trace.TraceError(err)
	}

	return c.Commit(msg, opts...)
}

func (c *GitClient) GetLogs() (logs []GitLog, err error) {
	iter, err := c.r.Log(&git.LogOptions{
		All: true,
	})
	if err != nil {
		return nil, trace.TraceError(err)
	}
	if err := iter.ForEach(func(commit *object.Commit) error {
		gitLog := GitLog{
			Hash:        commit.Hash.String(),
			Msg:         commit.Message,
			AuthorName:  commit.Author.Name,
			AuthorEmail: commit.Author.Email,
			Timestamp:   commit.Author.When,
		}
		logs = append(logs, gitLog)
		return nil
	}); err != nil {
		return nil, trace.TraceError(err)
	}
	return
}

func (c *GitClient) GetLogsWithRefs() (logs []GitLog, err error) {
	// logs without tags
	logs, err = c.GetLogs()
	if err != nil {
		return nil, err
	}

	// branches
	branches, err := c.GetBranches()
	if err != nil {
		return nil, err
	}

	// tags
	tags, err := c.GetTags()
	if err != nil {
		return nil, err
	}

	// refs
	refs := append(branches, tags...)

	// refs map
	refsMap := map[string][]GitRef{}
	for _, ref := range refs {
		_, ok := refsMap[ref.Hash]
		if !ok {
			refsMap[ref.Hash] = []GitRef{}
		}
		refsMap[ref.Hash] = append(refsMap[ref.Hash], ref)
	}

	// iterate logs
	for i, l := range logs {
		refs, ok := refsMap[l.Hash]
		if ok {
			logs[i].Refs = refs
		}
	}

	return logs, nil
}

func (c *GitClient) GetRepository() (r *git.Repository) {
	return c.r
}

func (c *GitClient) GetPath() (path string) {
	return c.path
}

func (c *GitClient) SetPath(path string) {
	c.path = path
}

func (c *GitClient) GetRemoteUrl() (path string) {
	return c.remoteUrl
}

func (c *GitClient) SetRemoteUrl(url string) {
	c.remoteUrl = url
}

func (c *GitClient) GetIsMem() (isMem bool) {
	return c.isMem
}

func (c *GitClient) SetIsMem(isMem bool) {
	c.isMem = isMem
}

func (c *GitClient) GetAuthType() (authType GitAuthType) {
	return c.authType
}

func (c *GitClient) SetAuthType(authType GitAuthType) {
	c.authType = authType
}

func (c *GitClient) GetUsername() (username string) {
	return c.username
}

func (c *GitClient) SetUsername(username string) {
	c.username = username
}

func (c *GitClient) GetPassword() (password string) {
	return c.password
}

func (c *GitClient) SetPassword(password string) {
	c.password = password
}

func (c *GitClient) GetPrivateKey() (key string) {
	return c.privateKey
}

func (c *GitClient) SetPrivateKey(key string) {
	c.privateKey = key
}

func (c *GitClient) GetPrivateKeyPath() (path string) {
	return c.privateKeyPath
}

func (c *GitClient) SetPrivateKeyPath(path string) {
	c.privateKeyPath = path
}

func (c *GitClient) GetCurrentBranch() (branch string, err error) {
	// attempt to get branch from .git/HEAD
	headRefStr, err := c.getHeadRef()
	if err != nil {
		return "", err
	}

	// if .git/HEAD points to refs/heads/master, return branch as master
	if headRefStr == plumbing.Master.String() {
		return GitBranchNameMaster, nil
	}

	// attempt to get head ref
	headRef, err := c.r.Head()
	if err != nil {
		return "", trace.TraceError(err)
	}
	if !headRef.Name().IsBranch() {
		return "", trace.TraceError(ErrUnableToGetCurrentBranch)
	}

	return headRef.Name().Short(), nil
}

func (c *GitClient) GetCurrentBranchRef() (ref *GitRef, err error) {
	currentBranch, err := c.GetCurrentBranch()
	if err != nil {
		return nil, err
	}
	branches, err := c.GetBranches()
	if err != nil {
		return nil, err
	}
	for _, branch := range branches {
		if branch.Name == currentBranch {
			return &branch, nil
		}
	}
	return nil, trace.TraceError(ErrUnableToGetCurrentBranch)
}

func (c *GitClient) GetBranches() (branches []GitRef, err error) {
	iter, err := c.r.Branches()
	if err != nil {
		return nil, trace.TraceError(err)
	}

	_ = iter.ForEach(func(r *plumbing.Reference) error {
		branches = append(branches, GitRef{
			Type: GitRefTypeBranch,
			Name: r.Name().Short(),
			Hash: r.Hash().String(),
		})
		return nil
	})

	return branches, nil
}

func (c *GitClient) GetRemoteRefs(remoteName string) (gitRefs []GitRef, err error) {
	if remoteName == "" {
		remoteName = GitRemoteNameOrigin
	}

	// remote
	r, err := c.r.Remote(remoteName)
	if err != nil {
		if errors.Is(err, git.ErrRemoteNotFound) {
			return nil, nil
		}
		return nil, trace.TraceError(err)
	}

	// auth
	auth, err := c.getGitAuth()
	if err != nil {
		return nil, err
	}

	// refs
	refs, err := r.List(&git.ListOptions{Auth: auth})
	if err != nil {
		if !errors.Is(err, transport.ErrEmptyRemoteRepository) {
			return nil, trace.TraceError(err)
		}
		return nil, nil
	}

	// iterate refs
	for _, ref := range refs {
		// ref type
		var refType string
		var gitRef *GitRef
		if strings.HasPrefix(ref.Name().String(), "refs/heads") {
			refType = GitRefTypeBranch
			gitRef = &GitRef{
				Type:     refType,
				Name:     remoteName + "/" + ref.Name().Short(),
				FullName: ref.Name().String(),
				Hash:     ref.Hash().String(),
			}
		} else if strings.HasPrefix(ref.Name().String(), "refs/tags") {
			refType = GitRefTypeTag
			gitRef = &GitRef{
				Type:     refType,
				Name:     ref.Name().Short(),
				FullName: ref.Name().String(),
				Hash:     ref.Hash().String(),
			}
		} else {
			continue
		}
		gitRefs = append(gitRefs, *gitRef)
	}

	// logs without tags
	logs, err := c.GetLogs()
	if err != nil {
		return nil, err
	}

	// logs map
	logsMap := map[string]GitLog{}
	for _, l := range logs {
		logsMap[l.Hash] = l
	}

	// iterate git refs
	for i, gitRef := range gitRefs {
		l, ok := logsMap[gitRef.Hash]
		if !ok {
			continue
		}
		gitRefs[i].Timestamp = l.Timestamp
	}

	// sort git refs
	sort.Slice(gitRefs, func(i, j int) bool {
		return gitRefs[i].Timestamp.Unix() > gitRefs[j].Timestamp.Unix()
	})

	return gitRefs, nil
}

func (c *GitClient) GetTags() (tags []GitRef, err error) {
	iter, err := c.r.Tags()
	if err != nil {
		return nil, trace.TraceError(err)
	}

	_ = iter.ForEach(func(r *plumbing.Reference) error {
		tags = append(tags, GitRef{
			Type: GitRefTypeTag,
			Name: r.Name().Short(),
			Hash: r.Hash().String(),
		})
		return nil
	})

	return tags, nil
}

func (c *GitClient) GetStatus() (statusList []GitFileStatus, err error) {
	// worktree
	wt, err := c.r.Worktree()
	if err != nil {
		return nil, trace.TraceError(err)
	}

	// status
	status, err := wt.Status()
	if err != nil {
		log.Warnf("failed to get worktree status: %v", err)
	}

	// file status list
	var list []GitFileStatus
	for filePath, fileStatus := range status {
		// file name
		fileName := path.Base(filePath)

		// file status
		s := GitFileStatus{
			Path:     filePath,
			Name:     fileName,
			IsDir:    false,
			Staging:  c.getStatusString(fileStatus.Staging),
			Worktree: c.getStatusString(fileStatus.Worktree),
			Extra:    fileStatus.Extra,
		}

		// add to list
		list = append(list, s)
	}

	// sort list ascending
	sort.Slice(list, func(i, j int) bool {
		return list[i].Path < list[j].Path
	})

	return list, nil
}

func (c *GitClient) Add(filePath string) (err error) {
	// worktree
	wt, err := c.r.Worktree()
	if err != nil {
		return trace.TraceError(err)
	}

	if _, err := wt.Add(filePath); err != nil {
		return trace.TraceError(err)
	}

	return nil
}

func (c *GitClient) GetRemote(name string) (r *git.Remote, err error) {
	return c.r.Remote(name)
}

func (c *GitClient) CreateRemote(cfg *config.RemoteConfig) (r *git.Remote, err error) {
	return c.r.CreateRemote(cfg)
}

func (c *GitClient) DeleteRemote(name string) (err error) {
	return c.r.DeleteRemote(name)
}

func (c *GitClient) IsRemoteChanged() (ok bool, err error) {
	return c.isRemoteChanged()
}

func (c *GitClient) initMem() (err error) {
	// validate options
	if !c.isMem || c.path == "" {
		return trace.TraceError(ErrInvalidOptions)
	}

	// get storage and worktree
	storage, wt := c.getMemStorageAndMemFs(c.path)

	// attempt to init
	c.r, err = git.Init(storage, wt)
	if err != nil {
		if errors.Is(err, git.ErrRepositoryAlreadyExists) {
			// if already exists, attempt to open
			c.r, err = git.Open(storage, wt)
			if err != nil {
				return trace.TraceError(err)
			}
		} else {
			return trace.TraceError(err)
		}
	}

	return nil
}

func (c *GitClient) initFs() (err error) {
	// validate options
	if c.path == "" {
		return trace.TraceError(ErrInvalidOptions)
	}

	// create directory if not exists
	_, err = os.Stat(c.path)
	if err != nil {
		if err := os.MkdirAll(c.path, os.ModePerm); err != nil {
			return trace.TraceError(err)
		}
		err = nil
	}

	// try to open repo
	c.r, err = git.PlainOpen(c.path)
	if err != nil {
		return err
	}

	return nil
}

func (c *GitClient) getInitType() (res GitInitType) {
	if c.isMem {
		return GitInitTypeMem
	} else {
		return GitInitTypeFs
	}
}

func (c *GitClient) createRemote(remoteName string, url string) (err error) {
	_, err = c.r.CreateRemote(&config.RemoteConfig{
		Name: remoteName,
		URLs: []string{url},
	})
	if err != nil {
		return trace.TraceError(err)
	}
	return
}

func (c *GitClient) getMemStorageAndMemFs(key string) (storage *memory.Storage, fs billy.Filesystem) {
	// storage
	storageItem, ok := GitMemStorages.Load(key)
	if !ok {
		storage = memory.NewStorage()
		GitMemStorages.Store(key, storage)
	} else {
		switch storageItem.(type) {
		case *memory.Storage:
			storage = storageItem.(*memory.Storage)
		default:
			storage = memory.NewStorage()
			GitMemStorages.Store(key, storage)
		}
	}

	// file system
	fsItem, ok := GitMemFileSystem.Load(key)
	if !ok {
		fs = memfs.New()
		GitMemFileSystem.Store(key, fs)
	} else {
		switch fsItem.(type) {
		case billy.Filesystem:
			fs = fsItem.(billy.Filesystem)
		default:
			fs = memfs.New()
			GitMemFileSystem.Store(key, fs)
		}
	}

	return storage, fs
}

func (c *GitClient) getGitAuth() (auth transport.AuthMethod, err error) {
	switch c.authType {
	case GitAuthTypeNone:
		return nil, nil
	case GitAuthTypeHTTP:
		if c.username == "" && c.password == "" {
			return nil, nil
		}
		auth = &http.BasicAuth{
			Username: c.username,
			Password: c.password,
		}
		return auth, nil
	case GitAuthTypeSSH:
		var privateKeyData []byte
		if c.privateKey != "" {
			// private key content
			privateKeyData = []byte(c.privateKey)
		} else if c.privateKeyPath != "" {
			// read from private key file
			privateKeyData, err = os.ReadFile(c.privateKeyPath)
			if err != nil {
				return nil, trace.TraceError(err)
			}
		} else {
			// no private key
			return nil, nil
		}
		signer, err := ssh.ParsePrivateKey(privateKeyData)
		if err != nil {
			return nil, trace.TraceError(err)
		}
		auth = &gitssh.PublicKeys{
			User:   c.username,
			Signer: signer,
			HostKeyCallbackHelper: gitssh.HostKeyCallbackHelper{
				HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			},
		}
		return auth, nil
	default:
		return nil, trace.TraceError(ErrInvalidAuthType)
	}
}

func (c *GitClient) getHeadRef() (ref string, err error) {
	wt, err := c.r.Worktree()
	if err != nil {
		return "", trace.TraceError(err)
	}
	fh, err := wt.Filesystem.Open(path.Join(".git", "HEAD"))
	if err != nil {
		return "", trace.TraceError(err)
	}
	data, err := io.ReadAll(fh)
	if err != nil {
		return "", trace.TraceError(err)
	}
	m := headRefRegexp.FindStringSubmatch(string(data))
	if len(m) < 2 {
		return "", trace.TraceError(ErrInvalidHeadRef)
	}
	return m[1], nil
}

func (c *GitClient) getStatusString(statusCode git.StatusCode) (code string) {
	return string(statusCode)
	//switch statusCode {
	//}
	//Unmodified         StatusCode = ' '
	//Untracked          StatusCode = '?'
	//Modified           StatusCode = 'M'
	//Added              StatusCode = 'A'
	//Deleted            StatusCode = 'D'
	//Renamed            StatusCode = 'R'
	//Copied             StatusCode = 'C'
	//UpdatedButUnmerged StatusCode = 'U'
}

func (c *GitClient) getDirPaths(filePath string) (paths []string) {
	pathItems := strings.Split(filePath, "/")

	var items []string
	for i, pathItem := range pathItems {
		if i == len(pathItems)-1 {
			continue
		}
		items = append(items, pathItem)
		dirPath := strings.Join(items, "/")
		paths = append(paths, dirPath)
	}

	return paths
}

func (c *GitClient) createBranch(branch, remote string, ref *plumbing.Reference) (err error) {
	// create a new branch if it does not exist
	cfg := config.Branch{
		Name:   branch,
		Remote: remote,
	}
	if err := c.r.CreateBranch(&cfg); err != nil {
		return err
	}
	// if ref is nil
	if ref == nil {
		// try to set to remote ref of branch first
		ref, err = c.getBranchHashRef(branch, remote)

		// if no matched remote branch, set to HEAD
		if errors.Is(err, ErrNoMatchedRemoteBranch) {
			ref, err = c.r.Head()
			if err != nil {
				return trace.TraceError(err)
			}
		}

		// error
		if err != nil {
			return trace.TraceError(err)
		}
	}

	// branch reference name
	branchRefName := plumbing.NewBranchReferenceName(branch)

	// branch reference
	branchRef := plumbing.NewHashReference(branchRefName, ref.Hash())

	// set HEAD to branch reference
	if err := c.r.Storer.SetReference(branchRef); err != nil {
		return err
	}

	return nil
}

func (c *GitClient) getBranchHashRef(branch, remote string) (hashRef *plumbing.Reference, err error) {
	refs, err := c.GetRemoteRefs(remote)
	if err != nil {
		return nil, err
	}
	var branchRef *GitRef
	for _, r := range refs {
		if r.Name == branch {
			branchRef = &r
			break
		}
	}
	if branchRef == nil {
		return nil, ErrNoMatchedRemoteBranch
	}
	branchHashRef := plumbing.NewHashReference(plumbing.NewBranchReferenceName(branch), plumbing.NewHash(branchRef.Hash))
	return branchHashRef, nil
}

func (c *GitClient) isRemoteChanged() (ok bool, err error) {
	b, err := c.GetCurrentBranchRef()
	if err != nil {
		return false, err
	}
	refs, err := c.GetRemoteRefs(GitRemoteNameOrigin)
	if err != nil {
		return false, err
	}
	for _, r := range refs {
		if r.Name == b.Name {
			return r.Hash != b.Hash, nil
		}
	}
	return false, nil
}

func NewGitClient(opts ...GitOption) (c *GitClient, err error) {
	// client
	c = &GitClient{
		isMem:          false,
		authType:       GitAuthTypeNone,
		username:       "git",
		privateKeyPath: getDefaultPublicKeyPath(),
		defaultInit:    true,
	}

	// apply options
	for _, opt := range opts {
		opt(c)
	}

	// initialize
	if c.defaultInit {
		if err = c.Init(); err != nil {
			return nil, err
		}
	}

	return
}
