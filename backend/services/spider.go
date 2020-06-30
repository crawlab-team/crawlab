package services

import (
	"crawlab/constants"
	"crawlab/database"
	"crawlab/entity"
	"crawlab/lib/cron"
	"crawlab/model"
	"crawlab/services/spider_handler"
	"crawlab/utils"
	"errors"
	"fmt"
	"github.com/apex/log"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime/debug"
	"time"
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

// 从主节点上传爬虫到GridFS
func UploadSpiderToGridFsFromMaster(spider model.Spider) error {
	// 爬虫所在目录
	spiderDir := spider.Src

	// 打包为 zip 文件
	files, err := utils.GetFilesFromDir(spiderDir)
	if err != nil {
		return err
	}
	randomId := uuid.NewV4()
	tmpFilePath := filepath.Join(viper.GetString("other.tmppath"), spider.Name+"."+randomId.String()+".zip")
	spiderZipFileName := spider.Name + ".zip"
	if err := utils.Compress(files, tmpFilePath); err != nil {
		return err
	}

	// 获取 GridFS 实例
	s, gf := database.GetGridFs("files")
	defer s.Close()

	// 判断文件是否已经存在
	var gfFile model.GridFs
	if err := gf.Find(bson.M{"filename": spiderZipFileName}).One(&gfFile); err == nil {
		// 已经存在文件，则删除
		log.Errorf(gfFile.Id.Hex() + " already exists. removing...")
		if err := gf.RemoveId(gfFile.Id); err != nil {
			log.Errorf(err.Error())
			debug.PrintStack()
			return err
		}
	}

	// 上传到GridFs
	fid, err := RetryUploadToGridFs(spiderZipFileName, tmpFilePath)
	if err != nil {
		log.Errorf("upload to grid fs error: %s", err.Error())
	}

	// 保存爬虫 FileId
	spider.FileId = fid
	if err := spider.Save(); err != nil {
		return err
	}

	// 获取爬虫同步实例
	spiderSync := spider_handler.SpiderSync{
		Spider: spider,
	}

	// 获取gfFile
	gfFile2 := model.GetGridFs(spider.FileId)

	// 生成MD5
	spiderSync.CreateMd5File(gfFile2.Md5)

	// 检查是否为 Scrapy 爬虫
	spiderSync.CheckIsScrapy()

	return nil
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
		log.Errorf("create file error: " + err.Error())
		debug.PrintStack()
		return
	}

	// 分片读取爬虫zip文件
	err = ReadFileByStep(filePath, WriteToGridFS, f)
	if err != nil {
		log.Errorf("read file by step error: " + err.Error())
		debug.PrintStack()
		return "", err
	}

	// 删除zip文件
	if err = os.Remove(filePath); err != nil {
		log.Errorf("remove file error: " + err.Error())
		debug.PrintStack()
		return
	}

	// 关闭文件，提交写入
	if err = f.Close(); err != nil {
		log.Errorf("close file error: " + err.Error())
		debug.PrintStack()
		return "", err
	}

	// 文件ID
	fid = f.Id().(bson.ObjectId)

	return fid, nil
}

// 带重试功能的上传至 GridFS
func RetryUploadToGridFs(fileName string, filePath string) (fid bson.ObjectId, err error) {
	maxErrCount := 10
	errCount := 0
	for {
		if errCount > maxErrCount {
			break
		}
		fid, err = UploadToGridFs(fileName, filePath)
		if err != nil {
			errCount++
			log.Errorf("upload to grid fs error: %s", err.Error())
			time.Sleep(3 * time.Second)
			continue
		}
		return fid, nil
	}
	return fid, errors.New("unable to upload to gridfs, please re-upload the spider")
}

// 写入grid fs
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
	defer utils.Close(f)
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
	spiders, _, _ := model.GetSpiderList(nil, 0, constants.Infinite, "-_id")
	if len(spiders) == 0 {
		return
	}
	log.Infof("start sync spider to local, total: %d", len(spiders))
	// 遍历爬虫列表
	for _, spider := range spiders {
		// 异步发布爬虫
		go func(s model.Spider) {
			PublishSpider(s)
		}(spider)
	}
}

// 发布爬虫
func PublishSpider(spider model.Spider) {
	var gfFile *model.GridFs
	if spider.FileId.Hex() != constants.ObjectIdNull {
		// 查询gf file，不存在则标记为爬虫文件不存在
		gfFile = model.GetGridFs(spider.FileId)
		if gfFile == nil {
			log.Errorf("get grid fs file error: cannot find grid fs file")
			log.Errorf("grid fs file_id: " + spider.FileId.Hex())
			log.Errorf("spider_name: " + spider.Name)
			debug.PrintStack()
			//spider.FileId = constants.ObjectIdNull
			//if err := spider.Save(); err != nil {
			//	return
			//}
			return
		}
	}

	// 如果FileId为空，表示还没有上传爬虫到GridFS，则跳过
	if spider.FileId == bson.ObjectIdHex(constants.ObjectIdNull) {
		return
	}

	// 获取爬虫同步实例
	spiderSync := spider_handler.SpiderSync{
		Spider: spider,
	}

	// 安装依赖
	go spiderSync.InstallDeps()

	//目录不存在，则直接下载
	path := filepath.Join(viper.GetString("spider.path"), spider.Name)
	if !utils.Exists(path) {
		log.Infof("path not found: %s", path)
		spiderSync.Download()
		spiderSync.CreateMd5File(gfFile.Md5)
		spiderSync.CheckIsScrapy()
		return
	}

	// md5文件不存在，则下载
	md5 := filepath.Join(path, spider_handler.Md5File)
	if !utils.Exists(md5) {
		log.Infof("md5 file not found: %s", md5)
		spiderSync.RemoveDownCreate(gfFile.Md5)
		return
	}

	// md5值不一样，则下载
	md5Str := utils.GetSpiderMd5Str(md5)
	if gfFile.Md5 != md5Str {
		log.Infof("md5 is different, gf-md5:%s, file-md5:%s", gfFile.Md5, md5Str)
		spiderSync.RemoveDownCreate(gfFile.Md5)
		return
	}
}

func RemoveSpider(id string) error {
	// 获取该爬虫
	spider, err := model.GetSpider(bson.ObjectIdHex(id))
	if err != nil {
		return err
	}

	// 删除爬虫文件目录
	path := filepath.Join(viper.GetString("spider.path"), spider.Name)
	utils.RemoveFiles(path)

	// 删除其他节点的爬虫目录
	msg := entity.NodeMessage{
		Type:     constants.MsgTypeRemoveSpider,
		SpiderId: id,
	}
	if err := database.Pub(constants.ChannelAllNode, msg); err != nil {
		return err
	}

	// 从数据库中删除该爬虫
	if err := model.RemoveSpider(bson.ObjectIdHex(id)); err != nil {
		return err
	}

	// 删除日志文件
	if err := RemoveLogBySpiderId(spider.Id); err != nil {
		return err
	}

	// 删除爬虫对应的task任务
	if err := model.RemoveTaskBySpiderId(spider.Id); err != nil {
		return err
	}

	// TODO 定时任务如何处理
	return nil
}

func CancelSpider(id string) error {
	// 获取该爬虫
	spider, err := model.GetSpider(bson.ObjectIdHex(id))
	if err != nil {
		return err
	}

	// 获取该爬虫待定或运行中的任务列表
	query := bson.M{
		"spider_id": spider.Id,
		"status": bson.M{
			"$in": []string{
				constants.StatusPending,
				constants.StatusRunning,
			},
		},
	}
	tasks, err := model.GetTaskList(query, 0, constants.Infinite, "-create_ts")
	if err != nil {
		return err
	}

	// 遍历任务列表，依次停止
	for _, task := range tasks {
		if err := CancelTask(task.Id); err != nil {
			return err
		}
	}

	return nil
}

func cloneGridFsFile(spider model.Spider, newName string) (err error) {
	// 构造新爬虫
	newSpider := spider
	newSpider.Id = bson.NewObjectId()
	newSpider.Name = newName
	newSpider.DisplayName = newName
	newSpider.Src = path.Join(path.Dir(spider.Src), newName)
	newSpider.CreateTs = time.Now()
	newSpider.UpdateTs = time.Now()

	// GridFS连接实例
	s, gf := database.GetGridFs("files")
	defer s.Close()

	// 被克隆爬虫的GridFS文件
	f, err := gf.OpenId(spider.FileId)
	if err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return err
	}

	// 新爬虫的GridFS文件
	fNew, err := gf.Create(newSpider.Name + ".zip")
	if err != nil {
		log.Errorf(err.Error())
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
	tmpFile := utils.OpenFile(tmpFilePath)

	// 拷贝到临时文件
	if _, err := io.Copy(tmpFile, f); err != nil {
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

	// 读取内容
	fContent, err := ioutil.ReadFile(tmpFilePath)
	if err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return err
	}

	// 写入GridFS文件
	if _, err := fNew.Write(fContent); err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return err
	}

	// 关闭被克隆爬虫GridFS文件
	if err = f.Close(); err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return err
	}

	// 将新爬虫文件复制
	newSpider.FileId = fNew.Id().(bson.ObjectId)

	// 保存新爬虫
	if err := newSpider.Add(); err != nil {
		return err
	}

	// 关闭新爬虫GridFS文件
	if err := fNew.Close(); err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return err
	}

	// 删除临时文件
	if err := os.RemoveAll(tmpFilePath); err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return err
	}

	// 同步爬虫
	PublishSpider(newSpider)

	return nil
}

func CopySpider(spider model.Spider, newName string) error {
	// 克隆GridFS文件
	if err := cloneGridFsFile(spider, newName); err != nil {
		return err
	}

	return nil
}

func UpdateSpiderDedup(spider model.Spider) error {
	col := utils.GetSpiderCol(spider.Col, spider.Name)

	s, c := database.GetCol(col)
	defer s.Close()

	if !spider.IsDedup {
		_ = c.DropIndex(spider.DedupField)
		//if err := c.DropIndex(spider.DedupField); err != nil {
		//	return err
		//}
		return nil
	}

	if err := c.EnsureIndex(mgo.Index{
		Key:    []string{spider.DedupField},
		Unique: true,
	}); err != nil {
		return err
	}

	return nil
}

func InitDemoSpiders() {
	// 添加Demo爬虫
	templateSpidersDir := "./template/spiders"
	for _, info := range utils.ListDir(templateSpidersDir) {
		if !info.IsDir() {
			continue
		}
		spiderName := info.Name()

		// 如果爬虫在数据库中不存在，则添加
		spider := model.GetSpiderByName(spiderName)
		if spider.Name != "" {
			// 存在同名爬虫，跳过
			continue
		}

		// 拷贝爬虫
		templateSpiderPath := path.Join(templateSpidersDir, spiderName)
		spiderPath := path.Join(viper.GetString("spider.path"), spiderName)
		if utils.Exists(spiderPath) {
			utils.RemoveFiles(spiderPath)
		}
		if err := utils.CopyDir(templateSpiderPath, spiderPath); err != nil {
			log.Errorf("copy error: " + err.Error())
			debug.PrintStack()
			continue
		}

		// 构造配置数据
		configData := entity.ConfigSpiderData{}

		// 读取YAML文件
		yamlFile, err := ioutil.ReadFile(path.Join(spiderPath, "Spiderfile"))
		if err != nil {
			log.Errorf("read yaml error: " + err.Error())
			//debug.PrintStack()
			continue
		}

		// 反序列化
		if err := yaml.Unmarshal(yamlFile, &configData); err != nil {
			log.Errorf("unmarshal error: " + err.Error())
			debug.PrintStack()
			continue
		}

		if configData.Type == constants.Customized {
			// 添加该爬虫到数据库
			spider = model.Spider{
				Id:          bson.NewObjectId(),
				Name:        spiderName,
				DisplayName: configData.DisplayName,
				Type:        constants.Customized,
				Col:         configData.Col,
				Src:         spiderPath,
				Remark:      configData.Remark,
				ProjectId:   bson.ObjectIdHex(constants.ObjectIdNull),
				FileId:      bson.ObjectIdHex(constants.ObjectIdNull),
				Cmd:         configData.Cmd,
				UserId:      bson.ObjectIdHex(constants.ObjectIdNull),
			}
			if err := spider.Add(); err != nil {
				log.Errorf("add spider error: " + err.Error())
				debug.PrintStack()
				continue
			}

			// 上传爬虫到GridFS
			if err := UploadSpiderToGridFsFromMaster(spider); err != nil {
				log.Errorf("upload spider error: " + err.Error())
				debug.PrintStack()
				continue
			}
		} else if configData.Type == constants.Configurable || configData.Type == "config" {
			// 添加该爬虫到数据库
			spider = model.Spider{
				Id:          bson.NewObjectId(),
				Name:        configData.Name,
				DisplayName: configData.DisplayName,
				Type:        constants.Configurable,
				Col:         configData.Col,
				Src:         spiderPath,
				Remark:      configData.Remark,
				ProjectId:   bson.ObjectIdHex(constants.ObjectIdNull),
				FileId:      bson.ObjectIdHex(constants.ObjectIdNull),
				Config:      configData,
				UserId:      bson.ObjectIdHex(constants.ObjectIdNull),
			}
			if err := spider.Add(); err != nil {
				log.Errorf("add spider error: " + err.Error())
				debug.PrintStack()
				continue
			}

			// 根据序列化后的数据处理爬虫文件
			if err := ProcessSpiderFilesFromConfigData(spider, configData); err != nil {
				log.Errorf("add spider error: " + err.Error())
				debug.PrintStack()
				continue
			}
		}
	}

	// 发布所有爬虫
	PublishAllSpiders()
}

// 启动爬虫服务
func InitSpiderService() error {
	// 构造定时任务执行器
	cPub := cron.New(cron.WithSeconds())
	if _, err := cPub.AddFunc("0 * * * * *", PublishAllSpiders); err != nil {
		return err
	}

	// 启动定时任务
	cPub.Start()

	if model.IsMaster() && viper.GetString("setting.demoSpiders") == "Y" {
		// 初始化Demo爬虫
		InitDemoSpiders()
	}

	if model.IsMaster() {
		// 构造 Git 定时任务
		GitCron = &GitCronScheduler{
			cron: cron.New(cron.WithSeconds()),
		}

		// 启动 Git 定时任务
		if err := GitCron.Start(); err != nil {
			return err
		}

		// 清理UserId
		InitSpiderCleanUserIds()
	}

	return nil
}
