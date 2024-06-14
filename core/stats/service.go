package stats

import (
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/entity"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
}

func (svc *Service) GetOverviewStats(query bson.M) (data interface{}, err error) {
	stats := bson.M{}

	// nodes
	stats["nodes"], err = mongo.GetMongoCol(interfaces.ModelColNameNode).Count(bson.M{"active": true})
	if err != nil {
		if err.Error() != mongo2.ErrNoDocuments.Error() {
			return nil, err
		}
		stats["nodes"] = 0
	}

	// projects
	stats["projects"], err = mongo.GetMongoCol(interfaces.ModelColNameProject).Count(nil)
	if err != nil {
		if err.Error() != mongo2.ErrNoDocuments.Error() {
			return nil, err
		}
		stats["projects"] = 0
	}

	// spiders
	stats["spiders"], err = mongo.GetMongoCol(interfaces.ModelColNameSpider).Count(nil)
	if err != nil {
		if err.Error() != mongo2.ErrNoDocuments.Error() {
			return nil, err
		}
		stats["spiders"] = 0
	}

	// schedules
	stats["schedules"], err = mongo.GetMongoCol(interfaces.ModelColNameSchedule).Count(nil)
	if err != nil {
		if err.Error() != mongo2.ErrNoDocuments.Error() {
			return nil, err
		}
		stats["schedules"] = 0
	}

	// tasks
	stats["tasks"], err = mongo.GetMongoCol(interfaces.ModelColNameTask).Count(nil)
	if err != nil {
		if err.Error() != mongo2.ErrNoDocuments.Error() {
			return nil, err
		}
		stats["tasks"] = 0
	}

	// error tasks
	stats["error_tasks"], err = mongo.GetMongoCol(interfaces.ModelColNameTask).Count(bson.M{"status": constants.TaskStatusError})
	if err != nil {
		if err.Error() != mongo2.ErrNoDocuments.Error() {
			return nil, err
		}
		stats["error_tasks"] = 0
	}

	// results
	stats["results"], err = svc.getOverviewResults(query)
	if err != nil {
		if err.Error() != mongo2.ErrNoDocuments.Error() {
			return nil, err
		}
		stats["results"] = 0
	}

	// users
	stats["users"], err = mongo.GetMongoCol(interfaces.ModelColNameUser).Count(nil)
	if err != nil {
		if err.Error() != mongo2.ErrNoDocuments.Error() {
			return nil, err
		}
		stats["users"] = 0
	}

	return stats, nil
}

func (svc *Service) GetDailyStats(query bson.M) (data interface{}, err error) {
	tasksStats, err := svc.getDailyTasksStats(query)
	if err != nil {
		return nil, err
	}
	return tasksStats, nil
}

func (svc *Service) GetTaskStats(query bson.M) (data interface{}, err error) {
	stats := bson.M{}

	// by status
	stats["by_status"], err = svc.getTaskStatsByStatus(query)
	if err != nil {
		return nil, err
	}

	// by node
	stats["by_node"], err = svc.getTaskStatsByNode(query)
	if err != nil {
		return nil, err
	}

	// by spider
	stats["by_spider"], err = svc.getTaskStatsBySpider(query)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

func (svc *Service) getDailyTasksStats(query bson.M) (data interface{}, err error) {
	pipeline := mongo2.Pipeline{
		{{
			"$match", query,
		}},
		{{
			"$addFields",
			bson.M{
				"date": bson.M{
					"$dateToString": bson.M{
						"date":     bson.M{"$toDate": "$_id"},
						"format":   "%Y-%m-%d",
						"timezone": "Asia/Shanghai", // TODO: parameterization
					},
				},
			},
		}},
		{{
			"$group",
			bson.M{
				"_id":     "$date",
				"tasks":   bson.M{"$sum": 1},
				"results": bson.M{"$sum": "$result_count"},
			},
		}},
		{{
			"$sort",
			bson.D{{"_id", 1}},
		}},
	}
	var results []entity.StatsDailyItem
	if err := mongo.GetMongoCol(interfaces.ModelColNameTaskStat).Aggregate(pipeline, nil).All(&results); err != nil {
		return nil, err
	}
	return results, nil
}

func (svc *Service) getOverviewResults(query bson.M) (data interface{}, err error) {
	pipeline := mongo2.Pipeline{
		{{"$match", query}},
		{{
			"$group",
			bson.M{
				"_id":     nil,
				"results": bson.M{"$sum": "$result_count"},
			},
		}},
	}
	var res bson.M
	if err := mongo.GetMongoCol(interfaces.ModelColNameTaskStat).Aggregate(pipeline, nil).One(&res); err != nil {
		return nil, err
	}
	return res["results"], nil
}

func (svc *Service) getTaskStatsByStatus(query bson.M) (data interface{}, err error) {
	pipeline := mongo2.Pipeline{
		{{"$match", query}},
		{{
			"$group",
			bson.M{
				"_id":   "$status",
				"tasks": bson.M{"$sum": 1},
			},
		}},
		{{
			"$project",
			bson.M{
				"status": "$_id",
				"tasks":  "$tasks",
			},
		}},
	}
	var results []bson.M
	if err := mongo.GetMongoCol(interfaces.ModelColNameTask).Aggregate(pipeline, nil).All(&results); err != nil {
		return nil, err
	}
	return results, nil
}

func (svc *Service) getTaskStatsByNode(query bson.M) (data interface{}, err error) {
	pipeline := mongo2.Pipeline{
		{{"$match", query}},
		{{
			"$group",
			bson.M{
				"_id":   "$node_id",
				"tasks": bson.M{"$sum": 1},
			},
		}},
		{{
			"$lookup",
			bson.M{
				"from":         interfaces.ModelColNameNode,
				"localField":   "_id",
				"foreignField": "_id",
				"as":           "_n",
			},
		}},
		{{
			"$project",
			bson.M{
				"node_id":   "$node_id",
				"node":      bson.M{"$arrayElemAt": bson.A{"$_n", 0}},
				"node_name": bson.M{"$arrayElemAt": bson.A{"$_n.name", 0}},
				"tasks":     "$tasks",
			},
		}},
	}
	var results []bson.M
	if err := mongo.GetMongoCol(interfaces.ModelColNameTask).Aggregate(pipeline, nil).All(&results); err != nil {
		return nil, err
	}
	return results, nil
}

func (svc *Service) getTaskStatsBySpider(query bson.M) (data interface{}, err error) {
	pipeline := mongo2.Pipeline{
		{{"$match", query}},
		{{
			"$group",
			bson.M{
				"_id":   "$spider_id",
				"tasks": bson.M{"$sum": 1},
			},
		}},
		{{
			"$lookup",
			bson.M{
				"from":         interfaces.ModelColNameSpider,
				"localField":   "_id",
				"foreignField": "_id",
				"as":           "_s",
			},
		}},
		{{
			"$project",
			bson.M{
				"spider_id":   "$spider_id",
				"spider":      bson.M{"$arrayElemAt": bson.A{"$_s", 0}},
				"spider_name": bson.M{"$arrayElemAt": bson.A{"$_s.name", 0}},
				"tasks":       "$tasks",
			},
		}},
		{{"$limit", 10}},
	}
	var results []bson.M
	if err := mongo.GetMongoCol(interfaces.ModelColNameTask).Aggregate(pipeline, nil).All(&results); err != nil {
		return nil, err
	}
	return results, nil
}

func (svc *Service) getTaskStatsHistogram(query bson.M) (data interface{}, err error) {
	pipeline := mongo2.Pipeline{
		{{"$match", query}},
		{{
			"$lookup",
			bson.M{
				"from":         interfaces.ModelColNameTaskStat,
				"localField":   "_id",
				"foreignField": "_id",
				"as":           "_ts",
			},
		}},
		{{
			"$facet",
			bson.M{
				"total_duration": bson.A{
					bson.M{
						"$bucketAuto": bson.M{
							"groupBy":     "$_ts.td",
							"buckets":     10,
							"granularity": "1-2-5",
						},
					},
				},
			},
		}},
	}
	var res bson.M
	if err := mongo.GetMongoCol(interfaces.ModelColNameTask).Aggregate(pipeline, nil).One(&res); err != nil {
		return nil, err
	}
	return res, nil
}

var svc interfaces.StatsService

func GetStatsService() interfaces.StatsService {
	if svc != nil {
		return svc
	}

	// service
	svc = &Service{}

	return svc
}
