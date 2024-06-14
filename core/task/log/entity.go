package log

import "time"

type Message struct {
	Id  int64     `json:"id" bson:"id"`
	Msg string    `json:"msg" bson:"msg"`
	Ts  time.Time `json:"ts" bson:"ts"`
}

type Metadata struct {
	Size       int64  `json:"size,omitempty" bson:"size"`
	TotalLines int64  `json:"total_lines,omitempty" bson:"total_lines"`
	TotalBytes int64  `json:"total_bytes,omitempty" bson:"total_bytes"`
	Md5        string `json:"md5,omitempty" bson:"md5"`
}
