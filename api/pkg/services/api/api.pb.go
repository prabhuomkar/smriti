// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.12.4
// source: api.proto

package api

import (
	empty "github.com/golang/protobuf/ptypes/empty"
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

type MediaItemMetadataRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId          string   `protobuf:"bytes,1,opt,name=userId,proto3" json:"userId,omitempty"`
	Id              string   `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	Status          string   `protobuf:"bytes,3,opt,name=status,proto3" json:"status,omitempty"`
	MimeType        *string  `protobuf:"bytes,4,opt,name=mimeType,proto3,oneof" json:"mimeType,omitempty"`
	SourceUrl       *string  `protobuf:"bytes,5,opt,name=sourceUrl,proto3,oneof" json:"sourceUrl,omitempty"`
	PreviewUrl      *string  `protobuf:"bytes,6,opt,name=previewUrl,proto3,oneof" json:"previewUrl,omitempty"`
	ThumbnailUrl    *string  `protobuf:"bytes,7,opt,name=thumbnailUrl,proto3,oneof" json:"thumbnailUrl,omitempty"`
	Type            *string  `protobuf:"bytes,8,opt,name=type,proto3,oneof" json:"type,omitempty"`
	Width           *int32   `protobuf:"varint,9,opt,name=width,proto3,oneof" json:"width,omitempty"`
	Height          *int32   `protobuf:"varint,10,opt,name=height,proto3,oneof" json:"height,omitempty"`
	CreationTime    *string  `protobuf:"bytes,11,opt,name=creationTime,proto3,oneof" json:"creationTime,omitempty"`
	CameraMake      *string  `protobuf:"bytes,12,opt,name=cameraMake,proto3,oneof" json:"cameraMake,omitempty"`
	CameraModel     *string  `protobuf:"bytes,13,opt,name=cameraModel,proto3,oneof" json:"cameraModel,omitempty"`
	FocalLength     *string  `protobuf:"bytes,14,opt,name=focalLength,proto3,oneof" json:"focalLength,omitempty"`
	ApertureFNumber *string  `protobuf:"bytes,15,opt,name=apertureFNumber,proto3,oneof" json:"apertureFNumber,omitempty"`
	IsoEquivalent   *string  `protobuf:"bytes,16,opt,name=isoEquivalent,proto3,oneof" json:"isoEquivalent,omitempty"`
	ExposureTime    *string  `protobuf:"bytes,17,opt,name=exposureTime,proto3,oneof" json:"exposureTime,omitempty"`
	Fps             *string  `protobuf:"bytes,18,opt,name=fps,proto3,oneof" json:"fps,omitempty"`
	Latitude        *float64 `protobuf:"fixed64,19,opt,name=latitude,proto3,oneof" json:"latitude,omitempty"`
	Longitude       *float64 `protobuf:"fixed64,20,opt,name=longitude,proto3,oneof" json:"longitude,omitempty"`
}

func (x *MediaItemMetadataRequest) Reset() {
	*x = MediaItemMetadataRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MediaItemMetadataRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MediaItemMetadataRequest) ProtoMessage() {}

func (x *MediaItemMetadataRequest) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use MediaItemMetadataRequest.ProtoReflect.Descriptor instead.
func (*MediaItemMetadataRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_rawDescGZIP(), []int{0}
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

func (x *MediaItemMetadataRequest) GetSourceUrl() string {
	if x != nil && x.SourceUrl != nil {
		return *x.SourceUrl
	}
	return ""
}

func (x *MediaItemMetadataRequest) GetPreviewUrl() string {
	if x != nil && x.PreviewUrl != nil {
		return *x.PreviewUrl
	}
	return ""
}

func (x *MediaItemMetadataRequest) GetThumbnailUrl() string {
	if x != nil && x.ThumbnailUrl != nil {
		return *x.ThumbnailUrl
	}
	return ""
}

func (x *MediaItemMetadataRequest) GetType() string {
	if x != nil && x.Type != nil {
		return *x.Type
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
		mi := &file_api_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MediaItemPlaceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MediaItemPlaceRequest) ProtoMessage() {}

func (x *MediaItemPlaceRequest) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use MediaItemPlaceRequest.ProtoReflect.Descriptor instead.
func (*MediaItemPlaceRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_rawDescGZIP(), []int{1}
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

var File_api_proto protoreflect.FileDescriptor

var file_api_proto_rawDesc = []byte{
	0x0a, 0x09, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70,
	0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xaa, 0x07, 0x0a, 0x18, 0x4d, 0x65, 0x64,
	0x69, 0x61, 0x49, 0x74, 0x65, 0x6d, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x16, 0x0a,
	0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1f, 0x0a, 0x08, 0x6d, 0x69, 0x6d, 0x65, 0x54, 0x79, 0x70,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x08, 0x6d, 0x69, 0x6d, 0x65, 0x54,
	0x79, 0x70, 0x65, 0x88, 0x01, 0x01, 0x12, 0x21, 0x0a, 0x09, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x55, 0x72, 0x6c, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x48, 0x01, 0x52, 0x09, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x55, 0x72, 0x6c, 0x88, 0x01, 0x01, 0x12, 0x23, 0x0a, 0x0a, 0x70, 0x72, 0x65,
	0x76, 0x69, 0x65, 0x77, 0x55, 0x72, 0x6c, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x48, 0x02, 0x52,
	0x0a, 0x70, 0x72, 0x65, 0x76, 0x69, 0x65, 0x77, 0x55, 0x72, 0x6c, 0x88, 0x01, 0x01, 0x12, 0x27,
	0x0a, 0x0c, 0x74, 0x68, 0x75, 0x6d, 0x62, 0x6e, 0x61, 0x69, 0x6c, 0x55, 0x72, 0x6c, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x09, 0x48, 0x03, 0x52, 0x0c, 0x74, 0x68, 0x75, 0x6d, 0x62, 0x6e, 0x61, 0x69,
	0x6c, 0x55, 0x72, 0x6c, 0x88, 0x01, 0x01, 0x12, 0x17, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18,
	0x08, 0x20, 0x01, 0x28, 0x09, 0x48, 0x04, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x88, 0x01, 0x01,
	0x12, 0x19, 0x0a, 0x05, 0x77, 0x69, 0x64, 0x74, 0x68, 0x18, 0x09, 0x20, 0x01, 0x28, 0x05, 0x48,
	0x05, 0x52, 0x05, 0x77, 0x69, 0x64, 0x74, 0x68, 0x88, 0x01, 0x01, 0x12, 0x1b, 0x0a, 0x06, 0x68,
	0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x05, 0x48, 0x06, 0x52, 0x06, 0x68,
	0x65, 0x69, 0x67, 0x68, 0x74, 0x88, 0x01, 0x01, 0x12, 0x27, 0x0a, 0x0c, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x48, 0x07,
	0x52, 0x0c, 0x63, 0x72, 0x65, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x88, 0x01,
	0x01, 0x12, 0x23, 0x0a, 0x0a, 0x63, 0x61, 0x6d, 0x65, 0x72, 0x61, 0x4d, 0x61, 0x6b, 0x65, 0x18,
	0x0c, 0x20, 0x01, 0x28, 0x09, 0x48, 0x08, 0x52, 0x0a, 0x63, 0x61, 0x6d, 0x65, 0x72, 0x61, 0x4d,
	0x61, 0x6b, 0x65, 0x88, 0x01, 0x01, 0x12, 0x25, 0x0a, 0x0b, 0x63, 0x61, 0x6d, 0x65, 0x72, 0x61,
	0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x09, 0x48, 0x09, 0x52, 0x0b, 0x63,
	0x61, 0x6d, 0x65, 0x72, 0x61, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x88, 0x01, 0x01, 0x12, 0x25, 0x0a,
	0x0b, 0x66, 0x6f, 0x63, 0x61, 0x6c, 0x4c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x18, 0x0e, 0x20, 0x01,
	0x28, 0x09, 0x48, 0x0a, 0x52, 0x0b, 0x66, 0x6f, 0x63, 0x61, 0x6c, 0x4c, 0x65, 0x6e, 0x67, 0x74,
	0x68, 0x88, 0x01, 0x01, 0x12, 0x2d, 0x0a, 0x0f, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65,
	0x46, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x09, 0x48, 0x0b, 0x52,
	0x0f, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x46, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72,
	0x88, 0x01, 0x01, 0x12, 0x29, 0x0a, 0x0d, 0x69, 0x73, 0x6f, 0x45, 0x71, 0x75, 0x69, 0x76, 0x61,
	0x6c, 0x65, 0x6e, 0x74, 0x18, 0x10, 0x20, 0x01, 0x28, 0x09, 0x48, 0x0c, 0x52, 0x0d, 0x69, 0x73,
	0x6f, 0x45, 0x71, 0x75, 0x69, 0x76, 0x61, 0x6c, 0x65, 0x6e, 0x74, 0x88, 0x01, 0x01, 0x12, 0x27,
	0x0a, 0x0c, 0x65, 0x78, 0x70, 0x6f, 0x73, 0x75, 0x72, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x11,
	0x20, 0x01, 0x28, 0x09, 0x48, 0x0d, 0x52, 0x0c, 0x65, 0x78, 0x70, 0x6f, 0x73, 0x75, 0x72, 0x65,
	0x54, 0x69, 0x6d, 0x65, 0x88, 0x01, 0x01, 0x12, 0x15, 0x0a, 0x03, 0x66, 0x70, 0x73, 0x18, 0x12,
	0x20, 0x01, 0x28, 0x09, 0x48, 0x0e, 0x52, 0x03, 0x66, 0x70, 0x73, 0x88, 0x01, 0x01, 0x12, 0x1f,
	0x0a, 0x08, 0x6c, 0x61, 0x74, 0x69, 0x74, 0x75, 0x64, 0x65, 0x18, 0x13, 0x20, 0x01, 0x28, 0x01,
	0x48, 0x0f, 0x52, 0x08, 0x6c, 0x61, 0x74, 0x69, 0x74, 0x75, 0x64, 0x65, 0x88, 0x01, 0x01, 0x12,
	0x21, 0x0a, 0x09, 0x6c, 0x6f, 0x6e, 0x67, 0x69, 0x74, 0x75, 0x64, 0x65, 0x18, 0x14, 0x20, 0x01,
	0x28, 0x01, 0x48, 0x10, 0x52, 0x09, 0x6c, 0x6f, 0x6e, 0x67, 0x69, 0x74, 0x75, 0x64, 0x65, 0x88,
	0x01, 0x01, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x6d, 0x69, 0x6d, 0x65, 0x54, 0x79, 0x70, 0x65, 0x42,
	0x0c, 0x0a, 0x0a, 0x5f, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x55, 0x72, 0x6c, 0x42, 0x0d, 0x0a,
	0x0b, 0x5f, 0x70, 0x72, 0x65, 0x76, 0x69, 0x65, 0x77, 0x55, 0x72, 0x6c, 0x42, 0x0f, 0x0a, 0x0d,
	0x5f, 0x74, 0x68, 0x75, 0x6d, 0x62, 0x6e, 0x61, 0x69, 0x6c, 0x55, 0x72, 0x6c, 0x42, 0x07, 0x0a,
	0x05, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x77, 0x69, 0x64, 0x74, 0x68,
	0x42, 0x09, 0x0a, 0x07, 0x5f, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x42, 0x0f, 0x0a, 0x0d, 0x5f,
	0x63, 0x72, 0x65, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x42, 0x0d, 0x0a, 0x0b,
	0x5f, 0x63, 0x61, 0x6d, 0x65, 0x72, 0x61, 0x4d, 0x61, 0x6b, 0x65, 0x42, 0x0e, 0x0a, 0x0c, 0x5f,
	0x63, 0x61, 0x6d, 0x65, 0x72, 0x61, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x42, 0x0e, 0x0a, 0x0c, 0x5f,
	0x66, 0x6f, 0x63, 0x61, 0x6c, 0x4c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x42, 0x12, 0x0a, 0x10, 0x5f,
	0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x46, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x42,
	0x10, 0x0a, 0x0e, 0x5f, 0x69, 0x73, 0x6f, 0x45, 0x71, 0x75, 0x69, 0x76, 0x61, 0x6c, 0x65, 0x6e,
	0x74, 0x42, 0x0f, 0x0a, 0x0d, 0x5f, 0x65, 0x78, 0x70, 0x6f, 0x73, 0x75, 0x72, 0x65, 0x54, 0x69,
	0x6d, 0x65, 0x42, 0x06, 0x0a, 0x04, 0x5f, 0x66, 0x70, 0x73, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x6c,
	0x61, 0x74, 0x69, 0x74, 0x75, 0x64, 0x65, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x6c, 0x6f, 0x6e, 0x67,
	0x69, 0x74, 0x75, 0x64, 0x65, 0x22, 0x81, 0x02, 0x0a, 0x15, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x49,
	0x74, 0x65, 0x6d, 0x50, 0x6c, 0x61, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1f, 0x0a, 0x08, 0x70, 0x6f, 0x73, 0x74, 0x63,
	0x6f, 0x64, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x08, 0x70, 0x6f, 0x73,
	0x74, 0x63, 0x6f, 0x64, 0x65, 0x88, 0x01, 0x01, 0x12, 0x1d, 0x0a, 0x07, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x72, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x48, 0x01, 0x52, 0x07, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x72, 0x79, 0x88, 0x01, 0x01, 0x12, 0x19, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x48, 0x02, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x88,
	0x01, 0x01, 0x12, 0x17, 0x0a, 0x04, 0x63, 0x69, 0x74, 0x79, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09,
	0x48, 0x03, 0x52, 0x04, 0x63, 0x69, 0x74, 0x79, 0x88, 0x01, 0x01, 0x12, 0x17, 0x0a, 0x04, 0x74,
	0x6f, 0x77, 0x6e, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x48, 0x04, 0x52, 0x04, 0x74, 0x6f, 0x77,
	0x6e, 0x88, 0x01, 0x01, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x70, 0x6f, 0x73, 0x74, 0x63, 0x6f, 0x64,
	0x65, 0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x79, 0x42, 0x08, 0x0a,
	0x06, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x63, 0x69, 0x74, 0x79,
	0x42, 0x07, 0x0a, 0x05, 0x5f, 0x74, 0x6f, 0x77, 0x6e, 0x32, 0x9b, 0x01, 0x0a, 0x03, 0x41, 0x50,
	0x49, 0x12, 0x4c, 0x0a, 0x15, 0x53, 0x61, 0x76, 0x65, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x49, 0x74,
	0x65, 0x6d, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x19, 0x2e, 0x4d, 0x65, 0x64,
	0x69, 0x61, 0x49, 0x74, 0x65, 0x6d, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12,
	0x46, 0x0a, 0x12, 0x53, 0x61, 0x76, 0x65, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x49, 0x74, 0x65, 0x6d,
	0x50, 0x6c, 0x61, 0x63, 0x65, 0x12, 0x16, 0x2e, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x49, 0x74, 0x65,
	0x6d, 0x50, 0x6c, 0x61, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x42, 0x06, 0x5a, 0x04, 0x61, 0x70, 0x69, 0x2f, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
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

var file_api_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_api_proto_goTypes = []interface{}{
	(*MediaItemMetadataRequest)(nil), // 0: MediaItemMetadataRequest
	(*MediaItemPlaceRequest)(nil),    // 1: MediaItemPlaceRequest
	(*empty.Empty)(nil),              // 2: google.protobuf.Empty
}
var file_api_proto_depIdxs = []int32{
	0, // 0: API.SaveMediaItemMetadata:input_type -> MediaItemMetadataRequest
	1, // 1: API.SaveMediaItemPlace:input_type -> MediaItemPlaceRequest
	2, // 2: API.SaveMediaItemMetadata:output_type -> google.protobuf.Empty
	2, // 3: API.SaveMediaItemPlace:output_type -> google.protobuf.Empty
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_api_proto_init() }
func file_api_proto_init() {
	if File_api_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
		file_api_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
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
	}
	file_api_proto_msgTypes[0].OneofWrappers = []interface{}{}
	file_api_proto_msgTypes[1].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
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
