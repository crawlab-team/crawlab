// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.2
// source: services/metrics_service_v2.proto

package grpc

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type MetricsServiceV2SendRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type                 string  `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	NodeKey              string  `protobuf:"bytes,2,opt,name=node_key,json=nodeKey,proto3" json:"node_key,omitempty"`
	Timestamp            int64   `protobuf:"varint,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	CpuUsagePercent      float32 `protobuf:"fixed32,4,opt,name=cpu_usage_percent,json=cpuUsagePercent,proto3" json:"cpu_usage_percent,omitempty"`
	TotalMemory          uint64  `protobuf:"varint,5,opt,name=total_memory,json=totalMemory,proto3" json:"total_memory,omitempty"`
	AvailableMemory      uint64  `protobuf:"varint,6,opt,name=available_memory,json=availableMemory,proto3" json:"available_memory,omitempty"`
	UsedMemory           uint64  `protobuf:"varint,7,opt,name=used_memory,json=usedMemory,proto3" json:"used_memory,omitempty"`
	UsedMemoryPercent    float32 `protobuf:"fixed32,8,opt,name=used_memory_percent,json=usedMemoryPercent,proto3" json:"used_memory_percent,omitempty"`
	TotalDisk            uint64  `protobuf:"varint,9,opt,name=total_disk,json=totalDisk,proto3" json:"total_disk,omitempty"`
	AvailableDisk        uint64  `protobuf:"varint,10,opt,name=available_disk,json=availableDisk,proto3" json:"available_disk,omitempty"`
	UsedDisk             uint64  `protobuf:"varint,11,opt,name=used_disk,json=usedDisk,proto3" json:"used_disk,omitempty"`
	UsedDiskPercent      float32 `protobuf:"fixed32,12,opt,name=used_disk_percent,json=usedDiskPercent,proto3" json:"used_disk_percent,omitempty"`
	DiskReadBytesRate    float32 `protobuf:"fixed32,15,opt,name=disk_read_bytes_rate,json=diskReadBytesRate,proto3" json:"disk_read_bytes_rate,omitempty"`
	DiskWriteBytesRate   float32 `protobuf:"fixed32,16,opt,name=disk_write_bytes_rate,json=diskWriteBytesRate,proto3" json:"disk_write_bytes_rate,omitempty"`
	NetworkBytesSentRate float32 `protobuf:"fixed32,17,opt,name=network_bytes_sent_rate,json=networkBytesSentRate,proto3" json:"network_bytes_sent_rate,omitempty"`
	NetworkBytesRecvRate float32 `protobuf:"fixed32,18,opt,name=network_bytes_recv_rate,json=networkBytesRecvRate,proto3" json:"network_bytes_recv_rate,omitempty"`
}

func (x *MetricsServiceV2SendRequest) Reset() {
	*x = MetricsServiceV2SendRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_metrics_service_v2_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MetricsServiceV2SendRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MetricsServiceV2SendRequest) ProtoMessage() {}

func (x *MetricsServiceV2SendRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_metrics_service_v2_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MetricsServiceV2SendRequest.ProtoReflect.Descriptor instead.
func (*MetricsServiceV2SendRequest) Descriptor() ([]byte, []int) {
	return file_services_metrics_service_v2_proto_rawDescGZIP(), []int{0}
}

func (x *MetricsServiceV2SendRequest) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *MetricsServiceV2SendRequest) GetNodeKey() string {
	if x != nil {
		return x.NodeKey
	}
	return ""
}

func (x *MetricsServiceV2SendRequest) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *MetricsServiceV2SendRequest) GetCpuUsagePercent() float32 {
	if x != nil {
		return x.CpuUsagePercent
	}
	return 0
}

func (x *MetricsServiceV2SendRequest) GetTotalMemory() uint64 {
	if x != nil {
		return x.TotalMemory
	}
	return 0
}

func (x *MetricsServiceV2SendRequest) GetAvailableMemory() uint64 {
	if x != nil {
		return x.AvailableMemory
	}
	return 0
}

func (x *MetricsServiceV2SendRequest) GetUsedMemory() uint64 {
	if x != nil {
		return x.UsedMemory
	}
	return 0
}

func (x *MetricsServiceV2SendRequest) GetUsedMemoryPercent() float32 {
	if x != nil {
		return x.UsedMemoryPercent
	}
	return 0
}

func (x *MetricsServiceV2SendRequest) GetTotalDisk() uint64 {
	if x != nil {
		return x.TotalDisk
	}
	return 0
}

func (x *MetricsServiceV2SendRequest) GetAvailableDisk() uint64 {
	if x != nil {
		return x.AvailableDisk
	}
	return 0
}

func (x *MetricsServiceV2SendRequest) GetUsedDisk() uint64 {
	if x != nil {
		return x.UsedDisk
	}
	return 0
}

func (x *MetricsServiceV2SendRequest) GetUsedDiskPercent() float32 {
	if x != nil {
		return x.UsedDiskPercent
	}
	return 0
}

func (x *MetricsServiceV2SendRequest) GetDiskReadBytesRate() float32 {
	if x != nil {
		return x.DiskReadBytesRate
	}
	return 0
}

func (x *MetricsServiceV2SendRequest) GetDiskWriteBytesRate() float32 {
	if x != nil {
		return x.DiskWriteBytesRate
	}
	return 0
}

func (x *MetricsServiceV2SendRequest) GetNetworkBytesSentRate() float32 {
	if x != nil {
		return x.NetworkBytesSentRate
	}
	return 0
}

func (x *MetricsServiceV2SendRequest) GetNetworkBytesRecvRate() float32 {
	if x != nil {
		return x.NetworkBytesRecvRate
	}
	return 0
}

var File_services_metrics_service_v2_proto protoreflect.FileDescriptor

var file_services_metrics_service_v2_proto_rawDesc = []byte{
	0x0a, 0x21, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x6d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x73, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x76, 0x32, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x04, 0x67, 0x72, 0x70, 0x63, 0x1a, 0x15, 0x65, 0x6e, 0x74, 0x69, 0x74,
	0x79, 0x2f, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x96, 0x05, 0x0a, 0x1b, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x56, 0x32, 0x53, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x6b, 0x65, 0x79,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6e, 0x6f, 0x64, 0x65, 0x4b, 0x65, 0x79, 0x12,
	0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x2a, 0x0a,
	0x11, 0x63, 0x70, 0x75, 0x5f, 0x75, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x70, 0x65, 0x72, 0x63, 0x65,
	0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x02, 0x52, 0x0f, 0x63, 0x70, 0x75, 0x55, 0x73, 0x61,
	0x67, 0x65, 0x50, 0x65, 0x72, 0x63, 0x65, 0x6e, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x74, 0x6f, 0x74,
	0x61, 0x6c, 0x5f, 0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x0b, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x4d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x12, 0x29, 0x0a, 0x10,
	0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x5f, 0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0f, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c,
	0x65, 0x4d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x12, 0x1f, 0x0a, 0x0b, 0x75, 0x73, 0x65, 0x64, 0x5f,
	0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x18, 0x07, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0a, 0x75, 0x73,
	0x65, 0x64, 0x4d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x12, 0x2e, 0x0a, 0x13, 0x75, 0x73, 0x65, 0x64,
	0x5f, 0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x5f, 0x70, 0x65, 0x72, 0x63, 0x65, 0x6e, 0x74, 0x18,
	0x08, 0x20, 0x01, 0x28, 0x02, 0x52, 0x11, 0x75, 0x73, 0x65, 0x64, 0x4d, 0x65, 0x6d, 0x6f, 0x72,
	0x79, 0x50, 0x65, 0x72, 0x63, 0x65, 0x6e, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x74, 0x6f, 0x74, 0x61,
	0x6c, 0x5f, 0x64, 0x69, 0x73, 0x6b, 0x18, 0x09, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x74, 0x6f,
	0x74, 0x61, 0x6c, 0x44, 0x69, 0x73, 0x6b, 0x12, 0x25, 0x0a, 0x0e, 0x61, 0x76, 0x61, 0x69, 0x6c,
	0x61, 0x62, 0x6c, 0x65, 0x5f, 0x64, 0x69, 0x73, 0x6b, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x0d, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x44, 0x69, 0x73, 0x6b, 0x12, 0x1b,
	0x0a, 0x09, 0x75, 0x73, 0x65, 0x64, 0x5f, 0x64, 0x69, 0x73, 0x6b, 0x18, 0x0b, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x08, 0x75, 0x73, 0x65, 0x64, 0x44, 0x69, 0x73, 0x6b, 0x12, 0x2a, 0x0a, 0x11, 0x75,
	0x73, 0x65, 0x64, 0x5f, 0x64, 0x69, 0x73, 0x6b, 0x5f, 0x70, 0x65, 0x72, 0x63, 0x65, 0x6e, 0x74,
	0x18, 0x0c, 0x20, 0x01, 0x28, 0x02, 0x52, 0x0f, 0x75, 0x73, 0x65, 0x64, 0x44, 0x69, 0x73, 0x6b,
	0x50, 0x65, 0x72, 0x63, 0x65, 0x6e, 0x74, 0x12, 0x2f, 0x0a, 0x14, 0x64, 0x69, 0x73, 0x6b, 0x5f,
	0x72, 0x65, 0x61, 0x64, 0x5f, 0x62, 0x79, 0x74, 0x65, 0x73, 0x5f, 0x72, 0x61, 0x74, 0x65, 0x18,
	0x0f, 0x20, 0x01, 0x28, 0x02, 0x52, 0x11, 0x64, 0x69, 0x73, 0x6b, 0x52, 0x65, 0x61, 0x64, 0x42,
	0x79, 0x74, 0x65, 0x73, 0x52, 0x61, 0x74, 0x65, 0x12, 0x31, 0x0a, 0x15, 0x64, 0x69, 0x73, 0x6b,
	0x5f, 0x77, 0x72, 0x69, 0x74, 0x65, 0x5f, 0x62, 0x79, 0x74, 0x65, 0x73, 0x5f, 0x72, 0x61, 0x74,
	0x65, 0x18, 0x10, 0x20, 0x01, 0x28, 0x02, 0x52, 0x12, 0x64, 0x69, 0x73, 0x6b, 0x57, 0x72, 0x69,
	0x74, 0x65, 0x42, 0x79, 0x74, 0x65, 0x73, 0x52, 0x61, 0x74, 0x65, 0x12, 0x35, 0x0a, 0x17, 0x6e,
	0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x5f, 0x62, 0x79, 0x74, 0x65, 0x73, 0x5f, 0x73, 0x65, 0x6e,
	0x74, 0x5f, 0x72, 0x61, 0x74, 0x65, 0x18, 0x11, 0x20, 0x01, 0x28, 0x02, 0x52, 0x14, 0x6e, 0x65,
	0x74, 0x77, 0x6f, 0x72, 0x6b, 0x42, 0x79, 0x74, 0x65, 0x73, 0x53, 0x65, 0x6e, 0x74, 0x52, 0x61,
	0x74, 0x65, 0x12, 0x35, 0x0a, 0x17, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x5f, 0x62, 0x79,
	0x74, 0x65, 0x73, 0x5f, 0x72, 0x65, 0x63, 0x76, 0x5f, 0x72, 0x61, 0x74, 0x65, 0x18, 0x12, 0x20,
	0x01, 0x28, 0x02, 0x52, 0x14, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x42, 0x79, 0x74, 0x65,
	0x73, 0x52, 0x65, 0x63, 0x76, 0x52, 0x61, 0x74, 0x65, 0x32, 0x4f, 0x0a, 0x10, 0x4d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x56, 0x32, 0x12, 0x3b, 0x0a,
	0x04, 0x53, 0x65, 0x6e, 0x64, 0x12, 0x21, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x4d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x56, 0x32, 0x53, 0x65, 0x6e,
	0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0e, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x08, 0x5a, 0x06, 0x2e, 0x3b,
	0x67, 0x72, 0x70, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_services_metrics_service_v2_proto_rawDescOnce sync.Once
	file_services_metrics_service_v2_proto_rawDescData = file_services_metrics_service_v2_proto_rawDesc
)

func file_services_metrics_service_v2_proto_rawDescGZIP() []byte {
	file_services_metrics_service_v2_proto_rawDescOnce.Do(func() {
		file_services_metrics_service_v2_proto_rawDescData = protoimpl.X.CompressGZIP(file_services_metrics_service_v2_proto_rawDescData)
	})
	return file_services_metrics_service_v2_proto_rawDescData
}

var file_services_metrics_service_v2_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_services_metrics_service_v2_proto_goTypes = []any{
	(*MetricsServiceV2SendRequest)(nil), // 0: grpc.MetricsServiceV2SendRequest
	(*Response)(nil),                    // 1: grpc.Response
}
var file_services_metrics_service_v2_proto_depIdxs = []int32{
	0, // 0: grpc.MetricsServiceV2.Send:input_type -> grpc.MetricsServiceV2SendRequest
	1, // 1: grpc.MetricsServiceV2.Send:output_type -> grpc.Response
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_services_metrics_service_v2_proto_init() }
func file_services_metrics_service_v2_proto_init() {
	if File_services_metrics_service_v2_proto != nil {
		return
	}
	file_entity_response_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_services_metrics_service_v2_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*MetricsServiceV2SendRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_services_metrics_service_v2_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_services_metrics_service_v2_proto_goTypes,
		DependencyIndexes: file_services_metrics_service_v2_proto_depIdxs,
		MessageInfos:      file_services_metrics_service_v2_proto_msgTypes,
	}.Build()
	File_services_metrics_service_v2_proto = out.File
	file_services_metrics_service_v2_proto_rawDesc = nil
	file_services_metrics_service_v2_proto_goTypes = nil
	file_services_metrics_service_v2_proto_depIdxs = nil
}
