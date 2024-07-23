package notification

import (
	"fmt"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/entity"
	"github.com/crawlab-team/crawlab/core/models/models/v2"
	"github.com/crawlab-team/crawlab/core/models/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"regexp"
	"strings"
	"sync"
	"time"
)

type ServiceV2 struct {
}

func (svc *ServiceV2) Send(s *models.NotificationSettingV2, args ...any) {
	content := svc.getContent(s, args...)

	wg := sync.WaitGroup{}
	wg.Add(len(s.ChannelIds))
	for _, chId := range s.ChannelIds {
		go func(chId primitive.ObjectID) {
			defer wg.Done()
			ch, err := service.NewModelServiceV2[models.NotificationChannelV2]().GetById(chId)
			if err != nil {
				log.Errorf("[NotificationServiceV2] get channel error: %v", err)
				return
			}
			switch ch.Type {
			case TypeMail:
				svc.SendMail(s, ch, content)
			case TypeIM:
				svc.SendIM(s, ch, content)
			}
		}(chId)
	}
	wg.Wait()
}

func (svc *ServiceV2) SendMail(s *models.NotificationSettingV2, ch *models.NotificationChannelV2, content string) {
	// TODO: parse to/cc/bcc
	mailTo := s.MailTo
	mailCc := s.MailCc
	mailBcc := s.MailBcc

	// send mail
	err := SendMail(s, ch, mailTo, mailCc, mailBcc, s.Title, content)
	if err != nil {
		log.Errorf("[NotificationServiceV2] send mail error: %v", err)
	}
}

func (svc *ServiceV2) SendIM(s *models.NotificationSettingV2, ch *models.NotificationChannelV2, content string) {
	err := SendIMNotification(s, ch, content)
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
			case "cmd":
				content = strings.ReplaceAll(content, v.GetKey(), vd.Task.Cmd)
			case "param":
				content = strings.ReplaceAll(content, v.GetKey(), vd.Task.Param)
			case "error":
				content = strings.ReplaceAll(content, v.GetKey(), vd.Task.Error)
			case "pid":
				content = strings.ReplaceAll(content, v.GetKey(), fmt.Sprintf("%d", vd.Task.Pid))
			case "type":
				content = strings.ReplaceAll(content, v.GetKey(), vd.Task.Type)
			case "mode":
				content = strings.ReplaceAll(content, v.GetKey(), vd.Task.Mode)
			case "priority":
				content = strings.ReplaceAll(content, v.GetKey(), fmt.Sprintf("%d", vd.Task.Priority))
			case "created_at":
				content = strings.ReplaceAll(content, v.GetKey(), vd.Task.CreatedAt.Format(time.DateTime))
			case "updated_at":
				content = strings.ReplaceAll(content, v.GetKey(), vd.Task.UpdatedAt.Format(time.DateTime))
			}

		case "spider":
			switch v.Name {
			case "id":
				content = strings.ReplaceAll(content, v.GetKey(), vd.Spider.Id.Hex())
			case "name":
				content = strings.ReplaceAll(content, v.GetKey(), vd.Spider.Name)
			case "type":
				content = strings.ReplaceAll(content, v.GetKey(), vd.Spider.Type)
			case "description":
				content = strings.ReplaceAll(content, v.GetKey(), vd.Spider.Description)
			case "mode":
				content = strings.ReplaceAll(content, v.GetKey(), vd.Spider.Mode)
			case "cmd":
				content = strings.ReplaceAll(content, v.GetKey(), vd.Spider.Cmd)
			case "param":
				content = strings.ReplaceAll(content, v.GetKey(), vd.Spider.Param)
			case "priority":
				content = strings.ReplaceAll(content, v.GetKey(), fmt.Sprintf("%d", vd.Spider.Priority))
			case "created_at":
				content = strings.ReplaceAll(content, v.GetKey(), vd.Spider.CreatedAt.Format(time.DateTime))
			case "updated_at":
				content = strings.ReplaceAll(content, v.GetKey(), vd.Spider.UpdatedAt.Format(time.DateTime))
			}

		case "node":
			switch v.Name {
			case "id":
				content = strings.ReplaceAll(content, v.GetKey(), vd.Node.Id.Hex())
			case "key":
				content = strings.ReplaceAll(content, v.GetKey(), vd.Node.Key)
			case "name":
				content = strings.ReplaceAll(content, v.GetKey(), vd.Node.Name)
			case "ip":
				content = strings.ReplaceAll(content, v.GetKey(), vd.Node.Ip)
			case "port":
				content = strings.ReplaceAll(content, v.GetKey(), vd.Node.Port)
			case "hostname":
				content = strings.ReplaceAll(content, v.GetKey(), vd.Node.Hostname)
			case "description":
				content = strings.ReplaceAll(content, v.GetKey(), vd.Node.Description)
			case "status":
				content = strings.ReplaceAll(content, v.GetKey(), vd.Node.Status)
			case "enabled":
				content = strings.ReplaceAll(content, v.GetKey(), fmt.Sprintf("%t", vd.Node.Enabled))
			case "active":
				content = strings.ReplaceAll(content, v.GetKey(), fmt.Sprintf("%t", vd.Node.Active))
			case "active_at":
				content = strings.ReplaceAll(content, v.GetKey(), vd.Node.ActiveAt.Format("2006-01-02 15:04:05"))
			case "available_runners":
				content = strings.ReplaceAll(content, v.GetKey(), fmt.Sprintf("%d", vd.Node.AvailableRunners))
			case "max_runners":
				content = strings.ReplaceAll(content, v.GetKey(), fmt.Sprintf("%d", vd.Node.MaxRunners))
			case "created_at":
				content = strings.ReplaceAll(content, v.GetKey(), vd.Node.CreatedAt.Format(time.DateTime))
			case "updated_at":
				content = strings.ReplaceAll(content, v.GetKey(), vd.Node.UpdatedAt.Format(time.DateTime))
			}

		case "schedule":
			switch v.Name {
			case "id":
				content = strings.ReplaceAll(content, v.GetKey(), vd.Schedule.Id.Hex())
			case "name":
				content = strings.ReplaceAll(content, v.GetKey(), vd.Schedule.Name)
			case "description":
				content = strings.ReplaceAll(content, v.GetKey(), vd.Schedule.Description)
			case "cron":
				content = strings.ReplaceAll(content, v.GetKey(), vd.Schedule.Cron)
			case "cmd":
				content = strings.ReplaceAll(content, v.GetKey(), vd.Schedule.Cmd)
			case "param":
				content = strings.ReplaceAll(content, v.GetKey(), vd.Schedule.Param)
			case "mode":
				content = strings.ReplaceAll(content, v.GetKey(), vd.Schedule.Mode)
			case "priority":
				content = strings.ReplaceAll(content, v.GetKey(), fmt.Sprintf("%d", vd.Schedule.Priority))
			case "enabled":
				content = strings.ReplaceAll(content, v.GetKey(), fmt.Sprintf("%t", vd.Schedule.Enabled))
			case "created_at":
				content = strings.ReplaceAll(content, v.GetKey(), vd.Schedule.CreatedAt.Format(time.DateTime))
			case "updated_at":
				content = strings.ReplaceAll(content, v.GetKey(), vd.Schedule.UpdatedAt.Format(time.DateTime))
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

func (svc *ServiceV2) getUserById(id primitive.ObjectID) {

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
