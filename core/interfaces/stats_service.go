package interfaces

import "go.mongodb.org/mongo-driver/bson"

type StatsService interface {
	GetOverviewStats(query bson.M) (data interface{}, err error)
	GetDailyStats(query bson.M) (data interface{}, err error)
	GetTaskStats(query bson.M) (data interface{}, err error)
}
