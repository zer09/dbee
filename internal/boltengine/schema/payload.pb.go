// Code generated by protoc-gen-go. DO NOT EDIT.
// source: payload.proto

package schema

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import timestamp "github.com/golang/protobuf/ptypes/timestamp"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Meta struct {
	CreatedOn            *timestamp.Timestamp `protobuf:"bytes,1,opt,name=createdOn,proto3" json:"createdOn,omitempty"`
	LastUpdate           *timestamp.Timestamp `protobuf:"bytes,2,opt,name=lastUpdate,proto3" json:"lastUpdate,omitempty"`
	Deleted              bool                 `protobuf:"varint,3,opt,name=deleted,proto3" json:"deleted,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Meta) Reset()         { *m = Meta{} }
func (m *Meta) String() string { return proto.CompactTextString(m) }
func (*Meta) ProtoMessage()    {}
func (*Meta) Descriptor() ([]byte, []int) {
	return fileDescriptor_payload_563679a0bb4d487b, []int{0}
}
func (m *Meta) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Meta.Unmarshal(m, b)
}
func (m *Meta) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Meta.Marshal(b, m, deterministic)
}
func (dst *Meta) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Meta.Merge(dst, src)
}
func (m *Meta) XXX_Size() int {
	return xxx_messageInfo_Meta.Size(m)
}
func (m *Meta) XXX_DiscardUnknown() {
	xxx_messageInfo_Meta.DiscardUnknown(m)
}

var xxx_messageInfo_Meta proto.InternalMessageInfo

func (m *Meta) GetCreatedOn() *timestamp.Timestamp {
	if m != nil {
		return m.CreatedOn
	}
	return nil
}

func (m *Meta) GetLastUpdate() *timestamp.Timestamp {
	if m != nil {
		return m.LastUpdate
	}
	return nil
}

func (m *Meta) GetDeleted() bool {
	if m != nil {
		return m.Deleted
	}
	return false
}

type Payload struct {
	Meta *Meta `protobuf:"bytes,1,opt,name=meta,proto3" json:"meta,omitempty"`
	// values of the payload the key will be the param index name.
	Values               map[uint64][]byte `protobuf:"bytes,2,rep,name=values,proto3" json:"values,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Payload) Reset()         { *m = Payload{} }
func (m *Payload) String() string { return proto.CompactTextString(m) }
func (*Payload) ProtoMessage()    {}
func (*Payload) Descriptor() ([]byte, []int) {
	return fileDescriptor_payload_563679a0bb4d487b, []int{1}
}
func (m *Payload) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Payload.Unmarshal(m, b)
}
func (m *Payload) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Payload.Marshal(b, m, deterministic)
}
func (dst *Payload) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Payload.Merge(dst, src)
}
func (m *Payload) XXX_Size() int {
	return xxx_messageInfo_Payload.Size(m)
}
func (m *Payload) XXX_DiscardUnknown() {
	xxx_messageInfo_Payload.DiscardUnknown(m)
}

var xxx_messageInfo_Payload proto.InternalMessageInfo

func (m *Payload) GetMeta() *Meta {
	if m != nil {
		return m.Meta
	}
	return nil
}

func (m *Payload) GetValues() map[uint64][]byte {
	if m != nil {
		return m.Values
	}
	return nil
}

type PayloadSint64 struct {
	Value                int64    `protobuf:"zigzag64,1,opt,name=Value,proto3" json:"Value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PayloadSint64) Reset()         { *m = PayloadSint64{} }
func (m *PayloadSint64) String() string { return proto.CompactTextString(m) }
func (*PayloadSint64) ProtoMessage()    {}
func (*PayloadSint64) Descriptor() ([]byte, []int) {
	return fileDescriptor_payload_563679a0bb4d487b, []int{2}
}
func (m *PayloadSint64) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PayloadSint64.Unmarshal(m, b)
}
func (m *PayloadSint64) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PayloadSint64.Marshal(b, m, deterministic)
}
func (dst *PayloadSint64) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PayloadSint64.Merge(dst, src)
}
func (m *PayloadSint64) XXX_Size() int {
	return xxx_messageInfo_PayloadSint64.Size(m)
}
func (m *PayloadSint64) XXX_DiscardUnknown() {
	xxx_messageInfo_PayloadSint64.DiscardUnknown(m)
}

var xxx_messageInfo_PayloadSint64 proto.InternalMessageInfo

func (m *PayloadSint64) GetValue() int64 {
	if m != nil {
		return m.Value
	}
	return 0
}

type PayloadUint64 struct {
	Value                uint64   `protobuf:"varint,1,opt,name=Value,proto3" json:"Value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PayloadUint64) Reset()         { *m = PayloadUint64{} }
func (m *PayloadUint64) String() string { return proto.CompactTextString(m) }
func (*PayloadUint64) ProtoMessage()    {}
func (*PayloadUint64) Descriptor() ([]byte, []int) {
	return fileDescriptor_payload_563679a0bb4d487b, []int{3}
}
func (m *PayloadUint64) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PayloadUint64.Unmarshal(m, b)
}
func (m *PayloadUint64) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PayloadUint64.Marshal(b, m, deterministic)
}
func (dst *PayloadUint64) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PayloadUint64.Merge(dst, src)
}
func (m *PayloadUint64) XXX_Size() int {
	return xxx_messageInfo_PayloadUint64.Size(m)
}
func (m *PayloadUint64) XXX_DiscardUnknown() {
	xxx_messageInfo_PayloadUint64.DiscardUnknown(m)
}

var xxx_messageInfo_PayloadUint64 proto.InternalMessageInfo

func (m *PayloadUint64) GetValue() uint64 {
	if m != nil {
		return m.Value
	}
	return 0
}

type PayloadBool struct {
	Value                bool     `protobuf:"varint,1,opt,name=Value,proto3" json:"Value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PayloadBool) Reset()         { *m = PayloadBool{} }
func (m *PayloadBool) String() string { return proto.CompactTextString(m) }
func (*PayloadBool) ProtoMessage()    {}
func (*PayloadBool) Descriptor() ([]byte, []int) {
	return fileDescriptor_payload_563679a0bb4d487b, []int{4}
}
func (m *PayloadBool) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PayloadBool.Unmarshal(m, b)
}
func (m *PayloadBool) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PayloadBool.Marshal(b, m, deterministic)
}
func (dst *PayloadBool) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PayloadBool.Merge(dst, src)
}
func (m *PayloadBool) XXX_Size() int {
	return xxx_messageInfo_PayloadBool.Size(m)
}
func (m *PayloadBool) XXX_DiscardUnknown() {
	xxx_messageInfo_PayloadBool.DiscardUnknown(m)
}

var xxx_messageInfo_PayloadBool proto.InternalMessageInfo

func (m *PayloadBool) GetValue() bool {
	if m != nil {
		return m.Value
	}
	return false
}

type PayloadString struct {
	Value                string   `protobuf:"bytes,1,opt,name=Value,proto3" json:"Value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PayloadString) Reset()         { *m = PayloadString{} }
func (m *PayloadString) String() string { return proto.CompactTextString(m) }
func (*PayloadString) ProtoMessage()    {}
func (*PayloadString) Descriptor() ([]byte, []int) {
	return fileDescriptor_payload_563679a0bb4d487b, []int{5}
}
func (m *PayloadString) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PayloadString.Unmarshal(m, b)
}
func (m *PayloadString) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PayloadString.Marshal(b, m, deterministic)
}
func (dst *PayloadString) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PayloadString.Merge(dst, src)
}
func (m *PayloadString) XXX_Size() int {
	return xxx_messageInfo_PayloadString.Size(m)
}
func (m *PayloadString) XXX_DiscardUnknown() {
	xxx_messageInfo_PayloadString.DiscardUnknown(m)
}

var xxx_messageInfo_PayloadString proto.InternalMessageInfo

func (m *PayloadString) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

type PayloadBytes struct {
	Value                []byte   `protobuf:"bytes,1,opt,name=Value,proto3" json:"Value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PayloadBytes) Reset()         { *m = PayloadBytes{} }
func (m *PayloadBytes) String() string { return proto.CompactTextString(m) }
func (*PayloadBytes) ProtoMessage()    {}
func (*PayloadBytes) Descriptor() ([]byte, []int) {
	return fileDescriptor_payload_563679a0bb4d487b, []int{6}
}
func (m *PayloadBytes) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PayloadBytes.Unmarshal(m, b)
}
func (m *PayloadBytes) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PayloadBytes.Marshal(b, m, deterministic)
}
func (dst *PayloadBytes) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PayloadBytes.Merge(dst, src)
}
func (m *PayloadBytes) XXX_Size() int {
	return xxx_messageInfo_PayloadBytes.Size(m)
}
func (m *PayloadBytes) XXX_DiscardUnknown() {
	xxx_messageInfo_PayloadBytes.DiscardUnknown(m)
}

var xxx_messageInfo_PayloadBytes proto.InternalMessageInfo

func (m *PayloadBytes) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

func init() {
	proto.RegisterType((*Meta)(nil), "schema.Meta")
	proto.RegisterType((*Payload)(nil), "schema.Payload")
	proto.RegisterMapType((map[uint64][]byte)(nil), "schema.Payload.ValuesEntry")
	proto.RegisterType((*PayloadSint64)(nil), "schema.PayloadSint64")
	proto.RegisterType((*PayloadUint64)(nil), "schema.PayloadUint64")
	proto.RegisterType((*PayloadBool)(nil), "schema.PayloadBool")
	proto.RegisterType((*PayloadString)(nil), "schema.PayloadString")
	proto.RegisterType((*PayloadBytes)(nil), "schema.PayloadBytes")
}

func init() { proto.RegisterFile("payload.proto", fileDescriptor_payload_563679a0bb4d487b) }

var fileDescriptor_payload_563679a0bb4d487b = []byte{
	// 314 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0xd1, 0x4d, 0x4b, 0xf3, 0x40,
	0x10, 0x07, 0x70, 0xd2, 0xe6, 0xe9, 0xcb, 0x24, 0x85, 0x87, 0xc5, 0x43, 0xa8, 0x07, 0x43, 0x55,
	0xe8, 0x69, 0x0b, 0xad, 0x48, 0xed, 0x51, 0xf0, 0x28, 0xca, 0xd6, 0x7a, 0xdf, 0x36, 0x63, 0x2d,
	0x26, 0xd9, 0x90, 0x4c, 0x85, 0x7c, 0x11, 0x2f, 0x7e, 0x59, 0xe9, 0xbe, 0x60, 0x83, 0x82, 0xb7,
	0xce, 0xee, 0xef, 0xff, 0xef, 0x0e, 0x81, 0x41, 0x21, 0xeb, 0x54, 0xc9, 0x84, 0x17, 0xa5, 0x22,
	0xc5, 0x3a, 0xd5, 0xe6, 0x15, 0x33, 0x39, 0x3c, 0xdb, 0x2a, 0xb5, 0x4d, 0x71, 0xa2, 0x4f, 0xd7,
	0xfb, 0x97, 0x09, 0xed, 0x32, 0xac, 0x48, 0x66, 0x85, 0x81, 0xa3, 0x0f, 0x0f, 0xfc, 0x7b, 0x24,
	0xc9, 0xe6, 0xd0, 0xdf, 0x94, 0x28, 0x09, 0x93, 0x87, 0x3c, 0xf2, 0x62, 0x6f, 0x1c, 0x4c, 0x87,
	0xdc, 0xa4, 0xb9, 0x4b, 0xf3, 0x27, 0x97, 0x16, 0xdf, 0x98, 0x2d, 0x00, 0x52, 0x59, 0xd1, 0xaa,
	0x48, 0x24, 0x61, 0xd4, 0xfa, 0x33, 0x7a, 0xa4, 0x59, 0x04, 0xdd, 0x04, 0x53, 0x24, 0x4c, 0xa2,
	0x76, 0xec, 0x8d, 0x7b, 0xc2, 0x8d, 0xa3, 0x4f, 0x0f, 0xba, 0x8f, 0x66, 0x27, 0x16, 0x83, 0x9f,
	0x21, 0x49, 0xfb, 0xac, 0x90, 0x9b, 0xe5, 0xf8, 0xe1, 0xdd, 0x42, 0xdf, 0xb0, 0x19, 0x74, 0xde,
	0x65, 0xba, 0xc7, 0x2a, 0x6a, 0xc5, 0xed, 0x71, 0x30, 0x3d, 0x75, 0xc6, 0x56, 0xf0, 0x67, 0x7d,
	0x7b, 0x97, 0x53, 0x59, 0x0b, 0x4b, 0x87, 0x37, 0x10, 0x1c, 0x1d, 0xb3, 0xff, 0xd0, 0x7e, 0xc3,
	0x5a, 0xff, 0x89, 0x2f, 0x0e, 0x3f, 0xd9, 0x09, 0xfc, 0xd3, 0x54, 0x2f, 0x15, 0x0a, 0x33, 0x2c,
	0x5a, 0x73, 0x6f, 0x74, 0x09, 0x03, 0xdb, 0xbc, 0xdc, 0xe5, 0x74, 0x7d, 0x75, 0xa0, 0xba, 0x4b,
	0xc7, 0x99, 0x30, 0xc3, 0x11, 0x5b, 0xfd, 0xc2, 0x7c, 0xc7, 0xce, 0x21, 0xb0, 0xec, 0x56, 0xa9,
	0xb4, 0x89, 0x7a, 0x3f, 0xbb, 0x96, 0x54, 0xee, 0xf2, 0x6d, 0x93, 0xf5, 0x1d, 0xbb, 0x80, 0xd0,
	0x75, 0xd5, 0x84, 0x55, 0x53, 0x85, 0x56, 0xad, 0x3b, 0xfa, 0xbb, 0xcc, 0xbe, 0x02, 0x00, 0x00,
	0xff, 0xff, 0x80, 0xef, 0xd5, 0xd4, 0x37, 0x02, 0x00, 0x00,
}
