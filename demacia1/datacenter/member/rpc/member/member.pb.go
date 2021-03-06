// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.19.4
// source: member.proto

package member

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type IdReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *IdReq) Reset() {
	*x = IdReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_member_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IdReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IdReq) ProtoMessage() {}

func (x *IdReq) ProtoReflect() protoreflect.Message {
	mi := &file_member_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IdReq.ProtoReflect.Descriptor instead.
func (*IdReq) Descriptor() ([]byte, []int) {
	return file_member_proto_rawDescGZIP(), []int{0}
}

func (x *IdReq) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type UserNameReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserName string `protobuf:"bytes,1,opt,name=userName,proto3" json:"userName,omitempty"`
}

func (x *UserNameReq) Reset() {
	*x = UserNameReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_member_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserNameReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserNameReq) ProtoMessage() {}

func (x *UserNameReq) ProtoReflect() protoreflect.Message {
	mi := &file_member_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserNameReq.ProtoReflect.Descriptor instead.
func (*UserNameReq) Descriptor() ([]byte, []int) {
	return file_member_proto_rawDescGZIP(), []int{1}
}

func (x *UserNameReq) GetUserName() string {
	if x != nil {
		return x.UserName
	}
	return ""
}

type MemberInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	OrgId    int64  `protobuf:"varint,2,opt,name=orgId,proto3" json:"orgId,omitempty"`
	UserName string `protobuf:"bytes,3,opt,name=userName,proto3" json:"userName,omitempty"`
	TrueName string `protobuf:"bytes,4,opt,name=trueName,proto3" json:"trueName,omitempty"`
	Mobile   string `protobuf:"bytes,5,opt,name=mobile,proto3" json:"mobile,omitempty"`
	Password string `protobuf:"bytes,6,opt,name=password,proto3" json:"password,omitempty"`
	Avatar   string `protobuf:"bytes,7,opt,name=avatar,proto3" json:"avatar,omitempty"`
	Role     int64  `protobuf:"varint,8,opt,name=role,proto3" json:"role,omitempty"`
}

func (x *MemberInfo) Reset() {
	*x = MemberInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_member_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MemberInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MemberInfo) ProtoMessage() {}

func (x *MemberInfo) ProtoReflect() protoreflect.Message {
	mi := &file_member_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MemberInfo.ProtoReflect.Descriptor instead.
func (*MemberInfo) Descriptor() ([]byte, []int) {
	return file_member_proto_rawDescGZIP(), []int{2}
}

func (x *MemberInfo) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *MemberInfo) GetOrgId() int64 {
	if x != nil {
		return x.OrgId
	}
	return 0
}

func (x *MemberInfo) GetUserName() string {
	if x != nil {
		return x.UserName
	}
	return ""
}

func (x *MemberInfo) GetTrueName() string {
	if x != nil {
		return x.TrueName
	}
	return ""
}

func (x *MemberInfo) GetMobile() string {
	if x != nil {
		return x.Mobile
	}
	return ""
}

func (x *MemberInfo) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *MemberInfo) GetAvatar() string {
	if x != nil {
		return x.Avatar
	}
	return ""
}

func (x *MemberInfo) GetRole() int64 {
	if x != nil {
		return x.Role
	}
	return 0
}

type InsertReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserName string `protobuf:"bytes,1,opt,name=userName,proto3" json:"userName,omitempty"`
	OrgId    int64  `protobuf:"varint,2,opt,name=orgId,proto3" json:"orgId,omitempty"`
	TrueName string `protobuf:"bytes,3,opt,name=trueName,proto3" json:"trueName,omitempty"`
	Mobile   string `protobuf:"bytes,4,opt,name=mobile,proto3" json:"mobile,omitempty"`
	Password string `protobuf:"bytes,5,opt,name=password,proto3" json:"password,omitempty"`
}

func (x *InsertReq) Reset() {
	*x = InsertReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_member_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InsertReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InsertReq) ProtoMessage() {}

func (x *InsertReq) ProtoReflect() protoreflect.Message {
	mi := &file_member_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InsertReq.ProtoReflect.Descriptor instead.
func (*InsertReq) Descriptor() ([]byte, []int) {
	return file_member_proto_rawDescGZIP(), []int{3}
}

func (x *InsertReq) GetUserName() string {
	if x != nil {
		return x.UserName
	}
	return ""
}

func (x *InsertReq) GetOrgId() int64 {
	if x != nil {
		return x.OrgId
	}
	return 0
}

func (x *InsertReq) GetTrueName() string {
	if x != nil {
		return x.TrueName
	}
	return ""
}

func (x *InsertReq) GetMobile() string {
	if x != nil {
		return x.Mobile
	}
	return ""
}

func (x *InsertReq) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type NullResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *NullResp) Reset() {
	*x = NullResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_member_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NullResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NullResp) ProtoMessage() {}

func (x *NullResp) ProtoReflect() protoreflect.Message {
	mi := &file_member_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NullResp.ProtoReflect.Descriptor instead.
func (*NullResp) Descriptor() ([]byte, []int) {
	return file_member_proto_rawDescGZIP(), []int{4}
}

var File_member_proto protoreflect.FileDescriptor

var file_member_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06,
	0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x22, 0x17, 0x0a, 0x05, 0x49, 0x64, 0x52, 0x65, 0x71, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x22,
	0x29, 0x0a, 0x0b, 0x55, 0x73, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x71, 0x12, 0x1a,
	0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0xca, 0x01, 0x0a, 0x0a, 0x4d,
	0x65, 0x6d, 0x62, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x6f, 0x72, 0x67,
	0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x6f, 0x72, 0x67, 0x49, 0x64, 0x12,
	0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x74,
	0x72, 0x75, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x74,
	0x72, 0x75, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x6d, 0x6f, 0x62, 0x69, 0x6c,
	0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6d, 0x6f, 0x62, 0x69, 0x6c, 0x65, 0x12,
	0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x61,
	0x76, 0x61, 0x74, 0x61, 0x72, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61, 0x76, 0x61,
	0x74, 0x61, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x22, 0x8d, 0x01, 0x0a, 0x09, 0x49, 0x6e, 0x73, 0x65,
	0x72, 0x74, 0x52, 0x65, 0x71, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x4e, 0x61, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x4e, 0x61, 0x6d,
	0x65, 0x12, 0x14, 0x0a, 0x05, 0x6f, 0x72, 0x67, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x05, 0x6f, 0x72, 0x67, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x74, 0x72, 0x75, 0x65, 0x4e,
	0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x74, 0x72, 0x75, 0x65, 0x4e,
	0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x6d, 0x6f, 0x62, 0x69, 0x6c, 0x65, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x6d, 0x6f, 0x62, 0x69, 0x6c, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x70,
	0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70,
	0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x22, 0x0a, 0x0a, 0x08, 0x4e, 0x75, 0x6c, 0x6c, 0x52,
	0x65, 0x73, 0x70, 0x32, 0xd3, 0x01, 0x0a, 0x06, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x30,
	0x0a, 0x0b, 0x66, 0x69, 0x6e, 0x64, 0x4f, 0x6e, 0x65, 0x42, 0x79, 0x49, 0x64, 0x12, 0x0d, 0x2e,
	0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x2e, 0x49, 0x64, 0x52, 0x65, 0x71, 0x1a, 0x12, 0x2e, 0x6d,
	0x65, 0x6d, 0x62, 0x65, 0x72, 0x2e, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f,
	0x12, 0x3c, 0x0a, 0x11, 0x66, 0x69, 0x6e, 0x64, 0x4f, 0x6e, 0x65, 0x42, 0x79, 0x55, 0x73, 0x65,
	0x72, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x13, 0x2e, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x2e, 0x55,
	0x73, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x71, 0x1a, 0x12, 0x2e, 0x6d, 0x65, 0x6d,
	0x62, 0x65, 0x72, 0x2e, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x2a,
	0x0a, 0x06, 0x69, 0x6e, 0x73, 0x65, 0x72, 0x74, 0x12, 0x11, 0x2e, 0x6d, 0x65, 0x6d, 0x62, 0x65,
	0x72, 0x2e, 0x49, 0x6e, 0x73, 0x65, 0x72, 0x74, 0x52, 0x65, 0x71, 0x1a, 0x0d, 0x2e, 0x6d, 0x65,
	0x6d, 0x62, 0x65, 0x72, 0x2e, 0x49, 0x64, 0x52, 0x65, 0x71, 0x12, 0x2d, 0x0a, 0x0a, 0x64, 0x65,
	0x6c, 0x65, 0x74, 0x65, 0x42, 0x79, 0x49, 0x64, 0x12, 0x0d, 0x2e, 0x6d, 0x65, 0x6d, 0x62, 0x65,
	0x72, 0x2e, 0x49, 0x64, 0x52, 0x65, 0x71, 0x1a, 0x10, 0x2e, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72,
	0x2e, 0x4e, 0x75, 0x6c, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x42, 0x08, 0x5a, 0x06, 0x6d, 0x65, 0x6d,
	0x62, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_member_proto_rawDescOnce sync.Once
	file_member_proto_rawDescData = file_member_proto_rawDesc
)

func file_member_proto_rawDescGZIP() []byte {
	file_member_proto_rawDescOnce.Do(func() {
		file_member_proto_rawDescData = protoimpl.X.CompressGZIP(file_member_proto_rawDescData)
	})
	return file_member_proto_rawDescData
}

var file_member_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_member_proto_goTypes = []interface{}{
	(*IdReq)(nil),       // 0: member.IdReq
	(*UserNameReq)(nil), // 1: member.UserNameReq
	(*MemberInfo)(nil),  // 2: member.MemberInfo
	(*InsertReq)(nil),   // 3: member.InsertReq
	(*NullResp)(nil),    // 4: member.NullResp
}
var file_member_proto_depIdxs = []int32{
	0, // 0: member.member.findOneById:input_type -> member.IdReq
	1, // 1: member.member.findOneByUserName:input_type -> member.UserNameReq
	3, // 2: member.member.insert:input_type -> member.InsertReq
	0, // 3: member.member.deleteById:input_type -> member.IdReq
	2, // 4: member.member.findOneById:output_type -> member.MemberInfo
	2, // 5: member.member.findOneByUserName:output_type -> member.MemberInfo
	0, // 6: member.member.insert:output_type -> member.IdReq
	4, // 7: member.member.deleteById:output_type -> member.NullResp
	4, // [4:8] is the sub-list for method output_type
	0, // [0:4] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_member_proto_init() }
func file_member_proto_init() {
	if File_member_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_member_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IdReq); i {
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
		file_member_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserNameReq); i {
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
		file_member_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MemberInfo); i {
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
		file_member_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InsertReq); i {
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
		file_member_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NullResp); i {
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
			RawDescriptor: file_member_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_member_proto_goTypes,
		DependencyIndexes: file_member_proto_depIdxs,
		MessageInfos:      file_member_proto_msgTypes,
	}.Build()
	File_member_proto = out.File
	file_member_proto_rawDesc = nil
	file_member_proto_goTypes = nil
	file_member_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// MemberClient is the client API for Member service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MemberClient interface {
	FindOneById(ctx context.Context, in *IdReq, opts ...grpc.CallOption) (*MemberInfo, error)
	FindOneByUserName(ctx context.Context, in *UserNameReq, opts ...grpc.CallOption) (*MemberInfo, error)
	Insert(ctx context.Context, in *InsertReq, opts ...grpc.CallOption) (*IdReq, error)
	DeleteById(ctx context.Context, in *IdReq, opts ...grpc.CallOption) (*NullResp, error)
}

type memberClient struct {
	cc grpc.ClientConnInterface
}

func NewMemberClient(cc grpc.ClientConnInterface) MemberClient {
	return &memberClient{cc}
}

func (c *memberClient) FindOneById(ctx context.Context, in *IdReq, opts ...grpc.CallOption) (*MemberInfo, error) {
	out := new(MemberInfo)
	err := c.cc.Invoke(ctx, "/member.member/findOneById", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *memberClient) FindOneByUserName(ctx context.Context, in *UserNameReq, opts ...grpc.CallOption) (*MemberInfo, error) {
	out := new(MemberInfo)
	err := c.cc.Invoke(ctx, "/member.member/findOneByUserName", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *memberClient) Insert(ctx context.Context, in *InsertReq, opts ...grpc.CallOption) (*IdReq, error) {
	out := new(IdReq)
	err := c.cc.Invoke(ctx, "/member.member/insert", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *memberClient) DeleteById(ctx context.Context, in *IdReq, opts ...grpc.CallOption) (*NullResp, error) {
	out := new(NullResp)
	err := c.cc.Invoke(ctx, "/member.member/deleteById", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MemberServer is the server API for Member service.
type MemberServer interface {
	FindOneById(context.Context, *IdReq) (*MemberInfo, error)
	FindOneByUserName(context.Context, *UserNameReq) (*MemberInfo, error)
	Insert(context.Context, *InsertReq) (*IdReq, error)
	DeleteById(context.Context, *IdReq) (*NullResp, error)
}

// UnimplementedMemberServer can be embedded to have forward compatible implementations.
type UnimplementedMemberServer struct {
}

func (*UnimplementedMemberServer) FindOneById(context.Context, *IdReq) (*MemberInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindOneById not implemented")
}
func (*UnimplementedMemberServer) FindOneByUserName(context.Context, *UserNameReq) (*MemberInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindOneByUserName not implemented")
}
func (*UnimplementedMemberServer) Insert(context.Context, *InsertReq) (*IdReq, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Insert not implemented")
}
func (*UnimplementedMemberServer) DeleteById(context.Context, *IdReq) (*NullResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteById not implemented")
}

func RegisterMemberServer(s *grpc.Server, srv MemberServer) {
	s.RegisterService(&_Member_serviceDesc, srv)
}

func _Member_FindOneById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IdReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MemberServer).FindOneById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/member.member/FindOneById",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MemberServer).FindOneById(ctx, req.(*IdReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Member_FindOneByUserName_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserNameReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MemberServer).FindOneByUserName(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/member.member/FindOneByUserName",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MemberServer).FindOneByUserName(ctx, req.(*UserNameReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Member_Insert_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InsertReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MemberServer).Insert(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/member.member/Insert",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MemberServer).Insert(ctx, req.(*InsertReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Member_DeleteById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IdReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MemberServer).DeleteById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/member.member/DeleteById",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MemberServer).DeleteById(ctx, req.(*IdReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _Member_serviceDesc = grpc.ServiceDesc{
	ServiceName: "member.member",
	HandlerType: (*MemberServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "findOneById",
			Handler:    _Member_FindOneById_Handler,
		},
		{
			MethodName: "findOneByUserName",
			Handler:    _Member_FindOneByUserName_Handler,
		},
		{
			MethodName: "insert",
			Handler:    _Member_Insert_Handler,
		},
		{
			MethodName: "deleteById",
			Handler:    _Member_DeleteById_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "member.proto",
}
