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
}

func (s *SpiderSync) CreateMd5File(md5 string, spiderName string) {
	path := filepath.Join(viper.GetString("spider.path"), spiderName)
	utils.CreateFilePath(path)

	fileName := filepath.Join(path, Md5File)
	file := utils.OpenFile(fileName)
	defer file.Close()
	if file != nil {
		if _, err := file.WriteString(md5); err != nil {
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
func (s *SpiderSync) RemoveSpiderFile(spiderName string) {
	//爬虫文件有变化，先删除本地文件
	_ = os.Remove(filepath.Join(
		viper.GetString("spider.path"),
		spiderName,
	))
}

// 检测是否已经下载中
func (s *SpiderSync) CheckDownLoading(spiderId string, fileId string) (bool, string) {
	key := s.GetLockDownloadKey(spiderId)
	if _, err := database.RedisClient.HGet("spider", key); err == nil {
		log.Infof("downloading spider file, spider_id: %s, file_id:%s", spiderId, fileId)
		return true, key
	}
	return false, key
}

// 下载爬虫
func (s *SpiderSync) Download(spiderId string, fileId string) {

	session, gf := database.GetGridFs("files")
	defer session.Close()

	f, err := gf.OpenId(bson.ObjectIdHex(fileId))
	defer f.Close()
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
	tmpFile, err := os.OpenFile(tmpFilePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return
	}
	defer tmpFile.Close()

	// 将该文件写入临时文件
	if _, err := io.Copy(tmpFile, f); err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return
	}

	// 解压缩临时文件到目标文件夹
	dstPath := filepath.Join(
		viper.GetString("spider.path"),
	)
	if err := utils.DeCompress(tmpFile, dstPath); err != nil {
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
}
