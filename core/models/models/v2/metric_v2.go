package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MetricV2 struct {
	any                   `collection:"metrics"`
	BaseModelV2[MetricV2] `bson:",inline"`
	Type                  string             `json:"type" bson:"type"`
	NodeId                primitive.ObjectID `json:"node_id" bson:"node_id"`
	CpuUsagePercent       float32            `json:"cpu_usage_percent" bson:"cpu_usage_percent"`
	TotalMemory           uint64             `json:"total_memory" bson:"total_memory"`
	AvailableMemory       uint64             `json:"available_memory" bson:"available_memory"`
	UsedMemory            uint64             `json:"used_memory" bson:"used_memory"`
	UsedMemoryPercent     float32            `json:"used_memory_percent" bson:"used_memory_percent"`
	TotalDisk             uint64             `json:"total_disk" bson:"total_disk"`
	AvailableDisk         uint64             `json:"available_disk" bson:"available_disk"`
	UsedDisk              uint64             `json:"used_disk" bson:"used_disk"`
	UsedDiskPercent       float32            `json:"used_disk_percent" bson:"used_disk_percent"`
	DiskReadBytesRate     float32            `json:"disk_read_bytes_rate" bson:"disk_read_bytes_rate"`
	DiskWriteBytesRate    float32            `json:"disk_write_bytes_rate" bson:"disk_write_bytes_rate"`
	NetworkBytesSentRate  float32            `json:"network_bytes_sent_rate" bson:"network_bytes_sent_rate"`
	NetworkBytesRecvRate  float32            `json:"network_bytes_recv_rate" bson:"network_bytes_recv_rate"`
}
