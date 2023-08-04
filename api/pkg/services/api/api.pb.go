// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.23.4
// source: api.proto

package api

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ConfigResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Config []byte `protobuf:"bytes,1,opt,name=config,proto3" json:"config,omitempty"`
}

func (x *ConfigResponse) Reset() {
	*x = ConfigResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConfigResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConfigResponse) ProtoMessage() {}

func (x *ConfigResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConfigResponse.ProtoReflect.Descriptor instead.
func (*ConfigResponse) Descriptor() ([]byte, []int) {
	return file_api_proto_rawDescGZIP(), []int{0}
}

func (x *ConfigResponse) GetConfig() []byte {
	if x != nil {
		return x.Config
	}
	return nil
}

type MediaItemMetadataRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId          string   `protobuf:"bytes,1,opt,name=userId,proto3" json:"userId,omitempty"`
	Id              string   `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	Status          string   `protobuf:"bytes,3,opt,name=status,proto3" json:"status,omitempty"`
	MimeType        *string  `protobuf:"bytes,4,opt,name=mimeType,proto3,oneof" json:"mimeType,omitempty"`
	SourcePath      string   `protobuf:"bytes,5,opt,name=sourcePath,proto3" json:"sourcePath,omitempty"`
	PreviewPath     *string  `protobuf:"bytes,6,opt,name=previewPath,proto3,oneof" json:"previewPath,omitempty"`
	ThumbnailPath   *string  `protobuf:"bytes,7,opt,name=thumbnailPath,proto3,oneof" json:"thumbnailPath,omitempty"`
	Type            string   `protobuf:"bytes,8,opt,name=type,proto3" json:"type,omitempty"`
	Category        string   `protobuf:"bytes,9,opt,name=category,proto3" json:"category,omitempty"`
	Width           *int32   `protobuf:"varint,10,opt,name=width,proto3,oneof" json:"width,omitempty"`
	Height          *int32   `protobuf:"varint,11,opt,name=height,proto3,oneof" json:"height,omitempty"`
	CreationTime    *string  `protobuf:"bytes,12,opt,name=creationTime,proto3,oneof" json:"creationTime,omitempty"`
	CameraMake      *string  `protobuf:"bytes,13,opt,name=cameraMake,proto3,oneof" json:"cameraMake,omitempty"`
	CameraModel     *string  `protobuf:"bytes,14,opt,name=cameraModel,proto3,oneof" json:"cameraModel,omitempty"`
	FocalLength     *string  `protobuf:"bytes,15,opt,name=focalLength,proto3,oneof" json:"focalLength,omitempty"`
	ApertureFNumber *string  `protobuf:"bytes,16,opt,name=apertureFNumber,proto3,oneof" json:"apertureFNumber,omitempty"`
	IsoEquivalent   *string  `protobuf:"bytes,17,opt,name=isoEquivalent,proto3,oneof" json:"isoEquivalent,omitempty"`
	ExposureTime    *string  `protobuf:"bytes,18,opt,name=exposureTime,proto3,oneof" json:"exposureTime,omitempty"`
	Fps             *string  `protobuf:"bytes,19,opt,name=fps,proto3,oneof" json:"fps,omitempty"`
	Latitude        *float64 `protobuf:"fixed64,20,opt,name=latitude,proto3,oneof" json:"latitude,omitempty"`
	Longitude       *float64 `protobuf:"fixed64,21,opt,name=longitude,proto3,oneof" json:"longitude,omitempty"`
}

func (x *MediaItemMetadataRequest) Reset() {
	*x = MediaItemMetadataRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MediaItemMetadataRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MediaItemMetadataRequest) ProtoMessage() {}

func (x *MediaItemMetadataRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MediaItemMetadataRequest.ProtoReflect.Descriptor instead.
func (*MediaItemMetadataRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_rawDescGZIP(), []int{1}
}

func (x *MediaItemMetadataRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *MediaItemMetadataRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *MediaItemMetadataRequest) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *MediaItemMetadataRequest) GetMimeType() string {
	if x != nil && x.MimeType != nil {
		return *x.MimeType
	}
	return ""
}

func (x *MediaItemMetadataRequest) GetSourcePath() string {
	if x != nil {
		return x.SourcePath
	}
	return ""
}

func (x *MediaItemMetadataRequest) GetPreviewPath() string {
	if x != nil && x.PreviewPath != nil {
		return *x.PreviewPath
	}
	return ""
}

func (x *MediaItemMetadataRequest) GetThumbnailPath() string {
	if x != nil && x.ThumbnailPath != nil {
		return *x.ThumbnailPath
	}
	return ""
}

func (x *MediaItemMetadataRequest) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *MediaItemMetadataRequest) GetCategory() string {
	if x != nil {
		return x.Category
	}
	return ""
}

func (x *MediaItemMetadataRequest) GetWidth() int32 {
	if x != nil && x.Width != nil {
		return *x.Width
	}
	return 0
}

func (x *MediaItemMetadataRequest) GetHeight() int32 {
	if x != nil && x.Height != nil {
		return *x.Height
	}
	return 0
}

func (x *MediaItemMetadataRequest) GetCreationTime() string {
	if x != nil && x.CreationTime != nil {
		return *x.CreationTime
	}
	return ""
}

func (x *MediaItemMetadataRequest) GetCameraMake() string {
	if x != nil && x.CameraMake != nil {
		return *x.CameraMake
	}
	return ""
}

func (x *MediaItemMetadataRequest) GetCameraModel() string {
	if x != nil && x.CameraModel != nil {
		return *x.CameraModel
	}
	return ""
}

func (x *MediaItemMetadataRequest) GetFocalLength() string {
	if x != nil && x.FocalLength != nil {
		return *x.FocalLength
	}
	return ""
}

func (x *MediaItemMetadataRequest) GetApertureFNumber() string {
	if x != nil && x.ApertureFNumber != nil {
		return *x.ApertureFNumber
	}
	return ""
}

func (x *MediaItemMetadataRequest) GetIsoEquivalent() string {
	if x != nil && x.IsoEquivalent != nil {
		return *x.IsoEquivalent
	}
	return ""
}

func (x *MediaItemMetadataRequest) GetExposureTime() string {
	if x != nil && x.ExposureTime != nil {
		return *x.ExposureTime
	}
	return ""
}

func (x *MediaItemMetadataRequest) GetFps() string {
	if x != nil && x.Fps != nil {
		return *x.Fps
	}
	return ""
}

func (x *MediaItemMetadataRequest) GetLatitude() float64 {
	if x != nil && x.Latitude != nil {
		return *x.Latitude
	}
	return 0
}

func (x *MediaItemMetadataRequest) GetLongitude() float64 {
	if x != nil && x.Longitude != nil {
		return *x.Longitude
	}
	return 0
}

type MediaItemPlaceRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId   string  `protobuf:"bytes,1,opt,name=userId,proto3" json:"userId,omitempty"`
	Id       string  `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	Postcode *string `protobuf:"bytes,3,opt,name=postcode,proto3,oneof" json:"postcode,omitempty"`
	Country  *string `protobuf:"bytes,4,opt,name=country,proto3,oneof" json:"country,omitempty"`
	State    *string `protobuf:"bytes,5,opt,name=state,proto3,oneof" json:"state,omitempty"`
	City     *string `protobuf:"bytes,6,opt,name=city,proto3,oneof" json:"city,omitempty"`
	Town     *string `protobuf:"bytes,7,opt,name=town,proto3,oneof" json:"town,omitempty"`
}

func (x *MediaItemPlaceRequest) Reset() {
	*x = MediaItemPlaceRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MediaItemPlaceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MediaItemPlaceRequest) ProtoMessage() {}

func (x *MediaItemPlaceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MediaItemPlaceRequest.ProtoReflect.Descriptor instead.
func (*MediaItemPlaceRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_rawDescGZIP(), []int{2}
}

func (x *MediaItemPlaceRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *MediaItemPlaceRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *MediaItemPlaceRequest) GetPostcode() string {
	if x != nil && x.Postcode != nil {
		return *x.Postcode
	}
	return ""
}

func (x *MediaItemPlaceRequest) GetCountry() string {
	if x != nil && x.Country != nil {
		return *x.Country
	}
	return ""
}

func (x *MediaItemPlaceRequest) GetState() string {
	if x != nil && x.State != nil {
		return *x.State
	}
	return ""
}

func (x *MediaItemPlaceRequest) GetCity() string {
	if x != nil && x.City != nil {
		return *x.City
	}
	return ""
}

func (x *MediaItemPlaceRequest) GetTown() string {
	if x != nil && x.Town != nil {
		return *x.Town
	}
	return ""
}

type MediaItemThingRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId string `protobuf:"bytes,1,opt,name=userId,proto3" json:"userId,omitempty"`
	Id     string `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	Name   string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *MediaItemThingRequest) Reset() {
	*x = MediaItemThingRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MediaItemThingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MediaItemThingRequest) ProtoMessage() {}

func (x *MediaItemThingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MediaItemThingRequest.ProtoReflect.Descriptor instead.
func (*MediaItemThingRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_rawDescGZIP(), []int{3}
}

func (x *MediaItemThingRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *MediaItemThingRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *MediaItemThingRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type MediaItemEmbedding struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Embedding []float32 `protobuf:"fixed32,1,rep,packed,name=embedding,proto3" json:"embedding,omitempty"`
}

func (x *MediaItemEmbedding) Reset() {
	*x = MediaItemEmbedding{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MediaItemEmbedding) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MediaItemEmbedding) ProtoMessage() {}

func (x *MediaItemEmbedding) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MediaItemEmbedding.ProtoReflect.Descriptor instead.
func (*MediaItemEmbedding) Descriptor() ([]byte, []int) {
	return file_api_proto_rawDescGZIP(), []int{4}
}

func (x *MediaItemEmbedding) GetEmbedding() []float32 {
	if x != nil {
		return x.Embedding
	}
	return nil
}

type MediaItemFacesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId     string                `protobuf:"bytes,1,opt,name=userId,proto3" json:"userId,omitempty"`
	Id         string                `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	Embeddings []*MediaItemEmbedding `protobuf:"bytes,3,rep,name=embeddings,proto3" json:"embeddings,omitempty"`
}

func (x *MediaItemFacesRequest) Reset() {
	*x = MediaItemFacesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MediaItemFacesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MediaItemFacesRequest) ProtoMessage() {}

func (x *MediaItemFacesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MediaItemFacesRequest.ProtoReflect.Descriptor instead.
func (*MediaItemFacesRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_rawDescGZIP(), []int{5}
}

func (x *MediaItemFacesRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *MediaItemFacesRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *MediaItemFacesRequest) GetEmbeddings() []*MediaItemEmbedding {
	if x != nil {
		return x.Embeddings
	}
	return nil
}

type MediaItemFinalResultRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId     string                `protobuf:"bytes,1,opt,name=userId,proto3" json:"userId,omitempty"`
	Id         string                `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	Keywords   string                `protobuf:"bytes,3,opt,name=keywords,proto3" json:"keywords,omitempty"`
	Embeddings []*MediaItemEmbedding `protobuf:"bytes,4,rep,name=embeddings,proto3" json:"embeddings,omitempty"`
}

func (x *MediaItemFinalResultRequest) Reset() {
	*x = MediaItemFinalResultRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MediaItemFinalResultRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MediaItemFinalResultRequest) ProtoMessage() {}

func (x *MediaItemFinalResultRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MediaItemFinalResultRequest.ProtoReflect.Descriptor instead.
func (*MediaItemFinalResultRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_rawDescGZIP(), []int{6}
}

func (x *MediaItemFinalResultRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *MediaItemFinalResultRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *MediaItemFinalResultRequest) GetKeywords() string {
	if x != nil {
		return x.Keywords
	}
	return ""
}

func (x *MediaItemFinalResultRequest) GetEmbeddings() []*MediaItemEmbedding {
	if x != nil {
		return x.Embeddings
	}
	return nil
}

var File_api_proto protoreflect.FileDescriptor

var file_api_proto_rawDesc = []byte{
	0x0a, 0x09, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70,
	0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x28, 0x0a, 0x0e, 0x43, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x63, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x63, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x22, 0xad, 0x07, 0x0a, 0x18, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x49, 0x74, 0x65, 0x6d,
	0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12,
	0x1f, 0x0a, 0x08, 0x6d, 0x69, 0x6d, 0x65, 0x54, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x48, 0x00, 0x52, 0x08, 0x6d, 0x69, 0x6d, 0x65, 0x54, 0x79, 0x70, 0x65, 0x88, 0x01, 0x01,
	0x12, 0x1e, 0x0a, 0x0a, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x50, 0x61, 0x74, 0x68, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x50, 0x61, 0x74, 0x68,
	0x12, 0x25, 0x0a, 0x0b, 0x70, 0x72, 0x65, 0x76, 0x69, 0x65, 0x77, 0x50, 0x61, 0x74, 0x68, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x09, 0x48, 0x01, 0x52, 0x0b, 0x70, 0x72, 0x65, 0x76, 0x69, 0x65, 0x77,
	0x50, 0x61, 0x74, 0x68, 0x88, 0x01, 0x01, 0x12, 0x29, 0x0a, 0x0d, 0x74, 0x68, 0x75, 0x6d, 0x62,
	0x6e, 0x61, 0x69, 0x6c, 0x50, 0x61, 0x74, 0x68, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x48, 0x02,
	0x52, 0x0d, 0x74, 0x68, 0x75, 0x6d, 0x62, 0x6e, 0x61, 0x69, 0x6c, 0x50, 0x61, 0x74, 0x68, 0x88,
	0x01, 0x01, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f,
	0x72, 0x79, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f,
	0x72, 0x79, 0x12, 0x19, 0x0a, 0x05, 0x77, 0x69, 0x64, 0x74, 0x68, 0x18, 0x0a, 0x20, 0x01, 0x28,
	0x05, 0x48, 0x03, 0x52, 0x05, 0x77, 0x69, 0x64, 0x74, 0x68, 0x88, 0x01, 0x01, 0x12, 0x1b, 0x0a,
	0x06, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x05, 0x48, 0x04, 0x52,
	0x06, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x88, 0x01, 0x01, 0x12, 0x27, 0x0a, 0x0c, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09,
	0x48, 0x05, 0x52, 0x0c, 0x63, 0x72, 0x65, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x6d, 0x65,
	0x88, 0x01, 0x01, 0x12, 0x23, 0x0a, 0x0a, 0x63, 0x61, 0x6d, 0x65, 0x72, 0x61, 0x4d, 0x61, 0x6b,
	0x65, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x09, 0x48, 0x06, 0x52, 0x0a, 0x63, 0x61, 0x6d, 0x65, 0x72,
	0x61, 0x4d, 0x61, 0x6b, 0x65, 0x88, 0x01, 0x01, 0x12, 0x25, 0x0a, 0x0b, 0x63, 0x61, 0x6d, 0x65,
	0x72, 0x61, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x09, 0x48, 0x07, 0x52,
	0x0b, 0x63, 0x61, 0x6d, 0x65, 0x72, 0x61, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x88, 0x01, 0x01, 0x12,
	0x25, 0x0a, 0x0b, 0x66, 0x6f, 0x63, 0x61, 0x6c, 0x4c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x18, 0x0f,
	0x20, 0x01, 0x28, 0x09, 0x48, 0x08, 0x52, 0x0b, 0x66, 0x6f, 0x63, 0x61, 0x6c, 0x4c, 0x65, 0x6e,
	0x67, 0x74, 0x68, 0x88, 0x01, 0x01, 0x12, 0x2d, 0x0a, 0x0f, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75,
	0x72, 0x65, 0x46, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x10, 0x20, 0x01, 0x28, 0x09, 0x48,
	0x09, 0x52, 0x0f, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x46, 0x4e, 0x75, 0x6d, 0x62,
	0x65, 0x72, 0x88, 0x01, 0x01, 0x12, 0x29, 0x0a, 0x0d, 0x69, 0x73, 0x6f, 0x45, 0x71, 0x75, 0x69,
	0x76, 0x61, 0x6c, 0x65, 0x6e, 0x74, 0x18, 0x11, 0x20, 0x01, 0x28, 0x09, 0x48, 0x0a, 0x52, 0x0d,
	0x69, 0x73, 0x6f, 0x45, 0x71, 0x75, 0x69, 0x76, 0x61, 0x6c, 0x65, 0x6e, 0x74, 0x88, 0x01, 0x01,
	0x12, 0x27, 0x0a, 0x0c, 0x65, 0x78, 0x70, 0x6f, 0x73, 0x75, 0x72, 0x65, 0x54, 0x69, 0x6d, 0x65,
	0x18, 0x12, 0x20, 0x01, 0x28, 0x09, 0x48, 0x0b, 0x52, 0x0c, 0x65, 0x78, 0x70, 0x6f, 0x73, 0x75,
	0x72, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x88, 0x01, 0x01, 0x12, 0x15, 0x0a, 0x03, 0x66, 0x70, 0x73,
	0x18, 0x13, 0x20, 0x01, 0x28, 0x09, 0x48, 0x0c, 0x52, 0x03, 0x66, 0x70, 0x73, 0x88, 0x01, 0x01,
	0x12, 0x1f, 0x0a, 0x08, 0x6c, 0x61, 0x74, 0x69, 0x74, 0x75, 0x64, 0x65, 0x18, 0x14, 0x20, 0x01,
	0x28, 0x01, 0x48, 0x0d, 0x52, 0x08, 0x6c, 0x61, 0x74, 0x69, 0x74, 0x75, 0x64, 0x65, 0x88, 0x01,
	0x01, 0x12, 0x21, 0x0a, 0x09, 0x6c, 0x6f, 0x6e, 0x67, 0x69, 0x74, 0x75, 0x64, 0x65, 0x18, 0x15,
	0x20, 0x01, 0x28, 0x01, 0x48, 0x0e, 0x52, 0x09, 0x6c, 0x6f, 0x6e, 0x67, 0x69, 0x74, 0x75, 0x64,
	0x65, 0x88, 0x01, 0x01, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x6d, 0x69, 0x6d, 0x65, 0x54, 0x79, 0x70,
	0x65, 0x42, 0x0e, 0x0a, 0x0c, 0x5f, 0x70, 0x72, 0x65, 0x76, 0x69, 0x65, 0x77, 0x50, 0x61, 0x74,
	0x68, 0x42, 0x10, 0x0a, 0x0e, 0x5f, 0x74, 0x68, 0x75, 0x6d, 0x62, 0x6e, 0x61, 0x69, 0x6c, 0x50,
	0x61, 0x74, 0x68, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x77, 0x69, 0x64, 0x74, 0x68, 0x42, 0x09, 0x0a,
	0x07, 0x5f, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x42, 0x0f, 0x0a, 0x0d, 0x5f, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x63, 0x61,
	0x6d, 0x65, 0x72, 0x61, 0x4d, 0x61, 0x6b, 0x65, 0x42, 0x0e, 0x0a, 0x0c, 0x5f, 0x63, 0x61, 0x6d,
	0x65, 0x72, 0x61, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x42, 0x0e, 0x0a, 0x0c, 0x5f, 0x66, 0x6f, 0x63,
	0x61, 0x6c, 0x4c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x42, 0x12, 0x0a, 0x10, 0x5f, 0x61, 0x70, 0x65,
	0x72, 0x74, 0x75, 0x72, 0x65, 0x46, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x42, 0x10, 0x0a, 0x0e,
	0x5f, 0x69, 0x73, 0x6f, 0x45, 0x71, 0x75, 0x69, 0x76, 0x61, 0x6c, 0x65, 0x6e, 0x74, 0x42, 0x0f,
	0x0a, 0x0d, 0x5f, 0x65, 0x78, 0x70, 0x6f, 0x73, 0x75, 0x72, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x42,
	0x06, 0x0a, 0x04, 0x5f, 0x66, 0x70, 0x73, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x6c, 0x61, 0x74, 0x69,
	0x74, 0x75, 0x64, 0x65, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x6c, 0x6f, 0x6e, 0x67, 0x69, 0x74, 0x75,
	0x64, 0x65, 0x22, 0x81, 0x02, 0x0a, 0x15, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x49, 0x74, 0x65, 0x6d,
	0x50, 0x6c, 0x61, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06,
	0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73,
	0x65, 0x72, 0x49, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x1f, 0x0a, 0x08, 0x70, 0x6f, 0x73, 0x74, 0x63, 0x6f, 0x64, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x08, 0x70, 0x6f, 0x73, 0x74, 0x63, 0x6f,
	0x64, 0x65, 0x88, 0x01, 0x01, 0x12, 0x1d, 0x0a, 0x07, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x79,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x48, 0x01, 0x52, 0x07, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x72,
	0x79, 0x88, 0x01, 0x01, 0x12, 0x19, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x48, 0x02, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x88, 0x01, 0x01, 0x12,
	0x17, 0x0a, 0x04, 0x63, 0x69, 0x74, 0x79, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x48, 0x03, 0x52,
	0x04, 0x63, 0x69, 0x74, 0x79, 0x88, 0x01, 0x01, 0x12, 0x17, 0x0a, 0x04, 0x74, 0x6f, 0x77, 0x6e,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x48, 0x04, 0x52, 0x04, 0x74, 0x6f, 0x77, 0x6e, 0x88, 0x01,
	0x01, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x70, 0x6f, 0x73, 0x74, 0x63, 0x6f, 0x64, 0x65, 0x42, 0x0a,
	0x0a, 0x08, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x79, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x73,
	0x74, 0x61, 0x74, 0x65, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x63, 0x69, 0x74, 0x79, 0x42, 0x07, 0x0a,
	0x05, 0x5f, 0x74, 0x6f, 0x77, 0x6e, 0x22, 0x53, 0x0a, 0x15, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x49,
	0x74, 0x65, 0x6d, 0x54, 0x68, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x32, 0x0a, 0x12, 0x4d,
	0x65, 0x64, 0x69, 0x61, 0x49, 0x74, 0x65, 0x6d, 0x45, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x69, 0x6e,
	0x67, 0x12, 0x1c, 0x0a, 0x09, 0x65, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x02, 0x52, 0x09, 0x65, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x22,
	0x74, 0x0a, 0x15, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x49, 0x74, 0x65, 0x6d, 0x46, 0x61, 0x63, 0x65,
	0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72,
	0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x33, 0x0a, 0x0a, 0x65, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x73, 0x18, 0x03,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x49, 0x74, 0x65, 0x6d,
	0x45, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x52, 0x0a, 0x65, 0x6d, 0x62, 0x65, 0x64,
	0x64, 0x69, 0x6e, 0x67, 0x73, 0x22, 0x96, 0x01, 0x0a, 0x1b, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x49,
	0x74, 0x65, 0x6d, 0x46, 0x69, 0x6e, 0x61, 0x6c, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1a, 0x0a,
	0x08, 0x6b, 0x65, 0x79, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x6b, 0x65, 0x79, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x12, 0x33, 0x0a, 0x0a, 0x65, 0x6d, 0x62,
	0x65, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e,
	0x4d, 0x65, 0x64, 0x69, 0x61, 0x49, 0x74, 0x65, 0x6d, 0x45, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x69,
	0x6e, 0x67, 0x52, 0x0a, 0x65, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x73, 0x32, 0xbd,
	0x03, 0x0a, 0x03, 0x41, 0x50, 0x49, 0x12, 0x3c, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x57, 0x6f, 0x72,
	0x6b, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x1a, 0x0f, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x00, 0x12, 0x4c, 0x0a, 0x15, 0x53, 0x61, 0x76, 0x65, 0x4d, 0x65, 0x64, 0x69,
	0x61, 0x49, 0x74, 0x65, 0x6d, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x19, 0x2e,
	0x4d, 0x65, 0x64, 0x69, 0x61, 0x49, 0x74, 0x65, 0x6d, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x22, 0x00, 0x12, 0x46, 0x0a, 0x12, 0x53, 0x61, 0x76, 0x65, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x49,
	0x74, 0x65, 0x6d, 0x50, 0x6c, 0x61, 0x63, 0x65, 0x12, 0x16, 0x2e, 0x4d, 0x65, 0x64, 0x69, 0x61,
	0x49, 0x74, 0x65, 0x6d, 0x50, 0x6c, 0x61, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x46, 0x0a, 0x12, 0x53, 0x61,
	0x76, 0x65, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x49, 0x74, 0x65, 0x6d, 0x54, 0x68, 0x69, 0x6e, 0x67,
	0x12, 0x16, 0x2e, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x49, 0x74, 0x65, 0x6d, 0x54, 0x68, 0x69, 0x6e,
	0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x22, 0x00, 0x12, 0x46, 0x0a, 0x12, 0x53, 0x61, 0x76, 0x65, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x49,
	0x74, 0x65, 0x6d, 0x46, 0x61, 0x63, 0x65, 0x73, 0x12, 0x16, 0x2e, 0x4d, 0x65, 0x64, 0x69, 0x61,
	0x49, 0x74, 0x65, 0x6d, 0x46, 0x61, 0x63, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x52, 0x0a, 0x18, 0x53, 0x61,
	0x76, 0x65, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x49, 0x74, 0x65, 0x6d, 0x46, 0x69, 0x6e, 0x61, 0x6c,
	0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x1c, 0x2e, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x49, 0x74,
	0x65, 0x6d, 0x46, 0x69, 0x6e, 0x61, 0x6c, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x42, 0x06,
	0x5a, 0x04, 0x61, 0x70, 0x69, 0x2f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_proto_rawDescOnce sync.Once
	file_api_proto_rawDescData = file_api_proto_rawDesc
)

func file_api_proto_rawDescGZIP() []byte {
	file_api_proto_rawDescOnce.Do(func() {
		file_api_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_proto_rawDescData)
	})
	return file_api_proto_rawDescData
}

var file_api_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_api_proto_goTypes = []interface{}{
	(*ConfigResponse)(nil),              // 0: ConfigResponse
	(*MediaItemMetadataRequest)(nil),    // 1: MediaItemMetadataRequest
	(*MediaItemPlaceRequest)(nil),       // 2: MediaItemPlaceRequest
	(*MediaItemThingRequest)(nil),       // 3: MediaItemThingRequest
	(*MediaItemEmbedding)(nil),          // 4: MediaItemEmbedding
	(*MediaItemFacesRequest)(nil),       // 5: MediaItemFacesRequest
	(*MediaItemFinalResultRequest)(nil), // 6: MediaItemFinalResultRequest
	(*emptypb.Empty)(nil),               // 7: google.protobuf.Empty
}
var file_api_proto_depIdxs = []int32{
	4, // 0: MediaItemFacesRequest.embeddings:type_name -> MediaItemEmbedding
	4, // 1: MediaItemFinalResultRequest.embeddings:type_name -> MediaItemEmbedding
	7, // 2: API.GetWorkerConfig:input_type -> google.protobuf.Empty
	1, // 3: API.SaveMediaItemMetadata:input_type -> MediaItemMetadataRequest
	2, // 4: API.SaveMediaItemPlace:input_type -> MediaItemPlaceRequest
	3, // 5: API.SaveMediaItemThing:input_type -> MediaItemThingRequest
	5, // 6: API.SaveMediaItemFaces:input_type -> MediaItemFacesRequest
	6, // 7: API.SaveMediaItemFinalResult:input_type -> MediaItemFinalResultRequest
	0, // 8: API.GetWorkerConfig:output_type -> ConfigResponse
	7, // 9: API.SaveMediaItemMetadata:output_type -> google.protobuf.Empty
	7, // 10: API.SaveMediaItemPlace:output_type -> google.protobuf.Empty
	7, // 11: API.SaveMediaItemThing:output_type -> google.protobuf.Empty
	7, // 12: API.SaveMediaItemFaces:output_type -> google.protobuf.Empty
	7, // 13: API.SaveMediaItemFinalResult:output_type -> google.protobuf.Empty
	8, // [8:14] is the sub-list for method output_type
	2, // [2:8] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_api_proto_init() }
func file_api_proto_init() {
	if File_api_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConfigResponse); i {
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
		file_api_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MediaItemMetadataRequest); i {
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
		file_api_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MediaItemPlaceRequest); i {
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
		file_api_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MediaItemThingRequest); i {
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
		file_api_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MediaItemEmbedding); i {
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
		file_api_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MediaItemFacesRequest); i {
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
		file_api_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MediaItemFinalResultRequest); i {
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
	file_api_proto_msgTypes[1].OneofWrappers = []interface{}{}
	file_api_proto_msgTypes[2].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_proto_goTypes,
		DependencyIndexes: file_api_proto_depIdxs,
		MessageInfos:      file_api_proto_msgTypes,
	}.Build()
	File_api_proto = out.File
	file_api_proto_rawDesc = nil
	file_api_proto_goTypes = nil
	file_api_proto_depIdxs = nil
}
