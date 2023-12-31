// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.23.4
// source: otp.proto

package v1

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

type OTPGenRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *OTPGenRequest) Reset() {
	*x = OTPGenRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_otp_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OTPGenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OTPGenRequest) ProtoMessage() {}

func (x *OTPGenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_otp_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OTPGenRequest.ProtoReflect.Descriptor instead.
func (*OTPGenRequest) Descriptor() ([]byte, []int) {
	return file_otp_proto_rawDescGZIP(), []int{0}
}

type OTPGenResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SecretKey string `protobuf:"bytes,1,opt,name=secret_key,json=secretKey,proto3" json:"secret_key,omitempty"`
	AuthUrl   string `protobuf:"bytes,2,opt,name=auth_url,json=authUrl,proto3" json:"auth_url,omitempty"`
}

func (x *OTPGenResponse) Reset() {
	*x = OTPGenResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_otp_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OTPGenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OTPGenResponse) ProtoMessage() {}

func (x *OTPGenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_otp_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OTPGenResponse.ProtoReflect.Descriptor instead.
func (*OTPGenResponse) Descriptor() ([]byte, []int) {
	return file_otp_proto_rawDescGZIP(), []int{1}
}

func (x *OTPGenResponse) GetSecretKey() string {
	if x != nil {
		return x.SecretKey
	}
	return ""
}

func (x *OTPGenResponse) GetAuthUrl() string {
	if x != nil {
		return x.AuthUrl
	}
	return ""
}

type OTPVerifyRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *OTPVerifyRequest) Reset() {
	*x = OTPVerifyRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_otp_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OTPVerifyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OTPVerifyRequest) ProtoMessage() {}

func (x *OTPVerifyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_otp_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OTPVerifyRequest.ProtoReflect.Descriptor instead.
func (*OTPVerifyRequest) Descriptor() ([]byte, []int) {
	return file_otp_proto_rawDescGZIP(), []int{2}
}

func (x *OTPVerifyRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type OTPVerifyResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AccessToken string `protobuf:"bytes,1,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
}

func (x *OTPVerifyResponse) Reset() {
	*x = OTPVerifyResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_otp_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OTPVerifyResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OTPVerifyResponse) ProtoMessage() {}

func (x *OTPVerifyResponse) ProtoReflect() protoreflect.Message {
	mi := &file_otp_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OTPVerifyResponse.ProtoReflect.Descriptor instead.
func (*OTPVerifyResponse) Descriptor() ([]byte, []int) {
	return file_otp_proto_rawDescGZIP(), []int{3}
}

func (x *OTPVerifyResponse) GetAccessToken() string {
	if x != nil {
		return x.AccessToken
	}
	return ""
}

type OTPValidateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *OTPValidateRequest) Reset() {
	*x = OTPValidateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_otp_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OTPValidateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OTPValidateRequest) ProtoMessage() {}

func (x *OTPValidateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_otp_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OTPValidateRequest.ProtoReflect.Descriptor instead.
func (*OTPValidateRequest) Descriptor() ([]byte, []int) {
	return file_otp_proto_rawDescGZIP(), []int{4}
}

func (x *OTPValidateRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type OTPValidateResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AccessToken string `protobuf:"bytes,1,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
}

func (x *OTPValidateResponse) Reset() {
	*x = OTPValidateResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_otp_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OTPValidateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OTPValidateResponse) ProtoMessage() {}

func (x *OTPValidateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_otp_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OTPValidateResponse.ProtoReflect.Descriptor instead.
func (*OTPValidateResponse) Descriptor() ([]byte, []int) {
	return file_otp_proto_rawDescGZIP(), []int{5}
}

func (x *OTPValidateResponse) GetAccessToken() string {
	if x != nil {
		return x.AccessToken
	}
	return ""
}

type OTPDisableRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PasswordHash string `protobuf:"bytes,1,opt,name=password_hash,json=passwordHash,proto3" json:"password_hash,omitempty"`
}

func (x *OTPDisableRequest) Reset() {
	*x = OTPDisableRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_otp_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OTPDisableRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OTPDisableRequest) ProtoMessage() {}

func (x *OTPDisableRequest) ProtoReflect() protoreflect.Message {
	mi := &file_otp_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OTPDisableRequest.ProtoReflect.Descriptor instead.
func (*OTPDisableRequest) Descriptor() ([]byte, []int) {
	return file_otp_proto_rawDescGZIP(), []int{6}
}

func (x *OTPDisableRequest) GetPasswordHash() string {
	if x != nil {
		return x.PasswordHash
	}
	return ""
}

type OTPDisableResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *OTPDisableResponse) Reset() {
	*x = OTPDisableResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_otp_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OTPDisableResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OTPDisableResponse) ProtoMessage() {}

func (x *OTPDisableResponse) ProtoReflect() protoreflect.Message {
	mi := &file_otp_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OTPDisableResponse.ProtoReflect.Descriptor instead.
func (*OTPDisableResponse) Descriptor() ([]byte, []int) {
	return file_otp_proto_rawDescGZIP(), []int{7}
}

var File_otp_proto protoreflect.FileDescriptor

var file_otp_proto_rawDesc = []byte{
	0x0a, 0x09, 0x6f, 0x74, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x70, 0x61, 0x73,
	0x73, 0x5f, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x22, 0x0f, 0x0a, 0x0d, 0x4f, 0x54, 0x50, 0x47,
	0x65, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x4a, 0x0a, 0x0e, 0x4f, 0x54, 0x50,
	0x47, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x73,
	0x65, 0x63, 0x72, 0x65, 0x74, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x4b, 0x65, 0x79, 0x12, 0x19, 0x0a, 0x08, 0x61, 0x75,
	0x74, 0x68, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x75,
	0x74, 0x68, 0x55, 0x72, 0x6c, 0x22, 0x28, 0x0a, 0x10, 0x4f, 0x54, 0x50, 0x56, 0x65, 0x72, 0x69,
	0x66, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b,
	0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22,
	0x36, 0x0a, 0x11, 0x4f, 0x54, 0x50, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x74,
	0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x61, 0x63, 0x63, 0x65,
	0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x2a, 0x0a, 0x12, 0x4f, 0x54, 0x50, 0x56, 0x61,
	0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a,
	0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f,
	0x6b, 0x65, 0x6e, 0x22, 0x38, 0x0a, 0x13, 0x4f, 0x54, 0x50, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61,
	0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x61, 0x63,
	0x63, 0x65, 0x73, 0x73, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0b, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x38, 0x0a,
	0x11, 0x4f, 0x54, 0x50, 0x44, 0x69, 0x73, 0x61, 0x62, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x23, 0x0a, 0x0d, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x5f, 0x68,
	0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x70, 0x61, 0x73, 0x73, 0x77,
	0x6f, 0x72, 0x64, 0x48, 0x61, 0x73, 0x68, 0x22, 0x14, 0x0a, 0x12, 0x4f, 0x54, 0x50, 0x44, 0x69,
	0x73, 0x61, 0x62, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0xc1, 0x02,
	0x0a, 0x0a, 0x4f, 0x74, 0x70, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x46, 0x0a, 0x0b,
	0x4f, 0x54, 0x50, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x12, 0x1a, 0x2e, 0x70, 0x61,
	0x73, 0x73, 0x5f, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x4f, 0x54, 0x50, 0x47, 0x65, 0x6e,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x70, 0x61, 0x73, 0x73, 0x5f, 0x6b,
	0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x4f, 0x54, 0x50, 0x47, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4a, 0x0a, 0x09, 0x4f, 0x54, 0x50, 0x56, 0x65, 0x72, 0x69, 0x66,
	0x79, 0x12, 0x1d, 0x2e, 0x70, 0x61, 0x73, 0x73, 0x5f, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e,
	0x4f, 0x54, 0x50, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x1e, 0x2e, 0x70, 0x61, 0x73, 0x73, 0x5f, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x4f,
	0x54, 0x50, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x50, 0x0a, 0x0b, 0x4f, 0x54, 0x50, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x12,
	0x1f, 0x2e, 0x70, 0x61, 0x73, 0x73, 0x5f, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x4f, 0x54,
	0x50, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x20, 0x2e, 0x70, 0x61, 0x73, 0x73, 0x5f, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x4f,
	0x54, 0x50, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x4d, 0x0a, 0x0a, 0x4f, 0x54, 0x50, 0x44, 0x69, 0x73, 0x61, 0x62, 0x6c, 0x65,
	0x12, 0x1e, 0x2e, 0x70, 0x61, 0x73, 0x73, 0x5f, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x4f,
	0x54, 0x50, 0x44, 0x69, 0x73, 0x61, 0x62, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x1f, 0x2e, 0x70, 0x61, 0x73, 0x73, 0x5f, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x4f,
	0x54, 0x50, 0x44, 0x69, 0x73, 0x61, 0x62, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x42, 0x33, 0x5a, 0x31, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x75, 0x6e, 0x62, 0x65, 0x6d, 0x61, 0x6e, 0x2f, 0x79, 0x61, 0x2d, 0x70, 0x72, 0x61, 0x63, 0x2d,
	0x67, 0x6f, 0x2d, 0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x2d, 0x67, 0x72, 0x61, 0x64, 0x65, 0x2f,
	0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_otp_proto_rawDescOnce sync.Once
	file_otp_proto_rawDescData = file_otp_proto_rawDesc
)

func file_otp_proto_rawDescGZIP() []byte {
	file_otp_proto_rawDescOnce.Do(func() {
		file_otp_proto_rawDescData = protoimpl.X.CompressGZIP(file_otp_proto_rawDescData)
	})
	return file_otp_proto_rawDescData
}

var file_otp_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_otp_proto_goTypes = []interface{}{
	(*OTPGenRequest)(nil),       // 0: pass_keeper.OTPGenRequest
	(*OTPGenResponse)(nil),      // 1: pass_keeper.OTPGenResponse
	(*OTPVerifyRequest)(nil),    // 2: pass_keeper.OTPVerifyRequest
	(*OTPVerifyResponse)(nil),   // 3: pass_keeper.OTPVerifyResponse
	(*OTPValidateRequest)(nil),  // 4: pass_keeper.OTPValidateRequest
	(*OTPValidateResponse)(nil), // 5: pass_keeper.OTPValidateResponse
	(*OTPDisableRequest)(nil),   // 6: pass_keeper.OTPDisableRequest
	(*OTPDisableResponse)(nil),  // 7: pass_keeper.OTPDisableResponse
}
var file_otp_proto_depIdxs = []int32{
	0, // 0: pass_keeper.OtpService.OTPGenerate:input_type -> pass_keeper.OTPGenRequest
	2, // 1: pass_keeper.OtpService.OTPVerify:input_type -> pass_keeper.OTPVerifyRequest
	4, // 2: pass_keeper.OtpService.OTPValidate:input_type -> pass_keeper.OTPValidateRequest
	6, // 3: pass_keeper.OtpService.OTPDisable:input_type -> pass_keeper.OTPDisableRequest
	1, // 4: pass_keeper.OtpService.OTPGenerate:output_type -> pass_keeper.OTPGenResponse
	3, // 5: pass_keeper.OtpService.OTPVerify:output_type -> pass_keeper.OTPVerifyResponse
	5, // 6: pass_keeper.OtpService.OTPValidate:output_type -> pass_keeper.OTPValidateResponse
	7, // 7: pass_keeper.OtpService.OTPDisable:output_type -> pass_keeper.OTPDisableResponse
	4, // [4:8] is the sub-list for method output_type
	0, // [0:4] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_otp_proto_init() }
func file_otp_proto_init() {
	if File_otp_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_otp_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OTPGenRequest); i {
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
		file_otp_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OTPGenResponse); i {
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
		file_otp_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OTPVerifyRequest); i {
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
		file_otp_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OTPVerifyResponse); i {
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
		file_otp_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OTPValidateRequest); i {
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
		file_otp_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OTPValidateResponse); i {
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
		file_otp_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OTPDisableRequest); i {
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
		file_otp_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OTPDisableResponse); i {
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
			RawDescriptor: file_otp_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_otp_proto_goTypes,
		DependencyIndexes: file_otp_proto_depIdxs,
		MessageInfos:      file_otp_proto_msgTypes,
	}.Build()
	File_otp_proto = out.File
	file_otp_proto_rawDesc = nil
	file_otp_proto_goTypes = nil
	file_otp_proto_depIdxs = nil
}
