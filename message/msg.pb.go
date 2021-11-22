// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.6.1
// source: msg.proto

package message

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

type Msg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PID  string `protobuf:"bytes,1,opt,name=PID,proto3" json:"PID,omitempty"`
	Type int32  `protobuf:"varint,2,opt,name=type,proto3" json:"type,omitempty"`
	Data []byte `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *Msg) Reset() {
	*x = Msg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Msg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Msg) ProtoMessage() {}

func (x *Msg) ProtoReflect() protoreflect.Message {
	mi := &file_msg_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Msg.ProtoReflect.Descriptor instead.
func (*Msg) Descriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{0}
}

func (x *Msg) GetPID() string {
	if x != nil {
		return x.PID
	}
	return ""
}

func (x *Msg) GetType() int32 {
	if x != nil {
		return x.Type
	}
	return 0
}

func (x *Msg) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type VoteMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Signature  []byte `protobuf:"bytes,1,opt,name=Signature,proto3" json:"Signature,omitempty"`
	Round      uint64 `protobuf:"varint,2,opt,name=round,proto3" json:"round,omitempty"`
	Event      string `protobuf:"bytes,3,opt,name=event,proto3" json:"event,omitempty"`
	VRF        []byte `protobuf:"bytes,4,opt,name=VRF,proto3" json:"VRF,omitempty"`
	Proof      []byte `protobuf:"bytes,5,opt,name=Proof,proto3" json:"Proof,omitempty"`
	ParentHash []byte `protobuf:"bytes,6,opt,name=ParentHash,proto3" json:"ParentHash,omitempty"`
	Hash       []byte `protobuf:"bytes,7,opt,name=Hash,proto3" json:"Hash,omitempty"`
}

func (x *VoteMessage) Reset() {
	*x = VoteMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VoteMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VoteMessage) ProtoMessage() {}

func (x *VoteMessage) ProtoReflect() protoreflect.Message {
	mi := &file_msg_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VoteMessage.ProtoReflect.Descriptor instead.
func (*VoteMessage) Descriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{1}
}

func (x *VoteMessage) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

func (x *VoteMessage) GetRound() uint64 {
	if x != nil {
		return x.Round
	}
	return 0
}

func (x *VoteMessage) GetEvent() string {
	if x != nil {
		return x.Event
	}
	return ""
}

func (x *VoteMessage) GetVRF() []byte {
	if x != nil {
		return x.VRF
	}
	return nil
}

func (x *VoteMessage) GetProof() []byte {
	if x != nil {
		return x.Proof
	}
	return nil
}

func (x *VoteMessage) GetParentHash() []byte {
	if x != nil {
		return x.ParentHash
	}
	return nil
}

func (x *VoteMessage) GetHash() []byte {
	if x != nil {
		return x.Hash
	}
	return nil
}

type Proposal struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Round  uint64 `protobuf:"varint,1,opt,name=Round,proto3" json:"Round,omitempty"`
	Hash   []byte `protobuf:"bytes,2,opt,name=Hash,proto3" json:"Hash,omitempty"`
	Prior  []byte `protobuf:"bytes,3,opt,name=Prior,proto3" json:"Prior,omitempty"`
	VRF    []byte `protobuf:"bytes,4,opt,name=VRF,proto3" json:"VRF,omitempty"`
	Proof  []byte `protobuf:"bytes,5,opt,name=Proof,proto3" json:"Proof,omitempty"`
	Pubkey []byte `protobuf:"bytes,6,opt,name=Pubkey,proto3" json:"Pubkey,omitempty"`
	Block  []byte `protobuf:"bytes,7,opt,name=Block,proto3" json:"Block,omitempty"`
}

func (x *Proposal) Reset() {
	*x = Proposal{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Proposal) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Proposal) ProtoMessage() {}

func (x *Proposal) ProtoReflect() protoreflect.Message {
	mi := &file_msg_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Proposal.ProtoReflect.Descriptor instead.
func (*Proposal) Descriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{2}
}

func (x *Proposal) GetRound() uint64 {
	if x != nil {
		return x.Round
	}
	return 0
}

func (x *Proposal) GetHash() []byte {
	if x != nil {
		return x.Hash
	}
	return nil
}

func (x *Proposal) GetPrior() []byte {
	if x != nil {
		return x.Prior
	}
	return nil
}

func (x *Proposal) GetVRF() []byte {
	if x != nil {
		return x.VRF
	}
	return nil
}

func (x *Proposal) GetProof() []byte {
	if x != nil {
		return x.Proof
	}
	return nil
}

func (x *Proposal) GetPubkey() []byte {
	if x != nil {
		return x.Pubkey
	}
	return nil
}

func (x *Proposal) GetBlock() []byte {
	if x != nil {
		return x.Block
	}
	return nil
}

type Transaction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	From      []byte `protobuf:"bytes,1,opt,name=From,proto3" json:"From,omitempty"`
	To        []byte `protobuf:"bytes,2,opt,name=To,proto3" json:"To,omitempty"`
	Nonce     uint64 `protobuf:"varint,4,opt,name=Nonce,proto3" json:"Nonce,omitempty"`
	Type      uint64 `protobuf:"varint,5,opt,name=Type,proto3" json:"Type,omitempty"`
	Data      []byte `protobuf:"bytes,6,opt,name=Data,proto3" json:"Data,omitempty"`
	Signature []byte `protobuf:"bytes,7,opt,name=Signature,proto3" json:"Signature,omitempty"`
	Id        uint64 `protobuf:"varint,8,opt,name=Id,proto3" json:"Id,omitempty"`
}

func (x *Transaction) Reset() {
	*x = Transaction{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Transaction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Transaction) ProtoMessage() {}

func (x *Transaction) ProtoReflect() protoreflect.Message {
	mi := &file_msg_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Transaction.ProtoReflect.Descriptor instead.
func (*Transaction) Descriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{3}
}

func (x *Transaction) GetFrom() []byte {
	if x != nil {
		return x.From
	}
	return nil
}

func (x *Transaction) GetTo() []byte {
	if x != nil {
		return x.To
	}
	return nil
}

func (x *Transaction) GetNonce() uint64 {
	if x != nil {
		return x.Nonce
	}
	return 0
}

func (x *Transaction) GetType() uint64 {
	if x != nil {
		return x.Type
	}
	return 0
}

func (x *Transaction) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *Transaction) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

func (x *Transaction) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type PendingRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Nonce uint64 `protobuf:"varint,1,opt,name=Nonce,proto3" json:"Nonce,omitempty"`
	URL   string `protobuf:"bytes,2,opt,name=URL,proto3" json:"URL,omitempty"`
	Id    uint64 `protobuf:"varint,3,opt,name=Id,proto3" json:"Id,omitempty"`
}

func (x *PendingRequest) Reset() {
	*x = PendingRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PendingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PendingRequest) ProtoMessage() {}

func (x *PendingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_msg_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PendingRequest.ProtoReflect.Descriptor instead.
func (*PendingRequest) Descriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{4}
}

func (x *PendingRequest) GetNonce() uint64 {
	if x != nil {
		return x.Nonce
	}
	return 0
}

func (x *PendingRequest) GetURL() string {
	if x != nil {
		return x.URL
	}
	return ""
}

func (x *PendingRequest) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type Block struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Round       uint64            `protobuf:"varint,1,opt,name=Round,proto3" json:"Round,omitempty"`
	ParentHash  []byte            `protobuf:"bytes,2,opt,name=ParentHash,proto3" json:"ParentHash,omitempty"`
	Author      []byte            `protobuf:"bytes,3,opt,name=Author,proto3" json:"Author,omitempty"`
	AuthorVRF   []byte            `protobuf:"bytes,4,opt,name=AuthorVRF,proto3" json:"AuthorVRF,omitempty"`
	AuthorProof []byte            `protobuf:"bytes,5,opt,name=AuthorProof,proto3" json:"AuthorProof,omitempty"`
	Time        int64             `protobuf:"varint,6,opt,name=Time,proto3" json:"Time,omitempty"`
	Seed        []byte            `protobuf:"bytes,7,opt,name=Seed,proto3" json:"Seed,omitempty"`
	Proof       []byte            `protobuf:"bytes,8,opt,name=Proof,proto3" json:"Proof,omitempty"`
	Txs         []*Transaction    `protobuf:"bytes,9,rep,name=Txs,proto3" json:"Txs,omitempty"`
	StateHash   []byte            `protobuf:"bytes,10,opt,name=StateHash,proto3" json:"StateHash,omitempty"`
	Signature   []byte            `protobuf:"bytes,11,opt,name=Signature,proto3" json:"Signature,omitempty"`
	Data        []byte            `protobuf:"bytes,12,opt,name=Data,proto3" json:"Data,omitempty"`
	Reqs        []*PendingRequest `protobuf:"bytes,13,rep,name=Reqs,proto3" json:"Reqs,omitempty"`
	ResHash     [][]byte          `protobuf:"bytes,14,rep,name=ResHash,proto3" json:"ResHash,omitempty"`
}

func (x *Block) Reset() {
	*x = Block{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Block) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Block) ProtoMessage() {}

func (x *Block) ProtoReflect() protoreflect.Message {
	mi := &file_msg_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Block.ProtoReflect.Descriptor instead.
func (*Block) Descriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{5}
}

func (x *Block) GetRound() uint64 {
	if x != nil {
		return x.Round
	}
	return 0
}

func (x *Block) GetParentHash() []byte {
	if x != nil {
		return x.ParentHash
	}
	return nil
}

func (x *Block) GetAuthor() []byte {
	if x != nil {
		return x.Author
	}
	return nil
}

func (x *Block) GetAuthorVRF() []byte {
	if x != nil {
		return x.AuthorVRF
	}
	return nil
}

func (x *Block) GetAuthorProof() []byte {
	if x != nil {
		return x.AuthorProof
	}
	return nil
}

func (x *Block) GetTime() int64 {
	if x != nil {
		return x.Time
	}
	return 0
}

func (x *Block) GetSeed() []byte {
	if x != nil {
		return x.Seed
	}
	return nil
}

func (x *Block) GetProof() []byte {
	if x != nil {
		return x.Proof
	}
	return nil
}

func (x *Block) GetTxs() []*Transaction {
	if x != nil {
		return x.Txs
	}
	return nil
}

func (x *Block) GetStateHash() []byte {
	if x != nil {
		return x.StateHash
	}
	return nil
}

func (x *Block) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

func (x *Block) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *Block) GetReqs() []*PendingRequest {
	if x != nil {
		return x.Reqs
	}
	return nil
}

func (x *Block) GetResHash() [][]byte {
	if x != nil {
		return x.ResHash
	}
	return nil
}

type AccessResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ParentHash  []byte `protobuf:"bytes,1,opt,name=ParentHash,proto3" json:"ParentHash,omitempty"`
	RequestHash []byte `protobuf:"bytes,2,opt,name=RequestHash,proto3" json:"RequestHash,omitempty"`
	Data        []byte `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	Signature   []byte `protobuf:"bytes,4,opt,name=Signature,proto3" json:"Signature,omitempty"`
	VRF         []byte `protobuf:"bytes,5,opt,name=VRF,proto3" json:"VRF,omitempty"`
	Proof       []byte `protobuf:"bytes,6,opt,name=Proof,proto3" json:"Proof,omitempty"`
	Id          uint64 `protobuf:"varint,7,opt,name=Id,proto3" json:"Id,omitempty"`
}

func (x *AccessResponse) Reset() {
	*x = AccessResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AccessResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AccessResponse) ProtoMessage() {}

func (x *AccessResponse) ProtoReflect() protoreflect.Message {
	mi := &file_msg_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AccessResponse.ProtoReflect.Descriptor instead.
func (*AccessResponse) Descriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{6}
}

func (x *AccessResponse) GetParentHash() []byte {
	if x != nil {
		return x.ParentHash
	}
	return nil
}

func (x *AccessResponse) GetRequestHash() []byte {
	if x != nil {
		return x.RequestHash
	}
	return nil
}

func (x *AccessResponse) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *AccessResponse) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

func (x *AccessResponse) GetVRF() []byte {
	if x != nil {
		return x.VRF
	}
	return nil
}

func (x *AccessResponse) GetProof() []byte {
	if x != nil {
		return x.Proof
	}
	return nil
}

func (x *AccessResponse) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type ProposedTx struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Epoch uint64         `protobuf:"varint,1,opt,name=Epoch,proto3" json:"Epoch,omitempty"`
	Txs   []*Transaction `protobuf:"bytes,2,rep,name=Txs,proto3" json:"Txs,omitempty"`
}

func (x *ProposedTx) Reset() {
	*x = ProposedTx{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProposedTx) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProposedTx) ProtoMessage() {}

func (x *ProposedTx) ProtoReflect() protoreflect.Message {
	mi := &file_msg_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProposedTx.ProtoReflect.Descriptor instead.
func (*ProposedTx) Descriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{7}
}

func (x *ProposedTx) GetEpoch() uint64 {
	if x != nil {
		return x.Epoch
	}
	return 0
}

func (x *ProposedTx) GetTxs() []*Transaction {
	if x != nil {
		return x.Txs
	}
	return nil
}

type StateHash struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Epoch     uint64 `protobuf:"varint,1,opt,name=Epoch,proto3" json:"Epoch,omitempty"`
	StateHash []byte `protobuf:"bytes,2,opt,name=StateHash,proto3" json:"StateHash,omitempty"`
}

func (x *StateHash) Reset() {
	*x = StateHash{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StateHash) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StateHash) ProtoMessage() {}

func (x *StateHash) ProtoReflect() protoreflect.Message {
	mi := &file_msg_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StateHash.ProtoReflect.Descriptor instead.
func (*StateHash) Descriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{8}
}

func (x *StateHash) GetEpoch() uint64 {
	if x != nil {
		return x.Epoch
	}
	return 0
}

func (x *StateHash) GetStateHash() []byte {
	if x != nil {
		return x.StateHash
	}
	return nil
}

type OraclePeerProposal struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pubkey []byte `protobuf:"bytes,1,opt,name=Pubkey,proto3" json:"Pubkey,omitempty"`
	Proof  []byte `protobuf:"bytes,2,opt,name=Proof,proto3" json:"Proof,omitempty"`
	Vrf    []byte `protobuf:"bytes,3,opt,name=Vrf,proto3" json:"Vrf,omitempty"`
	Epoch  uint64 `protobuf:"varint,4,opt,name=Epoch,proto3" json:"Epoch,omitempty"`
	Weight uint64 `protobuf:"varint,5,opt,name=Weight,proto3" json:"Weight,omitempty"`
}

func (x *OraclePeerProposal) Reset() {
	*x = OraclePeerProposal{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OraclePeerProposal) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OraclePeerProposal) ProtoMessage() {}

func (x *OraclePeerProposal) ProtoReflect() protoreflect.Message {
	mi := &file_msg_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OraclePeerProposal.ProtoReflect.Descriptor instead.
func (*OraclePeerProposal) Descriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{9}
}

func (x *OraclePeerProposal) GetPubkey() []byte {
	if x != nil {
		return x.Pubkey
	}
	return nil
}

func (x *OraclePeerProposal) GetProof() []byte {
	if x != nil {
		return x.Proof
	}
	return nil
}

func (x *OraclePeerProposal) GetVrf() []byte {
	if x != nil {
		return x.Vrf
	}
	return nil
}

func (x *OraclePeerProposal) GetEpoch() uint64 {
	if x != nil {
		return x.Epoch
	}
	return 0
}

func (x *OraclePeerProposal) GetWeight() uint64 {
	if x != nil {
		return x.Weight
	}
	return 0
}

var File_msg_proto protoreflect.FileDescriptor

var file_msg_proto_rawDesc = []byte{
	0x0a, 0x09, 0x6d, 0x73, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x3f, 0x0a, 0x03, 0x4d,
	0x73, 0x67, 0x12, 0x10, 0x0a, 0x03, 0x50, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x50, 0x49, 0x44, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0xb3, 0x01, 0x0a,
	0x0b, 0x56, 0x6f, 0x74, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1c, 0x0a, 0x09,
	0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x09, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x6f,
	0x75, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64,
	0x12, 0x14, 0x0a, 0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x56, 0x52, 0x46, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x03, 0x56, 0x52, 0x46, 0x12, 0x14, 0x0a, 0x05, 0x50, 0x72, 0x6f, 0x6f,
	0x66, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x50, 0x72, 0x6f, 0x6f, 0x66, 0x12, 0x1e,
	0x0a, 0x0a, 0x50, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x48, 0x61, 0x73, 0x68, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x0a, 0x50, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x48, 0x61, 0x73, 0x68, 0x12, 0x12,
	0x0a, 0x04, 0x48, 0x61, 0x73, 0x68, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x48, 0x61,
	0x73, 0x68, 0x22, 0xa0, 0x01, 0x0a, 0x08, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x12,
	0x14, 0x0a, 0x05, 0x52, 0x6f, 0x75, 0x6e, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05,
	0x52, 0x6f, 0x75, 0x6e, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x48, 0x61, 0x73, 0x68, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x04, 0x48, 0x61, 0x73, 0x68, 0x12, 0x14, 0x0a, 0x05, 0x50, 0x72, 0x69,
	0x6f, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x50, 0x72, 0x69, 0x6f, 0x72, 0x12,
	0x10, 0x0a, 0x03, 0x56, 0x52, 0x46, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x56, 0x52,
	0x46, 0x12, 0x14, 0x0a, 0x05, 0x50, 0x72, 0x6f, 0x6f, 0x66, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x05, 0x50, 0x72, 0x6f, 0x6f, 0x66, 0x12, 0x16, 0x0a, 0x06, 0x50, 0x75, 0x62, 0x6b, 0x65,
	0x79, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x50, 0x75, 0x62, 0x6b, 0x65, 0x79, 0x12,
	0x14, 0x0a, 0x05, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05,
	0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x22, 0x9d, 0x01, 0x0a, 0x0b, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x46, 0x72, 0x6f, 0x6d, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x04, 0x46, 0x72, 0x6f, 0x6d, 0x12, 0x0e, 0x0a, 0x02, 0x54, 0x6f, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x02, 0x54, 0x6f, 0x12, 0x14, 0x0a, 0x05, 0x4e, 0x6f, 0x6e,
	0x63, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x4e, 0x6f, 0x6e, 0x63, 0x65, 0x12,
	0x12, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x04, 0x44, 0x61, 0x74, 0x61, 0x12, 0x1c, 0x0a, 0x09, 0x53, 0x69, 0x67, 0x6e, 0x61,
	0x74, 0x75, 0x72, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x53, 0x69, 0x67, 0x6e,
	0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x64, 0x18, 0x08, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x02, 0x49, 0x64, 0x22, 0x48, 0x0a, 0x0e, 0x50, 0x65, 0x6e, 0x64, 0x69, 0x6e, 0x67,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x4e, 0x6f, 0x6e, 0x63, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x4e, 0x6f, 0x6e, 0x63, 0x65, 0x12, 0x10, 0x0a,
	0x03, 0x55, 0x52, 0x4c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x55, 0x52, 0x4c, 0x12,
	0x0e, 0x0a, 0x02, 0x49, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x49, 0x64, 0x22,
	0x82, 0x03, 0x0a, 0x05, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x14, 0x0a, 0x05, 0x52, 0x6f, 0x75,
	0x6e, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x52, 0x6f, 0x75, 0x6e, 0x64, 0x12,
	0x1e, 0x0a, 0x0a, 0x50, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x48, 0x61, 0x73, 0x68, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x0a, 0x50, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x48, 0x61, 0x73, 0x68, 0x12,
	0x16, 0x0a, 0x06, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x06, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x12, 0x1c, 0x0a, 0x09, 0x41, 0x75, 0x74, 0x68, 0x6f,
	0x72, 0x56, 0x52, 0x46, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x41, 0x75, 0x74, 0x68,
	0x6f, 0x72, 0x56, 0x52, 0x46, 0x12, 0x20, 0x0a, 0x0b, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x50,
	0x72, 0x6f, 0x6f, 0x66, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0b, 0x41, 0x75, 0x74, 0x68,
	0x6f, 0x72, 0x50, 0x72, 0x6f, 0x6f, 0x66, 0x12, 0x12, 0x0a, 0x04, 0x54, 0x69, 0x6d, 0x65, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x53,
	0x65, 0x65, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x53, 0x65, 0x65, 0x64, 0x12,
	0x14, 0x0a, 0x05, 0x50, 0x72, 0x6f, 0x6f, 0x66, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05,
	0x50, 0x72, 0x6f, 0x6f, 0x66, 0x12, 0x1e, 0x0a, 0x03, 0x54, 0x78, 0x73, 0x18, 0x09, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x03, 0x54, 0x78, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x53, 0x74, 0x61, 0x74, 0x65, 0x48, 0x61,
	0x73, 0x68, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x53, 0x74, 0x61, 0x74, 0x65, 0x48,
	0x61, 0x73, 0x68, 0x12, 0x1c, 0x0a, 0x09, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65,
	0x18, 0x0b, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x04, 0x44, 0x61, 0x74, 0x61, 0x12, 0x23, 0x0a, 0x04, 0x52, 0x65, 0x71, 0x73, 0x18, 0x0d, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x50, 0x65, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x52, 0x04, 0x52, 0x65, 0x71, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x52, 0x65,
	0x73, 0x48, 0x61, 0x73, 0x68, 0x18, 0x0e, 0x20, 0x03, 0x28, 0x0c, 0x52, 0x07, 0x52, 0x65, 0x73,
	0x48, 0x61, 0x73, 0x68, 0x22, 0xbc, 0x01, 0x0a, 0x0e, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x50, 0x61, 0x72, 0x65, 0x6e,
	0x74, 0x48, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0a, 0x50, 0x61, 0x72,
	0x65, 0x6e, 0x74, 0x48, 0x61, 0x73, 0x68, 0x12, 0x20, 0x0a, 0x0b, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x48, 0x61, 0x73, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0b, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x48, 0x61, 0x73, 0x68, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x1c, 0x0a,
	0x09, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x09, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x56,
	0x52, 0x46, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x56, 0x52, 0x46, 0x12, 0x14, 0x0a,
	0x05, 0x50, 0x72, 0x6f, 0x6f, 0x66, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x50, 0x72,
	0x6f, 0x6f, 0x66, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x02, 0x49, 0x64, 0x22, 0x42, 0x0a, 0x0a, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x65, 0x64, 0x54,
	0x78, 0x12, 0x14, 0x0a, 0x05, 0x45, 0x70, 0x6f, 0x63, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x05, 0x45, 0x70, 0x6f, 0x63, 0x68, 0x12, 0x1e, 0x0a, 0x03, 0x54, 0x78, 0x73, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x52, 0x03, 0x54, 0x78, 0x73, 0x22, 0x3f, 0x0a, 0x09, 0x53, 0x74, 0x61, 0x74, 0x65,
	0x48, 0x61, 0x73, 0x68, 0x12, 0x14, 0x0a, 0x05, 0x45, 0x70, 0x6f, 0x63, 0x68, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x05, 0x45, 0x70, 0x6f, 0x63, 0x68, 0x12, 0x1c, 0x0a, 0x09, 0x53, 0x74,
	0x61, 0x74, 0x65, 0x48, 0x61, 0x73, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x53,
	0x74, 0x61, 0x74, 0x65, 0x48, 0x61, 0x73, 0x68, 0x22, 0x82, 0x01, 0x0a, 0x12, 0x4f, 0x72, 0x61,
	0x63, 0x6c, 0x65, 0x50, 0x65, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x12,
	0x16, 0x0a, 0x06, 0x50, 0x75, 0x62, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x06, 0x50, 0x75, 0x62, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x50, 0x72, 0x6f, 0x6f, 0x66,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x50, 0x72, 0x6f, 0x6f, 0x66, 0x12, 0x10, 0x0a,
	0x03, 0x56, 0x72, 0x66, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x56, 0x72, 0x66, 0x12,
	0x14, 0x0a, 0x05, 0x45, 0x70, 0x6f, 0x63, 0x68, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05,
	0x45, 0x70, 0x6f, 0x63, 0x68, 0x12, 0x16, 0x0a, 0x06, 0x57, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x57, 0x65, 0x69, 0x67, 0x68, 0x74, 0x42, 0x0a, 0x5a,
	0x08, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_msg_proto_rawDescOnce sync.Once
	file_msg_proto_rawDescData = file_msg_proto_rawDesc
)

func file_msg_proto_rawDescGZIP() []byte {
	file_msg_proto_rawDescOnce.Do(func() {
		file_msg_proto_rawDescData = protoimpl.X.CompressGZIP(file_msg_proto_rawDescData)
	})
	return file_msg_proto_rawDescData
}

var file_msg_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_msg_proto_goTypes = []interface{}{
	(*Msg)(nil),                // 0: Msg
	(*VoteMessage)(nil),        // 1: VoteMessage
	(*Proposal)(nil),           // 2: Proposal
	(*Transaction)(nil),        // 3: Transaction
	(*PendingRequest)(nil),     // 4: PendingRequest
	(*Block)(nil),              // 5: Block
	(*AccessResponse)(nil),     // 6: AccessResponse
	(*ProposedTx)(nil),         // 7: ProposedTx
	(*StateHash)(nil),          // 8: StateHash
	(*OraclePeerProposal)(nil), // 9: OraclePeerProposal
}
var file_msg_proto_depIdxs = []int32{
	3, // 0: Block.Txs:type_name -> Transaction
	4, // 1: Block.Reqs:type_name -> PendingRequest
	3, // 2: ProposedTx.Txs:type_name -> Transaction
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_msg_proto_init() }
func file_msg_proto_init() {
	if File_msg_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_msg_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Msg); i {
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
		file_msg_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VoteMessage); i {
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
		file_msg_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Proposal); i {
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
		file_msg_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Transaction); i {
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
		file_msg_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PendingRequest); i {
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
		file_msg_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Block); i {
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
		file_msg_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AccessResponse); i {
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
		file_msg_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProposedTx); i {
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
		file_msg_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StateHash); i {
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
		file_msg_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OraclePeerProposal); i {
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
			RawDescriptor: file_msg_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_msg_proto_goTypes,
		DependencyIndexes: file_msg_proto_depIdxs,
		MessageInfos:      file_msg_proto_msgTypes,
	}.Build()
	File_msg_proto = out.File
	file_msg_proto_rawDesc = nil
	file_msg_proto_goTypes = nil
	file_msg_proto_depIdxs = nil
}
