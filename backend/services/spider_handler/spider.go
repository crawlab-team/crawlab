package spider_handler

import (
	"crawlab/database"
	"crawlab/model"
	"crawlab/utils"
	"github.com/apex/log"
	"github.com/globalsign/mgo/bson"
	"github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"io"
	"os"
	"path/filepath"
	"runtime/debug"
)

const (
	Md5File = "md5.txt"
)

type SpiderSync struct {
	Spider model.Spider
}

func (s *SpiderSync) CreateMd5File(md5 string) {
	path := filepath.Join(viper.GetString("spider.path"), s.Spider.Name)
	utils.CreateFilePath(path)

	fileName := filepath.Join(path, Md5File)
	file := utils.OpenFile(fileName)
	defer utils.Close(file)
	if file != nil {
		if _, err := file.WriteString(md5 + "\n"); err != nil {
			log.Errorf("file write string error: %s", err.Error())
			debug.PrintStack()
		}
	}
}

// 获得下载锁的key
func (s *SpiderSync) GetLockDownloadKey(spiderId string) string {
	node, _ := model.GetCurrentNode()
	return node.Id.Hex() + "#" + spiderId
}

// 删除本地文件
func (s *SpiderSync) RemoveSpiderFile() {
	path := filepath.Join(
		viper.GetString("spider.path"),
		s.Spider.Name,
	)
	//爬虫文件有变化，先删除本地文件
	if err := os.RemoveAll(path); err != nil {
		log.Errorf("remove spider files error: %s, path: %s", err.Error(), path)
		debug.PrintStack()
	}
}

// 检测是否已经下载中
func (s *SpiderSync) CheckDownLoading(spiderId string, fileId string) (bool, string) {
	key := s.GetLockDownloadKey(spiderId)
	if _, err := database.RedisClient.HGet("spider", key); err == nil {
		return true, key
	}
	return false, key
}

// 下载爬虫
func (s *SpiderSync) Download() {
	spiderId := s.Spider.Id.Hex()
	fileId := s.Spider.FileId.Hex()
	isDownloading, key := s.CheckDownLoading(spiderId, fileId)
	if isDownloading {
		return
	} else {
		_ = database.RedisClient.HSet("spider", key, key)
	}

	session, gf := database.GetGridFs("files")
	defer session.Close()

	f, err := gf.OpenId(bson.ObjectIdHex(fileId))
	defer utils.Close(f)
	if err != nil {
		log.Errorf("open file id: " + fileId + ", spider id:" + spiderId + ", error: " + err.Error())
		debug.PrintStack()
		return
	}

	// 生成唯一ID
	randomId := uuid.NewV4()
	tmpPath := viper.GetString("other.tmppath")
	if !utils.Exists(tmpPath) {
		if err := os.MkdirAll(tmpPath, 0777); err != nil {
			log.Errorf("mkdir other.tmppath error: %v", err.Error())
			return
		}
	}
	// 创建临时文件
	tmpFilePath := filepath.Join(tmpPath, randomId.String()+".zip")
	tmpFile := utils.OpenFile(tmpFilePath)
	defer utils.Close(tmpFile)

	// 将该文件写入临时文件
	if _, err := io.Copy(tmpFile, f); err != nil {
		log.Errorf("copy file error: %s, file_id: %s", err.Error(), f.Id())
		debug.PrintStack()
		return
	}

	// 解压缩临时文件到目标文件夹
	dstPath := filepath.Join(
		viper.GetString("spider.path"),
		s.Spider.Name,
	)
	if err := utils.DeCompress(tmpFile, dstPath); err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return
	}

	//修改目标文件夹权限
	// 解决scrapy.setting中开启LOG_ENABLED 和 LOG_FILE时不能创建log文件的问题
	if err := os.Chmod(dstPath, 0777); err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return
	}

	// 关闭临时文件
	if err := tmpFile.Close(); err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return
	}

	// 删除临时文件
	if err := os.Remove(tmpFilePath); err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return
	}

	_ = database.RedisClient.HDel("spider", key)
}
