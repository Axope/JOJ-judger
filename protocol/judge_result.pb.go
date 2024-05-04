// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v5.26.1
// source: judge_result.proto

package protocol

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

type StatusSet int32

const (
	StatusSet_CE  StatusSet = 0
	StatusSet_AC  StatusSet = 1
	StatusSet_WA  StatusSet = 2
	StatusSet_TLE StatusSet = 3
	StatusSet_MLE StatusSet = 4
	StatusSet_RE  StatusSet = 5
	StatusSet_OLE StatusSet = 6
	StatusSet_UKE StatusSet = 7
)

// Enum value maps for StatusSet.
var (
	StatusSet_name = map[int32]string{
		0: "CE",
		1: "AC",
		2: "WA",
		3: "TLE",
		4: "MLE",
		5: "RE",
		6: "OLE",
		7: "UKE",
	}
	StatusSet_value = map[string]int32{
		"CE":  0,
		"AC":  1,
		"WA":  2,
		"TLE": 3,
		"MLE": 4,
		"RE":  5,
		"OLE": 6,
		"UKE": 7,
	}
)

func (x StatusSet) Enum() *StatusSet {
	p := new(StatusSet)
	*p = x
	return p
}

func (x StatusSet) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (StatusSet) Descriptor() protoreflect.EnumDescriptor {
	return file_judge_result_proto_enumTypes[0].Descriptor()
}

func (StatusSet) Type() protoreflect.EnumType {
	return &file_judge_result_proto_enumTypes[0]
}

func (x StatusSet) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use StatusSet.Descriptor instead.
func (StatusSet) EnumDescriptor() ([]byte, []int) {
	return file_judge_result_proto_rawDescGZIP(), []int{0}
}

type JudgeResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sid             string    `protobuf:"bytes,1,opt,name=sid,proto3" json:"sid,omitempty"`
	Pid             string    `protobuf:"bytes,2,opt,name=pid,proto3" json:"pid,omitempty"`
	Uid             string    `protobuf:"bytes,3,opt,name=uid,proto3" json:"uid,omitempty"`
	Cid             string    `protobuf:"bytes,4,opt,name=cid,proto3" json:"cid,omitempty"`
	Status          StatusSet `protobuf:"varint,5,opt,name=status,proto3,enum=protocol.StatusSet" json:"status,omitempty"`
	SubmitTimestamp int64     `protobuf:"varint,6,opt,name=submitTimestamp,proto3" json:"submitTimestamp,omitempty"`
}

func (x *JudgeResult) Reset() {
	*x = JudgeResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_judge_result_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JudgeResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JudgeResult) ProtoMessage() {}

func (x *JudgeResult) ProtoReflect() protoreflect.Message {
	mi := &file_judge_result_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JudgeResult.ProtoReflect.Descriptor instead.
func (*JudgeResult) Descriptor() ([]byte, []int) {
	return file_judge_result_proto_rawDescGZIP(), []int{0}
}

func (x *JudgeResult) GetSid() string {
	if x != nil {
		return x.Sid
	}
	return ""
}

func (x *JudgeResult) GetPid() string {
	if x != nil {
		return x.Pid
	}
	return ""
}

func (x *JudgeResult) GetUid() string {
	if x != nil {
		return x.Uid
	}
	return ""
}

func (x *JudgeResult) GetCid() string {
	if x != nil {
		return x.Cid
	}
	return ""
}

func (x *JudgeResult) GetStatus() StatusSet {
	if x != nil {
		return x.Status
	}
	return StatusSet_CE
}

func (x *JudgeResult) GetSubmitTimestamp() int64 {
	if x != nil {
		return x.SubmitTimestamp
	}
	return 0
}

var File_judge_result_proto protoreflect.FileDescriptor

var file_judge_result_proto_rawDesc = []byte{
	0x0a, 0x12, 0x6a, 0x75, 0x64, 0x67, 0x65, 0x5f, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x22, 0xac,
	0x01, 0x0a, 0x0b, 0x4a, 0x75, 0x64, 0x67, 0x65, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x10,
	0x0a, 0x03, 0x73, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x73, 0x69, 0x64,
	0x12, 0x10, 0x0a, 0x03, 0x70, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x70,
	0x69, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x75, 0x69, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x63, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x63, 0x69, 0x64, 0x12, 0x2b, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f,
	0x6c, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x53, 0x65, 0x74, 0x52, 0x06, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x12, 0x28, 0x0a, 0x0f, 0x73, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0f, 0x73, 0x75,
	0x62, 0x6d, 0x69, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2a, 0x4f, 0x0a,
	0x09, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x53, 0x65, 0x74, 0x12, 0x06, 0x0a, 0x02, 0x43, 0x45,
	0x10, 0x00, 0x12, 0x06, 0x0a, 0x02, 0x41, 0x43, 0x10, 0x01, 0x12, 0x06, 0x0a, 0x02, 0x57, 0x41,
	0x10, 0x02, 0x12, 0x07, 0x0a, 0x03, 0x54, 0x4c, 0x45, 0x10, 0x03, 0x12, 0x07, 0x0a, 0x03, 0x4d,
	0x4c, 0x45, 0x10, 0x04, 0x12, 0x06, 0x0a, 0x02, 0x52, 0x45, 0x10, 0x05, 0x12, 0x07, 0x0a, 0x03,
	0x4f, 0x4c, 0x45, 0x10, 0x06, 0x12, 0x07, 0x0a, 0x03, 0x55, 0x4b, 0x45, 0x10, 0x07, 0x42, 0x0c,
	0x5a, 0x0a, 0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_judge_result_proto_rawDescOnce sync.Once
	file_judge_result_proto_rawDescData = file_judge_result_proto_rawDesc
)

func file_judge_result_proto_rawDescGZIP() []byte {
	file_judge_result_proto_rawDescOnce.Do(func() {
		file_judge_result_proto_rawDescData = protoimpl.X.CompressGZIP(file_judge_result_proto_rawDescData)
	})
	return file_judge_result_proto_rawDescData
}

var file_judge_result_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_judge_result_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_judge_result_proto_goTypes = []interface{}{
	(StatusSet)(0),      // 0: protocol.StatusSet
	(*JudgeResult)(nil), // 1: protocol.JudgeResult
}
var file_judge_result_proto_depIdxs = []int32{
	0, // 0: protocol.JudgeResult.status:type_name -> protocol.StatusSet
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_judge_result_proto_init() }
func file_judge_result_proto_init() {
	if File_judge_result_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_judge_result_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JudgeResult); i {
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
			RawDescriptor: file_judge_result_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_judge_result_proto_goTypes,
		DependencyIndexes: file_judge_result_proto_depIdxs,
		EnumInfos:         file_judge_result_proto_enumTypes,
		MessageInfos:      file_judge_result_proto_msgTypes,
	}.Build()
	File_judge_result_proto = out.File
	file_judge_result_proto_rawDesc = nil
	file_judge_result_proto_goTypes = nil
	file_judge_result_proto_depIdxs = nil
}