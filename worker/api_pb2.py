# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: api.proto
"""Generated protocol buffer code."""
from google.protobuf.internal import builder as _builder
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import empty_pb2 as google_dot_protobuf_dot_empty__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\tapi.proto\x1a\x1bgoogle/protobuf/empty.proto\"\xe8\x05\n\x16MediaItemResultRequest\x12\n\n\x02id\x18\x01 \x01(\t\x12\x0e\n\x06status\x18\x02 \x01(\t\x12\x15\n\x08\x66ilename\x18\x03 \x01(\tH\x00\x88\x01\x01\x12\x15\n\x08mimeType\x18\x04 \x01(\tH\x01\x88\x01\x01\x12\x16\n\tsourceUrl\x18\x05 \x01(\tH\x02\x88\x01\x01\x12\x17\n\npreviewUrl\x18\x06 \x01(\tH\x03\x88\x01\x01\x12\x19\n\x0cthumbnailUrl\x18\x07 \x01(\tH\x04\x88\x01\x01\x12\x11\n\x04type\x18\x08 \x01(\tH\x05\x88\x01\x01\x12\x12\n\x05width\x18\t \x01(\x05H\x06\x88\x01\x01\x12\x13\n\x06height\x18\n \x01(\x05H\x07\x88\x01\x01\x12\x19\n\x0c\x63reationTime\x18\x0b \x01(\tH\x08\x88\x01\x01\x12\x17\n\ncameraMake\x18\x0c \x01(\tH\t\x88\x01\x01\x12\x18\n\x0b\x63\x61meraModel\x18\r \x01(\tH\n\x88\x01\x01\x12\x18\n\x0b\x66ocalLength\x18\x0e \x01(\tH\x0b\x88\x01\x01\x12\x1c\n\x0f\x61pertureFNumber\x18\x0f \x01(\tH\x0c\x88\x01\x01\x12\x1a\n\risoEquivalent\x18\x10 \x01(\tH\r\x88\x01\x01\x12\x19\n\x0c\x65xposureTime\x18\x11 \x01(\tH\x0e\x88\x01\x01\x12\x10\n\x03\x66ps\x18\x12 \x01(\tH\x0f\x88\x01\x01\x12\x15\n\x08latitude\x18\x13 \x01(\x01H\x10\x88\x01\x01\x12\x16\n\tlongitude\x18\x14 \x01(\x01H\x11\x88\x01\x01\x42\x0b\n\t_filenameB\x0b\n\t_mimeTypeB\x0c\n\n_sourceUrlB\r\n\x0b_previewUrlB\x0f\n\r_thumbnailUrlB\x07\n\x05_typeB\x08\n\x06_widthB\t\n\x07_heightB\x0f\n\r_creationTimeB\r\n\x0b_cameraMakeB\x0e\n\x0c_cameraModelB\x0e\n\x0c_focalLengthB\x12\n\x10_apertureFNumberB\x10\n\x0e_isoEquivalentB\x0f\n\r_exposureTimeB\x06\n\x04_fpsB\x0b\n\t_latitudeB\x0c\n\n_longitude\"\xbf\x01\n\x15MediaItemPlaceRequest\x12\n\n\x02id\x18\x01 \x01(\t\x12\x15\n\x08postcode\x18\x02 \x01(\tH\x00\x88\x01\x01\x12\x14\n\x07\x63ountry\x18\x03 \x01(\tH\x01\x88\x01\x01\x12\x12\n\x05state\x18\x04 \x01(\tH\x02\x88\x01\x01\x12\x11\n\x04\x63ity\x18\x05 \x01(\tH\x03\x88\x01\x01\x12\x11\n\x04town\x18\x06 \x01(\tH\x04\x88\x01\x01\x42\x0b\n\t_postcodeB\n\n\x08_countryB\x08\n\x06_stateB\x07\n\x05_cityB\x07\n\x05_town2\x93\x01\n\x03\x41PI\x12\x46\n\x13SaveMediaItemResult\x12\x17.MediaItemResultRequest\x1a\x16.google.protobuf.Empty\x12\x44\n\x12SaveMediaItemPlace\x12\x16.MediaItemPlaceRequest\x1a\x16.google.protobuf.EmptyB\x06Z\x04\x61pi/b\x06proto3')

_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, globals())
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'api_pb2', globals())
if _descriptor._USE_C_DESCRIPTORS == False:

  DESCRIPTOR._options = None
  DESCRIPTOR._serialized_options = b'Z\004api/'
  _MEDIAITEMRESULTREQUEST._serialized_start=43
  _MEDIAITEMRESULTREQUEST._serialized_end=787
  _MEDIAITEMPLACEREQUEST._serialized_start=790
  _MEDIAITEMPLACEREQUEST._serialized_end=981
  _API._serialized_start=984
  _API._serialized_end=1131
# @@protoc_insertion_point(module_scope)
