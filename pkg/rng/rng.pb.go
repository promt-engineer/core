// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.21.9
// source: pkg/rng/rng.proto

package rng

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

type Status struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status string `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *Status) Reset() {
	*x = Status{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_rng_rng_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Status) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Status) ProtoMessage() {}

func (x *Status) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_rng_rng_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Status.ProtoReflect.Descriptor instead.
func (*Status) Descriptor() ([]byte, []int) {
	return file_pkg_rng_rng_proto_rawDescGZIP(), []int{0}
}

func (x *Status) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

type RandRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Max []uint64 `protobuf:"varint,1,rep,packed,name=max,proto3" json:"max,omitempty"`
}

func (x *RandRequest) Reset() {
	*x = RandRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_rng_rng_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RandRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RandRequest) ProtoMessage() {}

func (x *RandRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_rng_rng_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RandRequest.ProtoReflect.Descriptor instead.
func (*RandRequest) Descriptor() ([]byte, []int) {
	return file_pkg_rng_rng_proto_rawDescGZIP(), []int{1}
}

func (x *RandRequest) GetMax() []uint64 {
	if x != nil {
		return x.Max
	}
	return nil
}

type RandResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result []uint64 `protobuf:"varint,1,rep,packed,name=result,proto3" json:"result,omitempty"`
}

func (x *RandResponse) Reset() {
	*x = RandResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_rng_rng_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RandResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RandResponse) ProtoMessage() {}

func (x *RandResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_rng_rng_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RandResponse.ProtoReflect.Descriptor instead.
func (*RandResponse) Descriptor() ([]byte, []int) {
	return file_pkg_rng_rng_proto_rawDescGZIP(), []int{2}
}

func (x *RandResponse) GetResult() []uint64 {
	if x != nil {
		return x.Result
	}
	return nil
}

type RandRequestFloat struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Max uint64 `protobuf:"varint,1,opt,name=max,proto3" json:"max,omitempty"`
}

func (x *RandRequestFloat) Reset() {
	*x = RandRequestFloat{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_rng_rng_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RandRequestFloat) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RandRequestFloat) ProtoMessage() {}

func (x *RandRequestFloat) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_rng_rng_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RandRequestFloat.ProtoReflect.Descriptor instead.
func (*RandRequestFloat) Descriptor() ([]byte, []int) {
	return file_pkg_rng_rng_proto_rawDescGZIP(), []int{3}
}

func (x *RandRequestFloat) GetMax() uint64 {
	if x != nil {
		return x.Max
	}
	return 0
}

type RandResponseFloat struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result []float64 `protobuf:"fixed64,1,rep,packed,name=result,proto3" json:"result,omitempty"`
}

func (x *RandResponseFloat) Reset() {
	*x = RandResponseFloat{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_rng_rng_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RandResponseFloat) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RandResponseFloat) ProtoMessage() {}

func (x *RandResponseFloat) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_rng_rng_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RandResponseFloat.ProtoReflect.Descriptor instead.
func (*RandResponseFloat) Descriptor() ([]byte, []int) {
	return file_pkg_rng_rng_proto_rawDescGZIP(), []int{4}
}

func (x *RandResponseFloat) GetResult() []float64 {
	if x != nil {
		return x.Result
	}
	return nil
}

var File_pkg_rng_rng_proto protoreflect.FileDescriptor

var file_pkg_rng_rng_proto_rawDesc = []byte{
	0x0a, 0x11, 0x70, 0x6b, 0x67, 0x2f, 0x72, 0x6e, 0x67, 0x2f, 0x72, 0x6e, 0x67, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x20, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x16, 0x0a,
	0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x1f, 0x0a, 0x0b, 0x52, 0x61, 0x6e, 0x64, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x61, 0x78, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x04, 0x52, 0x03, 0x6d, 0x61, 0x78, 0x22, 0x26, 0x0a, 0x0c, 0x52, 0x61, 0x6e, 0x64, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x04, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x24,
	0x0a, 0x10, 0x52, 0x61, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x46, 0x6c, 0x6f,
	0x61, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x61, 0x78, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x03, 0x6d, 0x61, 0x78, 0x22, 0x2b, 0x0a, 0x11, 0x52, 0x61, 0x6e, 0x64, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x46, 0x6c, 0x6f, 0x61, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x73,
	0x75, 0x6c, 0x74, 0x18, 0x01, 0x20, 0x03, 0x28, 0x01, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c,
	0x74, 0x32, 0x89, 0x01, 0x0a, 0x03, 0x52, 0x4e, 0x47, 0x12, 0x25, 0x0a, 0x04, 0x52, 0x61, 0x6e,
	0x64, 0x12, 0x0c, 0x2e, 0x52, 0x61, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x0d, 0x2e, 0x52, 0x61, 0x6e, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x12, 0x34, 0x0a, 0x09, 0x52, 0x61, 0x6e, 0x64, 0x46, 0x6c, 0x6f, 0x61, 0x74, 0x12, 0x11, 0x2e,
	0x52, 0x61, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x46, 0x6c, 0x6f, 0x61, 0x74,
	0x1a, 0x12, 0x2e, 0x52, 0x61, 0x6e, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x46,
	0x6c, 0x6f, 0x61, 0x74, 0x22, 0x00, 0x12, 0x25, 0x0a, 0x0b, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68,
	0x43, 0x68, 0x65, 0x63, 0x6b, 0x12, 0x07, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x1a, 0x07,
	0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x42, 0x07, 0x5a,
	0x05, 0x2e, 0x2f, 0x72, 0x6e, 0x67, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_rng_rng_proto_rawDescOnce sync.Once
	file_pkg_rng_rng_proto_rawDescData = file_pkg_rng_rng_proto_rawDesc
)

func file_pkg_rng_rng_proto_rawDescGZIP() []byte {
	file_pkg_rng_rng_proto_rawDescOnce.Do(func() {
		file_pkg_rng_rng_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_rng_rng_proto_rawDescData)
	})
	return file_pkg_rng_rng_proto_rawDescData
}

var file_pkg_rng_rng_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_pkg_rng_rng_proto_goTypes = []interface{}{
	(*Status)(nil),            // 0: Status
	(*RandRequest)(nil),       // 1: RandRequest
	(*RandResponse)(nil),      // 2: RandResponse
	(*RandRequestFloat)(nil),  // 3: RandRequestFloat
	(*RandResponseFloat)(nil), // 4: RandResponseFloat
}
var file_pkg_rng_rng_proto_depIdxs = []int32{
	1, // 0: RNG.Rand:input_type -> RandRequest
	3, // 1: RNG.RandFloat:input_type -> RandRequestFloat
	0, // 2: RNG.HealthCheck:input_type -> Status
	2, // 3: RNG.Rand:output_type -> RandResponse
	4, // 4: RNG.RandFloat:output_type -> RandResponseFloat
	0, // 5: RNG.HealthCheck:output_type -> Status
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_pkg_rng_rng_proto_init() }
func file_pkg_rng_rng_proto_init() {
	if File_pkg_rng_rng_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_rng_rng_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Status); i {
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
		file_pkg_rng_rng_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RandRequest); i {
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
		file_pkg_rng_rng_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RandResponse); i {
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
		file_pkg_rng_rng_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RandRequestFloat); i {
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
		file_pkg_rng_rng_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RandResponseFloat); i {
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
			RawDescriptor: file_pkg_rng_rng_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_rng_rng_proto_goTypes,
		DependencyIndexes: file_pkg_rng_rng_proto_depIdxs,
		MessageInfos:      file_pkg_rng_rng_proto_msgTypes,
	}.Build()
	File_pkg_rng_rng_proto = out.File
	file_pkg_rng_rng_proto_rawDesc = nil
	file_pkg_rng_rng_proto_goTypes = nil
	file_pkg_rng_rng_proto_depIdxs = nil
}
