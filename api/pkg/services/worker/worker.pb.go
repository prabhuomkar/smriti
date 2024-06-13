// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v5.26.1
// source: worker.proto

package worker

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

type MediaItemComponent int32

const (
	MediaItemComponent_METADATA          MediaItemComponent = 0
	MediaItemComponent_PREVIEW_THUMBNAIL MediaItemComponent = 1
	MediaItemComponent_PLACES            MediaItemComponent = 2
	MediaItemComponent_CLASSIFICATION    MediaItemComponent = 3
	MediaItemComponent_FACES             MediaItemComponent = 4
	MediaItemComponent_OCR               MediaItemComponent = 5
	MediaItemComponent_SEARCH            MediaItemComponent = 6
)

// Enum value maps for MediaItemComponent.
var (
	MediaItemComponent_name = map[int32]string{
		0: "METADATA",
		1: "PREVIEW_THUMBNAIL",
		2: "PLACES",
		3: "CLASSIFICATION",
		4: "FACES",
		5: "OCR",
		6: "SEARCH",
	}
	MediaItemComponent_value = map[string]int32{
		"METADATA":          0,
		"PREVIEW_THUMBNAIL": 1,
		"PLACES":            2,
		"CLASSIFICATION":    3,
		"FACES":             4,
		"OCR":               5,
		"SEARCH":            6,
	}
)

func (x MediaItemComponent) Enum() *MediaItemComponent {
	p := new(MediaItemComponent)
	*p = x
	return p
}

func (x MediaItemComponent) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MediaItemComponent) Descriptor() protoreflect.EnumDescriptor {
	return file_worker_proto_enumTypes[0].Descriptor()
}

func (MediaItemComponent) Type() protoreflect.EnumType {
	return &file_worker_proto_enumTypes[0]
}

func (x MediaItemComponent) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MediaItemComponent.Descriptor instead.
func (MediaItemComponent) EnumDescriptor() ([]byte, []int) {
	return file_worker_proto_rawDescGZIP(), []int{0}
}

type MediaItemProcessRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId     string               `protobuf:"bytes,1,opt,name=userId,proto3" json:"userId,omitempty"`
	Id         string               `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	FilePath   string               `protobuf:"bytes,3,opt,name=filePath,proto3" json:"filePath,omitempty"`
	Components []MediaItemComponent `protobuf:"varint,4,rep,packed,name=components,proto3,enum=MediaItemComponent" json:"components,omitempty"`
	Payload    map[string]string    `protobuf:"bytes,5,rep,name=payload,proto3" json:"payload,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *MediaItemProcessRequest) Reset() {
	*x = MediaItemProcessRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_worker_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MediaItemProcessRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MediaItemProcessRequest) ProtoMessage() {}

func (x *MediaItemProcessRequest) ProtoReflect() protoreflect.Message {
	mi := &file_worker_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MediaItemProcessRequest.ProtoReflect.Descriptor instead.
func (*MediaItemProcessRequest) Descriptor() ([]byte, []int) {
	return file_worker_proto_rawDescGZIP(), []int{0}
}

func (x *MediaItemProcessRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *MediaItemProcessRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *MediaItemProcessRequest) GetFilePath() string {
	if x != nil {
		return x.FilePath
	}
	return ""
}

func (x *MediaItemProcessRequest) GetComponents() []MediaItemComponent {
	if x != nil {
		return x.Components
	}
	return nil
}

func (x *MediaItemProcessRequest) GetPayload() map[string]string {
	if x != nil {
		return x.Payload
	}
	return nil
}

type MediaItemProcessResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ok bool `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`
}

func (x *MediaItemProcessResponse) Reset() {
	*x = MediaItemProcessResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_worker_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MediaItemProcessResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MediaItemProcessResponse) ProtoMessage() {}

func (x *MediaItemProcessResponse) ProtoReflect() protoreflect.Message {
	mi := &file_worker_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MediaItemProcessResponse.ProtoReflect.Descriptor instead.
func (*MediaItemProcessResponse) Descriptor() ([]byte, []int) {
	return file_worker_proto_rawDescGZIP(), []int{1}
}

func (x *MediaItemProcessResponse) GetOk() bool {
	if x != nil {
		return x.Ok
	}
	return false
}

type GenerateEmbeddingRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Text string `protobuf:"bytes,1,opt,name=text,proto3" json:"text,omitempty"`
}

func (x *GenerateEmbeddingRequest) Reset() {
	*x = GenerateEmbeddingRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_worker_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GenerateEmbeddingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GenerateEmbeddingRequest) ProtoMessage() {}

func (x *GenerateEmbeddingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_worker_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GenerateEmbeddingRequest.ProtoReflect.Descriptor instead.
func (*GenerateEmbeddingRequest) Descriptor() ([]byte, []int) {
	return file_worker_proto_rawDescGZIP(), []int{2}
}

func (x *GenerateEmbeddingRequest) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

type GenerateEmbeddingResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Embedding []float32 `protobuf:"fixed32,1,rep,packed,name=embedding,proto3" json:"embedding,omitempty"`
}

func (x *GenerateEmbeddingResponse) Reset() {
	*x = GenerateEmbeddingResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_worker_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GenerateEmbeddingResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GenerateEmbeddingResponse) ProtoMessage() {}

func (x *GenerateEmbeddingResponse) ProtoReflect() protoreflect.Message {
	mi := &file_worker_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GenerateEmbeddingResponse.ProtoReflect.Descriptor instead.
func (*GenerateEmbeddingResponse) Descriptor() ([]byte, []int) {
	return file_worker_proto_rawDescGZIP(), []int{3}
}

func (x *GenerateEmbeddingResponse) GetEmbedding() []float32 {
	if x != nil {
		return x.Embedding
	}
	return nil
}

var File_worker_proto protoreflect.FileDescriptor

var file_worker_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x8f,
	0x02, 0x0a, 0x17, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x49, 0x74, 0x65, 0x6d, 0x50, 0x72, 0x6f, 0x63,
	0x65, 0x73, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73,
	0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72,
	0x49, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x50, 0x61, 0x74, 0x68, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x50, 0x61, 0x74, 0x68, 0x12, 0x33,
	0x0a, 0x0a, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x04, 0x20, 0x03,
	0x28, 0x0e, 0x32, 0x13, 0x2e, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x49, 0x74, 0x65, 0x6d, 0x43, 0x6f,
	0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x52, 0x0a, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65,
	0x6e, 0x74, 0x73, 0x12, 0x3f, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x05,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x49, 0x74, 0x65, 0x6d,
	0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x50,
	0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x70, 0x61, 0x79,
	0x6c, 0x6f, 0x61, 0x64, 0x1a, 0x3a, 0x0a, 0x0c, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01,
	0x22, 0x2a, 0x0a, 0x18, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x49, 0x74, 0x65, 0x6d, 0x50, 0x72, 0x6f,
	0x63, 0x65, 0x73, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02,
	0x6f, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x02, 0x6f, 0x6b, 0x22, 0x2e, 0x0a, 0x18,
	0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x45, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x69, 0x6e,
	0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x78, 0x74,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74, 0x22, 0x39, 0x0a, 0x19,
	0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x45, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x69, 0x6e,
	0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x65, 0x6d, 0x62,
	0x65, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x18, 0x01, 0x20, 0x03, 0x28, 0x02, 0x52, 0x09, 0x65, 0x6d,
	0x62, 0x65, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x2a, 0x79, 0x0a, 0x12, 0x4d, 0x65, 0x64, 0x69, 0x61,
	0x49, 0x74, 0x65, 0x6d, 0x43, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x12, 0x0c, 0x0a,
	0x08, 0x4d, 0x45, 0x54, 0x41, 0x44, 0x41, 0x54, 0x41, 0x10, 0x00, 0x12, 0x15, 0x0a, 0x11, 0x50,
	0x52, 0x45, 0x56, 0x49, 0x45, 0x57, 0x5f, 0x54, 0x48, 0x55, 0x4d, 0x42, 0x4e, 0x41, 0x49, 0x4c,
	0x10, 0x01, 0x12, 0x0a, 0x0a, 0x06, 0x50, 0x4c, 0x41, 0x43, 0x45, 0x53, 0x10, 0x02, 0x12, 0x12,
	0x0a, 0x0e, 0x43, 0x4c, 0x41, 0x53, 0x53, 0x49, 0x46, 0x49, 0x43, 0x41, 0x54, 0x49, 0x4f, 0x4e,
	0x10, 0x03, 0x12, 0x09, 0x0a, 0x05, 0x46, 0x41, 0x43, 0x45, 0x53, 0x10, 0x04, 0x12, 0x07, 0x0a,
	0x03, 0x4f, 0x43, 0x52, 0x10, 0x05, 0x12, 0x0a, 0x0a, 0x06, 0x53, 0x45, 0x41, 0x52, 0x43, 0x48,
	0x10, 0x06, 0x32, 0xa1, 0x01, 0x0a, 0x06, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x12, 0x49, 0x0a,
	0x10, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x49, 0x74, 0x65, 0x6d, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73,
	0x73, 0x12, 0x18, 0x2e, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x49, 0x74, 0x65, 0x6d, 0x50, 0x72, 0x6f,
	0x63, 0x65, 0x73, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x4d, 0x65,
	0x64, 0x69, 0x61, 0x49, 0x74, 0x65, 0x6d, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x4c, 0x0a, 0x11, 0x47, 0x65, 0x6e, 0x65,
	0x72, 0x61, 0x74, 0x65, 0x45, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x12, 0x19, 0x2e,
	0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x45, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x69, 0x6e,
	0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x47, 0x65, 0x6e, 0x65, 0x72,
	0x61, 0x74, 0x65, 0x45, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x09, 0x5a, 0x07, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72,
	0x2f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_worker_proto_rawDescOnce sync.Once
	file_worker_proto_rawDescData = file_worker_proto_rawDesc
)

func file_worker_proto_rawDescGZIP() []byte {
	file_worker_proto_rawDescOnce.Do(func() {
		file_worker_proto_rawDescData = protoimpl.X.CompressGZIP(file_worker_proto_rawDescData)
	})
	return file_worker_proto_rawDescData
}

var file_worker_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_worker_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_worker_proto_goTypes = []interface{}{
	(MediaItemComponent)(0),           // 0: MediaItemComponent
	(*MediaItemProcessRequest)(nil),   // 1: MediaItemProcessRequest
	(*MediaItemProcessResponse)(nil),  // 2: MediaItemProcessResponse
	(*GenerateEmbeddingRequest)(nil),  // 3: GenerateEmbeddingRequest
	(*GenerateEmbeddingResponse)(nil), // 4: GenerateEmbeddingResponse
	nil,                               // 5: MediaItemProcessRequest.PayloadEntry
}
var file_worker_proto_depIdxs = []int32{
	0, // 0: MediaItemProcessRequest.components:type_name -> MediaItemComponent
	5, // 1: MediaItemProcessRequest.payload:type_name -> MediaItemProcessRequest.PayloadEntry
	1, // 2: Worker.MediaItemProcess:input_type -> MediaItemProcessRequest
	3, // 3: Worker.GenerateEmbedding:input_type -> GenerateEmbeddingRequest
	2, // 4: Worker.MediaItemProcess:output_type -> MediaItemProcessResponse
	4, // 5: Worker.GenerateEmbedding:output_type -> GenerateEmbeddingResponse
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_worker_proto_init() }
func file_worker_proto_init() {
	if File_worker_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_worker_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MediaItemProcessRequest); i {
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
		file_worker_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MediaItemProcessResponse); i {
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
		file_worker_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GenerateEmbeddingRequest); i {
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
		file_worker_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GenerateEmbeddingResponse); i {
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
			RawDescriptor: file_worker_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_worker_proto_goTypes,
		DependencyIndexes: file_worker_proto_depIdxs,
		EnumInfos:         file_worker_proto_enumTypes,
		MessageInfos:      file_worker_proto_msgTypes,
	}.Build()
	File_worker_proto = out.File
	file_worker_proto_rawDesc = nil
	file_worker_proto_goTypes = nil
	file_worker_proto_depIdxs = nil
}
