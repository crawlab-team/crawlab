package services

import (
	"crawlab/constants"
	"crawlab/database"
	"crawlab/entity"
	"crawlab/model"
	"crawlab/model/config_spider"
	"crawlab/services/spider_handler"
	"crawlab/utils"
	"errors"
	"fmt"
	"github.com/apex/log"
	"github.com/globalsign/mgo/bson"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
)

func GenerateConfigSpiderFiles(spider model.Spider, configData entity.ConfigSpiderData) error {
	// 校验Spiderfile正确性
	if err := ValidateSpiderfile(configData); err != nil {
		return err
	}

	// 构造代码生成器
	generator := config_spider.ScrapyGenerator{
		Spider:     spider,
		ConfigData: configData,
	}

	// 生成代码
	if err := generator.Generate(); err != nil {
		return err
	}

	return nil
}

// 验证Spiderfile
func ValidateSpiderfile(configData entity.ConfigSpiderData) error {
	// 获取所有字段
	fields := config_spider.GetAllFields(configData)

	// 校验是否存在 start_url
	if configData.StartUrl == "" {
		return errors.New("spiderfile invalid: start_url is empty")
	}

	// 校验是否存在 start_stage
	if configData.StartStage == "" {
		return errors.New("spiderfile invalid: start_stage is empty")
	}

	// 校验是否存在 stages
	if len(configData.Stages) == 0 {
		return errors.New("spiderfile invalid: stages is empty")
	}

	// 校验stages
	dict := map[string]int{}
	for _, stage := range configData.Stages {
		stageName := stage.Name

		// stage 名称不能为空
		if stageName == "" {
			return errors.New("spiderfile invalid: stage name is empty")
		}

		// stage 名称不能为保留字符串
		// NOTE: 如果有其他Engine，可以扩展，默认为Scrapy
		if configData.Engine == "" || configData.Engine == constants.EngineScrapy {
			if strings.Contains(constants.ScrapyProtectedStageNames, stageName) {
				return errors.New(fmt.Sprintf("spiderfile invalid: stage name '%s' is protected", stageName))
			}
		} else {
			return errors.New(fmt.Sprintf("spiderfile invalid: engine '%s' is not implemented", configData.Engine))
		}

		// stage 名称不能重复
		if dict[stageName] == 1 {
			return errors.New(fmt.Sprintf("spiderfile invalid: stage name '%s' is duplicated", stageName))
		}
		dict[stageName] = 1

		// stage 字段不能为空
		if len(stage.Fields) == 0 {
			return errors.New(fmt.Sprintf("spiderfile invalid: stage '%s' has no fields", stageName))
		}

		// 是否包含 next_stage
		hasNextStage := false

		// 遍历字段列表
		for _, field := range stage.Fields {
			// stage 的 next stage 只能有一个
			if field.NextStage != "" {
				if hasNextStage {
					return errors.New(fmt.Sprintf("spiderfile invalid: stage '%s' has more than 1 next_stage", stageName))
				}
				hasNextStage = true
			}

			// 字段里 css 和 xpath 只能包含一个
			if field.Css != "" && field.Xpath != "" {
				return errors.New(fmt.Sprintf("spiderfile invalid: field '%s' in stage '%s' has both css and xpath set which is prohibited", field.Name, stageName))
			}
		}

		// stage 里 page_css 和 page_xpath 只能包含一个
		if stage.PageCss != "" && stage.PageXpath != "" {
			return errors.New(fmt.Sprintf("spiderfile invalid: stage '%s' has both page_css and page_xpath set which is prohibited", stageName))
		}

		// stage 里 list_css 和 list_xpath 只能包含一个
		if stage.ListCss != "" && stage.ListXpath != "" {
			return errors.New(fmt.Sprintf("spiderfile invalid: stage '%s' has both list_css and list_xpath set which is prohibited", stageName))
		}

		// 如果 stage 的 is_list 为 true 但 list_css 为空，报错
		if stage.IsList && (stage.ListCss == "" && stage.ListXpath == "") {
			return errors.New("spiderfile invalid: stage with is_list = true should have either list_css or list_xpath being set")
		}
	}

	// 校验字段唯一性
	if !IsUniqueConfigSpiderFields(fields) {
		return errors.New("spiderfile invalid: fields not unique")
	}

	// 字段名称不能为保留字符串
	for _, field := range fields {
		if strings.Contains(constants.ScrapyProtectedFieldNames, field.Name) {
			return errors.New(fmt.Sprintf("spiderfile invalid: field name '%s' is protected", field.Name))
		}
	}

	return nil
}

func IsUniqueConfigSpiderFields(fields []entity.Field) bool {
	dict := map[string]int{}
	for _, field := range fields {
		if dict[field.Name] == 1 {
			return false
		}
		dict[field.Name] = 1
	}
	return true
}

func ProcessSpiderFilesFromConfigData(spider model.Spider, configData entity.ConfigSpiderData) error {
	spiderDir := spider.Src

	// 删除已有的爬虫文件
	for _, fInfo := range utils.ListDir(spiderDir) {
		// 不删除Spiderfile
		if fInfo.Name() == "Spiderfile" {
			continue
		}

		// 删除其他文件
		if err := os.RemoveAll(filepath.Join(spiderDir, fInfo.Name())); err != nil {
			return err
		}
	}

	// 拷贝爬虫文件
	tplDir := "./template/scrapy"
	for _, fInfo := range utils.ListDir(tplDir) {
		// 跳过Spiderfile
		if fInfo.Name() == "Spiderfile" {
			continue
		}

		srcPath := filepath.Join(tplDir, fInfo.Name())
		if fInfo.IsDir() {
			dirPath := filepath.Join(spiderDir, fInfo.Name())
			if err := utils.CopyDir(srcPath, dirPath); err != nil {
				return err
			}
		} else {
			if err := utils.CopyFile(srcPath, filepath.Join(spiderDir, fInfo.Name())); err != nil {
				return err
			}
		}
	}

	// 更改爬虫文件
	if err := GenerateConfigSpiderFiles(spider, configData); err != nil {
		return err
	}

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
		if err := gf.RemoveId(gfFile.Id); err != nil {
			log.Errorf("remove grid fs error: %s", err.Error())
			debug.PrintStack()
			return err
		}
	}

	// 上传到GridFs
	fid, err := RetryUploadToGridFs(spiderZipFileName, tmpFilePath)
	if err != nil {
		log.Errorf("upload to grid fs error: %s", err.Error())
		return err
	}

	// 保存爬虫 FileId
	spider.FileId = fid
	_ = spider.Save()

	// 获取爬虫同步实例
	spiderSync := spider_handler.SpiderSync{
		Spider: spider,
	}

	// 获取gfFile
	gfFile2 := model.GetGridFs(spider.FileId)

	// 生成MD5
	spiderSync.CreateMd5File(gfFile2.Md5)

	return nil
}

func GenerateSpiderfileFromConfigData(spider model.Spider, configData entity.ConfigSpiderData) error {
	// Spiderfile 路径
	sfPath := filepath.Join(spider.Src, "Spiderfile")

	// 生成Yaml内容
	sfContentByte, err := yaml.Marshal(configData)
	if err != nil {
		return err
	}

	// 打开文件
	var f *os.File
	if utils.Exists(sfPath) {
		f, err = os.OpenFile(sfPath, os.O_WRONLY|os.O_TRUNC, 0777)
	} else {
		f, err = os.OpenFile(sfPath, os.O_CREATE, 0777)
	}
	if err != nil {
		return err
	}
	defer f.Close()

	// 写入内容
	if _, err := f.Write(sfContentByte); err != nil {
		return err
	}

	return nil
}
