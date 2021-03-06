// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1-devel
// 	protoc        v3.19.4
// source: web_service.proto

package api_v1

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

type GetAlertRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AlertDetails string `protobuf:"bytes,1,opt,name=alert_details,json=alertDetails,proto3" json:"alert_details,omitempty"`
}

func (x *GetAlertRequest) Reset() {
	*x = GetAlertRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_web_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAlertRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAlertRequest) ProtoMessage() {}

func (x *GetAlertRequest) ProtoReflect() protoreflect.Message {
	mi := &file_web_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAlertRequest.ProtoReflect.Descriptor instead.
func (*GetAlertRequest) Descriptor() ([]byte, []int) {
	return file_web_service_proto_rawDescGZIP(), []int{0}
}

func (x *GetAlertRequest) GetAlertDetails() string {
	if x != nil {
		return x.AlertDetails
	}
	return ""
}

type GetAlertResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Alert map[string]string `protobuf:"bytes,1,rep,name=alert,proto3" json:"alert,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *GetAlertResponse) Reset() {
	*x = GetAlertResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_web_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAlertResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAlertResponse) ProtoMessage() {}

func (x *GetAlertResponse) ProtoReflect() protoreflect.Message {
	mi := &file_web_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAlertResponse.ProtoReflect.Descriptor instead.
func (*GetAlertResponse) Descriptor() ([]byte, []int) {
	return file_web_service_proto_rawDescGZIP(), []int{1}
}

func (x *GetAlertResponse) GetAlert() map[string]string {
	if x != nil {
		return x.Alert
	}
	return nil
}

var File_web_service_proto protoreflect.FileDescriptor

var file_web_service_proto_rawDesc = []byte{
	0x0a, 0x11, 0x77, 0x65, 0x62, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x62, 0x75, 0x74, 0x6c, 0x65, 0x72, 0x2e, 0x61, 0x70, 0x69, 0x5f,
	0x76, 0x31, 0x22, 0x36, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x65, 0x72, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x23, 0x0a, 0x0d, 0x61, 0x6c, 0x65, 0x72, 0x74, 0x5f, 0x64,
	0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x61, 0x6c,
	0x65, 0x72, 0x74, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x22, 0x8e, 0x01, 0x0a, 0x10, 0x47,
	0x65, 0x74, 0x41, 0x6c, 0x65, 0x72, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x40, 0x0a, 0x05, 0x61, 0x6c, 0x65, 0x72, 0x74, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2a,
	0x2e, 0x62, 0x75, 0x74, 0x6c, 0x65, 0x72, 0x2e, 0x61, 0x70, 0x69, 0x5f, 0x76, 0x31, 0x2e, 0x47,
	0x65, 0x74, 0x41, 0x6c, 0x65, 0x72, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e,
	0x41, 0x6c, 0x65, 0x72, 0x74, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x05, 0x61, 0x6c, 0x65, 0x72,
	0x74, 0x1a, 0x38, 0x0a, 0x0a, 0x41, 0x6c, 0x65, 0x72, 0x74, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12,
	0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65,
	0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x32, 0x5b, 0x0a, 0x0a, 0x57,
	0x65, 0x62, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x4d, 0x0a, 0x08, 0x47, 0x65, 0x74,
	0x41, 0x6c, 0x65, 0x72, 0x74, 0x12, 0x1e, 0x2e, 0x62, 0x75, 0x74, 0x6c, 0x65, 0x72, 0x2e, 0x61,
	0x70, 0x69, 0x5f, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x65, 0x72, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x62, 0x75, 0x74, 0x6c, 0x65, 0x72, 0x2e, 0x61,
	0x70, 0x69, 0x5f, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x65, 0x72, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x0a, 0x5a, 0x08, 0x2e, 0x2f, 0x61, 0x70,
	0x69, 0x5f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_web_service_proto_rawDescOnce sync.Once
	file_web_service_proto_rawDescData = file_web_service_proto_rawDesc
)

func file_web_service_proto_rawDescGZIP() []byte {
	file_web_service_proto_rawDescOnce.Do(func() {
		file_web_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_web_service_proto_rawDescData)
	})
	return file_web_service_proto_rawDescData
}

var file_web_service_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_web_service_proto_goTypes = []interface{}{
	(*GetAlertRequest)(nil),  // 0: butler.api_v1.GetAlertRequest
	(*GetAlertResponse)(nil), // 1: butler.api_v1.GetAlertResponse
	nil,                      // 2: butler.api_v1.GetAlertResponse.AlertEntry
}
var file_web_service_proto_depIdxs = []int32{
	2, // 0: butler.api_v1.GetAlertResponse.alert:type_name -> butler.api_v1.GetAlertResponse.AlertEntry
	0, // 1: butler.api_v1.WebService.GetAlert:input_type -> butler.api_v1.GetAlertRequest
	1, // 2: butler.api_v1.WebService.GetAlert:output_type -> butler.api_v1.GetAlertResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_web_service_proto_init() }
func file_web_service_proto_init() {
	if File_web_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_web_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAlertRequest); i {
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
		file_web_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAlertResponse); i {
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
			RawDescriptor: file_web_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_web_service_proto_goTypes,
		DependencyIndexes: file_web_service_proto_depIdxs,
		MessageInfos:      file_web_service_proto_msgTypes,
	}.Build()
	File_web_service_proto = out.File
	file_web_service_proto_rawDesc = nil
	file_web_service_proto_goTypes = nil
	file_web_service_proto_depIdxs = nil
}
