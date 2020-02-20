package services

import (
	"crawlab/entity"
	"github.com/apex/log"
	"github.com/imroc/req"
	"runtime/debug"
	"sort"
)

func GetLatestRelease() (release entity.Release, err error) {
	res, err := req.Get("https://api.github.com/repos/crawlab-team/crawlab/releases")
	if err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return release, err
	}

	var releaseDataList entity.ReleaseSlices
	if err := res.ToJSON(&releaseDataList); err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return release, err
	}

	sort.Sort(releaseDataList)

	return releaseDataList[len(releaseDataList)-1], nil
}
