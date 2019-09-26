package services

import (
	"context"
	"crawlab/constants"
	"crawlab/database"
	"crawlab/lib/cron"
	"crawlab/model"
	"crawlab/utils"
	"encoding/json"
	"fmt"
	"github.com/apex/log"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/gomodule/redigo/redis"
	"github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"io"
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
func UploadToGridFs(fileName string, filePath string) (fid bson.ObjectId, err error) {
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
		return "", err
	}

	// 删除zip文件
	if err = os.Remove(filePath); err != nil {
		debug.PrintStack()
		return
	}
	// 关闭文件，提交写入
	if err = f.Close(); err != nil {
		return "", err
	}
	// 文件ID
	fid = f.Id().(bson.ObjectId)

	return fid, nil
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
func PublishAllSpiders() error {
	// 获取爬虫列表
	spiders, err := model.GetSpiderList(nil, 0, constants.Infinite)
	if err != nil {
		log.Errorf(err.Error())
		return err
	}

	// 遍历爬虫列表
	for _, spider := range spiders {
		// 发布爬虫
		if err := PublishSpider(spider); err != nil {
			log.Errorf("publish spider error:" + err.Error())
			// return err
		}
	}

	return nil
}

func PublishAllSpidersJob() {
	if err := PublishAllSpiders(); err != nil {
		log.Errorf(err.Error())
	}
}

// 发布爬虫
// 1. 将源文件夹打包为zip文件
// 2. 上传zip文件到GridFS
// 3. 发布消息给工作节点
func PublishSpider(spider model.Spider) (err error) {

	s, gf := database.GetGridFs("files")
	defer s.Close()

	f, err := gf.OpenId(spider.FileId)
	defer f.Close()
	if err != nil {
		log.Errorf("open file id: " + spider.FileId.Hex() + ", spider id:" + spider.Id.Hex() + ", error: " + err.Error())
		debug.PrintStack()
		// 爬虫和文件没有对应，则删除爬虫
		_ = model.RemoveSpider(spider.Id)
		return err
	}

	// 发布消息给工作节点
	msg := SpiderUploadMessage{
		FileId:   spider.FileId.Hex(),
		FileName: f.Name(),
		SpiderId: spider.Id.Hex(),
	}
	msgStr, err := json.Marshal(msg)
	if err != nil {
		return
	}
	channel := "files:upload"
	if _, err = database.RedisClient.Publish(channel, utils.BytesToString(msgStr)); err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return
	}

	return
}

// 上传爬虫回调
func OnFileUpload(message redis.Message) (err error) {
	s, gf := database.GetGridFs("files")
	defer s.Close()

	// 反序列化消息
	var msg SpiderUploadMessage
	if err := json.Unmarshal(message.Data, &msg); err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return err
	}

	// 从GridFS获取该文件
	f, err := gf.OpenId(bson.ObjectIdHex(msg.FileId))
	defer f.Close()
	if err != nil {
		log.Errorf("open file id: " + msg.FileId + ", spider id:" + msg.SpiderId + ", error: " + err.Error())
		debug.PrintStack()
		return err
	}

	// 生成唯一ID
	randomId := uuid.NewV4()
	tmpPath := viper.GetString("other.tmppath")
	if !utils.Exists(tmpPath) {
		if err := os.MkdirAll(tmpPath, 0777); err != nil {
			log.Errorf("mkdir other.tmppath error: %v", err.Error())
			return err
		}
	}
	// 创建临时文件
	tmpFilePath := filepath.Join(tmpPath, randomId.String()+".zip")

	tmpFile, err := os.OpenFile(tmpFilePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return err
	}
	defer tmpFile.Close()

	// 将该文件写入临时文件
	if _, err := io.Copy(tmpFile, f); err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return err
	}

	// 解压缩临时文件到目标文件夹
	dstPath := filepath.Join(
		viper.GetString("spider.path"),
	)
	if err := utils.DeCompress(tmpFile, dstPath); err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return err
	}

	// 关闭临时文件
	if err := tmpFile.Close(); err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return err
	}

	// 删除临时文件
	if err := os.Remove(tmpFilePath); err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return err
	}
	return nil
}

// 启动爬虫服务
func InitSpiderService() error {
	// 构造定时任务执行器
	c := cron.New(cron.WithSeconds())

	if IsMaster() {
		// 主节点

		// 每60秒同步爬虫给工作节点
		if _, err := c.AddFunc("0 * * * * *", PublishAllSpidersJob); err != nil {
			return err
		}
	} else {
		// 非主节点

		// 订阅文件上传
		channel := "files:upload"
		ctx := context.Background()
		return database.RedisClient.Subscribe(ctx, OnFileUpload, channel)

	}

	// 启动定时任务
	c.Start()

	return nil
}
