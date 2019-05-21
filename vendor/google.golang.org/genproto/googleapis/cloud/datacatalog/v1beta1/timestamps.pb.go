// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/cloud/datacatalog/v1beta1/timestamps.proto

package datacatalog // import "google.golang.org/genproto/googleapis/cloud/datacatalog/v1beta1"

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

// Timestamps about this resource according to a particular system.
type SystemTimestamps struct {
	// Output only. The creation time of the resource within the given system.
	CreateTime *timestamp.Timestamp `protobuf:"bytes,1,opt,name=create_time,json=createTime,proto3" json:"create_time,omitempty"`
	// Output only. The last-modified time of the resource within the given
	// system.
	UpdateTime *timestamp.Timestamp `protobuf:"bytes,2,opt,name=update_time,json=updateTime,proto3" json:"update_time,omitempty"`
	// Output only. The expiration time of the resource within the given system.
	ExpireTime           *timestamp.Timestamp `protobuf:"bytes,3,opt,name=expire_time,json=expireTime,proto3" json:"expire_time,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *SystemTimestamps) Reset()         { *m = SystemTimestamps{} }
func (m *SystemTimestamps) String() string { return proto.CompactTextString(m) }
func (*SystemTimestamps) ProtoMessage()    {}
func (*SystemTimestamps) Descriptor() ([]byte, []int) {
	return fileDescriptor_timestamps_d8e8b54cbb1fb3a6, []int{0}
}
func (m *SystemTimestamps) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SystemTimestamps.Unmarshal(m, b)
}
func (m *SystemTimestamps) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SystemTimestamps.Marshal(b, m, deterministic)
}
func (dst *SystemTimestamps) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SystemTimestamps.Merge(dst, src)
}
func (m *SystemTimestamps) XXX_Size() int {
	return xxx_messageInfo_SystemTimestamps.Size(m)
}
func (m *SystemTimestamps) XXX_DiscardUnknown() {
	xxx_messageInfo_SystemTimestamps.DiscardUnknown(m)
}

var xxx_messageInfo_SystemTimestamps proto.InternalMessageInfo

func (m *SystemTimestamps) GetCreateTime() *timestamp.Timestamp {
	if m != nil {
		return m.CreateTime
	}
	return nil
}

func (m *SystemTimestamps) GetUpdateTime() *timestamp.Timestamp {
	if m != nil {
		return m.UpdateTime
	}
	return nil
}

func (m *SystemTimestamps) GetExpireTime() *timestamp.Timestamp {
	if m != nil {
		return m.ExpireTime
	}
	return nil
}

func init() {
	proto.RegisterType((*SystemTimestamps)(nil), "google.cloud.datacatalog.v1beta1.SystemTimestamps")
}

func init() {
	proto.RegisterFile("google/cloud/datacatalog/v1beta1/timestamps.proto", fileDescriptor_timestamps_d8e8b54cbb1fb3a6)
}

var fileDescriptor_timestamps_d8e8b54cbb1fb3a6 = []byte{
	// 229 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0xd1, 0x3f, 0x4b, 0x03, 0x41,
	0x10, 0x05, 0x70, 0x56, 0xc1, 0x62, 0xd3, 0x48, 0x2a, 0x09, 0x82, 0xc1, 0xca, 0x6a, 0x96, 0xd3,
	0x32, 0x9d, 0xad, 0x8d, 0xa8, 0x95, 0x8d, 0xcc, 0xdd, 0x8d, 0xcb, 0xc1, 0x6d, 0x66, 0xd9, 0x9d,
	0x13, 0xfd, 0x88, 0x7e, 0x23, 0x4b, 0xd9, 0x3f, 0xb9, 0xa4, 0x09, 0x69, 0xdf, 0xbc, 0xdf, 0x3b,
	0xb8, 0xd5, 0x8d, 0x65, 0xb6, 0x23, 0x99, 0x6e, 0xe4, 0xa9, 0x37, 0x3d, 0x0a, 0x76, 0x28, 0x38,
	0xb2, 0x35, 0x5f, 0x4d, 0x4b, 0x82, 0x8d, 0x91, 0xc1, 0x51, 0x14, 0x74, 0x3e, 0x82, 0x0f, 0x2c,
	0xbc, 0x5c, 0x17, 0x02, 0x99, 0xc0, 0x01, 0x81, 0x4a, 0x56, 0x37, 0x75, 0x34, 0xf7, 0xdb, 0xe9,
	0x73, 0xbf, 0x51, 0x26, 0x6e, 0x7f, 0x95, 0xbe, 0x7c, 0xfd, 0x89, 0x42, 0xee, 0x6d, 0x5e, 0x5f,
	0x6e, 0xf4, 0xa2, 0x0b, 0x84, 0x42, 0x1f, 0xa9, 0x7e, 0xa5, 0xd6, 0xea, 0x6e, 0x71, 0xbf, 0x82,
	0xfa, 0xb5, 0xdd, 0x16, 0xcc, 0xe2, 0x45, 0x97, 0x7a, 0x0a, 0x12, 0x9e, 0x7c, 0x3f, 0xe3, 0xb3,
	0xd3, 0xb8, 0xd4, 0x77, 0x98, 0xbe, 0xfd, 0x10, 0x2a, 0x3e, 0x3f, 0x8d, 0x4b, 0x3d, 0x05, 0x8f,
	0x5e, 0x5f, 0x77, 0xec, 0xe0, 0xd8, 0x4f, 0x79, 0x56, 0xef, 0x4f, 0xf5, 0x66, 0x79, 0xc4, 0xad,
	0x05, 0x0e, 0xd6, 0x58, 0xda, 0xe6, 0x59, 0x53, 0x4e, 0xe8, 0x87, 0x78, 0xfc, 0x09, 0x36, 0x07,
	0xd9, 0x9f, 0x52, 0xed, 0x45, 0xa6, 0x0f, 0xff, 0x01, 0x00, 0x00, 0xff, 0xff, 0xef, 0x7f, 0xec,
	0xb8, 0xbc, 0x01, 0x00, 0x00,
}