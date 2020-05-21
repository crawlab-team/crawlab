package services

import (
	"bytes"
	"crawlab/lib/cron"
	"crawlab/model"
	"crawlab/services/spider_handler"
	"crawlab/utils"
	"fmt"
	"github.com/apex/log"
	"github.com/globalsign/mgo/bson"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"path"
	"regexp"
	"runtime/debug"
	"strings"
	"time"
)

var GitCron *GitCronScheduler

type GitCronScheduler struct {
	cron *cron.Cron
}

type GitBranch struct {
	Hash  string `json:"hash"`
	Name  string `json:"name"`
	Label string `json:"label"`
}

type GitTag struct {
	Hash  string `json:"hash"`
	Name  string `json:"name"`
	Label string `json:"label"`
}

type GitCommit struct {
	Hash           string      `json:"hash"`
	TreeHash       string      `json:"tree_hash"`
	Author         string      `json:"author"`
	Email          string      `json:"email"`
	Message        string      `json:"message"`
	IsHead         bool        `json:"is_head"`
	Ts             time.Time   `json:"ts"`
	Branches       []GitBranch `json:"branches"`
	RemoteBranches []GitBranch `json:"remote_branches"`
	Tags           []GitTag    `json:"tags"`
}

func (g *GitCronScheduler) Start() error {
	c := cron.New(cron.WithSeconds())

	// 启动cron服务
	g.cron.Start()

	// 更新任务列表
	if err := g.Update(); err != nil {
		log.Errorf("update scheduler error: %s", err.Error())
		debug.PrintStack()
		return err
	}

	// 每30秒更新一次任务列表
	spec := "*/30 * * * * *"
	if _, err := c.AddFunc(spec, UpdateGitCron); err != nil {
		log.Errorf("add func update schedulers error: %s", err.Error())
		debug.PrintStack()
		return err
	}

	return nil
}

func (g *GitCronScheduler) RemoveAll() {
	entries := g.cron.Entries()
	for i := 0; i < len(entries); i++ {
		g.cron.Remove(entries[i].ID)
	}
}

func (g *GitCronScheduler) Update() error {
	// 删除所有定时任务
	g.RemoveAll()

	// 获取开启 Git 自动同步的爬虫
	spiders, err := model.GetSpiderAllList(bson.M{"git_auto_sync": true})
	if err != nil {
		log.Errorf("get spider list error: %s", err.Error())
		debug.PrintStack()
		return err
	}

	// 遍历任务列表
	for _, s := range spiders {
		// 添加到定时任务
		if err := g.AddJob(s); err != nil {
			log.Errorf("add job error: %s, job: %s, cron: %s", err.Error(), s.Name, s.GitSyncFrequency)
			debug.PrintStack()
			return err
		}
	}

	return nil
}

func (g *GitCronScheduler) AddJob(s model.Spider) error {
	spec := s.GitSyncFrequency

	// 添加定时任务
	_, err := g.cron.AddFunc(spec, AddGitCronJob(s))
	if err != nil {
		log.Errorf("add func task error: %s", err.Error())
		debug.PrintStack()
		return err
	}

	return nil
}

// 保存爬虫Git同步错误
func SaveSpiderGitSyncError(s model.Spider, errMsg string) {
	s, _ = model.GetSpider(s.Id)
	s.GitSyncError = errMsg
	if err := s.Save(); err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return
	}
}

// 获得Git分支
func GetGitRemoteBranchesPlain(url string) (branches []string, err error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("git", "ls-remote", url)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return branches, err
	}

	for _, line := range strings.Split(stdout.String(), "\n") {
		regex := regexp.MustCompile("refs/heads/(.*)$")
		res := regex.FindStringSubmatch(line)
		if len(res) > 1 {
			branches = append(branches, res[1])
		}
	}

	return branches, nil
}

// 重置爬虫Git
func ResetSpiderGit(s model.Spider) (err error) {
	// 删除文件夹
	if err := os.RemoveAll(s.Src); err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return err
	}

	// 创建空文件夹
	if err := os.MkdirAll(s.Src, os.ModePerm); err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return err
	}

	// 同步到GridFS
	if err := UploadSpiderToGridFsFromMaster(s); err != nil {
		return err
	}

	return nil
}

// 同步爬虫Git
func SyncSpiderGit(s model.Spider) (err error) {
	// 如果 .git 不存在，初始化一个仓库
	if !utils.Exists(path.Join(s.Src, ".git")) {
		_, err := git.PlainInit(s.Src, false)
		if err != nil {
			log.Error(err.Error())
			debug.PrintStack()
			return err
		}
	}

	// 打开 repo
	repo, err := git.PlainOpen(s.Src)
	if err != nil {
		log.Error(err.Error())
		debug.PrintStack()
		SaveSpiderGitSyncError(s, err.Error())
		return err
	}

	// 生成 URL
	gitUrl := s.GitUrl
	if s.GitUsername != "" && s.GitPassword != "" {
		u, err := url.Parse(s.GitUrl)
		if err != nil {
			SaveSpiderGitSyncError(s, err.Error())
			return err
		}
		gitUrl = fmt.Sprintf(
			"%s://%s:%s@%s%s",
			u.Scheme,
			s.GitUsername,
			s.GitPassword,
			u.Hostname(),
			u.Path,
		)
	}

	// 创建 remote
	_ = repo.DeleteRemote("origin")
	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{gitUrl},
	})
	if err != nil {
		log.Error(err.Error())
		debug.PrintStack()
		SaveSpiderGitSyncError(s, err.Error())
		return err
	}

	// 生成验证信息
	var auth ssh.AuthMethod
	if !strings.HasPrefix(s.GitUrl, "http") {
		// 为 SSH
		regex := regexp.MustCompile("^(?:ssh://?)?([0-9a-zA-Z_]+)@")
		res := regex.FindStringSubmatch(s.GitUrl)
		username := s.GitUsername
		if username == "" {
			if len(res) > 1 {
				username = res[1]
			} else {
				username = "git"
			}
		}
		auth, err = ssh.NewPublicKeysFromFile(username, path.Join(os.Getenv("HOME"), ".ssh", "id_rsa"), "")
		if err != nil {
			log.Error(err.Error())
			debug.PrintStack()
			SaveSpiderGitSyncError(s, err.Error())
			return err
		}
	}

	// 获取 repo
	_ = repo.Fetch(&git.FetchOptions{
		RemoteName: "origin",
		Force:      true,
		Auth:       auth,
		Tags:       git.AllTags,
	})

	// 获得 WorkTree
	wt, err := repo.Worktree()
	if err != nil {
		log.Error(err.Error())
		debug.PrintStack()
		SaveSpiderGitSyncError(s, err.Error())
		return err
	}

	// 拉取 repo
	if err := wt.Pull(&git.PullOptions{
		RemoteName:    "origin",
		Auth:          auth,
		ReferenceName: plumbing.HEAD,
		SingleBranch:  false,
	}); err != nil {
		if err.Error() == "already up-to-date" {
			// 检查是否为 Scrapy
			sync := spider_handler.SpiderSync{Spider: s}
			sync.CheckIsScrapy()

			// 同步到GridFS
			if err := UploadSpiderToGridFsFromMaster(s); err != nil {
				SaveSpiderGitSyncError(s, err.Error())
				return err
			}

			// 如果没有错误，则保存空字符串
			SaveSpiderGitSyncError(s, "")

			return nil
		}
		log.Error(err.Error())
		debug.PrintStack()
		SaveSpiderGitSyncError(s, err.Error())
		return err
	}

	// 切换分支
	if err := wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(s.GitBranch),
	}); err != nil {
		log.Error(err.Error())
		debug.PrintStack()
		SaveSpiderGitSyncError(s, err.Error())
		return err
	}

	// 同步到GridFS
	if err := UploadSpiderToGridFsFromMaster(s); err != nil {
		SaveSpiderGitSyncError(s, err.Error())
		return err
	}

	// 获取更新后的爬虫
	s, err = model.GetSpider(s.Id)
	if err != nil {
		SaveSpiderGitSyncError(s, err.Error())
		return err
	}

	// 检查是否为 Scrapy
	sync := spider_handler.SpiderSync{Spider: s}
	sync.CheckIsScrapy()

	// 如果没有错误，则保存空字符串
	SaveSpiderGitSyncError(s, "")

	return nil
}

// 添加Git定时任务
func AddGitCronJob(s model.Spider) func() {
	return func() {
		if err := SyncSpiderGit(s); err != nil {
			log.Errorf(err.Error())
			debug.PrintStack()
			return
		}
	}
}

// 更新Git定时任务
func UpdateGitCron() {
	if err := GitCron.Update(); err != nil {
		log.Errorf(err.Error())
		return
	}
}

// 获取SSH公钥
func GetGitSshPublicKey() string {
	if !utils.Exists(path.Join(os.Getenv("HOME"), ".ssh")) ||
		!utils.Exists(path.Join(os.Getenv("HOME"), ".ssh", "id_rsa")) ||
		!utils.Exists(path.Join(os.Getenv("HOME"), ".ssh", "id_rsa.pub")) {
		log.Errorf("no ssh public key")
		debug.PrintStack()
		return ""
	}
	content, err := ioutil.ReadFile(path.Join(os.Getenv("HOME"), ".ssh", "id_rsa.pub"))
	if err != nil {
		return ""
	}
	return string(content)
}

// 获取Git分支
func GetGitBranches(s model.Spider) (branches []GitBranch, err error) {
	// 打开 repo
	repo, err := git.PlainOpen(s.Src)
	if err != nil {
		log.Error(err.Error())
		debug.PrintStack()
		return branches, err
	}

	iter, err := repo.Branches()
	if iter == nil {
		return branches, nil
	}
	if err := iter.ForEach(func(reference *plumbing.Reference) error {
		branches = append(branches, GitBranch{
			Hash:  reference.Hash().String(),
			Name:  reference.Name().String(),
			Label: reference.Name().Short(),
		})
		return nil
	}); err != nil {
		return branches, err
	}

	return branches, nil
}

// 获取Git Tags
func GetGitTags(s model.Spider) (tags []GitTag, err error) {
	// 打开 repo
	repo, err := git.PlainOpen(s.Src)
	if err != nil {
		log.Error(err.Error())
		debug.PrintStack()
		return tags, err
	}

	iter, err := repo.Tags()
	if iter == nil {
		return tags, nil
	}
	if err := iter.ForEach(func(reference *plumbing.Reference) error {
		tags = append(tags, GitTag{
			Hash:  reference.Hash().String(),
			Name:  reference.Name().String(),
			Label: reference.Name().Short(),
		})
		return nil
	}); err != nil {
		return tags, err
	}

	return tags, nil
}

// 获取Git Head Hash
func GetGitHeadHash(repo *git.Repository) string {
	head, _ := repo.Head()
	return head.Hash().String()
}

// 获取Git远端分支
func GetGitRemoteBranches(s model.Spider) (branches []GitBranch, err error) {
	// 打开 repo
	repo, err := git.PlainOpen(s.Src)
	if err != nil {
		log.Error(err.Error())
		debug.PrintStack()
		return branches, err
	}

	iter, err := repo.References()
	if err != nil {
		log.Error(err.Error())
		debug.PrintStack()
		return branches, err
	}
	if err := iter.ForEach(func(reference *plumbing.Reference) error {
		if reference.Name().IsRemote() {
			log.Infof(reference.Hash().String())
			log.Infof(reference.Name().String())
			branches = append(branches, GitBranch{
				Hash:  reference.Hash().String(),
				Name:  reference.Name().String(),
				Label: reference.Name().Short(),
			})
		}
		return nil
	}); err != nil {
		log.Error(err.Error())
		debug.PrintStack()
		return branches, err
	}
	return branches, err
}

// 获取Git Commits
func GetGitCommits(s model.Spider) (commits []GitCommit, err error) {
	// 打开 repo
	repo, err := git.PlainOpen(s.Src)
	if err != nil {
		log.Error(err.Error())
		debug.PrintStack()
		return commits, err
	}

	// 获取分支列表
	branches, err := GetGitBranches(s)
	branchesDict := map[string][]GitBranch{}
	for _, b := range branches {
		branchesDict[b.Hash] = append(branchesDict[b.Hash], b)
	}

	// 获取分支列表
	remoteBranches, err := GetGitRemoteBranches(s)
	remoteBranchesDict := map[string][]GitBranch{}
	for _, b := range remoteBranches {
		remoteBranchesDict[b.Hash] = append(remoteBranchesDict[b.Hash], b)
	}

	// 获取标签列表
	tags, err := GetGitTags(s)
	tagsDict := map[string][]GitTag{}
	for _, t := range tags {
		tagsDict[t.Hash] = append(tagsDict[t.Hash], t)
	}

	// 获取日志遍历器
	iter, err := repo.Log(&git.LogOptions{
		All: true,
	})
	if err != nil {
		log.Error(err.Error())
		debug.PrintStack()
		return commits, err
	}

	// 遍历日志
	if err := iter.ForEach(func(commit *object.Commit) error {
		gc := GitCommit{
			Hash:           commit.Hash.String(),
			TreeHash:       commit.TreeHash.String(),
			Message:        commit.Message,
			Author:         commit.Author.Name,
			Email:          commit.Author.Email,
			Ts:             commit.Author.When,
			IsHead:         commit.Hash.String() == GetGitHeadHash(repo),
			Branches:       branchesDict[commit.Hash.String()],
			RemoteBranches: remoteBranchesDict[commit.Hash.String()],
			Tags:           tagsDict[commit.Hash.String()],
		}
		commits = append(commits, gc)
		return nil
	}); err != nil {
		log.Error(err.Error())
		debug.PrintStack()
		return commits, err
	}

	return commits, nil
}

func GitCheckout(s model.Spider, hash string) (err error) {
	// 打开 repo
	repo, err := git.PlainOpen(s.Src)
	if err != nil {
		log.Error(err.Error())
		debug.PrintStack()
		return err
	}

	// 获取worktree
	wt, err := repo.Worktree()
	if err != nil {
		log.Error(err.Error())
		debug.PrintStack()
		return err
	}
	//判断远程origin路径是否和当前的GitUrl是同一个，如果不是删掉原来的路径，重新拉取远程代码
	remote, err := repo.Remote("origin")
	if err != nil {
		log.Error(err.Error())
		debug.PrintStack()
		return err
	}
	if remote.String() != s.GitUrl {
		utils.RemoveFiles(s.Src)
		return SyncSpiderGit(s)
	}

	// Checkout
	if err := wt.Checkout(&git.CheckoutOptions{
		Hash:   plumbing.NewHash(hash),
		Create: false,
		Force:  true,
		Keep:   false,
	}); err != nil {
		log.Error(err.Error())
		debug.PrintStack()
		return err
	}

	return nil
}
