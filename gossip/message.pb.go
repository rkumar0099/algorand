// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.6.1
// source: message.proto

package gossip

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

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_message_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_message_proto_rawDescGZIP(), []int{0}
}

type NeighborReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Topic  string `protobuf:"bytes,1,opt,name=topic,proto3" json:"topic,omitempty"`
	NodeId string `protobuf:"bytes,2,opt,name=nodeId,proto3" json:"nodeId,omitempty"`
	MaxNum int32  `protobuf:"varint,3,opt,name=maxNum,proto3" json:"maxNum,omitempty"`
}

func (x *NeighborReq) Reset() {
	*x = NeighborReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NeighborReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NeighborReq) ProtoMessage() {}

func (x *NeighborReq) ProtoReflect() protoreflect.Message {
	mi := &file_message_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NeighborReq.ProtoReflect.Descriptor instead.
func (*NeighborReq) Descriptor() ([]byte, []int) {
	return file_message_proto_rawDescGZIP(), []int{1}
}

func (x *NeighborReq) GetTopic() string {
	if x != nil {
		return x.Topic
	}
	return ""
}

func (x *NeighborReq) GetNodeId() string {
	if x != nil {
		return x.NodeId
	}
	return ""
}

func (x *NeighborReq) GetMaxNum() int32 {
	if x != nil {
		return x.MaxNum
	}
	return 0
}

type NeighborRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Topic     string   `protobuf:"bytes,1,opt,name=topic,proto3" json:"topic,omitempty"`
	NodeId    string   `protobuf:"bytes,2,opt,name=nodeId,proto3" json:"nodeId,omitempty"`
	Neighbors []string `protobuf:"bytes,3,rep,name=neighbors,proto3" json:"neighbors,omitempty"`
}

func (x *NeighborRes) Reset() {
	*x = NeighborRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NeighborRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NeighborRes) ProtoMessage() {}

func (x *NeighborRes) ProtoReflect() protoreflect.Message {
	mi := &file_message_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NeighborRes.ProtoReflect.Descriptor instead.
func (*NeighborRes) Descriptor() ([]byte, []int) {
	return file_message_proto_rawDescGZIP(), []int{2}
}

func (x *NeighborRes) GetTopic() string {
	if x != nil {
		return x.Topic
	}
	return ""
}

func (x *NeighborRes) GetNodeId() string {
	if x != nil {
		return x.NodeId
	}
	return ""
}

func (x *NeighborRes) GetNeighbors() []string {
	if x != nil {
		return x.Neighbors
	}
	return nil
}

type GossipData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Topic   string `protobuf:"bytes,1,opt,name=topic,proto3" json:"topic,omitempty"`
	NodeId  string `protobuf:"bytes,2,opt,name=nodeId,proto3" json:"nodeId,omitempty"`
	Nonce   uint64 `protobuf:"varint,3,opt,name=nonce,proto3" json:"nonce,omitempty"`
	Payload []byte `protobuf:"bytes,4,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (x *GossipData) Reset() {
	*x = GossipData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GossipData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GossipData) ProtoMessage() {}

func (x *GossipData) ProtoReflect() protoreflect.Message {
	mi := &file_message_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GossipData.ProtoReflect.Descriptor instead.
func (*GossipData) Descriptor() ([]byte, []int) {
	return file_message_proto_rawDescGZIP(), []int{3}
}

func (x *GossipData) GetTopic() string {
	if x != nil {
		return x.Topic
	}
	return ""
}

func (x *GossipData) GetNodeId() string {
	if x != nil {
		return x.NodeId
	}
	return ""
}

func (x *GossipData) GetNonce() uint64 {
	if x != nil {
		return x.Nonce
	}
	return 0
}

func (x *GossipData) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

var File_message_proto protoreflect.FileDescriptor

var file_message_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x07, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x53, 0x0a, 0x0b, 0x4e, 0x65, 0x69, 0x67,
	0x68, 0x62, 0x6f, 0x72, 0x52, 0x65, 0x71, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x12, 0x16, 0x0a,
	0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6e,
	0x6f, 0x64, 0x65, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x6d, 0x61, 0x78, 0x4e, 0x75, 0x6d, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x6d, 0x61, 0x78, 0x4e, 0x75, 0x6d, 0x22, 0x59, 0x0a,
	0x0b, 0x4e, 0x65, 0x69, 0x67, 0x68, 0x62, 0x6f, 0x72, 0x52, 0x65, 0x73, 0x12, 0x14, 0x0a, 0x05,
	0x74, 0x6f, 0x70, 0x69, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x70,
	0x69, 0x63, 0x12, 0x16, 0x0a, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x6e, 0x65,
	0x69, 0x67, 0x68, 0x62, 0x6f, 0x72, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x09, 0x6e,
	0x65, 0x69, 0x67, 0x68, 0x62, 0x6f, 0x72, 0x73, 0x22, 0x6a, 0x0a, 0x0a, 0x47, 0x6f, 0x73, 0x73,
	0x69, 0x70, 0x44, 0x61, 0x74, 0x61, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x12, 0x16, 0x0a, 0x06,
	0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6e, 0x6f,
	0x64, 0x65, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x6e, 0x6f, 0x6e, 0x63, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x05, 0x6e, 0x6f, 0x6e, 0x63, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x61,
	0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x70, 0x61, 0x79,
	0x6c, 0x6f, 0x61, 0x64, 0x32, 0x55, 0x0a, 0x06, 0x47, 0x6f, 0x73, 0x73, 0x69, 0x70, 0x12, 0x28,
	0x0a, 0x08, 0x47, 0x65, 0x74, 0x50, 0x65, 0x65, 0x72, 0x73, 0x12, 0x0c, 0x2e, 0x4e, 0x65, 0x69,
	0x67, 0x68, 0x62, 0x6f, 0x72, 0x52, 0x65, 0x71, 0x1a, 0x0c, 0x2e, 0x4e, 0x65, 0x69, 0x67, 0x68,
	0x62, 0x6f, 0x72, 0x52, 0x65, 0x73, 0x22, 0x00, 0x12, 0x21, 0x0a, 0x08, 0x53, 0x65, 0x6e, 0x64,
	0x44, 0x61, 0x74, 0x61, 0x12, 0x0b, 0x2e, 0x47, 0x6f, 0x73, 0x73, 0x69, 0x70, 0x44, 0x61, 0x74,
	0x61, 0x1a, 0x06, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x42, 0x09, 0x5a, 0x07, 0x2f,
	0x67, 0x6f, 0x73, 0x73, 0x69, 0x70, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_message_proto_rawDescOnce sync.Once
	file_message_proto_rawDescData = file_message_proto_rawDesc
)

func file_message_proto_rawDescGZIP() []byte {
	file_message_proto_rawDescOnce.Do(func() {
		file_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_message_proto_rawDescData)
	})
	return file_message_proto_rawDescData
}

var file_message_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_message_proto_goTypes = []interface{}{
	(*Empty)(nil),       // 0: Empty
	(*NeighborReq)(nil), // 1: NeighborReq
	(*NeighborRes)(nil), // 2: NeighborRes
	(*GossipData)(nil),  // 3: GossipData
}
var file_message_proto_depIdxs = []int32{
	1, // 0: Gossip.GetPeers:input_type -> NeighborReq
	3, // 1: Gossip.SendData:input_type -> GossipData
	2, // 2: Gossip.GetPeers:output_type -> NeighborRes
	0, // 3: Gossip.SendData:output_type -> Empty
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_message_proto_init() }
func file_message_proto_init() {
	if File_message_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_message_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Empty); i {
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
		file_message_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NeighborReq); i {
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
		file_message_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NeighborRes); i {
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
		file_message_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GossipData); i {
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
			RawDescriptor: file_message_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_message_proto_goTypes,
		DependencyIndexes: file_message_proto_depIdxs,
		MessageInfos:      file_message_proto_msgTypes,
	}.Build()
	File_message_proto = out.File
	file_message_proto_rawDesc = nil
	file_message_proto_goTypes = nil
	file_message_proto_depIdxs = nil
}
