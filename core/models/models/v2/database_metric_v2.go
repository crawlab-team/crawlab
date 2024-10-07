package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type DatabaseMetricV2 struct {
	any                           `collection:"database_metrics"`
	BaseModelV2[DatabaseMetricV2] `bson:",inline"`
	DatabaseId                    primitive.ObjectID `json:"database_id" bson:"database_id"`
	CpuUsagePercent               float32            `json:"cpu_usage_percent" bson:"cpu_usage_percent"`
	TotalMemory                   uint64             `json:"total_memory" bson:"total_memory"`
	AvailableMemory               uint64             `json:"available_memory" bson:"available_memory"`
	UsedMemory                    uint64             `json:"used_memory" bson:"used_memory"`
	UsedMemoryPercent             float32            `json:"used_memory_percent" bson:"used_memory_percent"`
	TotalDisk                     uint64             `json:"total_disk" bson:"total_disk"`
	AvailableDisk                 uint64             `json:"available_disk" bson:"available_disk"`
	UsedDisk                      uint64             `json:"used_disk" bson:"used_disk"`
	UsedDiskPercent               float32            `json:"used_disk_percent" bson:"used_disk_percent"`
	Connections                   int                `json:"connections" bson:"connections"`
	QueryPerSecond                float64            `json:"query_per_second" bson:"query_per_second"`
	TotalQuery                    uint64             `json:"total_query,omitempty" bson:"total_query,omitempty"`
	CacheHitRatio                 float64            `json:"cache_hit_ratio" bson:"cache_hit_ratio"`
	ReplicationLag                float64            `json:"replication_lag" bson:"replication_lag"`
	LockWaitTime                  float64            `json:"lock_wait_time" bson:"lock_wait_time"`
}
