package services

import (
	"crawlab/constants"
	"crawlab/database"
	"crawlab/lib/cron"
	"crawlab/model"
	"crawlab/services/spider_handler"
	"crawlab/utils"
	"fmt"
	"github.com/apex/log"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"runtime/debug"
)

type SpiderFileData struct {
	FileName string
	File     []byte
}

type SpiderUploadMessage struct {
	FileId   string
	FileName string
	SpiderId string
}

// 上传zip文件到GridFS
func UploadToGridFs(fileName string, filePath string) (fid bson.ObjectId, md5 string, err error) {
	fid = ""

	// 获取MongoDB GridFS连接
	s, gf := database.GetGridFs("files")
	defer s.Close()

	// 创建一个新GridFS文件
	f, err := gf.Create(fileName)
	if err != nil {
		debug.PrintStack()
		return
	}

	//分片读取爬虫zip文件
	err = ReadFileByStep(filePath, WriteToGridFS, f)
	if err != nil {
		debug.PrintStack()
		return "", "", err
	}

	// 删除zip文件
	if err = os.Remove(filePath); err != nil {
		debug.PrintStack()
		return
	}
	// 关闭文件，提交写入
	if err = f.Close(); err != nil {
		return "", "", err
	}
	// 文件ID
	fid = f.Id().(bson.ObjectId)

	return fid, f.MD5(), nil
}

func WriteToGridFS(content []byte, f *mgo.GridFile) {
	if _, err := f.Write(content); err != nil {
		debug.PrintStack()
		return
	}
}

//分片读取大文件
func ReadFileByStep(filePath string, handle func([]byte, *mgo.GridFile), fileCreate *mgo.GridFile) error {
	f, err := os.OpenFile(filePath, os.O_RDONLY, 0777)
	if err != nil {
		log.Infof("can't opened this file")
		return err
	}
	defer f.Close()
	s := make([]byte, 4096)
	for {
		switch nr, err := f.Read(s[:]); true {
		case nr < 0:
			_, _ = fmt.Fprintf(os.Stderr, "cat: error reading: %s\n", err.Error())
			debug.PrintStack()
		case nr == 0: // EOF
			return nil
		case nr > 0:
			handle(s[0:nr], fileCreate)
		}
	}
}

// 发布所有爬虫
func PublishAllSpiders() {
	// 获取爬虫列表
	spiders, _ := model.GetSpiderList(nil, 0, constants.Infinite)
	if len(spiders) == 0 {
		return
	}
	log.Infof("start sync spider to local, total: %d", len(spiders))
	// 遍历爬虫列表
	for _, spider := range spiders {
		// 异步发布爬虫
		go func() {
			PublishSpider(spider)
		}()
	}
}

// 发布爬虫
func PublishSpider(spider model.Spider) {
	// 查询gf file，不存在则删除
	gfFile := model.GetGridFs(spider.FileId)
	if gfFile == nil {
		_ = model.RemoveSpider(spider.FileId)
		return
	}
	spiderSync := spider_handler.SpiderSync{}
	defer spiderSync.CreateMd5File(gfFile.Md5, spider.Name)

	//目录不存在，则直接下载
	path := filepath.Join(viper.GetString("spider.path"), spider.Name)
	if !utils.Exists(path) {
		log.Infof("path not found: %s", path)
		spiderSync.Download(spider.Id.Hex(), spider.FileId.Hex())
		return
	}
	// md5文件不存在，则下载
	md5 := filepath.Join(path, spider_handler.Md5File)
	if !utils.Exists(md5) {
		log.Infof("md5.txt file not found: %s", md5)
		spiderSync.RemoveSpiderFile(spider.Name)
		spiderSync.Download(spider.Id.Hex(), spider.FileId.Hex())
		return
	}
	// md5值不一样，则下载
	md5Str := utils.ReadFile(md5)
	if spider.Md5 != md5Str {
		log.Infof("md5 is different: %s:%s ", md5Str, md5)
		spiderSync.RemoveSpiderFile(spider.Name)
		spiderSync.Download(spider.Id.Hex(), spider.FileId.Hex())
		return
	}
}

// 启动爬虫服务
func InitSpiderService() error {
	// 构造定时任务执行器
	c := cron.New(cron.WithSeconds())
	if _, err := c.AddFunc("0/15 * * * * *", PublishAllSpiders); err != nil {
		return err
	}
	// 启动定时任务
	c.Start()

	return nil
}
