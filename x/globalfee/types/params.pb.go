// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: gaia/globalfee/v1beta1/params.proto

package types

import (
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
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

// Params defines the set of module parameters.
type Params struct {
	// minimum_gas_prices stores the minimum gas price(s) for all TX on the chain.
	// When multiple coins are defined then they are accepted alternatively.
	// The list must be sorted by denoms asc. No duplicate denoms or zero amount
	// values allowed. For more information see
	// https://docs.cosmos.network/main/modules/auth#concepts
	MinimumGasPrices github_com_cosmos_cosmos_sdk_types.DecCoins `protobuf:"bytes,1,rep,name=minimum_gas_prices,json=minimumGasPrices,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.DecCoins" json:"minimum_gas_prices,omitempty" yaml:"minimum_gas_prices"`
	// bypass_min_fee_msg_types defines a list of message type urls
	// that are free of fee charge.
	BypassMinFeeMsgTypes []string `protobuf:"bytes,2,rep,name=bypass_min_fee_msg_types,json=bypassMinFeeMsgTypes,proto3" json:"bypass_min_fee_msg_types,omitempty" yaml:"bypass_min_fee_msg_types"`
	// max_total_bypass_min_fee_msg_gas_usage defines the total maximum gas usage
	// allowed for a transaction containing only messages of types in bypass_min_fee_msg_types
	// to bypass fee charge.
	MaxTotalBypassMinFeeMsgGasUsage uint64 `protobuf:"varint,3,opt,name=max_total_bypass_min_fee_msg_gas_usage,json=maxTotalBypassMinFeeMsgGasUsage,proto3" json:"max_total_bypass_min_fee_msg_gas_usage,omitempty"`
}

func (m *Params) Reset()         { *m = Params{} }
func (m *Params) String() string { return proto.CompactTextString(m) }
func (*Params) ProtoMessage()    {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_f135cd41f9af437e, []int{0}
}
func (m *Params) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Params) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Params.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Params) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Params.Merge(m, src)
}
func (m *Params) XXX_Size() int {
	return m.Size()
}
func (m *Params) XXX_DiscardUnknown() {
	xxx_messageInfo_Params.DiscardUnknown(m)
}

var xxx_messageInfo_Params proto.InternalMessageInfo

func (m *Params) GetMinimumGasPrices() github_com_cosmos_cosmos_sdk_types.DecCoins {
	if m != nil {
		return m.MinimumGasPrices
	}
	return nil
}

func (m *Params) GetBypassMinFeeMsgTypes() []string {
	if m != nil {
		return m.BypassMinFeeMsgTypes
	}
	return nil
}

func (m *Params) GetMaxTotalBypassMinFeeMsgGasUsage() uint64 {
	if m != nil {
		return m.MaxTotalBypassMinFeeMsgGasUsage
	}
	return 0
}

func init() {
	proto.RegisterType((*Params)(nil), "gaia.globalfee.v1beta1.Params")
}

func init() {
	proto.RegisterFile("gaia/globalfee/v1beta1/params.proto", fileDescriptor_f135cd41f9af437e)
}

var fileDescriptor_f135cd41f9af437e = []byte{
	// 402 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x92, 0xc1, 0xaa, 0xd3, 0x40,
	0x14, 0x86, 0x13, 0x23, 0x17, 0x8c, 0x1b, 0x09, 0x17, 0x89, 0x97, 0x4b, 0x52, 0x22, 0x48, 0x41,
	0x6f, 0x86, 0xeb, 0xc5, 0x8d, 0xcb, 0x28, 0x16, 0x17, 0xc5, 0x52, 0xea, 0xc6, 0xcd, 0x70, 0x12,
	0xa7, 0xe3, 0x60, 0x26, 0x13, 0x72, 0x26, 0xa5, 0x59, 0xfa, 0x06, 0x3e, 0x80, 0x4f, 0xe0, 0x33,
	0xf8, 0x00, 0x5d, 0x76, 0xe9, 0x2a, 0x4a, 0xbb, 0xeb, 0xd2, 0x27, 0x90, 0x24, 0xad, 0xb6, 0xf4,
	0x76, 0x95, 0xc3, 0xc9, 0xf7, 0xff, 0xe7, 0x3f, 0x33, 0x63, 0x3f, 0xe6, 0x20, 0x80, 0xf0, 0x54,
	0xc5, 0x90, 0x4e, 0x19, 0x23, 0xb3, 0xeb, 0x98, 0x69, 0xb8, 0x26, 0x39, 0x14, 0x20, 0x31, 0xcc,
	0x0b, 0xa5, 0x95, 0xf3, 0xb0, 0x81, 0xc2, 0x7f, 0x50, 0xb8, 0x85, 0x2e, 0xbc, 0x44, 0xa1, 0x54,
	0x48, 0x62, 0xc0, 0xff, 0xca, 0x44, 0x89, 0xac, 0xd3, 0x5d, 0x9c, 0x73, 0xc5, 0x55, 0x5b, 0x92,
	0xa6, 0xea, 0xba, 0xc1, 0x37, 0xcb, 0x3e, 0x1b, 0xb5, 0xf6, 0xce, 0x0f, 0xd3, 0x76, 0xa4, 0xc8,
	0x84, 0x2c, 0x25, 0xe5, 0x80, 0x34, 0x2f, 0x44, 0xc2, 0xd0, 0x35, 0x7b, 0x56, 0xff, 0xfe, 0xf3,
	0xcb, 0xb0, 0xb3, 0x0f, 0x1b, 0xfb, 0xdd, 0xcc, 0xf0, 0x35, 0x4b, 0x5e, 0x29, 0x91, 0x45, 0xf9,
	0xa2, 0xf6, 0x8d, 0x4d, 0xed, 0x5f, 0x1e, 0xeb, 0x9f, 0x29, 0x29, 0x34, 0x93, 0xb9, 0xae, 0xfe,
	0xd4, 0xfe, 0xa3, 0x0a, 0x64, 0xfa, 0x32, 0x38, 0xa6, 0x82, 0xef, 0xbf, 0xfc, 0xa7, 0x5c, 0xe8,
	0x4f, 0x65, 0x1c, 0x26, 0x4a, 0x92, 0xed, 0x2e, 0xdd, 0xe7, 0x0a, 0x3f, 0x7e, 0x26, 0xba, 0xca,
	0x19, 0xee, 0x06, 0xe2, 0xf8, 0xc1, 0xd6, 0x63, 0x00, 0x38, 0x6a, 0x1d, 0x9c, 0x2f, 0xa6, 0xed,
	0xc6, 0x55, 0x0e, 0x88, 0x54, 0x8a, 0x8c, 0x4e, 0x19, 0xa3, 0x12, 0x39, 0x6d, 0x75, 0xee, 0x9d,
	0x9e, 0xd5, 0xbf, 0x17, 0xbd, 0xdd, 0xd4, 0x7e, 0x70, 0x8a, 0x39, 0x08, 0xea, 0x77, 0x41, 0x4f,
	0xb1, 0xc1, 0xf8, 0xbc, 0xfb, 0x35, 0x14, 0xd9, 0x1b, 0xc6, 0x86, 0xc8, 0x27, 0x4d, 0xdb, 0x79,
	0x67, 0x3f, 0x91, 0x30, 0xa7, 0x5a, 0x69, 0x48, 0xe9, 0x2d, 0xe2, 0x66, 0xe1, 0x12, 0x81, 0x33,
	0xd7, 0xea, 0x99, 0xfd, 0xbb, 0x63, 0x5f, 0xc2, 0x7c, 0xd2, 0xc0, 0xd1, 0xa1, 0xdb, 0x00, 0xf0,
	0x7d, 0x83, 0x45, 0xc3, 0xc5, 0xca, 0x33, 0x97, 0x2b, 0xcf, 0xfc, 0xbd, 0xf2, 0xcc, 0xaf, 0x6b,
	0xcf, 0x58, 0xae, 0x3d, 0xe3, 0xe7, 0xda, 0x33, 0x3e, 0xdc, 0xec, 0x9d, 0x56, 0xc6, 0x4a, 0x5d,
	0xa8, 0xec, 0x4a, 0x15, 0x7c, 0x57, 0x93, 0xd9, 0x0b, 0x32, 0xdf, 0x7b, 0x4a, 0x6d, 0xec, 0xf8,
	0xac, 0xbd, 0xf4, 0x9b, 0xbf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x0d, 0x68, 0x4b, 0x82, 0x69, 0x02,
	0x00, 0x00,
}

func (m *Params) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Params) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Params) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.MaxTotalBypassMinFeeMsgGasUsage != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.MaxTotalBypassMinFeeMsgGasUsage))
		i--
		dAtA[i] = 0x18
	}
	if len(m.BypassMinFeeMsgTypes) > 0 {
		for iNdEx := len(m.BypassMinFeeMsgTypes) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.BypassMinFeeMsgTypes[iNdEx])
			copy(dAtA[i:], m.BypassMinFeeMsgTypes[iNdEx])
			i = encodeVarintParams(dAtA, i, uint64(len(m.BypassMinFeeMsgTypes[iNdEx])))
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.MinimumGasPrices) > 0 {
		for iNdEx := len(m.MinimumGasPrices) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.MinimumGasPrices[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintParams(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintParams(dAtA []byte, offset int, v uint64) int {
	offset -= sovParams(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Params) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.MinimumGasPrices) > 0 {
		for _, e := range m.MinimumGasPrices {
			l = e.Size()
			n += 1 + l + sovParams(uint64(l))
		}
	}
	if len(m.BypassMinFeeMsgTypes) > 0 {
		for _, s := range m.BypassMinFeeMsgTypes {
			l = len(s)
			n += 1 + l + sovParams(uint64(l))
		}
	}
	if m.MaxTotalBypassMinFeeMsgGasUsage != 0 {
		n += 1 + sovParams(uint64(m.MaxTotalBypassMinFeeMsgGasUsage))
	}
	return n
}

func sovParams(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozParams(x uint64) (n int) {
	return sovParams(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Params) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowParams
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
			return fmt.Errorf("proto: Params: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Params: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinimumGasPrices", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.MinimumGasPrices = append(m.MinimumGasPrices, types.DecCoin{})
			if err := m.MinimumGasPrices[len(m.MinimumGasPrices)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BypassMinFeeMsgTypes", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.BypassMinFeeMsgTypes = append(m.BypassMinFeeMsgTypes, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxTotalBypassMinFeeMsgGasUsage", wireType)
			}
			m.MaxTotalBypassMinFeeMsgGasUsage = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MaxTotalBypassMinFeeMsgGasUsage |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipParams(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthParams
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
func skipParams(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowParams
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
					return 0, ErrIntOverflowParams
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
					return 0, ErrIntOverflowParams
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
				return 0, ErrInvalidLengthParams
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupParams
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthParams
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthParams        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowParams          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupParams = fmt.Errorf("proto: unexpected end of group")
)
