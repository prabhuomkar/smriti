# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: src/protos/api.proto
# Protobuf Python Version: 5.29.0
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import runtime_version as _runtime_version
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
_runtime_version.ValidateProtobufRuntimeVersion(
    _runtime_version.Domain.PUBLIC,
    5,
    29,
    0,
    '',
    'src/protos/api.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import empty_pb2 as google_dot_protobuf_dot_empty__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x14src/protos/api.proto\x1a\x1bgoogle/protobuf/empty.proto\" \n\x0e\x43onfigResponse\x12\x0e\n\x06\x63onfig\x18\x01 \x01(\x0c\"\x84\x05\n\x18MediaItemMetadataRequest\x12\x0e\n\x06userId\x18\x01 \x01(\t\x12\n\n\x02id\x18\x02 \x01(\t\x12\x0e\n\x06status\x18\x03 \x01(\t\x12\x15\n\x08mimeType\x18\x04 \x01(\tH\x00\x88\x01\x01\x12\x0c\n\x04type\x18\x05 \x01(\t\x12\x10\n\x08\x63\x61tegory\x18\x06 \x01(\t\x12\x12\n\x05width\x18\x07 \x01(\x05H\x01\x88\x01\x01\x12\x13\n\x06height\x18\x08 \x01(\x05H\x02\x88\x01\x01\x12\x19\n\x0c\x63reationTime\x18\t \x01(\tH\x03\x88\x01\x01\x12\x17\n\ncameraMake\x18\n \x01(\tH\x04\x88\x01\x01\x12\x18\n\x0b\x63\x61meraModel\x18\x0b \x01(\tH\x05\x88\x01\x01\x12\x18\n\x0b\x66ocalLength\x18\x0c \x01(\tH\x06\x88\x01\x01\x12\x1c\n\x0f\x61pertureFNumber\x18\r \x01(\tH\x07\x88\x01\x01\x12\x1a\n\risoEquivalent\x18\x0e \x01(\tH\x08\x88\x01\x01\x12\x19\n\x0c\x65xposureTime\x18\x0f \x01(\tH\t\x88\x01\x01\x12\x10\n\x03\x66ps\x18\x10 \x01(\tH\n\x88\x01\x01\x12\x15\n\x08latitude\x18\x11 \x01(\x01H\x0b\x88\x01\x01\x12\x16\n\tlongitude\x18\x12 \x01(\x01H\x0c\x88\x01\x01\x12\x15\n\x08\x65xifData\x18\x13 \x01(\tH\r\x88\x01\x01\x42\x0b\n\t_mimeTypeB\x08\n\x06_widthB\t\n\x07_heightB\x0f\n\r_creationTimeB\r\n\x0b_cameraMakeB\x0e\n\x0c_cameraModelB\x0e\n\x0c_focalLengthB\x12\n\x10_apertureFNumberB\x10\n\x0e_isoEquivalentB\x0f\n\r_exposureTimeB\x06\n\x04_fpsB\x0b\n\t_latitudeB\x0c\n\n_longitudeB\x0b\n\t_exifData\"\xf8\x01\n MediaItemPreviewThumbnailRequest\x12\x0e\n\x06userId\x18\x01 \x01(\t\x12\n\n\x02id\x18\x02 \x01(\t\x12\x0e\n\x06status\x18\x03 \x01(\t\x12\x17\n\nsourcePath\x18\x04 \x01(\tH\x00\x88\x01\x01\x12\x18\n\x0bpreviewPath\x18\x05 \x01(\tH\x01\x88\x01\x01\x12\x1a\n\rthumbnailPath\x18\x06 \x01(\tH\x02\x88\x01\x01\x12\x18\n\x0bplaceholder\x18\x07 \x01(\tH\x03\x88\x01\x01\x42\r\n\x0b_sourcePathB\x0e\n\x0c_previewPathB\x10\n\x0e_thumbnailPathB\x0e\n\x0c_placeholder\"\xcf\x01\n\x15MediaItemPlaceRequest\x12\x0e\n\x06userId\x18\x01 \x01(\t\x12\n\n\x02id\x18\x02 \x01(\t\x12\x15\n\x08postcode\x18\x03 \x01(\tH\x00\x88\x01\x01\x12\x14\n\x07\x63ountry\x18\x04 \x01(\tH\x01\x88\x01\x01\x12\x12\n\x05state\x18\x05 \x01(\tH\x02\x88\x01\x01\x12\x11\n\x04\x63ity\x18\x06 \x01(\tH\x03\x88\x01\x01\x12\x11\n\x04town\x18\x07 \x01(\tH\x04\x88\x01\x01\x42\x0b\n\t_postcodeB\n\n\x08_countryB\x08\n\x06_stateB\x07\n\x05_cityB\x07\n\x05_town\"A\n\x15MediaItemThingRequest\x12\x0e\n\x06userId\x18\x01 \x01(\t\x12\n\n\x02id\x18\x02 \x01(\t\x12\x0c\n\x04name\x18\x03 \x01(\t\"\'\n\x12MediaItemEmbedding\x12\x11\n\tembedding\x18\x01 \x03(\x02\"p\n\x15MediaItemFacesRequest\x12\x0e\n\x06userId\x18\x01 \x01(\t\x12\n\n\x02id\x18\x02 \x01(\t\x12\'\n\nembeddings\x18\x03 \x03(\x0b\x32\x13.MediaItemEmbedding\x12\x12\n\nthumbnails\x18\x04 \x03(\t\"t\n\x1bMediaItemFinalResultRequest\x12\x0e\n\x06userId\x18\x01 \x01(\t\x12\n\n\x02id\x18\x02 \x01(\t\x12\x10\n\x08keywords\x18\x03 \x01(\t\x12\'\n\nembeddings\x18\x04 \x03(\x0b\x32\x13.MediaItemEmbedding\"0\n\x1eMediaItemFaceEmbeddingsRequest\x12\x0e\n\x06userId\x18\x01 \x01(\t\"s\n\x16MediaItemFaceEmbedding\x12\n\n\x02id\x18\x01 \x01(\t\x12\x13\n\x0bmediaItemId\x18\x02 \x01(\t\x12\x10\n\x08peopleId\x18\x03 \x01(\t\x12&\n\tembedding\x18\x04 \x01(\x0b\x32\x13.MediaItemEmbedding\"[\n\x1fMediaItemFaceEmbeddingsResponse\x12\x38\n\x17mediaItemFaceEmbeddings\x18\x01 \x03(\x0b\x32\x17.MediaItemFaceEmbedding\"!\n\x10GetUsersResponse\x12\r\n\x05users\x18\x01 \x03(\t\"\x82\x01\n\x13MediaItemFacePeople\x12\x38\n\nfacePeople\x18\x01 \x03(\x0b\x32$.MediaItemFacePeople.FacePeopleEntry\x1a\x31\n\x0f\x46\x61\x63\x65PeopleEntry\x12\x0b\n\x03key\x18\x01 \x01(\t\x12\r\n\x05value\x18\x02 \x01(\t:\x02\x38\x01\"\xc9\x01\n\x16MediaItemPeopleRequest\x12\x0e\n\x06userId\x18\x01 \x01(\t\x12M\n\x13mediaItemFacePeople\x18\x02 \x03(\x0b\x32\x30.MediaItemPeopleRequest.MediaItemFacePeopleEntry\x1aP\n\x18MediaItemFacePeopleEntry\x12\x0b\n\x03key\x18\x01 \x01(\t\x12#\n\x05value\x18\x02 \x01(\x0b\x32\x14.MediaItemFacePeople:\x02\x38\x01\x32\x81\x06\n\x03\x41PI\x12<\n\x0fGetWorkerConfig\x12\x16.google.protobuf.Empty\x1a\x0f.ConfigResponse\"\x00\x12\x61\n\x1aGetMediaItemFaceEmbeddings\x12\x1f.MediaItemFaceEmbeddingsRequest\x1a .MediaItemFaceEmbeddingsResponse\"\x00\x12\x37\n\x08GetUsers\x12\x16.google.protobuf.Empty\x1a\x11.GetUsersResponse\"\x00\x12L\n\x15SaveMediaItemMetadata\x12\x19.MediaItemMetadataRequest\x1a\x16.google.protobuf.Empty\"\x00\x12\\\n\x1dSaveMediaItemPreviewThumbnail\x12!.MediaItemPreviewThumbnailRequest\x1a\x16.google.protobuf.Empty\"\x00\x12\x46\n\x12SaveMediaItemPlace\x12\x16.MediaItemPlaceRequest\x1a\x16.google.protobuf.Empty\"\x00\x12\x46\n\x12SaveMediaItemThing\x12\x16.MediaItemThingRequest\x1a\x16.google.protobuf.Empty\"\x00\x12\x46\n\x12SaveMediaItemFaces\x12\x16.MediaItemFacesRequest\x1a\x16.google.protobuf.Empty\"\x00\x12H\n\x13SaveMediaItemPeople\x12\x17.MediaItemPeopleRequest\x1a\x16.google.protobuf.Empty\"\x00\x12R\n\x18SaveMediaItemFinalResult\x12\x1c.MediaItemFinalResultRequest\x1a\x16.google.protobuf.Empty\"\x00\x42\x06Z\x04\x61pi/b\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'src.protos.api_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z\004api/'
  _globals['_MEDIAITEMFACEPEOPLE_FACEPEOPLEENTRY']._loaded_options = None
  _globals['_MEDIAITEMFACEPEOPLE_FACEPEOPLEENTRY']._serialized_options = b'8\001'
  _globals['_MEDIAITEMPEOPLEREQUEST_MEDIAITEMFACEPEOPLEENTRY']._loaded_options = None
  _globals['_MEDIAITEMPEOPLEREQUEST_MEDIAITEMFACEPEOPLEENTRY']._serialized_options = b'8\001'
  _globals['_CONFIGRESPONSE']._serialized_start=53
  _globals['_CONFIGRESPONSE']._serialized_end=85
  _globals['_MEDIAITEMMETADATAREQUEST']._serialized_start=88
  _globals['_MEDIAITEMMETADATAREQUEST']._serialized_end=732
  _globals['_MEDIAITEMPREVIEWTHUMBNAILREQUEST']._serialized_start=735
  _globals['_MEDIAITEMPREVIEWTHUMBNAILREQUEST']._serialized_end=983
  _globals['_MEDIAITEMPLACEREQUEST']._serialized_start=986
  _globals['_MEDIAITEMPLACEREQUEST']._serialized_end=1193
  _globals['_MEDIAITEMTHINGREQUEST']._serialized_start=1195
  _globals['_MEDIAITEMTHINGREQUEST']._serialized_end=1260
  _globals['_MEDIAITEMEMBEDDING']._serialized_start=1262
  _globals['_MEDIAITEMEMBEDDING']._serialized_end=1301
  _globals['_MEDIAITEMFACESREQUEST']._serialized_start=1303
  _globals['_MEDIAITEMFACESREQUEST']._serialized_end=1415
  _globals['_MEDIAITEMFINALRESULTREQUEST']._serialized_start=1417
  _globals['_MEDIAITEMFINALRESULTREQUEST']._serialized_end=1533
  _globals['_MEDIAITEMFACEEMBEDDINGSREQUEST']._serialized_start=1535
  _globals['_MEDIAITEMFACEEMBEDDINGSREQUEST']._serialized_end=1583
  _globals['_MEDIAITEMFACEEMBEDDING']._serialized_start=1585
  _globals['_MEDIAITEMFACEEMBEDDING']._serialized_end=1700
  _globals['_MEDIAITEMFACEEMBEDDINGSRESPONSE']._serialized_start=1702
  _globals['_MEDIAITEMFACEEMBEDDINGSRESPONSE']._serialized_end=1793
  _globals['_GETUSERSRESPONSE']._serialized_start=1795
  _globals['_GETUSERSRESPONSE']._serialized_end=1828
  _globals['_MEDIAITEMFACEPEOPLE']._serialized_start=1831
  _globals['_MEDIAITEMFACEPEOPLE']._serialized_end=1961
  _globals['_MEDIAITEMFACEPEOPLE_FACEPEOPLEENTRY']._serialized_start=1912
  _globals['_MEDIAITEMFACEPEOPLE_FACEPEOPLEENTRY']._serialized_end=1961
  _globals['_MEDIAITEMPEOPLEREQUEST']._serialized_start=1964
  _globals['_MEDIAITEMPEOPLEREQUEST']._serialized_end=2165
  _globals['_MEDIAITEMPEOPLEREQUEST_MEDIAITEMFACEPEOPLEENTRY']._serialized_start=2085
  _globals['_MEDIAITEMPEOPLEREQUEST_MEDIAITEMFACEPEOPLEENTRY']._serialized_end=2165
  _globals['_API']._serialized_start=2168
  _globals['_API']._serialized_end=2937
# @@protoc_insertion_point(module_scope)
