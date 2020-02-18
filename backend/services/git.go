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
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"path"
	"regexp"
	"runtime/debug"
	"strings"
)

var GitCron *GitCronScheduler

type GitCronScheduler struct {
	cron *cron.Cron
}

func SaveSpiderGitSyncError(s model.Spider, errMsg string) {
	s.GitSyncError = errMsg
	if err := s.Save(); err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return
	}
}

func GetGitBranches(url string) (branches []string, err error) {
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
		RemoteName: "origin",
		Auth:       auth,
	}); err != nil {
		if err.Error() == "already up-to-date" {
			// 检查是否为 Scrapy
			sync := spider_handler.SpiderSync{Spider: s}
			sync.CheckIsScrapy()
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

	// 检查是否为 Scrapy
	sync := spider_handler.SpiderSync{Spider: s}
	sync.CheckIsScrapy()

	// 如果没有错误，则保存空字符串
	SaveSpiderGitSyncError(s, "")

	return nil
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

func AddGitCronJob(s model.Spider) func() {
	return func() {
		if err := SyncSpiderGit(s); err != nil {
			log.Errorf(err.Error())
			debug.PrintStack()
			return
		}
	}
}

func UpdateGitCron() {
	if err := GitCron.Update(); err != nil {
		log.Errorf(err.Error())
		return
	}
}

func GetGitSshPublicKey() string {
	if !utils.Exists(path.Join(os.Getenv("HOME"), ".ssh")) ||
		!utils.Exists(path.Join(os.Getenv("HOME"), ".ssh", "id_rsa")) ||
		!utils.Exists(path.Join(os.Getenv("HOME"), ".ssh", "id_rsa.pub")) {
		cmd := exec.Command("ssh-keygen -q -t rsa -N \"\" -f $HOME/.ssh/id_rsa")
		if err := cmd.Start(); err != nil {
			log.Errorf(err.Error())
			debug.PrintStack()
			return ""
		}
	}
	content, err := ioutil.ReadFile(path.Join(os.Getenv("HOME"), ".ssh", "id_rsa.pub"))
	if err != nil {
		return ""
	}
	return string(content)
}
