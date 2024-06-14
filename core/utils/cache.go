package utils

import (
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"time"
)

func GetFromDbCache(key string, getFn func() (string, error)) (res string, err error) {
	col := mongo.GetMongoCol(constants.CacheColName)

	var d bson.M
	if err := col.Find(bson.M{
		constants.CacheColKey: key,
	}, nil).One(&d); err != nil {
		if err != mongo2.ErrNoDocuments {
			return "", err
		}

		// get cache value
		res, err = getFn()
		if err != nil {
			return "", err
		}

		// save cache
		d = bson.M{
			constants.CacheColKey:   key,
			constants.CacheColValue: res,
			constants.CacheColTime:  time.Now(),
		}
		if _, err := col.Insert(d); err != nil {
			return "", err
		}
		return res, nil
	}

	// type conversion
	r, ok := d[constants.CacheColValue]
	if !ok {
		if err := col.Delete(bson.M{constants.CacheColKey: key}); err != nil {
			return "", err
		}
		return GetFromDbCache(key, getFn)
	}
	res, ok = r.(string)
	if !ok {
		if err := col.Delete(bson.M{constants.CacheColKey: key}); err != nil {
			return "", err
		}
		return GetFromDbCache(key, getFn)
	}

	return res, nil
}
