# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: worker.proto
"""Generated protocol buffer code."""
from google.protobuf.internal import builder as _builder
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import empty_pb2 as google_dot_protobuf_dot_empty__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x0cworker.proto\x1a\x1bgoogle/protobuf/empty.proto\"W\n\x17MediaItemProcessRequest\x12\n\n\x02id\x18\x01 \x01(\t\x12\x0e\n\x06offset\x18\x02 \x01(\x03\x12\x0f\n\x07\x63ommand\x18\x03 \x01(\t\x12\x0f\n\x07\x63ontent\x18\x04 \x01(\x0c\x32P\n\x06Worker\x12\x46\n\x10MediaItemProcess\x12\x18.MediaItemProcessRequest\x1a\x16.google.protobuf.Empty(\x01\x42\tZ\x07worker/b\x06proto3')

_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, globals())
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'worker_pb2', globals())
if _descriptor._USE_C_DESCRIPTORS == False:

  DESCRIPTOR._options = None
  DESCRIPTOR._serialized_options = b'Z\007worker/'
  _MEDIAITEMPROCESSREQUEST._serialized_start=45
  _MEDIAITEMPROCESSREQUEST._serialized_end=132
  _WORKER._serialized_start=134
  _WORKER._serialized_end=214
# @@protoc_insertion_point(module_scope)
