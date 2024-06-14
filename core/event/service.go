package event

import (
	"fmt"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab/core/entity"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/utils"
	"github.com/crawlab-team/crawlab/trace"
	"github.com/thoas/go-funk"
	"regexp"
)

var S interfaces.EventService

type Service struct {
	keys     []string
	includes []string
	excludes []string
	chs      []*chan interfaces.EventData
}

func (svc *Service) Register(key, include, exclude string, ch *chan interfaces.EventData) {
	svc.keys = append(svc.keys, key)
	svc.includes = append(svc.includes, include)
	svc.excludes = append(svc.excludes, exclude)
	svc.chs = append(svc.chs, ch)
}

func (svc *Service) Unregister(key string) {
	idx := funk.IndexOfString(svc.keys, key)
	if idx != -1 {
		svc.keys = append(svc.keys[:idx], svc.keys[(idx+1):]...)
		svc.includes = append(svc.includes[:idx], svc.includes[(idx+1):]...)
		svc.excludes = append(svc.excludes[:idx], svc.excludes[(idx+1):]...)
		svc.chs = append(svc.chs[:idx], svc.chs[(idx+1):]...)
		log.Infof("[EventService] unregistered %s", key)
	}
}

func (svc *Service) SendEvent(eventName string, data ...interface{}) {
	for i, key := range svc.keys {
		// include
		include := svc.includes[i]
		matchedInclude, err := regexp.MatchString(include, eventName)
		if err != nil {
			trace.PrintError(err)
			continue
		}
		if !matchedInclude {
			continue
		}

		// exclude
		exclude := svc.excludes[i]
		matchedExclude, err := regexp.MatchString(exclude, eventName)
		if err != nil {
			trace.PrintError(err)
			continue
		}
		if matchedExclude {
			continue
		}

		// send event
		utils.LogDebug(fmt.Sprintf("key %s matches event %s", key, eventName))
		ch := svc.chs[i]
		go func(ch *chan interfaces.EventData) {
			for _, d := range data {
				*ch <- &entity.EventData{
					Event: eventName,
					Data:  d,
				}
			}
		}(ch)
	}
}

func NewEventService() (svc interfaces.EventService) {
	if S != nil {
		return S
	}

	svc = &Service{
		chs:  []*chan interfaces.EventData{},
		keys: []string{},
	}

	S = svc

	return svc
}
