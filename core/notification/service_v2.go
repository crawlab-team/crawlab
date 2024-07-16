package notification

import (
	"fmt"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/entity"
	"github.com/crawlab-team/crawlab/core/models/models/v2"
	"github.com/crawlab-team/crawlab/core/models/service"
	"regexp"
	"strings"
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

func (svc *ServiceV2) Send(s *models.NotificationSettingV2, args ...any) {
	content := svc.getContent(s, args...)
	switch s.Type {
	case TypeMail:
		svc.SendMail(s, content)
	case TypeMobile:
		svc.SendMobile(s, content)
	}
}

func (svc *ServiceV2) SendMail(s *models.NotificationSettingV2, content string) {
	// TODO: parse to/cc/bcc

	// send mail
	err := SendMail(&s.Mail, s.Mail.To, s.Mail.Cc, s.Title, content)
	if err != nil {
		log.Errorf("[NotificationServiceV2] send mail error: %v", err)
	}
}

func (svc *ServiceV2) SendMobile(s *models.NotificationSettingV2, content string) {
	err := SendMobileNotification(s.Mobile.Webhook, s.Title, content)
	if err != nil {
		log.Errorf("[NotificationServiceV2] send mobile notification error: %v", err)
	}
}

func (svc *ServiceV2) getContent(s *models.NotificationSettingV2, args ...any) (content string) {
	switch s.TriggerTarget {
	case constants.NotificationTriggerTargetTask:
		vd := svc.getTaskVariableData(args...)
		switch s.TemplateMode {
		case constants.NotificationTemplateModeMarkdown:
			variables := svc.parseTemplateVariables(s.TemplateMarkdown)
			return svc.getTaskContent(s.TemplateMarkdown, variables, vd)
		case constants.NotificationTemplateModeRichText:
			variables := svc.parseTemplateVariables(s.TemplateRichText)
			return svc.getTaskContent(s.TemplateRichText, variables, vd)
		}

	case constants.NotificationTriggerTargetNode:
		// TODO: implement

	}

	return content
}

func (svc *ServiceV2) getTaskContent(template string, variables []entity.NotificationVariable, vd VariableDataTask) (content string) {
	content = template
	for _, v := range variables {
		switch v.Category {
		case "task":
			switch v.Name {
			case "id":
				content = strings.ReplaceAll(content, v.GetKey(), vd.Task.Id.Hex())
			case "status":
				content = strings.ReplaceAll(content, v.GetKey(), vd.Task.Status)
			case "priority":
				content = strings.ReplaceAll(content, v.GetKey(), fmt.Sprintf("%d", vd.Task.Priority))
			case "mode":
				content = strings.ReplaceAll(content, v.GetKey(), vd.Task.Mode)
			case "cmd":
				content = strings.ReplaceAll(content, v.GetKey(), vd.Task.Cmd)
			case "param":
				content = strings.ReplaceAll(content, v.GetKey(), vd.Task.Param)
			}
		}
	}
	return content
}

func (svc *ServiceV2) getTaskVariableData(args ...any) (vd VariableDataTask) {
	for _, arg := range args {
		switch arg.(type) {
		case *models.TaskV2:
			vd.Task = arg.(*models.TaskV2)
		case *models.TaskStatV2:
			vd.TaskStat = arg.(*models.TaskStatV2)
		case *models.SpiderV2:
			vd.Spider = arg.(*models.SpiderV2)
		case *models.NodeV2:
			vd.Node = arg.(*models.NodeV2)
		case *models.ScheduleV2:
			vd.Schedule = arg.(*models.ScheduleV2)
		}
	}
	return vd
}

func (svc *ServiceV2) parseTemplateVariables(template string) (variables []entity.NotificationVariable) {
	// regex pattern
	regex := regexp.MustCompile("\\$\\{(\\w+):(\\w+)}")

	// find all matches
	matches := regex.FindAllStringSubmatch(template, -1)

	// variables map
	variablesMap := make(map[string]entity.NotificationVariable)

	// iterate over matches
	for _, match := range matches {
		variable := entity.NotificationVariable{
			Category: match[1],
			Name:     match[2],
		}
		key := fmt.Sprintf("%s:%s", variable.Category, variable.Name)
		if _, ok := variablesMap[key]; !ok {
			variablesMap[key] = variable
		}
	}

	// convert map to slice
	for _, variable := range variablesMap {
		variables = append(variables, variable)
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
