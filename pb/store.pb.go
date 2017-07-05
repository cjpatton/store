// Code generated by protoc-gen-go. DO NOT EDIT.
// source: store.proto

/*
Package pb is a generated protocol buffer package.

It is generated from these files:
	store.proto

It has these top-level messages:
	StoreParams
	StoreTable
	ShareRequest
	ShareReply
	ParamsRequest
	ParamsReply
*/
package pb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Errors output by the remote procedure calls.
type StoreProviderError int32

const (
	StoreProviderError_OK       StoreProviderError = 0
	StoreProviderError_BAD_USER StoreProviderError = 1
	StoreProviderError_INDEX    StoreProviderError = 2
)

var StoreProviderError_name = map[int32]string{
	0: "OK",
	1: "BAD_USER",
	2: "INDEX",
}
var StoreProviderError_value = map[string]int32{
	"OK":       0,
	"BAD_USER": 1,
	"INDEX":    2,
}

func (x StoreProviderError) String() string {
	return proto.EnumName(StoreProviderError_name, int32(x))
}
func (StoreProviderError) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

// Parameters needed by store.PubStore and store.PrivStore.
type StoreParams struct {
	TableLen       int32  `protobuf:"varint,1,opt,name=table_len,json=tableLen" json:"table_len,omitempty"`
	MaxOutputBytes int32  `protobuf:"varint,2,opt,name=max_output_bytes,json=maxOutputBytes" json:"max_output_bytes,omitempty"`
	RowBytes       int32  `protobuf:"varint,3,opt,name=row_bytes,json=rowBytes" json:"row_bytes,omitempty"`
	TagBytes       int32  `protobuf:"varint,4,opt,name=tag_bytes,json=tagBytes" json:"tag_bytes,omitempty"`
	SaltBytes      int32  `protobuf:"varint,5,opt,name=salt_bytes,json=saltBytes" json:"salt_bytes,omitempty"`
	Salt           []byte `protobuf:"bytes,6,opt,name=salt,proto3" json:"salt,omitempty"`
}

func (m *StoreParams) Reset()                    { *m = StoreParams{} }
func (m *StoreParams) String() string            { return proto.CompactTextString(m) }
func (*StoreParams) ProtoMessage()               {}
func (*StoreParams) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *StoreParams) GetTableLen() int32 {
	if m != nil {
		return m.TableLen
	}
	return 0
}

func (m *StoreParams) GetMaxOutputBytes() int32 {
	if m != nil {
		return m.MaxOutputBytes
	}
	return 0
}

func (m *StoreParams) GetRowBytes() int32 {
	if m != nil {
		return m.RowBytes
	}
	return 0
}

func (m *StoreParams) GetTagBytes() int32 {
	if m != nil {
		return m.TagBytes
	}
	return 0
}

func (m *StoreParams) GetSaltBytes() int32 {
	if m != nil {
		return m.SaltBytes
	}
	return 0
}

func (m *StoreParams) GetSalt() []byte {
	if m != nil {
		return m.Salt
	}
	return nil
}

// A compressed representation of store.PubStore.
type StoreTable struct {
	Params *StoreParams `protobuf:"bytes,1,opt,name=params" json:"params,omitempty"`
	Table  []byte       `protobuf:"bytes,2,opt,name=table,proto3" json:"table,omitempty"`
	Idx    []int32      `protobuf:"varint,3,rep,packed,name=idx" json:"idx,omitempty"`
}

func (m *StoreTable) Reset()                    { *m = StoreTable{} }
func (m *StoreTable) String() string            { return proto.CompactTextString(m) }
func (*StoreTable) ProtoMessage()               {}
func (*StoreTable) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *StoreTable) GetParams() *StoreParams {
	if m != nil {
		return m.Params
	}
	return nil
}

func (m *StoreTable) GetTable() []byte {
	if m != nil {
		return m.Table
	}
	return nil
}

func (m *StoreTable) GetIdx() []int32 {
	if m != nil {
		return m.Idx
	}
	return nil
}

// The share request message.
type ShareRequest struct {
	UserId string `protobuf:"bytes,1,opt,name=user_id,json=userId" json:"user_id,omitempty"`
	X      int32  `protobuf:"varint,2,opt,name=x" json:"x,omitempty"`
	Y      int32  `protobuf:"varint,3,opt,name=y" json:"y,omitempty"`
}

func (m *ShareRequest) Reset()                    { *m = ShareRequest{} }
func (m *ShareRequest) String() string            { return proto.CompactTextString(m) }
func (*ShareRequest) ProtoMessage()               {}
func (*ShareRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *ShareRequest) GetUserId() string {
	if m != nil {
		return m.UserId
	}
	return ""
}

func (m *ShareRequest) GetX() int32 {
	if m != nil {
		return m.X
	}
	return 0
}

func (m *ShareRequest) GetY() int32 {
	if m != nil {
		return m.Y
	}
	return 0
}

// The share response message.
type ShareReply struct {
	PubShare []byte             `protobuf:"bytes,1,opt,name=pub_share,json=pubShare,proto3" json:"pub_share,omitempty"`
	Error    StoreProviderError `protobuf:"varint,2,opt,name=error,enum=pb.StoreProviderError" json:"error,omitempty"`
}

func (m *ShareReply) Reset()                    { *m = ShareReply{} }
func (m *ShareReply) String() string            { return proto.CompactTextString(m) }
func (*ShareReply) ProtoMessage()               {}
func (*ShareReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *ShareReply) GetPubShare() []byte {
	if m != nil {
		return m.PubShare
	}
	return nil
}

func (m *ShareReply) GetError() StoreProviderError {
	if m != nil {
		return m.Error
	}
	return StoreProviderError_OK
}

// The parameters request message.
type ParamsRequest struct {
	UserId string `protobuf:"bytes,1,opt,name=user_id,json=userId" json:"user_id,omitempty"`
}

func (m *ParamsRequest) Reset()                    { *m = ParamsRequest{} }
func (m *ParamsRequest) String() string            { return proto.CompactTextString(m) }
func (*ParamsRequest) ProtoMessage()               {}
func (*ParamsRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *ParamsRequest) GetUserId() string {
	if m != nil {
		return m.UserId
	}
	return ""
}

// The parameters response message.
type ParamsReply struct {
	Params *StoreParams       `protobuf:"bytes,1,opt,name=params" json:"params,omitempty"`
	Error  StoreProviderError `protobuf:"varint,2,opt,name=error,enum=pb.StoreProviderError" json:"error,omitempty"`
}

func (m *ParamsReply) Reset()                    { *m = ParamsReply{} }
func (m *ParamsReply) String() string            { return proto.CompactTextString(m) }
func (*ParamsReply) ProtoMessage()               {}
func (*ParamsReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *ParamsReply) GetParams() *StoreParams {
	if m != nil {
		return m.Params
	}
	return nil
}

func (m *ParamsReply) GetError() StoreProviderError {
	if m != nil {
		return m.Error
	}
	return StoreProviderError_OK
}

func init() {
	proto.RegisterType((*StoreParams)(nil), "pb.StoreParams")
	proto.RegisterType((*StoreTable)(nil), "pb.StoreTable")
	proto.RegisterType((*ShareRequest)(nil), "pb.ShareRequest")
	proto.RegisterType((*ShareReply)(nil), "pb.ShareReply")
	proto.RegisterType((*ParamsRequest)(nil), "pb.ParamsRequest")
	proto.RegisterType((*ParamsReply)(nil), "pb.ParamsReply")
	proto.RegisterEnum("pb.StoreProviderError", StoreProviderError_name, StoreProviderError_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for StoreProvider service

type StoreProviderClient interface {
	GetShare(ctx context.Context, in *ShareRequest, opts ...grpc.CallOption) (*ShareReply, error)
	GetParams(ctx context.Context, in *ParamsRequest, opts ...grpc.CallOption) (*ParamsReply, error)
}

type storeProviderClient struct {
	cc *grpc.ClientConn
}

func NewStoreProviderClient(cc *grpc.ClientConn) StoreProviderClient {
	return &storeProviderClient{cc}
}

func (c *storeProviderClient) GetShare(ctx context.Context, in *ShareRequest, opts ...grpc.CallOption) (*ShareReply, error) {
	out := new(ShareReply)
	err := grpc.Invoke(ctx, "/pb.StoreProvider/GetShare", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storeProviderClient) GetParams(ctx context.Context, in *ParamsRequest, opts ...grpc.CallOption) (*ParamsReply, error) {
	out := new(ParamsReply)
	err := grpc.Invoke(ctx, "/pb.StoreProvider/GetParams", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for StoreProvider service

type StoreProviderServer interface {
	GetShare(context.Context, *ShareRequest) (*ShareReply, error)
	GetParams(context.Context, *ParamsRequest) (*ParamsReply, error)
}

func RegisterStoreProviderServer(s *grpc.Server, srv StoreProviderServer) {
	s.RegisterService(&_StoreProvider_serviceDesc, srv)
}

func _StoreProvider_GetShare_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShareRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StoreProviderServer).GetShare(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.StoreProvider/GetShare",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StoreProviderServer).GetShare(ctx, req.(*ShareRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StoreProvider_GetParams_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ParamsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StoreProviderServer).GetParams(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.StoreProvider/GetParams",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StoreProviderServer).GetParams(ctx, req.(*ParamsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _StoreProvider_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.StoreProvider",
	HandlerType: (*StoreProviderServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetShare",
			Handler:    _StoreProvider_GetShare_Handler,
		},
		{
			MethodName: "GetParams",
			Handler:    _StoreProvider_GetParams_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "store.proto",
}

func init() { proto.RegisterFile("store.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 425 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x53, 0x61, 0x6b, 0xd4, 0x40,
	0x10, 0xed, 0xde, 0x35, 0xf1, 0x32, 0x97, 0x5e, 0xe3, 0x20, 0x7a, 0x54, 0x84, 0x23, 0x5f, 0x0c,
	0x22, 0x07, 0x9e, 0xf8, 0x03, 0xac, 0x3d, 0x4a, 0x51, 0xac, 0xec, 0x29, 0xfa, 0x45, 0xc2, 0x86,
	0x2c, 0xf5, 0x20, 0xd7, 0x5d, 0x37, 0x1b, 0x9b, 0xfc, 0x3c, 0xff, 0x99, 0xec, 0x64, 0x4f, 0x53,
	0xfc, 0xa0, 0xfd, 0xb6, 0xf3, 0xde, 0xec, 0x9b, 0xf7, 0x26, 0x59, 0x98, 0xd6, 0x56, 0x19, 0xb9,
	0xd4, 0x46, 0x59, 0x85, 0x23, 0x5d, 0xa4, 0x3f, 0x19, 0x4c, 0x37, 0x0e, 0xfb, 0x20, 0x8c, 0xd8,
	0xd5, 0xf8, 0x18, 0x22, 0x2b, 0x8a, 0x4a, 0xe6, 0x95, 0xbc, 0x9e, 0xb3, 0x05, 0xcb, 0x02, 0x3e,
	0x21, 0xe0, 0x9d, 0xbc, 0xc6, 0x0c, 0x92, 0x9d, 0x68, 0x73, 0xd5, 0x58, 0xdd, 0xd8, 0xbc, 0xe8,
	0xac, 0xac, 0xe7, 0x23, 0xea, 0x99, 0xed, 0x44, 0x7b, 0x49, 0xf0, 0xa9, 0x43, 0x9d, 0x8c, 0x51,
	0x37, 0xbe, 0x65, 0xdc, 0xcb, 0x18, 0x75, 0xf3, 0x9b, 0xb4, 0xe2, 0xca, 0x93, 0x87, 0xfb, 0x19,
	0x57, 0x3d, 0xf9, 0x04, 0xa0, 0x16, 0xd5, 0x5e, 0x3d, 0x20, 0x36, 0x72, 0x48, 0x4f, 0x23, 0x1c,
	0xba, 0x62, 0x1e, 0x2e, 0x58, 0x16, 0x73, 0x3a, 0xa7, 0x5f, 0x01, 0x28, 0xc2, 0x47, 0xe7, 0x13,
	0x9f, 0x42, 0xa8, 0x29, 0x0b, 0xd9, 0x9f, 0xae, 0x8e, 0x97, 0xba, 0x58, 0x0e, 0x22, 0x72, 0x4f,
	0xe3, 0x03, 0x08, 0x28, 0x19, 0x45, 0x88, 0x79, 0x5f, 0x60, 0x02, 0xe3, 0x6d, 0xd9, 0xce, 0xc7,
	0x8b, 0x71, 0x16, 0x70, 0x77, 0x4c, 0xdf, 0x40, 0xbc, 0xf9, 0x26, 0x8c, 0xe4, 0xf2, 0x7b, 0x23,
	0x6b, 0x8b, 0x8f, 0xe0, 0x5e, 0x53, 0x4b, 0x93, 0x6f, 0x4b, 0x9a, 0x10, 0xf1, 0xd0, 0x95, 0x17,
	0x25, 0xc6, 0xc0, 0x5a, 0xbf, 0x0f, 0xd6, 0xba, 0xaa, 0xf3, 0xd1, 0x59, 0x97, 0x7e, 0x06, 0xf0,
	0x22, 0xba, 0xea, 0xdc, 0x06, 0x74, 0x53, 0xe4, 0xb5, 0x43, 0x48, 0x24, 0xe6, 0x13, 0xdd, 0x14,
	0xd4, 0x81, 0xcf, 0x21, 0x90, 0xc6, 0x28, 0x43, 0x52, 0xb3, 0xd5, 0xc3, 0x3f, 0xfe, 0x8d, 0xfa,
	0xb1, 0x2d, 0xa5, 0x59, 0x3b, 0x96, 0xf7, 0x4d, 0x69, 0x06, 0x47, 0x3e, 0xd7, 0x3f, 0xec, 0xa5,
	0x25, 0x4c, 0xf7, 0x9d, 0xce, 0xc3, 0x7f, 0xef, 0xe9, 0x4e, 0x7e, 0x9e, 0xbd, 0x02, 0xfc, 0x9b,
	0xc4, 0x10, 0x46, 0x97, 0x6f, 0x93, 0x03, 0x8c, 0x61, 0x72, 0xfa, 0xfa, 0x2c, 0xff, 0xb4, 0x59,
	0xf3, 0x84, 0x61, 0x04, 0xc1, 0xc5, 0xfb, 0xb3, 0xf5, 0x97, 0x64, 0xb4, 0x32, 0x70, 0x74, 0xeb,
	0x1a, 0x2e, 0x61, 0x72, 0x2e, 0x6d, 0xbf, 0x91, 0x84, 0x46, 0x0e, 0xbe, 0xc1, 0xc9, 0x6c, 0x80,
	0xe8, 0xaa, 0x4b, 0x0f, 0xf0, 0x05, 0x44, 0xe7, 0xd2, 0xfa, 0xbf, 0xf8, 0xbe, 0xa3, 0x6f, 0xad,
	0xe5, 0xe4, 0x78, 0x08, 0xd1, 0x95, 0x22, 0xa4, 0x67, 0xf0, 0xf2, 0x57, 0x00, 0x00, 0x00, 0xff,
	0xff, 0x47, 0x5a, 0xc8, 0x2e, 0x15, 0x03, 0x00, 0x00,
}