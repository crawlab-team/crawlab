package notification

import (
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/entity"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/crawlab-team/crawlab/core/utils"
	mongo2 "github.com/crawlab-team/crawlab/db/mongo"
	parser "github.com/crawlab-team/crawlab/template-parser"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	col      *mongo2.Col // notification settings
	modelSvc service.ModelService
}

func (svc *Service) Init() (err error) {
	if !utils.IsPro() {
		return nil
	}

	return nil
}

func (svc *Service) Start() (err error) {
	// initialize data
	if err := svc.initData(); err != nil {
		return err
	}

	return nil
}

func (svc *Service) Stop() (err error) {
	return nil
}

func (svc *Service) initData() (err error) {
	total, err := svc.col.Count(nil)
	if err != nil {
		return err
	}
	if total > 0 {
		return nil
	}

	// data to initialize
	settings := []NotificationSetting{
		{
			Id:          primitive.NewObjectID(),
			Type:        TypeMail,
			Enabled:     true,
			Name:        "任务通知（邮件）",
			Description: "这是默认的邮件通知。您可以使用您自己的设置进行编辑。",
			TaskTrigger: constants.NotificationTriggerTaskFinish,
			Title:       "[Crawlab] 爬虫任务更新: {{$.status}}",
			Template: `尊敬的 {{$.user.username}},

请查看下面的任务数据。

|键|值|
|:-:|:--|
|任务状态|{{$.status}}|
|任务优先级|{{$.priority}}|
|任务模式|{{$.mode}}|
|执行命令|{{$.cmd}}|
|执行参数|{{$.param}}|
|错误信息|{{$.error}}|
|节点|{{$.node.name}}|
|爬虫|{{$.spider.name}}|
|项目|{{$.spider.project.name}}|
|定时任务|{{$.schedule.name}}|
|结果数|{{$.:task_stat.result_count}}|
|等待时间（秒）|{#{{$.:task_stat.wait_duration}}/1000#}|
|运行时间（秒）|{#{{$.:task_stat.runtime_duration}}/1000#}|
|总时间（秒）|{#{{$.:task_stat.total_duration}}/1000#}|
|平均结果数/秒|{#{{$.:task_stat.result_count}}/({{$.:task_stat.total_duration}}/1000)#}|
`,
			Mail: NotificationSettingMail{
				Server: "smtp.163.com",
				Port:   "465",
				To:     "{{$.user[create].email}}",
			},
		},
		{
			Id:          primitive.NewObjectID(),
			Type:        TypeMail,
			Enabled:     true,
			Name:        "Task Change (Mail)",
			Description: "This is the default mail notification. You can edit it with your own settings",
			TaskTrigger: constants.NotificationTriggerTaskFinish,
			Title:       "[Crawlab] Task Update: {{$.status}}",
			Template: `Dear {{$.user.username}},

Please find the task data as below.

|Key|Value|
|:-:|:--|
|Task Status|{{$.status}}|
|Task Priority|{{$.priority}}|
|Task Mode|{{$.mode}}|
|Task Command|{{$.cmd}}|
|Task Params|{{$.param}}|
|Error Message|{{$.error}}|
|Node|{{$.node.name}}|
|Spider|{{$.spider.name}}|
|Project|{{$.spider.project.name}}|
|Schedule|{{$.schedule.name}}|
|Result Count|{{$.:task_stat.result_count}}|
|Wait Duration (sec)|{#{{$.:task_stat.wait_duration}}/1000#}|
|Runtime Duration (sec)|{#{{$.:task_stat.runtime_duration}}/1000#}|
|Total Duration (sec)|{#{{$.:task_stat.total_duration}}/1000#}|
|Avg Results / Sec|{#{{$.:task_stat.result_count}}/({{$.:task_stat.total_duration}}/1000)#}|
`,
			Mail: NotificationSettingMail{
				Server: "smtp.163.com",
				Port:   "465",
				To:     "{{$.user[create].email}}",
			},
		},
		{
			Id:          primitive.NewObjectID(),
			Type:        TypeMobile,
			Enabled:     true,
			Name:        "任务通知（移动端）",
			Description: "这是默认的手机通知。您可以使用您自己的设置进行编辑。",
			TaskTrigger: constants.NotificationTriggerTaskFinish,
			Title:       "[Crawlab] 任务更新: {{$.status}}",
			Template: `尊敬的 {{$.user.username}},

请查看下面的任务数据。

- **任务状态**: {{$.status}}
- **任务优先级**: {{$.priority}}
- **任务模式**: {{$.mode}}
- **执行命令**: {{$.cmd}}
- **执行参数**: {{$.param}}
- **错误信息**: {{$.error}}
- **节点**: {{$.node.name}}
- **爬虫**: {{$.spider.name}}
- **项目**: {{$.spider.project.name}}
- **定时任务**: {{$.schedule.name}}
- **结果数**: {{$.:task_stat.result_count}}
- **等待时间（秒）**: {#{{$.:task_stat.wait_duration}}/1000#}
- **运行时间（秒）**: {#{{$.:task_stat.runtime_duration}}/1000#}
- **总时间（秒）**: {#{{$.:task_stat.total_duration}}/1000#}
- **平均结果数/秒**: {#{{$.:task_stat.result_count}}/({{$.:task_stat.total_duration}}/1000)#}`,
			Mobile: NotificationSettingMobile{},
		},
		{
			Id:          primitive.NewObjectID(),
			Type:        TypeMobile,
			Enabled:     true,
			Name:        "Task Change (Mobile)",
			Description: "This is the default mobile notification. You can edit it with your own settings",
			TaskTrigger: constants.NotificationTriggerTaskError,
			Title:       "[Crawlab] Task Update: {{$.status}}",
			Template: `Dear {{$.user.username}},

Please find the task data as below.

- **Task Status**: {{$.status}}
- **Task Priority**: {{$.priority}}
- **Task Mode**: {{$.mode}}
- **Task Command**: {{$.cmd}}
- **Task Params**: {{$.param}}
- **Error Message**: {{$.error}}
- **Node**: {{$.node.name}}
- **Spider**: {{$.spider.name}}
- **Project**: {{$.spider.project.name}}
- **Schedule**: {{$.schedule.name}}
- **Result Count**: {{$.:task_stat.result_count}}
- **Wait Duration (sec)**: {#{{$.:task_stat.wait_duration}}/1000#}
- **Runtime Duration (sec)**: {#{{$.:task_stat.runtime_duration}}/1000#}
- **Total Duration (sec)**: {#{{$.:task_stat.total_duration}}/1000#}
- **Avg Results / Sec**: {#{{$.:task_stat.result_count}}/({{$.:task_stat.total_duration}}/1000)#}`,
			Mobile: NotificationSettingMobile{},
		},
	}
	var data []interface{}
	for _, s := range settings {
		data = append(data, s)
	}
	_, err = svc.col.InsertMany(data)
	if err != nil {
		return err
	}
	return nil
}

func (svc *Service) Send(s NotificationSetting, entity bson.M) (err error) {
	switch s.Type {
	case TypeMail:
		return svc.SendMail(s, entity)
	case TypeMobile:
		return svc.SendMobile(s, entity)
	}
	return nil
}

func (svc *Service) SendMail(s NotificationSetting, entity bson.M) (err error) {
	// to
	to, err := parser.Parse(s.Mail.To, entity)
	if err != nil {
		log.Warnf("parsing 'to' error: %v", err)
	}
	if to == "" {
		return nil
	}

	// cc
	cc, err := parser.Parse(s.Mail.Cc, entity)
	if err != nil {
		log.Warnf("parsing 'cc' error: %v", err)
	}

	// title
	title, err := parser.Parse(s.Title, entity)
	if err != nil {
		log.Warnf("parsing 'title' error: %v", err)
	}

	// content
	content, err := parser.Parse(s.Template, entity)
	if err != nil {
		log.Warnf("parsing 'content' error: %v", err)
	}

	// send mail
	if err := SendMail(&models.NotificationSettingMail{
		Server:         s.Mail.Server,
		Port:           s.Mail.Port,
		User:           s.Mail.User,
		Password:       s.Mail.Password,
		SenderEmail:    s.Mail.SenderEmail,
		SenderIdentity: s.Mail.SenderIdentity,
		To:             s.Mail.To,
		Cc:             s.Mail.Cc,
	}, to, cc, title, content); err != nil {
		return err
	}

	return nil
}

func (svc *Service) SendMobile(s NotificationSetting, entity bson.M) (err error) {
	// webhook
	webhook, err := parser.Parse(s.Mobile.Webhook, entity)
	if err != nil {
		log.Warnf("parsing 'webhook' error: %v", err)
	}
	if webhook == "" {
		return nil
	}

	// title
	title, err := parser.Parse(s.Title, entity)
	if err != nil {
		log.Warnf("parsing 'title' error: %v", err)
	}

	// content
	content, err := parser.Parse(s.Template, entity)
	if err != nil {
		log.Warnf("parsing 'content' error: %v", err)
	}

	// send
	if err := SendMobileNotification(webhook, title, content); err != nil {
		return err
	}

	return nil
}

func (svc *Service) GetSettingList(query bson.M, pagination *entity.Pagination, sort bson.D) (res []NotificationSetting, total int, err error) {
	// options
	var options *mongo2.FindOptions
	if pagination != nil || sort != nil {
		options = new(mongo2.FindOptions)
		if pagination != nil {
			options.Skip = pagination.Size * (pagination.Page - 1)
			options.Limit = pagination.Size
		}
		if sort != nil {
			options.Sort = sort
		}
	}

	// get list
	var list []NotificationSetting
	if err := svc.col.Find(query, options).All(&list); err != nil {
		if err.Error() == mongo.ErrNoDocuments.Error() {
			return nil, 0, nil
		} else {
			return nil, 0, err
		}
	}

	// total count
	total, err = svc.col.Count(query)
	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (svc *Service) GetSetting(id primitive.ObjectID) (res *NotificationSetting, err error) {
	var s NotificationSetting
	if err := svc.col.FindId(id).One(&s); err != nil {
		return nil, err
	}
	return &s, nil
}

func (svc *Service) PosSetting(s *NotificationSetting) (err error) {
	s.Id = primitive.NewObjectID()
	if _, err := svc.col.Insert(s); err != nil {
		return err
	}
	return nil
}

func (svc *Service) PutSetting(id primitive.ObjectID, s NotificationSetting) (err error) {
	if err := svc.col.ReplaceId(id, s); err != nil {
		return err
	}

	return nil
}

func (svc *Service) DeleteSetting(id primitive.ObjectID) (err error) {
	if err := svc.col.DeleteId(id); err != nil {
		return err
	}

	return nil
}

func (svc *Service) EnableSetting(id primitive.ObjectID) (err error) {
	return svc._toggleSettingFunc(true)(id)
}

func (svc *Service) DisableSetting(id primitive.ObjectID) (err error) {
	return svc._toggleSettingFunc(false)(id)
}

func (svc *Service) _toggleSettingFunc(value bool) func(id primitive.ObjectID) error {
	return func(id primitive.ObjectID) (err error) {
		var s NotificationSetting
		if err := svc.col.FindId(id).One(&s); err != nil {
			return err
		}
		s.Enabled = value
		if err := svc.col.ReplaceId(id, s); err != nil {
			return err
		}
		return nil
	}
}

func NewService() *Service {
	// service
	svc := &Service{
		col: mongo2.GetMongoCol(SettingsColName),
	}

	// model service
	modelSvc, err := service.GetService()
	if err != nil {
		panic(err)
	}
	svc.modelSvc = modelSvc

	if err := svc.Init(); err != nil {
		panic(err)
	}

	return svc
}

var _service *Service

func GetService() *Service {
	if _service == nil {
		_service = NewService()
	}
	return _service
}
