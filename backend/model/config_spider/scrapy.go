package config_spider

import (
	"crawlab/constants"
	"crawlab/entity"
	"crawlab/model"
	"crawlab/utils"
	"errors"
	"fmt"
	"path/filepath"
)

type ScrapyGenerator struct {
	Spider     model.Spider
	ConfigData entity.ConfigSpiderData
}

// 生成爬虫文件
func (g ScrapyGenerator) Generate() error {
	// 生成 items.py
	if err := g.ProcessItems(); err != nil {
		return err
	}

	// 生成 spider.py
	if err := g.ProcessSpider(); err != nil {
		return err
	}
	return nil
}

// 生成 items.py
func (g ScrapyGenerator) ProcessItems() error {
	// 待处理文件名
	src := g.Spider.Src
	filePath := filepath.Join(src, "config_spider", "items.py")

	// 获取所有字段
	fields := g.GetAllFields()

	// 字段名列表（包含默认字段名）
	fieldNames := []string{
		"_id",
		"task_id",
		"ts",
	}

	// 加入字段
	for _, field := range fields {
		fieldNames = append(fieldNames, field.Name)
	}

	// 将字段名转化为python代码
	str := ""
	for _, fieldName := range fieldNames {
		line := g.PadCode(fmt.Sprintf("%s = scrapy.Field()", fieldName), 1)
		str += line
	}

	// 将占位符替换为代码
	if err := utils.SetFileVariable(filePath, constants.AnchorItems, str); err != nil {
		return err
	}

	return nil
}

// 生成 spider.py
func (g ScrapyGenerator) ProcessSpider() error {
	// 待处理文件名
	src := g.Spider.Src
	filePath := filepath.Join(src, "config_spider", "spiders", "spider.py")

	// 替换 start_stage
	if err := utils.SetFileVariable(filePath, constants.AnchorStartStage, "parse_"+GetStartStageName(g.ConfigData)); err != nil {
		return err
	}

	// 替换 start_url
	if err := utils.SetFileVariable(filePath, constants.AnchorStartUrl, g.ConfigData.StartUrl); err != nil {
		return err
	}

	// 替换 parsers
	strParser := ""
	for stageName, stage := range g.ConfigData.Stages {
		stageStr := g.GetParserString(stageName, stage)
		strParser += stageStr
	}
	if err := utils.SetFileVariable(filePath, constants.AnchorParsers, strParser); err != nil {
		return err
	}

	return nil
}

func (g ScrapyGenerator) GetParserString(stageName string, stage entity.Stage) string {
	// 构造函数定义行
	strDef := g.PadCode(fmt.Sprintf("def parse_%s(self, response):", stageName), 1)

	strParse := ""
	if stage.IsList {
		// 列表逻辑
		strParse = g.GetListParserString(stageName, stage)
	} else {
		// 非列表逻辑
		strParse = g.GetNonListParserString(stageName, stage)
	}

	// 构造
	str := fmt.Sprintf(`%s%s`, strDef, strParse)

	return str
}

func (g ScrapyGenerator) PadCode(str string, num int) string {
	res := ""
	for i := 0; i < num; i++ {
		res += "    "
	}
	res += str
	res += "\n"
	return res
}

func (g ScrapyGenerator) GetNonListParserString(stageName string, stage entity.Stage) string {
	str := ""

	// 获取或构造item
	str += g.PadCode("item = Item() if response.meta.get('item') is None else response.meta.get('item')", 2)

	// 遍历字段列表
	for _, f := range stage.Fields {
		line := fmt.Sprintf(`item['%s'] = response.%s.extract_first()`, f.Name, g.GetExtractStringFromField(f))
		line = g.PadCode(line, 2)
		str += line
	}

	// next stage 字段
	if f, err := g.GetNextStageField(stage); err == nil {
		// 如果找到 next stage 字段，进行下一个回调
		str += g.PadCode(fmt.Sprintf(`yield scrapy.Request(url="get_real_url(response, item['%s'])", callback=self.parse_%s, meta={'item': item})`, f.Name, f.NextStage), 2)
	} else {
		// 如果没找到 next stage 字段，返回 item
		str += g.PadCode(fmt.Sprintf(`yield item`), 2)
	}

	// 加入末尾换行
	str += g.PadCode("", 0)

	return str
}

func (g ScrapyGenerator) GetListParserString(stageName string, stage entity.Stage) string {
	str := ""

	// 获取前一个 stage 的 item
	str += g.PadCode(`prev_item = response.meta.get('item')`, 2)

	// for 循环遍历列表
	str += g.PadCode(fmt.Sprintf(`for elem in response.css('%s'):`, stage.ListCss), 2)

	// 构造item
	str += g.PadCode(`item = Item()`, 3)

	// 遍历字段列表
	for _, f := range stage.Fields {
		line := fmt.Sprintf(`item['%s'] = elem.%s.extract_first()`, f.Name, g.GetExtractStringFromField(f))
		line = g.PadCode(line, 3)
		str += line
	}

	// 把前一个 stage 的 item 值赋给当前 item
	str += g.PadCode(`if prev_item is not None:`, 3)
	str += g.PadCode(`for key, value in prev_item.items():`, 4)
	str += g.PadCode(`item[key] = value`, 5)

	// next stage 字段
	if f, err := g.GetNextStageField(stage); err == nil {
		// 如果找到 next stage 字段，进行下一个回调
		str += g.PadCode(fmt.Sprintf(`yield scrapy.Request(url=get_real_url(response, item['%s']), callback=self.parse_%s, meta={'item': item})`, f.Name, f.NextStage), 3)
	} else {
		// 如果没找到 next stage 字段，返回 item
		str += g.PadCode(fmt.Sprintf(`yield item`), 3)
	}

	// 分页
	if stage.PageCss != "" || stage.PageXpath != "" {
		str += g.PadCode(fmt.Sprintf(`next_url = response.%s.extract_first()`, g.GetExtractStringFromStage(stage)), 2)
		str += g.PadCode(fmt.Sprintf(`yield scrapy.Request(url=get_real_url(response, next_url), callback=self.parse_%s, meta={'item': item})`, stageName), 2)
	}

	// 加入末尾换行
	str += g.PadCode("", 0)

	return str
}

// 获取所有字段
func (g ScrapyGenerator) GetAllFields() []entity.Field {
	return GetAllFields(g.ConfigData)
}

// 获取包含 next stage 的字段
func (g ScrapyGenerator) GetNextStageField(stage entity.Stage) (entity.Field, error) {
	for _, field := range stage.Fields {
		if field.NextStage != "" {
			return field, nil
		}
	}
	return entity.Field{}, errors.New("cannot find next stage field")
}

func (g ScrapyGenerator) GetExtractStringFromField(f entity.Field) string {
	if f.Css != "" {
		// 如果为CSS
		if f.Attr == "" {
			// 文本
			return fmt.Sprintf(`css('%s::text')`, f.Css)
		} else {
			// 属性
			return fmt.Sprintf(`css('%s::attr("%s")')`, f.Css, f.Attr)
		}
	} else {
		// 如果为XPath
		if f.Attr == "" {
			// 文本
			return fmt.Sprintf(`xpath('%s/text()')`, f.Xpath)
		} else {
			// 属性
			return fmt.Sprintf(`xpath('%s/@%s')`, f.Xpath, f.Attr)
		}
	}
}

func (g ScrapyGenerator) GetExtractStringFromStage(stage entity.Stage) string {
	// 分页元素属性，默认为 href
	pageAttr := "href"
	if stage.PageAttr != "" {
		pageAttr = stage.PageAttr
	}

	if stage.PageCss != "" {
		// 如果为CSS
		return fmt.Sprintf(`css('%s::attr("%s")')`, stage.PageCss, pageAttr)
	} else {
		// 如果为XPath
		return fmt.Sprintf(`xpath('%s/@%s')`, stage.PageXpath, pageAttr)
	}
}
