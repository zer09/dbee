// Code generated by protoc-gen-go. DO NOT EDIT.
// source: set.proto

package schema

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Set struct {
	Partition            []*Partition `protobuf:"bytes,1,rep,name=partition,proto3" json:"partition,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Set) Reset()         { *m = Set{} }
func (m *Set) String() string { return proto.CompactTextString(m) }
func (*Set) ProtoMessage()    {}
func (*Set) Descriptor() ([]byte, []int) {
	return fileDescriptor_set_ad2406fdd1e32581, []int{0}
}
func (m *Set) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Set.Unmarshal(m, b)
}
func (m *Set) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Set.Marshal(b, m, deterministic)
}
func (dst *Set) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Set.Merge(dst, src)
}
func (m *Set) XXX_Size() int {
	return xxx_messageInfo_Set.Size(m)
}
func (m *Set) XXX_DiscardUnknown() {
	xxx_messageInfo_Set.DiscardUnknown(m)
}

var xxx_messageInfo_Set proto.InternalMessageInfo

func (m *Set) GetPartition() []*Partition {
	if m != nil {
		return m.Partition
	}
	return nil
}

type Partition struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Index                string   `protobuf:"bytes,2,opt,name=index,proto3" json:"index,omitempty"`
	Store                string   `protobuf:"bytes,3,opt,name=store,proto3" json:"store,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Partition) Reset()         { *m = Partition{} }
func (m *Partition) String() string { return proto.CompactTextString(m) }
func (*Partition) ProtoMessage()    {}
func (*Partition) Descriptor() ([]byte, []int) {
	return fileDescriptor_set_ad2406fdd1e32581, []int{1}
}
func (m *Partition) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Partition.Unmarshal(m, b)
}
func (m *Partition) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Partition.Marshal(b, m, deterministic)
}
func (dst *Partition) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Partition.Merge(dst, src)
}
func (m *Partition) XXX_Size() int {
	return xxx_messageInfo_Partition.Size(m)
}
func (m *Partition) XXX_DiscardUnknown() {
	xxx_messageInfo_Partition.DiscardUnknown(m)
}

var xxx_messageInfo_Partition proto.InternalMessageInfo

func (m *Partition) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Partition) GetIndex() string {
	if m != nil {
		return m.Index
	}
	return ""
}

func (m *Partition) GetStore() string {
	if m != nil {
		return m.Store
	}
	return ""
}

func init() {
	proto.RegisterType((*Set)(nil), "schema.Set")
	proto.RegisterType((*Partition)(nil), "schema.Partition")
}

func init() { proto.RegisterFile("set.proto", fileDescriptor_set_ad2406fdd1e32581) }

var fileDescriptor_set_ad2406fdd1e32581 = []byte{
	// 136 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2c, 0x4e, 0x2d, 0xd1,
	0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2b, 0x4e, 0xce, 0x48, 0xcd, 0x4d, 0x54, 0x32, 0xe3,
	0x62, 0x0e, 0x4e, 0x2d, 0x11, 0xd2, 0xe7, 0xe2, 0x2c, 0x48, 0x2c, 0x2a, 0xc9, 0x2c, 0xc9, 0xcc,
	0xcf, 0x93, 0x60, 0x54, 0x60, 0xd6, 0xe0, 0x36, 0x12, 0xd4, 0x83, 0x28, 0xd1, 0x0b, 0x80, 0x49,
	0x04, 0x21, 0xd4, 0x28, 0x79, 0x73, 0x71, 0xc2, 0xc5, 0x85, 0x84, 0xb8, 0x58, 0xf2, 0x12, 0x73,
	0x53, 0x25, 0x18, 0x15, 0x18, 0x35, 0x38, 0x83, 0xc0, 0x6c, 0x21, 0x11, 0x2e, 0xd6, 0xcc, 0xbc,
	0x94, 0xd4, 0x0a, 0x09, 0x26, 0xb0, 0x20, 0x84, 0x03, 0x12, 0x2d, 0x2e, 0xc9, 0x2f, 0x4a, 0x95,
	0x60, 0x86, 0x88, 0x82, 0x39, 0x49, 0x6c, 0x60, 0x37, 0x19, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff,
	0x1d, 0x62, 0x71, 0x78, 0xa0, 0x00, 0x00, 0x00,
}
