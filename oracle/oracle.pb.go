// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.6.1
// source: oracle.proto

package oracle

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

type Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []byte `protobuf:"bytes,1,opt,name=Data,proto3" json:"Data,omitempty"`
	Type int32  `protobuf:"varint,2,opt,name=Type,proto3" json:"Type,omitempty"`
	Id   uint64 `protobuf:"varint,3,opt,name=Id,proto3" json:"Id,omitempty"`
}

func (x *Response) Reset() {
	*x = Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_oracle_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_oracle_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response.ProtoReflect.Descriptor instead.
func (*Response) Descriptor() ([]byte, []int) {
	return file_oracle_proto_rawDescGZIP(), []int{0}
}

func (x *Response) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *Response) GetType() int32 {
	if x != nil {
		return x.Type
	}
	return 0
}

func (x *Response) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type ReqOPP struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Epoch uint64 `protobuf:"varint,1,opt,name=epoch,proto3" json:"epoch,omitempty"`
}

func (x *ReqOPP) Reset() {
	*x = ReqOPP{}
	if protoimpl.UnsafeEnabled {
		mi := &file_oracle_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReqOPP) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReqOPP) ProtoMessage() {}

func (x *ReqOPP) ProtoReflect() protoreflect.Message {
	mi := &file_oracle_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReqOPP.ProtoReflect.Descriptor instead.
func (*ReqOPP) Descriptor() ([]byte, []int) {
	return file_oracle_proto_rawDescGZIP(), []int{1}
}

func (x *ReqOPP) GetEpoch() uint64 {
	if x != nil {
		return x.Epoch
	}
	return 0
}

type ResOPP struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Proof  []byte `protobuf:"bytes,1,opt,name=Proof,proto3" json:"Proof,omitempty"`
	VRF    []byte `protobuf:"bytes,2,opt,name=VRF,proto3" json:"VRF,omitempty"`
	Pubkey []byte `protobuf:"bytes,3,opt,name=Pubkey,proto3" json:"Pubkey,omitempty"`
	Weight uint64 `protobuf:"varint,4,opt,name=weight,proto3" json:"weight,omitempty"`
}

func (x *ResOPP) Reset() {
	*x = ResOPP{}
	if protoimpl.UnsafeEnabled {
		mi := &file_oracle_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResOPP) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResOPP) ProtoMessage() {}

func (x *ResOPP) ProtoReflect() protoreflect.Message {
	mi := &file_oracle_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResOPP.ProtoReflect.Descriptor instead.
func (*ResOPP) Descriptor() ([]byte, []int) {
	return file_oracle_proto_rawDescGZIP(), []int{2}
}

func (x *ResOPP) GetProof() []byte {
	if x != nil {
		return x.Proof
	}
	return nil
}

func (x *ResOPP) GetVRF() []byte {
	if x != nil {
		return x.VRF
	}
	return nil
}

func (x *ResOPP) GetPubkey() []byte {
	if x != nil {
		return x.Pubkey
	}
	return nil
}

func (x *ResOPP) GetWeight() uint64 {
	if x != nil {
		return x.Weight
	}
	return 0
}

var File_oracle_proto protoreflect.FileDescriptor

var file_oracle_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x6f, 0x72, 0x61, 0x63, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x42,
	0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x44, 0x61,
	0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x44, 0x61, 0x74, 0x61, 0x12, 0x12,
	0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x54, 0x79,
	0x70, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02,
	0x49, 0x64, 0x22, 0x1e, 0x0a, 0x06, 0x52, 0x65, 0x71, 0x4f, 0x50, 0x50, 0x12, 0x14, 0x0a, 0x05,
	0x65, 0x70, 0x6f, 0x63, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x65, 0x70, 0x6f,
	0x63, 0x68, 0x22, 0x60, 0x0a, 0x06, 0x52, 0x65, 0x73, 0x4f, 0x50, 0x50, 0x12, 0x14, 0x0a, 0x05,
	0x50, 0x72, 0x6f, 0x6f, 0x66, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x50, 0x72, 0x6f,
	0x6f, 0x66, 0x12, 0x10, 0x0a, 0x03, 0x56, 0x52, 0x46, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x03, 0x56, 0x52, 0x46, 0x12, 0x16, 0x0a, 0x06, 0x50, 0x75, 0x62, 0x6b, 0x65, 0x79, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x50, 0x75, 0x62, 0x6b, 0x65, 0x79, 0x12, 0x16, 0x0a, 0x06,
	0x77, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x77, 0x65,
	0x69, 0x67, 0x68, 0x74, 0x32, 0x2b, 0x0a, 0x0a, 0x52, 0x50, 0x43, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x1d, 0x0a, 0x07, 0x73, 0x65, 0x6e, 0x64, 0x4f, 0x50, 0x50, 0x12, 0x07, 0x2e,
	0x52, 0x65, 0x71, 0x4f, 0x50, 0x50, 0x1a, 0x07, 0x2e, 0x52, 0x65, 0x73, 0x4f, 0x50, 0x50, 0x22,
	0x00, 0x42, 0x09, 0x5a, 0x07, 0x2f, 0x6f, 0x72, 0x61, 0x63, 0x6c, 0x65, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_oracle_proto_rawDescOnce sync.Once
	file_oracle_proto_rawDescData = file_oracle_proto_rawDesc
)

func file_oracle_proto_rawDescGZIP() []byte {
	file_oracle_proto_rawDescOnce.Do(func() {
		file_oracle_proto_rawDescData = protoimpl.X.CompressGZIP(file_oracle_proto_rawDescData)
	})
	return file_oracle_proto_rawDescData
}

var file_oracle_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_oracle_proto_goTypes = []interface{}{
	(*Response)(nil), // 0: Response
	(*ReqOPP)(nil),   // 1: ReqOPP
	(*ResOPP)(nil),   // 2: ResOPP
}
var file_oracle_proto_depIdxs = []int32{
	1, // 0: RPCService.sendOPP:input_type -> ReqOPP
	2, // 1: RPCService.sendOPP:output_type -> ResOPP
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_oracle_proto_init() }
func file_oracle_proto_init() {
	if File_oracle_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_oracle_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Response); i {
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
		file_oracle_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReqOPP); i {
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
		file_oracle_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResOPP); i {
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
			RawDescriptor: file_oracle_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_oracle_proto_goTypes,
		DependencyIndexes: file_oracle_proto_depIdxs,
		MessageInfos:      file_oracle_proto_msgTypes,
	}.Build()
	File_oracle_proto = out.File
	file_oracle_proto_rawDesc = nil
	file_oracle_proto_goTypes = nil
	file_oracle_proto_depIdxs = nil
}
