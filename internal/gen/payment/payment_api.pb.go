// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        (unknown)
// source: payment/payment_api.proto

package paymentv1

import (
	_ "github.com/grpc-buf/internal/gen/google/api"
	_type "github.com/grpc-buf/internal/gen/google/type"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Payment card types
type CardType int32

const (
	CardType_DebitCard  CardType = 0 // Debit Card
	CardType_CreditCard CardType = 1 // Credit Card
	CardType_MasterCard CardType = 2 // Master Card
	CardType_RewardCard CardType = 3 // Reward Card
)

// Enum value maps for CardType.
var (
	CardType_name = map[int32]string{
		0: "DebitCard",
		1: "CreditCard",
		2: "MasterCard",
		3: "RewardCard",
	}
	CardType_value = map[string]int32{
		"DebitCard":  0,
		"CreditCard": 1,
		"MasterCard": 2,
		"RewardCard": 3,
	}
)

func (x CardType) Enum() *CardType {
	p := new(CardType)
	*p = x
	return p
}

func (x CardType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CardType) Descriptor() protoreflect.EnumDescriptor {
	return file_payment_payment_api_proto_enumTypes[0].Descriptor()
}

func (CardType) Type() protoreflect.EnumType {
	return &file_payment_payment_api_proto_enumTypes[0]
}

func (x CardType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CardType.Descriptor instead.
func (CardType) EnumDescriptor() ([]byte, []int) {
	return file_payment_payment_api_proto_rawDescGZIP(), []int{0}
}

// Payment invoice information
type Invoice struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Invoice id
	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	// Invoice name
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// Amount
	Amount *_type.Money `protobuf:"bytes,3,opt,name=amount,proto3" json:"amount,omitempty"`
	// Paid
	Paid bool `protobuf:"varint,4,opt,name=paid,proto3" json:"paid,omitempty"`
}

func (x *Invoice) Reset() {
	*x = Invoice{}
	if protoimpl.UnsafeEnabled {
		mi := &file_payment_payment_api_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Invoice) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Invoice) ProtoMessage() {}

func (x *Invoice) ProtoReflect() protoreflect.Message {
	mi := &file_payment_payment_api_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Invoice.ProtoReflect.Descriptor instead.
func (*Invoice) Descriptor() ([]byte, []int) {
	return file_payment_payment_api_proto_rawDescGZIP(), []int{0}
}

func (x *Invoice) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Invoice) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Invoice) GetAmount() *_type.Money {
	if x != nil {
		return x.Amount
	}
	return nil
}

func (x *Invoice) GetPaid() bool {
	if x != nil {
		return x.Paid
	}
	return false
}

// Payment information request
type PaymentRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Card No
	CardNo int64 `protobuf:"varint,1,opt,name=CardNo,proto3" json:"CardNo,omitempty"`
	// Card types
	Card CardType `protobuf:"varint,2,opt,name=card,proto3,enum=rpc.payment.v1.CardType" json:"card,omitempty"`
	// Card holder name
	Name string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	// Card holder address
	Address []string `protobuf:"bytes,4,rep,name=address,proto3" json:"address,omitempty"`
	// Total amount
	Amount float32 `protobuf:"fixed32,5,opt,name=amount,proto3" json:"amount,omitempty"`
	// Payment created time
	PaymentCreated *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=payment_created,json=paymentCreated,proto3" json:"payment_created,omitempty"`
}

func (x *PaymentRequest) Reset() {
	*x = PaymentRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_payment_payment_api_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PaymentRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PaymentRequest) ProtoMessage() {}

func (x *PaymentRequest) ProtoReflect() protoreflect.Message {
	mi := &file_payment_payment_api_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PaymentRequest.ProtoReflect.Descriptor instead.
func (*PaymentRequest) Descriptor() ([]byte, []int) {
	return file_payment_payment_api_proto_rawDescGZIP(), []int{1}
}

func (x *PaymentRequest) GetCardNo() int64 {
	if x != nil {
		return x.CardNo
	}
	return 0
}

func (x *PaymentRequest) GetCard() CardType {
	if x != nil {
		return x.Card
	}
	return CardType_DebitCard
}

func (x *PaymentRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *PaymentRequest) GetAddress() []string {
	if x != nil {
		return x.Address
	}
	return nil
}

func (x *PaymentRequest) GetAmount() float32 {
	if x != nil {
		return x.Amount
	}
	return 0
}

func (x *PaymentRequest) GetPaymentCreated() *timestamppb.Timestamp {
	if x != nil {
		return x.PaymentCreated
	}
	return nil
}

// Payment information response
type PaymentResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Response:
	//
	//	*PaymentResponse_Paid
	//	*PaymentResponse_Error
	Response isPaymentResponse_Response `protobuf_oneof:"response"`
}

func (x *PaymentResponse) Reset() {
	*x = PaymentResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_payment_payment_api_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PaymentResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PaymentResponse) ProtoMessage() {}

func (x *PaymentResponse) ProtoReflect() protoreflect.Message {
	mi := &file_payment_payment_api_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PaymentResponse.ProtoReflect.Descriptor instead.
func (*PaymentResponse) Descriptor() ([]byte, []int) {
	return file_payment_payment_api_proto_rawDescGZIP(), []int{2}
}

func (m *PaymentResponse) GetResponse() isPaymentResponse_Response {
	if m != nil {
		return m.Response
	}
	return nil
}

func (x *PaymentResponse) GetPaid() bool {
	if x, ok := x.GetResponse().(*PaymentResponse_Paid); ok {
		return x.Paid
	}
	return false
}

func (x *PaymentResponse) GetError() string {
	if x, ok := x.GetResponse().(*PaymentResponse_Error); ok {
		return x.Error
	}
	return ""
}

type isPaymentResponse_Response interface {
	isPaymentResponse_Response()
}

type PaymentResponse_Paid struct {
	// payment id
	Paid bool `protobuf:"varint,1,opt,name=paid,proto3,oneof"`
}

type PaymentResponse_Error struct {
	// payment error
	Error string `protobuf:"bytes,2,opt,name=error,proto3,oneof"`
}

func (*PaymentResponse_Paid) isPaymentResponse_Response() {}

func (*PaymentResponse_Error) isPaymentResponse_Response() {}

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_payment_payment_api_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_payment_payment_api_proto_msgTypes[3]
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
	return file_payment_payment_api_proto_rawDescGZIP(), []int{3}
}

var File_payment_payment_api_proto protoreflect.FileDescriptor

var file_payment_payment_api_proto_rawDesc = []byte{
	0x0a, 0x19, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2f, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e,
	0x74, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x72, 0x70, 0x63,
	0x2e, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x1a, 0x1f, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x2f, 0x6d, 0x6f, 0x6e, 0x65, 0x79, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x6d, 0x0a, 0x07, 0x49, 0x6e, 0x76, 0x6f, 0x69, 0x63, 0x65, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12,
	0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x2a, 0x0a, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x12, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x74, 0x79, 0x70, 0x65,
	0x2e, 0x4d, 0x6f, 0x6e, 0x65, 0x79, 0x52, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x12,
	0x0a, 0x04, 0x70, 0x61, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x04, 0x70, 0x61,
	0x69, 0x64, 0x22, 0xe1, 0x01, 0x0a, 0x0e, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x43, 0x61, 0x72, 0x64, 0x4e, 0x6f, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x43, 0x61, 0x72, 0x64, 0x4e, 0x6f, 0x12, 0x2c, 0x0a,
	0x04, 0x63, 0x61, 0x72, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x18, 0x2e, 0x72, 0x70,
	0x63, 0x2e, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x61, 0x72,
	0x64, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x63, 0x61, 0x72, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x6d, 0x6f,
	0x75, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x02, 0x52, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e,
	0x74, 0x12, 0x43, 0x0a, 0x0f, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0e, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x22, 0x4b, 0x0a, 0x0f, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x04, 0x70, 0x61, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x48, 0x00, 0x52, 0x04, 0x70, 0x61, 0x69, 0x64, 0x12,
	0x16, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00,
	0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x42, 0x0a, 0x0a, 0x08, 0x72, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x07, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x2a, 0x49, 0x0a, 0x08,
	0x43, 0x61, 0x72, 0x64, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0d, 0x0a, 0x09, 0x44, 0x65, 0x62, 0x69,
	0x74, 0x43, 0x61, 0x72, 0x64, 0x10, 0x00, 0x12, 0x0e, 0x0a, 0x0a, 0x43, 0x72, 0x65, 0x64, 0x69,
	0x74, 0x43, 0x61, 0x72, 0x64, 0x10, 0x01, 0x12, 0x0e, 0x0a, 0x0a, 0x4d, 0x61, 0x73, 0x74, 0x65,
	0x72, 0x43, 0x61, 0x72, 0x64, 0x10, 0x02, 0x12, 0x0e, 0x0a, 0x0a, 0x52, 0x65, 0x77, 0x61, 0x72,
	0x64, 0x43, 0x61, 0x72, 0x64, 0x10, 0x03, 0x32, 0xb4, 0x02, 0x0a, 0x07, 0x50, 0x61, 0x79, 0x6d,
	0x65, 0x6e, 0x74, 0x12, 0x6b, 0x0a, 0x0b, 0x4d, 0x61, 0x6b, 0x65, 0x50, 0x61, 0x79, 0x6d, 0x65,
	0x6e, 0x74, 0x12, 0x1e, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74,
	0x2e, 0x76, 0x31, 0x2e, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74,
	0x2e, 0x76, 0x31, 0x2e, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x1b, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x15, 0x3a, 0x01, 0x2a, 0x22, 0x10,
	0x2f, 0x76, 0x31, 0x2f, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2f, 0x6d, 0x61, 0x6b, 0x65,
	0x12, 0x60, 0x0a, 0x0f, 0x4d, 0x61, 0x72, 0x6b, 0x49, 0x6e, 0x76, 0x6f, 0x69, 0x63, 0x65, 0x50,
	0x61, 0x69, 0x64, 0x12, 0x17, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e,
	0x74, 0x2e, 0x76, 0x31, 0x2e, 0x49, 0x6e, 0x76, 0x6f, 0x69, 0x63, 0x65, 0x1a, 0x17, 0x2e, 0x72,
	0x70, 0x63, 0x2e, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x49, 0x6e,
	0x76, 0x6f, 0x69, 0x63, 0x65, 0x22, 0x1b, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x15, 0x3a, 0x01, 0x2a,
	0x22, 0x10, 0x2f, 0x76, 0x31, 0x2f, 0x69, 0x6e, 0x76, 0x6f, 0x69, 0x63, 0x65, 0x2f, 0x6d, 0x61,
	0x72, 0x6b, 0x12, 0x5a, 0x0a, 0x0a, 0x50, 0x61, 0x79, 0x49, 0x6e, 0x76, 0x6f, 0x69, 0x63, 0x65,
	0x12, 0x17, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76,
	0x31, 0x2e, 0x49, 0x6e, 0x76, 0x6f, 0x69, 0x63, 0x65, 0x1a, 0x17, 0x2e, 0x72, 0x70, 0x63, 0x2e,
	0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x49, 0x6e, 0x76, 0x6f, 0x69,
	0x63, 0x65, 0x22, 0x1a, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x14, 0x3a, 0x01, 0x2a, 0x22, 0x0f, 0x2f,
	0x76, 0x31, 0x2f, 0x69, 0x6e, 0x76, 0x6f, 0x69, 0x63, 0x65, 0x2f, 0x70, 0x61, 0x79, 0x42, 0xb3,
	0x01, 0x0a, 0x12, 0x63, 0x6f, 0x6d, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x70, 0x61, 0x79, 0x6d, 0x65,
	0x6e, 0x74, 0x2e, 0x76, 0x31, 0x42, 0x0f, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x41, 0x70,
	0x69, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x32, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2d, 0x62, 0x75, 0x66, 0x2f, 0x69, 0x6e,
	0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x61, 0x79, 0x6d, 0x65,
	0x6e, 0x74, 0x3b, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x52,
	0x50, 0x58, 0xaa, 0x02, 0x0e, 0x52, 0x70, 0x63, 0x2e, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74,
	0x2e, 0x56, 0x31, 0xca, 0x02, 0x0e, 0x52, 0x70, 0x63, 0x5c, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e,
	0x74, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x1a, 0x52, 0x70, 0x63, 0x5c, 0x50, 0x61, 0x79, 0x6d, 0x65,
	0x6e, 0x74, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0xea, 0x02, 0x10, 0x52, 0x70, 0x63, 0x3a, 0x3a, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74,
	0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_payment_payment_api_proto_rawDescOnce sync.Once
	file_payment_payment_api_proto_rawDescData = file_payment_payment_api_proto_rawDesc
)

func file_payment_payment_api_proto_rawDescGZIP() []byte {
	file_payment_payment_api_proto_rawDescOnce.Do(func() {
		file_payment_payment_api_proto_rawDescData = protoimpl.X.CompressGZIP(file_payment_payment_api_proto_rawDescData)
	})
	return file_payment_payment_api_proto_rawDescData
}

var file_payment_payment_api_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_payment_payment_api_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_payment_payment_api_proto_goTypes = []interface{}{
	(CardType)(0),                 // 0: rpc.payment.v1.CardType
	(*Invoice)(nil),               // 1: rpc.payment.v1.Invoice
	(*PaymentRequest)(nil),        // 2: rpc.payment.v1.PaymentRequest
	(*PaymentResponse)(nil),       // 3: rpc.payment.v1.PaymentResponse
	(*Empty)(nil),                 // 4: rpc.payment.v1.Empty
	(*_type.Money)(nil),           // 5: google.type.Money
	(*timestamppb.Timestamp)(nil), // 6: google.protobuf.Timestamp
}
var file_payment_payment_api_proto_depIdxs = []int32{
	5, // 0: rpc.payment.v1.Invoice.amount:type_name -> google.type.Money
	0, // 1: rpc.payment.v1.PaymentRequest.card:type_name -> rpc.payment.v1.CardType
	6, // 2: rpc.payment.v1.PaymentRequest.payment_created:type_name -> google.protobuf.Timestamp
	2, // 3: rpc.payment.v1.Payment.MakePayment:input_type -> rpc.payment.v1.PaymentRequest
	1, // 4: rpc.payment.v1.Payment.MarkInvoicePaid:input_type -> rpc.payment.v1.Invoice
	1, // 5: rpc.payment.v1.Payment.PayInvoice:input_type -> rpc.payment.v1.Invoice
	3, // 6: rpc.payment.v1.Payment.MakePayment:output_type -> rpc.payment.v1.PaymentResponse
	1, // 7: rpc.payment.v1.Payment.MarkInvoicePaid:output_type -> rpc.payment.v1.Invoice
	1, // 8: rpc.payment.v1.Payment.PayInvoice:output_type -> rpc.payment.v1.Invoice
	6, // [6:9] is the sub-list for method output_type
	3, // [3:6] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_payment_payment_api_proto_init() }
func file_payment_payment_api_proto_init() {
	if File_payment_payment_api_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_payment_payment_api_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Invoice); i {
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
		file_payment_payment_api_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PaymentRequest); i {
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
		file_payment_payment_api_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PaymentResponse); i {
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
		file_payment_payment_api_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
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
	}
	file_payment_payment_api_proto_msgTypes[2].OneofWrappers = []interface{}{
		(*PaymentResponse_Paid)(nil),
		(*PaymentResponse_Error)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_payment_payment_api_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_payment_payment_api_proto_goTypes,
		DependencyIndexes: file_payment_payment_api_proto_depIdxs,
		EnumInfos:         file_payment_payment_api_proto_enumTypes,
		MessageInfos:      file_payment_payment_api_proto_msgTypes,
	}.Build()
	File_payment_payment_api_proto = out.File
	file_payment_payment_api_proto_rawDesc = nil
	file_payment_payment_api_proto_goTypes = nil
	file_payment_payment_api_proto_depIdxs = nil
}
