package services

import (
	"crawlab/constants"
	"crawlab/model"
	"crawlab/utils"
	"fmt"
	"github.com/apex/log"
	"github.com/globalsign/mgo/bson"
	"github.com/imroc/req"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"path"
	"path/filepath"
	"runtime/debug"
	"strings"
)

func DownloadRepo(fullName string, userId bson.ObjectId) (err error) {
	// 下载 zip 文件
	url := fmt.Sprintf("%s/%s.zip", viper.GetString("repo.ossUrl"), fullName)
	progress := func(current, total int64) {
		fmt.Println(float32(current)/float32(total)*100, "%")
	}
	res, err := req.Get(url, req.DownloadProgress(progress))
	if err != nil {
		log.Errorf("download repo error: " + err.Error())
		debug.PrintStack()
		return err
	}
	spiderName := strings.Replace(fullName, "/", "_", -1)
	randomId := uuid.NewV4()
	tmpFilePath := filepath.Join(viper.GetString("other.tmppath"), spiderName+"."+randomId.String()+".zip")
	if err := res.ToFile(tmpFilePath); err != nil {
		log.Errorf("to file error: " + err.Error())
		debug.PrintStack()
		return err
	}

	// 解压 zip 文件
	tmpFile := utils.OpenFile(tmpFilePath)
	if err := utils.DeCompress(tmpFile, viper.GetString("other.tmppath")); err != nil {
		log.Errorf("de-compress error: " + err.Error())
		debug.PrintStack()
		return err
	}

	// 拷贝文件
	spiderPath := path.Join(viper.GetString("spider.path"), spiderName)
	srcDirPath := fmt.Sprintf("%s/data/github.com/%s", viper.GetString("other.tmppath"), fullName)
	if err := utils.CopyDir(srcDirPath, spiderPath); err != nil {
		log.Errorf("copy error: " + err.Error())
		debug.PrintStack()
		return err
	}

	// 创建爬虫
	spider := model.GetSpiderByName(spiderName)
	if spider.Name == "" {
		// 新增
		spider = model.Spider{
			Id:          bson.NewObjectId(),
			Name:        spiderName,
			DisplayName: spiderName,
			Type:        constants.Customized,
			Src:         spiderPath,
			ProjectId:   bson.ObjectIdHex(constants.ObjectIdNull),
			FileId:      bson.ObjectIdHex(constants.ObjectIdNull),
			UserId:      userId,
		}
		if err := spider.Add(); err != nil {
			log.Error("add spider error: " + err.Error())
			debug.PrintStack()
			return err
		}
	} else {
		// 更新
		if err := spider.Save(); err != nil {
			log.Error("save spider error: " + err.Error())
			debug.PrintStack()
			return err
		}
	}

	// 上传爬虫
	if err := UploadSpiderToGridFsFromMaster(spider); err != nil {
		log.Error("upload spider error: " + err.Error())
		debug.PrintStack()
		return err
	}

	return nil
}
