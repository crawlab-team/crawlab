package entity

import (
	"encoding/json"
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/trace"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Result map[string]interface{}

func (r Result) Value() map[string]interface{} {
	return r
}

func (r Result) SetValue(key string, value interface{}) {
	r[key] = value
}

func (r Result) GetValue(key string) (value interface{}) {
	value, _ = r[key]
	return value
}

func (r Result) GetTaskId() (id primitive.ObjectID) {
	_tid, ok := r[constants.TaskKey]
	if !ok {
		return id
	}
	switch _tid.(type) {
	case string:
		oid, err := primitive.ObjectIDFromHex(_tid.(string))
		if err != nil {
			return id
		}
		return oid
	default:
		return id
	}
}

func (r Result) SetTaskId(id primitive.ObjectID) {
	r[constants.TaskKey] = id
}

func (r Result) DenormalizeObjectId() (res Result) {
	for k, v := range r {
		switch v.(type) {
		case primitive.ObjectID:
			r[k] = v.(primitive.ObjectID).Hex()
		case Result:
			r[k] = v.(Result).DenormalizeObjectId()
		}
	}
	return r
}

func (r Result) ToJSON() (res Result) {
	r = r.DenormalizeObjectId()
	for k, v := range r {
		switch v.(type) {
		case []byte:
			r[k] = string(v.([]byte))
		}
	}
	return r
}

func (r Result) Flatten() (res Result) {
	r = r.ToJSON()
	for k, v := range r {
		switch v.(type) {
		case string,
			bool,
			uint, uint8, uint16, uint32, uint64,
			int, int8, int16, int32, int64,
			float32, float64:
		default:
			bytes, err := json.Marshal(v)
			if err != nil {
				trace.PrintError(err)
				return nil
			}
			r[k] = string(bytes)
		}
	}
	return r
}

func (r Result) String() (s string) {
	return string(r.Bytes())
}

func (r Result) Bytes() (bytes []byte) {
	bytes, err := json.Marshal(r.ToJSON())
	if err != nil {
		return bytes
	}
	return bytes
}
