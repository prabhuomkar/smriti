# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: src/protos/api.proto
"""Generated protocol buffer code."""
from google.protobuf.internal import builder as _builder
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import empty_pb2 as google_dot_protobuf_dot_empty__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x14src/protos/api.proto\x1a\x1bgoogle/protobuf/empty.proto\" \n\x0e\x43onfigResponse\x12\x0e\n\x06\x63onfig\x18\x01 \x01(\x0c\"\xcc\x05\n\x18MediaItemMetadataRequest\x12\x0e\n\x06userId\x18\x01 \x01(\t\x12\n\n\x02id\x18\x02 \x01(\t\x12\x0e\n\x06status\x18\x03 \x01(\t\x12\x15\n\x08mimeType\x18\x04 \x01(\tH\x00\x88\x01\x01\x12\x12\n\nsourcePath\x18\x05 \x01(\t\x12\x18\n\x0bpreviewPath\x18\x06 \x01(\tH\x01\x88\x01\x01\x12\x1a\n\rthumbnailPath\x18\x07 \x01(\tH\x02\x88\x01\x01\x12\x0c\n\x04type\x18\x08 \x01(\t\x12\x10\n\x08\x63\x61tegory\x18\t \x01(\t\x12\x12\n\x05width\x18\n \x01(\x05H\x03\x88\x01\x01\x12\x13\n\x06height\x18\x0b \x01(\x05H\x04\x88\x01\x01\x12\x19\n\x0c\x63reationTime\x18\x0c \x01(\tH\x05\x88\x01\x01\x12\x17\n\ncameraMake\x18\r \x01(\tH\x06\x88\x01\x01\x12\x18\n\x0b\x63\x61meraModel\x18\x0e \x01(\tH\x07\x88\x01\x01\x12\x18\n\x0b\x66ocalLength\x18\x0f \x01(\tH\x08\x88\x01\x01\x12\x1c\n\x0f\x61pertureFNumber\x18\x10 \x01(\tH\t\x88\x01\x01\x12\x1a\n\risoEquivalent\x18\x11 \x01(\tH\n\x88\x01\x01\x12\x19\n\x0c\x65xposureTime\x18\x12 \x01(\tH\x0b\x88\x01\x01\x12\x10\n\x03\x66ps\x18\x13 \x01(\tH\x0c\x88\x01\x01\x12\x15\n\x08latitude\x18\x14 \x01(\x01H\r\x88\x01\x01\x12\x16\n\tlongitude\x18\x15 \x01(\x01H\x0e\x88\x01\x01\x42\x0b\n\t_mimeTypeB\x0e\n\x0c_previewPathB\x10\n\x0e_thumbnailPathB\x08\n\x06_widthB\t\n\x07_heightB\x0f\n\r_creationTimeB\r\n\x0b_cameraMakeB\x0e\n\x0c_cameraModelB\x0e\n\x0c_focalLengthB\x12\n\x10_apertureFNumberB\x10\n\x0e_isoEquivalentB\x0f\n\r_exposureTimeB\x06\n\x04_fpsB\x0b\n\t_latitudeB\x0c\n\n_longitude\"\xcf\x01\n\x15MediaItemPlaceRequest\x12\x0e\n\x06userId\x18\x01 \x01(\t\x12\n\n\x02id\x18\x02 \x01(\t\x12\x15\n\x08postcode\x18\x03 \x01(\tH\x00\x88\x01\x01\x12\x14\n\x07\x63ountry\x18\x04 \x01(\tH\x01\x88\x01\x01\x12\x12\n\x05state\x18\x05 \x01(\tH\x02\x88\x01\x01\x12\x11\n\x04\x63ity\x18\x06 \x01(\tH\x03\x88\x01\x01\x12\x11\n\x04town\x18\x07 \x01(\tH\x04\x88\x01\x01\x42\x0b\n\t_postcodeB\n\n\x08_countryB\x08\n\x06_stateB\x07\n\x05_cityB\x07\n\x05_town\"A\n\x15MediaItemThingRequest\x12\x0e\n\x06userId\x18\x01 \x01(\t\x12\n\n\x02id\x18\x02 \x01(\t\x12\x0c\n\x04name\x18\x03 \x01(\t\"S\n\x18MediaItemMLResultRequest\x12\x0e\n\x06userId\x18\x01 \x01(\t\x12\n\n\x02id\x18\x02 \x01(\t\x12\x0c\n\x04name\x18\x03 \x01(\t\x12\r\n\x05value\x18\x04 \x03(\t\"\\\n\x19\x46inalSaveMediaItemRequest\x12\x0e\n\x06userId\x18\x01 \x01(\t\x12\n\n\x02id\x18\x02 \x01(\t\x12\x10\n\x08keywords\x18\x03 \x01(\t\x12\x11\n\tembedding\x18\x04 \x03(\x02\x32\xbb\x03\n\x03\x41PI\x12<\n\x0fGetWorkerConfig\x12\x16.google.protobuf.Empty\x1a\x0f.ConfigResponse\"\x00\x12L\n\x15SaveMediaItemMetadata\x12\x19.MediaItemMetadataRequest\x1a\x16.google.protobuf.Empty\"\x00\x12\x46\n\x12SaveMediaItemPlace\x12\x16.MediaItemPlaceRequest\x1a\x16.google.protobuf.Empty\"\x00\x12\x46\n\x12SaveMediaItemThing\x12\x16.MediaItemThingRequest\x1a\x16.google.protobuf.Empty\"\x00\x12L\n\x15SaveMediaItemMLResult\x12\x19.MediaItemMLResultRequest\x1a\x16.google.protobuf.Empty\"\x00\x12J\n\x12\x46inalSaveMediaItem\x12\x1a.FinalSaveMediaItemRequest\x1a\x16.google.protobuf.Empty\"\x00\x42\x06Z\x04\x61pi/b\x06proto3')

_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, globals())
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'src.protos.api_pb2', globals())
if _descriptor._USE_C_DESCRIPTORS == False:

  DESCRIPTOR._options = None
  DESCRIPTOR._serialized_options = b'Z\004api/'
  _CONFIGRESPONSE._serialized_start=53
  _CONFIGRESPONSE._serialized_end=85
  _MEDIAITEMMETADATAREQUEST._serialized_start=88
  _MEDIAITEMMETADATAREQUEST._serialized_end=804
  _MEDIAITEMPLACEREQUEST._serialized_start=807
  _MEDIAITEMPLACEREQUEST._serialized_end=1014
  _MEDIAITEMTHINGREQUEST._serialized_start=1016
  _MEDIAITEMTHINGREQUEST._serialized_end=1081
  _MEDIAITEMMLRESULTREQUEST._serialized_start=1083
  _MEDIAITEMMLRESULTREQUEST._serialized_end=1166
  _FINALSAVEMEDIAITEMREQUEST._serialized_start=1168
  _FINALSAVEMEDIAITEMREQUEST._serialized_end=1260
  _API._serialized_start=1263
  _API._serialized_end=1706
# @@protoc_insertion_point(module_scope)
