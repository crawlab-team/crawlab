package model

import (
	"crawlab/database"
	"crawlab/utils"
	"github.com/apex/log"
	"github.com/globalsign/mgo/bson"
	"os"
	"runtime/debug"
	"time"
)

type LogItem struct {
	Id      bson.ObjectId `json:"_id" bson:"_id"`
	Message string        `json:"msg" bson:"msg"`
	TaskId  string        `json:"task_id" bson:"task_id"`
	IsError bool          `json:"is_error" bson:"is_error"`
	Seq     int64         `json:"seq" bson:"seq"`
	Ts      time.Time     `json:"ts" bson:"ts"`
}

// 获取本地日志
func GetLocalLog(logPath string) (fileBytes []byte, err error) {

	f, err := os.Open(logPath)
	if err != nil {
		log.Error(err.Error())
		debug.PrintStack()
		return nil, err
	}
	fi, err := f.Stat()
	if err != nil {
		log.Error(err.Error())
		debug.PrintStack()
		return nil, err
	}
	defer utils.Close(f)

	const bufLen = 2 * 1024 * 1024
	logBuf := make([]byte, bufLen)

	off := int64(0)
	if fi.Size() > int64(len(logBuf)) {
		off = fi.Size() - int64(len(logBuf))
	}
	n, err := f.ReadAt(logBuf, off)

	//到文件结尾会有EOF标识
	if err != nil && err.Error() != "EOF" {
		log.Error(err.Error())
		debug.PrintStack()
		return nil, err
	}
	logBuf = logBuf[:n]
	return logBuf, nil
}

func AddLogItem(l LogItem) error {
	s, c := database.GetCol("logs")
	defer s.Close()
	if err := c.Insert(l); err != nil {
		log.Errorf("insert log error: " + err.Error())
		debug.PrintStack()
		return err
	}
	return nil
}

func GetLogItemList(query bson.M, keyword string, skip int, limit int, sortStr string) ([]LogItem, error) {
	s, c := database.GetCol("logs")
	defer s.Close()

	filter := query

	var logItems []LogItem
	if keyword == "" {
		filter["seq"] = bson.M{
			"$gte": skip,
			"$lt":  skip + limit,
		}
		if err := c.Find(filter).Sort(sortStr).All(&logItems); err != nil {
			debug.PrintStack()
			return logItems, err
		}
	} else {
		filter["msg"] = bson.M{
			"$regex": bson.RegEx{
				Pattern: keyword,
				Options: "i",
			},
		}
		if err := c.Find(filter).Sort(sortStr).Skip(skip).Limit(limit).All(&logItems); err != nil {
			debug.PrintStack()
			return logItems, err
		}
	}

	return logItems, nil
}

func GetLogItemTotal(query bson.M, keyword string) (int, error) {
	s, c := database.GetCol("logs")
	defer s.Close()

	filter := query

	if keyword != "" {
		filter["msg"] = bson.M{
			"$regex": bson.RegEx{
				Pattern: keyword,
				Options: "i",
			},
		}
	}

	total, err := c.Find(filter).Count()
	if err != nil {
		debug.PrintStack()
		return total, err
	}

	return total, nil
}
