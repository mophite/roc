// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: proto/phello/pim.proto

package phello

import (
	fmt "fmt"
	parcel "github.com/go-roc/roc/parcel"
	context "github.com/go-roc/roc/parcel/context"
	service "github.com/go-roc/roc/service"
	client "github.com/go-roc/roc/service/client"
	handler "github.com/go-roc/roc/service/handler"
	invoke "github.com/go-roc/roc/service/invoke"
	server "github.com/go-roc/roc/service/server"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// ConnectReq is like handshake
type ConnectReq struct {
	UserName string `protobuf:"bytes,1,opt,name=user_name,json=userName,proto3" json:"user_name,omitempty"`
}

func (m *ConnectReq) Reset()         { *m = ConnectReq{} }
func (m *ConnectReq) String() string { return proto.CompactTextString(m) }
func (*ConnectReq) ProtoMessage()    {}
func (*ConnectReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_23025222ae894cc8, []int{0}
}
func (m *ConnectReq) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ConnectReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ConnectReq.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ConnectReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConnectReq.Merge(m, src)
}
func (m *ConnectReq) XXX_Size() int {
	return m.Size()
}
func (m *ConnectReq) XXX_DiscardUnknown() {
	xxx_messageInfo_ConnectReq.DiscardUnknown(m)
}

var xxx_messageInfo_ConnectReq proto.InternalMessageInfo

func (m *ConnectReq) GetUserName() string {
	if m != nil {
		return m.UserName
	}
	return ""
}

type ConnectRsp struct {
	IsConnect bool `protobuf:"varint,1,opt,name=is_connect,json=isConnect,proto3" json:"is_connect,omitempty"`
}

func (m *ConnectRsp) Reset()         { *m = ConnectRsp{} }
func (m *ConnectRsp) String() string { return proto.CompactTextString(m) }
func (*ConnectRsp) ProtoMessage()    {}
func (*ConnectRsp) Descriptor() ([]byte, []int) {
	return fileDescriptor_23025222ae894cc8, []int{1}
}
func (m *ConnectRsp) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ConnectRsp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ConnectRsp.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ConnectRsp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConnectRsp.Merge(m, src)
}
func (m *ConnectRsp) XXX_Size() int {
	return m.Size()
}
func (m *ConnectRsp) XXX_DiscardUnknown() {
	xxx_messageInfo_ConnectRsp.DiscardUnknown(m)
}

var xxx_messageInfo_ConnectRsp proto.InternalMessageInfo

func (m *ConnectRsp) GetIsConnect() bool {
	if m != nil {
		return m.IsConnect
	}
	return false
}

// CountReq is for count online member
type CountReq struct {
	Prefix string `protobuf:"bytes,1,opt,name=prefix,proto3" json:"prefix,omitempty"`
}

func (m *CountReq) Reset()         { *m = CountReq{} }
func (m *CountReq) String() string { return proto.CompactTextString(m) }
func (*CountReq) ProtoMessage()    {}
func (*CountReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_23025222ae894cc8, []int{2}
}
func (m *CountReq) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *CountReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_CountReq.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *CountReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CountReq.Merge(m, src)
}
func (m *CountReq) XXX_Size() int {
	return m.Size()
}
func (m *CountReq) XXX_DiscardUnknown() {
	xxx_messageInfo_CountReq.DiscardUnknown(m)
}

var xxx_messageInfo_CountReq proto.InternalMessageInfo

func (m *CountReq) GetPrefix() string {
	if m != nil {
		return m.Prefix
	}
	return ""
}

type CountRsp struct {
	Count uint32 `protobuf:"varint,1,opt,name=count,proto3" json:"count,omitempty"`
}

func (m *CountRsp) Reset()         { *m = CountRsp{} }
func (m *CountRsp) String() string { return proto.CompactTextString(m) }
func (*CountRsp) ProtoMessage()    {}
func (*CountRsp) Descriptor() ([]byte, []int) {
	return fileDescriptor_23025222ae894cc8, []int{3}
}
func (m *CountRsp) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *CountRsp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_CountRsp.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *CountRsp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CountRsp.Merge(m, src)
}
func (m *CountRsp) XXX_Size() int {
	return m.Size()
}
func (m *CountRsp) XXX_DiscardUnknown() {
	xxx_messageInfo_CountRsp.DiscardUnknown(m)
}

var xxx_messageInfo_CountRsp proto.InternalMessageInfo

func (m *CountRsp) GetCount() uint32 {
	if m != nil {
		return m.Count
	}
	return 0
}

// SendMessageReq send a message
type SendMessageReq struct {
	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (m *SendMessageReq) Reset()         { *m = SendMessageReq{} }
func (m *SendMessageReq) String() string { return proto.CompactTextString(m) }
func (*SendMessageReq) ProtoMessage()    {}
func (*SendMessageReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_23025222ae894cc8, []int{4}
}
func (m *SendMessageReq) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SendMessageReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SendMessageReq.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SendMessageReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SendMessageReq.Merge(m, src)
}
func (m *SendMessageReq) XXX_Size() int {
	return m.Size()
}
func (m *SendMessageReq) XXX_DiscardUnknown() {
	xxx_messageInfo_SendMessageReq.DiscardUnknown(m)
}

var xxx_messageInfo_SendMessageReq proto.InternalMessageInfo

func (m *SendMessageReq) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

// response content.
// SendMessageRsp usually use for broadcast
type SendMessageRsp struct {
	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (m *SendMessageRsp) Reset()         { *m = SendMessageRsp{} }
func (m *SendMessageRsp) String() string { return proto.CompactTextString(m) }
func (*SendMessageRsp) ProtoMessage()    {}
func (*SendMessageRsp) Descriptor() ([]byte, []int) {
	return fileDescriptor_23025222ae894cc8, []int{5}
}
func (m *SendMessageRsp) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SendMessageRsp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SendMessageRsp.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SendMessageRsp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SendMessageRsp.Merge(m, src)
}
func (m *SendMessageRsp) XXX_Size() int {
	return m.Size()
}
func (m *SendMessageRsp) XXX_DiscardUnknown() {
	xxx_messageInfo_SendMessageRsp.DiscardUnknown(m)
}

var xxx_messageInfo_SendMessageRsp proto.InternalMessageInfo

func (m *SendMessageRsp) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func init() {
	proto.RegisterType((*ConnectReq)(nil), "phello.ConnectReq")
	proto.RegisterType((*ConnectRsp)(nil), "phello.ConnectRsp")
	proto.RegisterType((*CountReq)(nil), "phello.CountReq")
	proto.RegisterType((*CountRsp)(nil), "phello.CountRsp")
	proto.RegisterType((*SendMessageReq)(nil), "phello.SendMessageReq")
	proto.RegisterType((*SendMessageRsp)(nil), "phello.SendMessageRsp")
}

func init() { proto.RegisterFile("proto/phello/pim.proto", fileDescriptor_23025222ae894cc8) }

var fileDescriptor_23025222ae894cc8 = []byte{
	// 288 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x91, 0x41, 0x4e, 0xb3, 0x40,
	0x14, 0xc7, 0x99, 0x2f, 0x29, 0x85, 0xf7, 0x45, 0x63, 0x5e, 0x0c, 0x21, 0x18, 0x27, 0xcd, 0xac,
	0xaa, 0x46, 0x6a, 0xec, 0x0d, 0x64, 0xe5, 0x42, 0x17, 0x78, 0x80, 0x06, 0xf1, 0xa9, 0x24, 0x05,
	0x46, 0xa6, 0x4d, 0x3c, 0x86, 0x17, 0xf1, 0x1e, 0x2e, 0xbb, 0x74, 0x69, 0xe0, 0x22, 0x86, 0x19,
	0x08, 0x5a, 0x75, 0xc7, 0xff, 0xc7, 0x2f, 0xff, 0xbc, 0xf7, 0x06, 0x3c, 0x59, 0x95, 0xab, 0x72,
	0x26, 0x1f, 0x69, 0xb9, 0x2c, 0x67, 0x32, 0xcb, 0x43, 0x0d, 0xd0, 0x36, 0x44, 0x1c, 0x01, 0x44,
	0x65, 0x51, 0x50, 0xba, 0x8a, 0xe9, 0x09, 0x0f, 0xc0, 0x5d, 0x2b, 0xaa, 0x16, 0x45, 0x92, 0x93,
	0xcf, 0x26, 0x6c, 0xea, 0xc6, 0x4e, 0x0b, 0xae, 0x93, 0x9c, 0xc4, 0xc9, 0xa0, 0x2a, 0x89, 0x87,
	0x00, 0x99, 0x5a, 0xa4, 0x06, 0x68, 0xd7, 0x89, 0xdd, 0x4c, 0x75, 0x86, 0x10, 0xe0, 0x44, 0xe5,
	0xba, 0xd0, 0xad, 0x1e, 0xd8, 0xb2, 0xa2, 0xfb, 0xec, 0xb9, 0xab, 0xec, 0x92, 0x98, 0xf4, 0x8e,
	0x92, 0xb8, 0x0f, 0xa3, 0xb4, 0xfd, 0xd6, 0xca, 0x4e, 0x6c, 0x82, 0x38, 0x86, 0xdd, 0x1b, 0x2a,
	0xee, 0xae, 0x48, 0xa9, 0xe4, 0x81, 0xda, 0x2e, 0x1f, 0xc6, 0xb9, 0x49, 0x5d, 0x59, 0x1f, 0xb7,
	0x5d, 0x25, 0xff, 0x76, 0xcf, 0x5f, 0x19, 0xfc, 0xbb, 0xcc, 0x71, 0x0e, 0xe3, 0x6e, 0x5e, 0xc4,
	0xd0, 0x1c, 0x24, 0x1c, 0xae, 0x11, 0xfc, 0x60, 0x4a, 0x0a, 0x0b, 0x4f, 0x61, 0xa4, 0xa7, 0xc6,
	0xbd, 0xe1, 0xb7, 0x59, 0x34, 0xd8, 0x22, 0x5a, 0x8f, 0xe0, 0xff, 0x97, 0xb1, 0xd0, 0xeb, 0x95,
	0xef, 0x7b, 0x05, 0xbf, 0xf2, 0xb6, 0x60, 0xca, 0xce, 0xd8, 0x85, 0xff, 0x56, 0x73, 0xb6, 0xa9,
	0x39, 0xfb, 0xa8, 0x39, 0x7b, 0x69, 0xb8, 0xb5, 0x69, 0xb8, 0xf5, 0xde, 0x70, 0xeb, 0xd6, 0xd6,
	0xcf, 0x39, 0xff, 0x0c, 0x00, 0x00, 0xff, 0xff, 0x35, 0x0b, 0x15, 0xbe, 0xe8, 0x01, 0x00, 0x00,
}

func (m *ConnectReq) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ConnectReq) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ConnectReq) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.UserName) > 0 {
		i -= len(m.UserName)
		copy(dAtA[i:], m.UserName)
		i = encodeVarintPim(dAtA, i, uint64(len(m.UserName)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *ConnectRsp) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ConnectRsp) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ConnectRsp) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.IsConnect {
		i--
		if m.IsConnect {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *CountReq) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *CountReq) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *CountReq) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Prefix) > 0 {
		i -= len(m.Prefix)
		copy(dAtA[i:], m.Prefix)
		i = encodeVarintPim(dAtA, i, uint64(len(m.Prefix)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *CountRsp) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *CountRsp) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *CountRsp) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Count != 0 {
		i = encodeVarintPim(dAtA, i, uint64(m.Count))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *SendMessageReq) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SendMessageReq) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SendMessageReq) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Message) > 0 {
		i -= len(m.Message)
		copy(dAtA[i:], m.Message)
		i = encodeVarintPim(dAtA, i, uint64(len(m.Message)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *SendMessageRsp) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SendMessageRsp) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SendMessageRsp) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Message) > 0 {
		i -= len(m.Message)
		copy(dAtA[i:], m.Message)
		i = encodeVarintPim(dAtA, i, uint64(len(m.Message)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintPim(dAtA []byte, offset int, v uint64) int {
	offset -= sovPim(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ invoke.Invoke
var _ handler.Handler
var _ service.Service
var _ parcel.RocPacket
var _ client.Client
var _ server.Server

// This is a compile-time assertion to ensure that this generated file
// is compatible with the roc package it is being compiled against.
const _ = service.SupportPackageIsVersion1

type ImClient interface {
	// Connect server for wait message
	Connect(c *context.Context, req *ConnectReq, opts ...invoke.InvokeOptions) (*ConnectRsp, error)
	// Count online member
	Count(c *context.Context, req *CountReq, opts ...invoke.InvokeOptions) (*CountRsp, error)
	// SendMessage is the im kernel
	SendMessage(c *context.Context, req chan *SendMessageReq, errIn chan error, opts ...invoke.InvokeOptions) (chan *SendMessageRsp, chan error)
}

type imClient struct {
	c *client.Client
}

func NewImClient(c *client.Client) ImClient {
	return &imClient{c}
}

func (cc *imClient) Connect(c *context.Context, req *ConnectReq, opts ...invoke.InvokeOptions) (*ConnectRsp, error) {
	rsp := &ConnectRsp{}
	err := cc.c.InvokeRR(c, service.GetApiPrefix()+"im/connect", req, rsp, opts...)
	return rsp, err
}

func (cc *imClient) Count(c *context.Context, req *CountReq, opts ...invoke.InvokeOptions) (*CountRsp, error) {
	rsp := &CountRsp{}
	err := cc.c.InvokeRR(c, service.GetApiPrefix()+"im/count", req, rsp, opts...)
	return rsp, err
}

func (cc *imClient) SendMessage(c *context.Context, req chan *SendMessageReq, errIn chan error, opts ...invoke.InvokeOptions) (chan *SendMessageRsp, chan error) {
	var in = make(chan []byte)
	go func() {
		for b := range req {
			v, err := c.Codec().Encode(b)
			if err != nil {
				errIn <- err
				break
			}
			in <- v
		}
		close(in)
	}()

	data, errs := cc.c.InvokeRC(c, service.GetApiPrefix()+"im/sendmessage", in, errIn, opts...)
	var rsp = make(chan *SendMessageRsp)
	go func() {
		for b := range data {
			v := &SendMessageRsp{}
			err := c.Codec().Decode(b, v)
			if err != nil {
				errs <- err
				break
			}
			rsp <- v
		}
		close(rsp)
	}()
	return rsp, errs
}

// ImServer is the server API for Im server.
type ImServer interface {
	// Connect server for wait message
	Connect(c *context.Context, req *ConnectReq, rsp *ConnectRsp) (err error)
	// Count online member
	Count(c *context.Context, req *CountReq, rsp *CountRsp) (err error)
	// SendMessage is the im kernel
	SendMessage(c *context.Context, req chan *SendMessageReq, errIn chan error) (chan *SendMessageRsp, chan error)
}

func RegisterImServer(s *server.Server, h ImServer) {
	var r = &imHandler{h: h, s: s}
	s.RegisterHandler(service.GetApiPrefix()+"im/connect", r.Connect)
	s.RegisterHandler(service.GetApiPrefix()+"im/count", r.Count)
	s.RegisterChannelHandler(service.GetApiPrefix()+"im/sendmessage", r.SendMessage)
}

type imHandler struct {
	h ImServer
	s *server.Server
}

func (r *imHandler) Connect(c *context.Context, req *parcel.RocPacket, interrupt handler.Interceptor) (rsp proto.Message, err error) {
	var in ConnectReq
	err = c.Codec().Decode(req.Bytes(), &in)
	if err != nil {
		return nil, err
	}
	var out = ConnectRsp{}
	if interrupt == nil {
		err = r.h.Connect(c, &in, &out)
		return &out, err
	}
	f := func(c *context.Context, req proto.Message) (proto.Message, error) {
		err = r.h.Connect(c, req.(*ConnectReq), &out)
		return &out, err
	}
	return interrupt(c, &in, f)
}

func (r *imHandler) Count(c *context.Context, req *parcel.RocPacket, interrupt handler.Interceptor) (rsp proto.Message, err error) {
	var in CountReq
	err = c.Codec().Decode(req.Bytes(), &in)
	if err != nil {
		return nil, err
	}
	var out = CountRsp{}
	if interrupt == nil {
		err = r.h.Count(c, &in, &out)
		return &out, err
	}
	f := func(c *context.Context, req proto.Message) (proto.Message, error) {
		err = r.h.Count(c, req.(*CountReq), &out)
		return &out, err
	}
	return interrupt(c, &in, f)
}

func (r *imHandler) SendMessage(c *context.Context, req chan *parcel.RocPacket, errIn chan error) (chan proto.Message, chan error) {
	var in = make(chan *SendMessageReq)
	go func() {
		for b := range req {
			var v = &SendMessageReq{}
			err := c.Codec().Decode(b.Bytes(), v)
			if err != nil {
				errIn <- err
				break
			}
			in <- v
		}
		close(in)
	}()

	out, outErrs := r.h.SendMessage(c, in, errIn)
	var rsp = make(chan proto.Message)

	go func() {
		for d := range out {
			rsp <- d
		}
		close(rsp)
	}()
	return rsp, outErrs
}

func (m *ConnectReq) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.UserName)
	if l > 0 {
		n += 1 + l + sovPim(uint64(l))
	}
	return n
}

func (m *ConnectRsp) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.IsConnect {
		n += 2
	}
	return n
}

func (m *CountReq) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Prefix)
	if l > 0 {
		n += 1 + l + sovPim(uint64(l))
	}
	return n
}

func (m *CountRsp) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Count != 0 {
		n += 1 + sovPim(uint64(m.Count))
	}
	return n
}

func (m *SendMessageReq) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Message)
	if l > 0 {
		n += 1 + l + sovPim(uint64(l))
	}
	return n
}

func (m *SendMessageRsp) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Message)
	if l > 0 {
		n += 1 + l + sovPim(uint64(l))
	}
	return n
}

func sovPim(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozPim(x uint64) (n int) {
	return sovPim(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ConnectReq) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPim
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ConnectReq: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ConnectReq: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UserName", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPim
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPim
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPim
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.UserName = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPim(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthPim
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *ConnectRsp) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPim
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ConnectRsp: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ConnectRsp: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IsConnect", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPim
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.IsConnect = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipPim(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthPim
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *CountReq) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPim
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: CountReq: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: CountReq: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Prefix", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPim
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPim
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPim
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Prefix = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPim(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthPim
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *CountRsp) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPim
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: CountRsp: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: CountRsp: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Count", wireType)
			}
			m.Count = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPim
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Count |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipPim(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthPim
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *SendMessageReq) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPim
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: SendMessageReq: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SendMessageReq: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Message", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPim
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPim
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPim
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Message = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPim(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthPim
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *SendMessageRsp) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPim
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: SendMessageRsp: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SendMessageRsp: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Message", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPim
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPim
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPim
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Message = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPim(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthPim
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipPim(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowPim
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowPim
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowPim
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthPim
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupPim
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthPim
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthPim        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowPim          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupPim = fmt.Errorf("proto: unexpected end of group")
)
