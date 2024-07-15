package notification

import (
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/entity"
	"github.com/crawlab-team/crawlab/core/models/models/v2"
	"github.com/crawlab-team/crawlab/core/models/service"
	"regexp"
	"sync"
)

type ServiceV2 struct {
}

func (svc *ServiceV2) Start() {
	// initialize data
	if err := svc.initData(); err != nil {
		log.Errorf("[NotificationServiceV2] initializing data error: %v", err)
		return
	}
}

func (svc *ServiceV2) initData() (err error) {
	total, err := service.NewModelServiceV2[models.NotificationSettingV2]().Count(nil)
	if err != nil {
		return err
	}
	if total > 0 {
		return nil
	}

	// data to initialize
	settings := []models.NotificationSettingV2{
		{
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
			Mail: models.NotificationSettingMail{
				Server: "smtp.163.com",
				Port:   "465",
				To:     "{{$.user[create].email}}",
			},
		},
		{
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
			Mail: models.NotificationSettingMail{
				Server: "smtp.163.com",
				Port:   "465",
				To:     "{{$.user[create].email}}",
			},
		},
		{
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
			Mobile: models.NotificationSettingMobile{},
		},
		{
			Type:        TypeMobile,
			Enabled:     true,
			Name:        "Task Change (Mobile)",
			Description: "This is the default mobile notification. You can edit it with your own settings",
			TaskTrigger: constants.NotificationTriggerTaskFinish,
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
			Mobile: models.NotificationSettingMobile{},
		},
	}
	_, err = service.NewModelServiceV2[models.NotificationSettingV2]().InsertMany(settings)
	if err != nil {
		return err
	}
	return nil
}

func (svc *ServiceV2) Send(s *models.NotificationSettingV2, args ...any) (err error) {
	content := svc.getContent(s, args...)
	switch s.Type {
	case TypeMail:
		return svc.SendMail(s, content)
	case TypeMobile:
		return svc.SendMobile(s, content)
	}
	return nil
}

func (svc *ServiceV2) SendMail(s *models.NotificationSettingV2, content string) (err error) {
	// TODO: parse to/cc/bcc

	// send mail
	if err := SendMail(&s.Mail, s.Mail.To, s.Mail.Cc, s.Title, content); err != nil {
		return err
	}

	return nil
}

func (svc *ServiceV2) SendMobile(s *models.NotificationSettingV2, content string) (err error) {
	// send
	if err := SendMobileNotification(s.Mobile.Webhook, s.Title, content); err != nil {
		return err
	}

	return nil
}

func (svc *ServiceV2) getContent(s *models.NotificationSettingV2, args ...any) (content string) {
	switch s.TriggerTarget {
	case constants.NotificationTriggerTargetTask:
		//task := new(models.TaskV2)
		//taskStat := new(models.TaskStatV2)
		//spider := new(models.SpiderV2)
		//node := new(models.NodeV2)
		//for _, arg := range args {
		//	switch arg.(type) {
		//	case models.TaskV2:
		//		task = arg.(*models.TaskV2)
		//	case models.TaskStatV2:
		//		taskStat = arg.(*models.TaskStatV2)
		//	case models.SpiderV2:
		//		spider = arg.(*models.SpiderV2)
		//	case models.NodeV2:
		//		node = arg.(*models.NodeV2)
		//	}
		//}
		switch s.TemplateMode {
		case constants.NotificationTemplateModeMarkdown:
			// TODO: implement
		case constants.NotificationTemplateModeRichText:
			//s.TemplateRichText
		}

	case constants.NotificationTriggerTargetNode:

	}

	return content
}

func (svc *ServiceV2) parseTemplateVariables(s *models.NotificationSettingV2) (variables []entity.NotificationVariable) {
	regex := regexp.MustCompile("\\$\\{(\\w+):(\\w+)}")

	// find all matches
	matches := regex.FindAllStringSubmatch(s.Template, -1)

	// iterate over matches
	for _, match := range matches {
		variables = append(variables, entity.NotificationVariable{
			Category: match[1],
			Name:     match[2],
		})
	}

	return variables
}

func newNotificationServiceV2() *ServiceV2 {
	return &ServiceV2{}
}

var _serviceV2 *ServiceV2
var _serviceV2Once = new(sync.Once)

func GetNotificationServiceV2() *ServiceV2 {
	_serviceV2Once.Do(func() {
		_serviceV2 = newNotificationServiceV2()
	})
	return _serviceV2
}
