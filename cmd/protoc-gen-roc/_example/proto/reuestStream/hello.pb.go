// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: hello.proto

package hello

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
	client "github.com/go-roc/roc/client"
	parcel "github.com/go-roc/roc/parcel"
	context "github.com/go-roc/roc/parcel/context"
	server "github.com/go-roc/roc/server"
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

type SayReq struct {
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (m *SayReq) Reset()         { *m = SayReq{} }
func (m *SayReq) String() string { return proto.CompactTextString(m) }
func (*SayReq) ProtoMessage()    {}
func (*SayReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_61ef911816e0a8ce, []int{0}
}
func (m *SayReq) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SayReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SayReq.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SayReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SayReq.Merge(m, src)
}
func (m *SayReq) XXX_Size() int {
	return m.Size()
}
func (m *SayReq) XXX_DiscardUnknown() {
	xxx_messageInfo_SayReq.DiscardUnknown(m)
}

var xxx_messageInfo_SayReq proto.InternalMessageInfo

func (m *SayReq) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type SayRsp struct {
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
}

func (m *SayRsp) Reset()         { *m = SayRsp{} }
func (m *SayRsp) String() string { return proto.CompactTextString(m) }
func (*SayRsp) ProtoMessage()    {}
func (*SayRsp) Descriptor() ([]byte, []int) {
	return fileDescriptor_61ef911816e0a8ce, []int{1}
}
func (m *SayRsp) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SayRsp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SayRsp.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SayRsp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SayRsp.Merge(m, src)
}
func (m *SayRsp) XXX_Size() int {
	return m.Size()
}
func (m *SayRsp) XXX_DiscardUnknown() {
	xxx_messageInfo_SayRsp.DiscardUnknown(m)
}

var xxx_messageInfo_SayRsp proto.InternalMessageInfo

func (m *SayRsp) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type RocReq struct {
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (m *RocReq) Reset()         { *m = RocReq{} }
func (m *RocReq) String() string { return proto.CompactTextString(m) }
func (*RocReq) ProtoMessage()    {}
func (*RocReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_61ef911816e0a8ce, []int{2}
}
func (m *RocReq) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *RocReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_RocReq.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *RocReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RocReq.Merge(m, src)
}
func (m *RocReq) XXX_Size() int {
	return m.Size()
}
func (m *RocReq) XXX_DiscardUnknown() {
	xxx_messageInfo_RocReq.DiscardUnknown(m)
}

var xxx_messageInfo_RocReq proto.InternalMessageInfo

func (m *RocReq) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type RocRsp struct {
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
}

func (m *RocRsp) Reset()         { *m = RocRsp{} }
func (m *RocRsp) String() string { return proto.CompactTextString(m) }
func (*RocRsp) ProtoMessage()    {}
func (*RocRsp) Descriptor() ([]byte, []int) {
	return fileDescriptor_61ef911816e0a8ce, []int{3}
}
func (m *RocRsp) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *RocRsp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_RocRsp.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *RocRsp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RocRsp.Merge(m, src)
}
func (m *RocRsp) XXX_Size() int {
	return m.Size()
}
func (m *RocRsp) XXX_DiscardUnknown() {
	xxx_messageInfo_RocRsp.DiscardUnknown(m)
}

var xxx_messageInfo_RocRsp proto.InternalMessageInfo

func (m *RocRsp) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func init() {
	proto.RegisterType((*SayReq)(nil), "SayReq")
	proto.RegisterType((*SayRsp)(nil), "SayRsp")
	proto.RegisterType((*RocReq)(nil), "RocReq")
	proto.RegisterType((*RocRsp)(nil), "RocRsp")
}

func init() { proto.RegisterFile("hello.proto", fileDescriptor_61ef911816e0a8ce) }

var fileDescriptor_61ef911816e0a8ce = []byte{
	// 158 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xce, 0x48, 0xcd, 0xc9,
	0xc9, 0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x57, 0x92, 0xe1, 0x62, 0x0b, 0x4e, 0xac, 0x0c, 0x4a,
	0x2d, 0x14, 0x12, 0xe2, 0x62, 0xc9, 0x4b, 0xcc, 0x4d, 0x95, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0c,
	0x02, 0xb3, 0x61, 0xb2, 0xc5, 0x05, 0x70, 0x59, 0x26, 0x54, 0xd9, 0xa0, 0xfc, 0x64, 0x3c, 0x7a,
	0x41, 0xb2, 0xd8, 0xf5, 0x1a, 0x39, 0x72, 0xb1, 0x7a, 0x80, 0x9c, 0x21, 0x24, 0xcd, 0xc5, 0x1c,
	0x9c, 0x58, 0x29, 0xc4, 0xae, 0x07, 0x71, 0x86, 0x14, 0x84, 0x51, 0x5c, 0xa0, 0xc4, 0xa0, 0xc1,
	0x08, 0x92, 0x0c, 0xca, 0x4f, 0x16, 0x62, 0xd7, 0x83, 0xd8, 0x23, 0x05, 0x61, 0x40, 0x24, 0x9d,
	0x24, 0x4e, 0x3c, 0x92, 0x63, 0xbc, 0xf0, 0x48, 0x8e, 0xf1, 0xc1, 0x23, 0x39, 0xc6, 0x09, 0x8f,
	0xe5, 0x18, 0x2e, 0x3c, 0x96, 0x63, 0xb8, 0xf1, 0x58, 0x8e, 0x21, 0x89, 0x0d, 0xec, 0x37, 0x63,
	0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0x8a, 0x20, 0x35, 0x61, 0xea, 0x00, 0x00, 0x00,
}

func (m *SayReq) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SayReq) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SayReq) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Name) > 0 {
		i -= len(m.Name)
		copy(dAtA[i:], m.Name)
		i = encodeVarintHello(dAtA, i, uint64(len(m.Name)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *SayRsp) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SayRsp) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SayRsp) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Name) > 0 {
		i -= len(m.Name)
		copy(dAtA[i:], m.Name)
		i = encodeVarintHello(dAtA, i, uint64(len(m.Name)))
		i--
		dAtA[i] = 0x12
	}
	return len(dAtA) - i, nil
}

func (m *RocReq) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RocReq) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *RocReq) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Name) > 0 {
		i -= len(m.Name)
		copy(dAtA[i:], m.Name)
		i = encodeVarintHello(dAtA, i, uint64(len(m.Name)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *RocRsp) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RocRsp) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *RocRsp) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Name) > 0 {
		i -= len(m.Name)
		copy(dAtA[i:], m.Name)
		i = encodeVarintHello(dAtA, i, uint64(len(m.Name)))
		i--
		dAtA[i] = 0x12
	}
	return len(dAtA) - i, nil
}

func encodeVarintHello(dAtA []byte, offset int, v uint64) int {
	offset -= sovHello(v)
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
var _ client.RocClient
var _ server.RocServer
var _ parcel.RocPacket

// This is a compile-time assertion to ensure that this generated file
// is compatible with the roc package it is being compiled against.
const _ = server.SupportPackageIsVersion1

type HelloClient interface {
	Say(c *context.Context, req *SayReq, opts ...client.InvokeOptions) (chan *SayRsp, chan error)
	Roc(c *context.Context, req *RocReq, opts ...client.InvokeOptions) (chan *RocRsp, chan error)
}

type helloClient struct {
	c *client.RocClient
}

func NewHelloClient(c *client.RocClient) HelloClient {
	return &helloClient{c}
}

func (cc *helloClient) Say(c *context.Context, req *SayReq, opts ...client.InvokeOptions) (chan *SayRsp, chan error) {
	data, errs := cc.c.InvokeRS(c, "Hello.Say", req, opts...)
	var rsp = make(chan *SayRsp)
	go func() {
		for b := range data {
			v := &SayRsp{}
			err := cc.c.Codec().Decode(b, v)
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

func (cc *helloClient) Roc(c *context.Context, req *RocReq, opts ...client.InvokeOptions) (chan *RocRsp, chan error) {
	data, errs := cc.c.InvokeRS(c, "Hello.Roc", req, opts...)
	var rsp = make(chan *RocRsp)
	go func() {
		for b := range data {
			v := &RocRsp{}
			err := cc.c.Codec().Decode(b, v)
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

// HelloServer is the server API for Hello service.
type HelloServer interface {
	Say(c *context.Context, req *SayReq) (chan *SayRsp, chan error)
	Roc(c *context.Context, req *RocReq) (chan *RocRsp, chan error)
}

func RegisterHelloServer(s *server.RocServer, h HelloServer) {
	var r = &helloHandler{h: h, s: s}
	s.RegisterStreamHandler("Hello.Say", r.Say)
	s.RegisterStreamHandler("Hello.Roc", r.Roc)
}

type helloHandler struct {
	h HelloServer
	s *server.RocServer
}

func (r *helloHandler) Say(c *context.Context, req *parcel.RocPacket) (chan proto.Message, chan error) {
	var errs = make(chan error)
	var in SayReq
	err := r.s.Codec().Decode(req.Bytes(), &in)
	if err != nil {
		errs <- err
		close(errs)
		return nil, errs
	}

	out, outErrs := r.h.Say(c, &in)
	var rsp = make(chan proto.Message, len(out))

	go func() {
	QUIT:
		for {
			select {
			case d, ok := <-out:
				if ok {
					rsp <- d
				} else {
					break QUIT
				}
			case err := <-outErrs:
				errs <- err
				break QUIT
			}
		}
		close(rsp)
		close(errs)
	}()
	return rsp, errs
}

func (r *helloHandler) Roc(c *context.Context, req *parcel.RocPacket) (chan proto.Message, chan error) {
	var errs = make(chan error)
	var in RocReq
	err := r.s.Codec().Decode(req.Bytes(), &in)
	if err != nil {
		errs <- err
		close(errs)
		return nil, errs
	}

	out, outErrs := r.h.Roc(c, &in)
	var rsp = make(chan proto.Message, len(out))

	go func() {
	QUIT:
		for {
			select {
			case d, ok := <-out:
				if ok {
					rsp <- d
				} else {
					break QUIT
				}
			case err := <-outErrs:
				errs <- err
				break QUIT
			}
		}
		close(rsp)
		close(errs)
	}()
	return rsp, errs
}

func (m *SayReq) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovHello(uint64(l))
	}
	return n
}

func (m *SayRsp) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovHello(uint64(l))
	}
	return n
}

func (m *RocReq) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovHello(uint64(l))
	}
	return n
}

func (m *RocRsp) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovHello(uint64(l))
	}
	return n
}

func sovHello(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozHello(x uint64) (n int) {
	return sovHello(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *SayReq) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowHello
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
			return fmt.Errorf("proto: SayReq: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SayReq: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHello
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
				return ErrInvalidLengthHello
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthHello
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipHello(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthHello
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
func (m *SayRsp) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowHello
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
			return fmt.Errorf("proto: SayRsp: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SayRsp: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHello
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
				return ErrInvalidLengthHello
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthHello
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipHello(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthHello
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
func (m *RocReq) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowHello
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
			return fmt.Errorf("proto: RocReq: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RocReq: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHello
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
				return ErrInvalidLengthHello
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthHello
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipHello(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthHello
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
func (m *RocRsp) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowHello
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
			return fmt.Errorf("proto: RocRsp: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RocRsp: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHello
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
				return ErrInvalidLengthHello
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthHello
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipHello(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthHello
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
func skipHello(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowHello
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
					return 0, ErrIntOverflowHello
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
					return 0, ErrIntOverflowHello
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
				return 0, ErrInvalidLengthHello
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupHello
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthHello
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthHello        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowHello          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupHello = fmt.Errorf("proto: unexpected end of group")
)