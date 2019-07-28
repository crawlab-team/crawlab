package services

import (
	"crawlab/constants"
	"crawlab/database"
	"crawlab/lib/cron"
	"crawlab/model"
	"crawlab/utils"
	"encoding/json"
	"github.com/apex/log"
	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
)

type SpiderFileData struct {
	FileName string
	File     []byte
}

type SpiderUploadMessage struct {
	FileId   string
	FileName string
}

// 从项目目录中获取爬虫列表
func GetSpidersFromDir() ([]model.Spider, error) {
	// 爬虫项目目录路径
	srcPath := viper.GetString("spider.path")

	// 如果爬虫项目目录不存在，则创建一个
	if !utils.Exists(srcPath) {
		if err := os.MkdirAll(srcPath, 0666); err != nil {
			debug.PrintStack()
			return []model.Spider{}, err
		}
	}

	// 获取爬虫项目目录下的所有子项
	items, err := ioutil.ReadDir(srcPath)
	if err != nil {
		debug.PrintStack()
		return []model.Spider{}, err
	}

	// 定义爬虫列表
	spiders := make([]model.Spider, 0)

	// 遍历所有子项
	for _, item := range items {
		// 忽略不为目录的子项
		if !item.IsDir() {
			continue
		}

		// 忽略隐藏目录
		if strings.HasPrefix(item.Name(), ".") {
			continue
		}

		// 构造爬虫
		spider := model.Spider{
			Name:        item.Name(),
			DisplayName: item.Name(),
			Type:        constants.Customized,
			Src:         filepath.Join(srcPath, item.Name()),
			FileId:      bson.ObjectIdHex(constants.ObjectIdNull),
		}

		// 将爬虫加入列表
		spiders = append(spiders, spider)
	}

	return spiders, nil
}

// 将爬虫保存到数据库
func SaveSpiders(spiders []model.Spider) error {
	// 遍历爬虫列表
	for _, spider := range spiders {
		// 忽略非自定义爬虫
		if spider.Type != constants.Customized {
			continue
		}

		// 如果该爬虫不存在于数据库，则保存爬虫到数据库
		s, c := database.GetCol("spiders")
		defer s.Close()
		var spider_ *model.Spider
		if err := c.Find(bson.M{"src": spider.Src}).One(&spider_); err != nil {
			// 不存在
			if err := spider.Add(); err != nil {
				debug.PrintStack()
				return err
			}
		} else {
			// 存在
		}
	}

	return nil
}

// 更新爬虫
func UpdateSpiders() {
	// 从项目目录获取爬虫列表
	spiders, err := GetSpidersFromDir()
	if err != nil {
		log.Errorf(err.Error())
		return
	}

	// 储存爬虫
	if err := SaveSpiders(spiders); err != nil {
		log.Errorf(err.Error())
		return
	}
}

// 打包爬虫目录为zip文件
func ZipSpider(spider model.Spider) (filePath string, err error) {
	// 如果源文件夹不存在，抛错
	if !utils.Exists(spider.Src) {
		debug.PrintStack()
		return "", errors.New("source path does not exist")
	}

	// 临时文件路径
	randomId := uuid.NewV4()
	if err != nil {
		debug.PrintStack()
		return "", err
	}
	filePath = filepath.Join(
		viper.GetString("other.tmppath"),
		randomId.String()+".zip",
	)

	// 将源文件夹打包为zip文件
	d, err := os.Open(spider.Src)
	if err != nil {
		debug.PrintStack()
		return filePath, err
	}
	var files []*os.File
	files = append(files, d)
	if err := utils.Compress(files, filePath); err != nil {
		return filePath, err
	}

	return filePath, nil
}

// 上传zip文件到GridFS
func UploadToGridFs(spider model.Spider, fileName string, filePath string) (fid bson.ObjectId, err error) {
	fid = ""

	// 读取zip文件
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0777)
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		debug.PrintStack()
		return
	}
	if err = file.Close(); err != nil {
		debug.PrintStack()
		return
	}

	// 删除zip文件
	if err = os.Remove(filePath); err != nil {
		debug.PrintStack()
		return
	}

	// 获取MongoDB GridFS连接
	s, gf := database.GetGridFs("files")
	defer s.Close()

	// 如果存在FileId删除GridFS上的老文件
	if !utils.IsObjectIdNull(spider.FileId) {
		if err = gf.RemoveId(spider.FileId); err != nil {
			debug.PrintStack()
		}
	}

	// 创建一个新GridFS文件
	f, err := gf.Create(fileName)
	if err != nil {
		debug.PrintStack()
		return
	}

	// 将文件写入到GridFS
	if _, err = f.Write(fileBytes); err != nil {
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
			log.Errorf(err.Error())
			return err
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
	// 将源文件夹打包为zip文件
	filePath, err := ZipSpider(spider)
	if err != nil {
		return err
	}

	// 上传zip文件到GridFS
	fileName := filepath.Base(spider.Src) + ".zip"
	fid, err := UploadToGridFs(spider, fileName, filePath)
	if err != nil {
		return err
	}

	// 保存FileId
	spider.FileId = fid
	if err := spider.Save(); err != nil {
		return err
	}

	// 发布消息给工作节点
	msg := SpiderUploadMessage{
		FileId:   fid.Hex(),
		FileName: fileName,
	}
	msgStr, err := json.Marshal(msg)
	if err != nil {
		return
	}
	channel := "files:upload"
	if err = database.Publish(channel, string(msgStr)); err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return
	}

	return
}

// 上传爬虫回调
func OnFileUpload(channel string, msgStr string) {
	s, gf := database.GetGridFs("files")
	defer s.Close()

	// 反序列化消息
	var msg SpiderUploadMessage
	if err := json.Unmarshal([]byte(msgStr), &msg); err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return
	}

	// 从GridFS获取该文件
	f, err := gf.OpenId(bson.ObjectIdHex(msg.FileId))
	if err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return
	}
	defer f.Close()

	// 生成唯一ID
	randomId := uuid.NewV4()

	// 创建临时文件
	tmpFilePath := filepath.Join(viper.GetString("other.tmppath"), randomId.String()+".zip")
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
		//strings.Replace(msg.FileName, ".zip", "", -1),
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

// 启动爬虫服务
func InitSpiderService() error {
	// 构造定时任务执行器
	c := cron.New(cron.WithSeconds())

	if IsMaster() {
		// 主节点

		// 每5秒更新一次爬虫信息
		if _, err := c.AddFunc("*/5 * * * * *", UpdateSpiders); err != nil {
			return err
		}

		// 每60秒同步爬虫给工作节点
		if _, err := c.AddFunc("0 * * * * *", PublishAllSpidersJob); err != nil {
			return err
		}
	} else {
		// 非主节点

		// 订阅文件上传
		channel := "files:upload"
		var sub database.Subscriber
		sub.Connect()
		sub.Subscribe(channel, OnFileUpload)
	}

	// 启动定时任务
	c.Start()

	return nil
}
