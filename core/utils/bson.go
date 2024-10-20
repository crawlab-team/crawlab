package utils

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NormalizeBsonMObjectId(m bson.M) (res bson.M) {
	for k, v := range m {
		switch v.(type) {
		case string:
			oid, err := primitive.ObjectIDFromHex(v.(string))
			if err == nil {
				m[k] = oid
			}
		case bson.M:
			m[k] = NormalizeBsonMObjectId(v.(bson.M))
		}
	}
	return m
}

func NormalizeObjectId(v interface{}) (res interface{}) {
	switch v.(type) {
	case string:
		oid, err := primitive.ObjectIDFromHex(v.(string))
		if err != nil {
			return v
		}
		return oid
	default:
		return v
	}
}
