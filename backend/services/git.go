package services

import (
	"bytes"
	"crawlab/model"
	"crawlab/utils"
	"fmt"
	"github.com/apex/log"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"net/url"
	"os"
	"os/exec"
	"path"
	"regexp"
	"runtime/debug"
	"strings"
)

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
		return err
	}

	// 生成 URL
	gitUrl := s.GitUrl
	if s.GitUsername != "" && s.GitPassword != "" {
		u, err := url.Parse(s.GitUrl)
		if err != nil {
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
		return err
	}

	// 获取 repo
	_ = repo.Fetch(&git.FetchOptions{
		RemoteName: "origin",
		Force:      true,
	})

	// 获得 WorkTree
	wt, err := repo.Worktree()
	if err != nil {
		log.Error(err.Error())
		debug.PrintStack()
		return err
	}

	// 拉取 repo
	if err := wt.Pull(&git.PullOptions{
		RemoteName: "origin",
	}); err != nil {
		log.Error(err.Error())
		debug.PrintStack()
		return err
	}

	// 切换分支
	if err := wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(s.GitBranch),
	}); err != nil {
		log.Error(err.Error())
		debug.PrintStack()
		return err
	}

	// 同步到GridFS
	if err := UploadSpiderToGridFsFromMaster(s); err != nil {
		return err
	}

	return nil
}
