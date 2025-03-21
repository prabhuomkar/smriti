# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc
import warnings

from google.protobuf import empty_pb2 as google_dot_protobuf_dot_empty__pb2
from src.protos import api_pb2 as src_dot_protos_dot_api__pb2

GRPC_GENERATED_VERSION = '1.70.0'
GRPC_VERSION = grpc.__version__
_version_not_supported = False

try:
    from grpc._utilities import first_version_is_lower
    _version_not_supported = first_version_is_lower(GRPC_VERSION, GRPC_GENERATED_VERSION)
except ImportError:
    _version_not_supported = True

if _version_not_supported:
    raise RuntimeError(
        f'The grpc package installed is at version {GRPC_VERSION},'
        + f' but the generated code in src/protos/api_pb2_grpc.py depends on'
        + f' grpcio>={GRPC_GENERATED_VERSION}.'
        + f' Please upgrade your grpc module to grpcio>={GRPC_GENERATED_VERSION}'
        + f' or downgrade your generated code using grpcio-tools<={GRPC_VERSION}.'
    )


class APIStub(object):
    """Missing associated documentation comment in .proto file."""

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.GetWorkerConfig = channel.unary_unary(
                '/API/GetWorkerConfig',
                request_serializer=google_dot_protobuf_dot_empty__pb2.Empty.SerializeToString,
                response_deserializer=src_dot_protos_dot_api__pb2.ConfigResponse.FromString,
                _registered_method=True)
        self.GetMediaItemFaceEmbeddings = channel.unary_unary(
                '/API/GetMediaItemFaceEmbeddings',
                request_serializer=src_dot_protos_dot_api__pb2.MediaItemFaceEmbeddingsRequest.SerializeToString,
                response_deserializer=src_dot_protos_dot_api__pb2.MediaItemFaceEmbeddingsResponse.FromString,
                _registered_method=True)
        self.GetUsers = channel.unary_unary(
                '/API/GetUsers',
                request_serializer=google_dot_protobuf_dot_empty__pb2.Empty.SerializeToString,
                response_deserializer=src_dot_protos_dot_api__pb2.GetUsersResponse.FromString,
                _registered_method=True)
        self.SaveMediaItemMetadata = channel.unary_unary(
                '/API/SaveMediaItemMetadata',
                request_serializer=src_dot_protos_dot_api__pb2.MediaItemMetadataRequest.SerializeToString,
                response_deserializer=google_dot_protobuf_dot_empty__pb2.Empty.FromString,
                _registered_method=True)
        self.SaveMediaItemPreviewThumbnail = channel.unary_unary(
                '/API/SaveMediaItemPreviewThumbnail',
                request_serializer=src_dot_protos_dot_api__pb2.MediaItemPreviewThumbnailRequest.SerializeToString,
                response_deserializer=google_dot_protobuf_dot_empty__pb2.Empty.FromString,
                _registered_method=True)
        self.SaveMediaItemPlace = channel.unary_unary(
                '/API/SaveMediaItemPlace',
                request_serializer=src_dot_protos_dot_api__pb2.MediaItemPlaceRequest.SerializeToString,
                response_deserializer=google_dot_protobuf_dot_empty__pb2.Empty.FromString,
                _registered_method=True)
        self.SaveMediaItemThing = channel.unary_unary(
                '/API/SaveMediaItemThing',
                request_serializer=src_dot_protos_dot_api__pb2.MediaItemThingRequest.SerializeToString,
                response_deserializer=google_dot_protobuf_dot_empty__pb2.Empty.FromString,
                _registered_method=True)
        self.SaveMediaItemFaces = channel.unary_unary(
                '/API/SaveMediaItemFaces',
                request_serializer=src_dot_protos_dot_api__pb2.MediaItemFacesRequest.SerializeToString,
                response_deserializer=google_dot_protobuf_dot_empty__pb2.Empty.FromString,
                _registered_method=True)
        self.SaveMediaItemPeople = channel.unary_unary(
                '/API/SaveMediaItemPeople',
                request_serializer=src_dot_protos_dot_api__pb2.MediaItemPeopleRequest.SerializeToString,
                response_deserializer=google_dot_protobuf_dot_empty__pb2.Empty.FromString,
                _registered_method=True)
        self.SaveMediaItemFinalResult = channel.unary_unary(
                '/API/SaveMediaItemFinalResult',
                request_serializer=src_dot_protos_dot_api__pb2.MediaItemFinalResultRequest.SerializeToString,
                response_deserializer=google_dot_protobuf_dot_empty__pb2.Empty.FromString,
                _registered_method=True)


class APIServicer(object):
    """Missing associated documentation comment in .proto file."""

    def GetWorkerConfig(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def GetMediaItemFaceEmbeddings(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def GetUsers(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def SaveMediaItemMetadata(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def SaveMediaItemPreviewThumbnail(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def SaveMediaItemPlace(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def SaveMediaItemThing(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def SaveMediaItemFaces(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def SaveMediaItemPeople(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def SaveMediaItemFinalResult(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_APIServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'GetWorkerConfig': grpc.unary_unary_rpc_method_handler(
                    servicer.GetWorkerConfig,
                    request_deserializer=google_dot_protobuf_dot_empty__pb2.Empty.FromString,
                    response_serializer=src_dot_protos_dot_api__pb2.ConfigResponse.SerializeToString,
            ),
            'GetMediaItemFaceEmbeddings': grpc.unary_unary_rpc_method_handler(
                    servicer.GetMediaItemFaceEmbeddings,
                    request_deserializer=src_dot_protos_dot_api__pb2.MediaItemFaceEmbeddingsRequest.FromString,
                    response_serializer=src_dot_protos_dot_api__pb2.MediaItemFaceEmbeddingsResponse.SerializeToString,
            ),
            'GetUsers': grpc.unary_unary_rpc_method_handler(
                    servicer.GetUsers,
                    request_deserializer=google_dot_protobuf_dot_empty__pb2.Empty.FromString,
                    response_serializer=src_dot_protos_dot_api__pb2.GetUsersResponse.SerializeToString,
            ),
            'SaveMediaItemMetadata': grpc.unary_unary_rpc_method_handler(
                    servicer.SaveMediaItemMetadata,
                    request_deserializer=src_dot_protos_dot_api__pb2.MediaItemMetadataRequest.FromString,
                    response_serializer=google_dot_protobuf_dot_empty__pb2.Empty.SerializeToString,
            ),
            'SaveMediaItemPreviewThumbnail': grpc.unary_unary_rpc_method_handler(
                    servicer.SaveMediaItemPreviewThumbnail,
                    request_deserializer=src_dot_protos_dot_api__pb2.MediaItemPreviewThumbnailRequest.FromString,
                    response_serializer=google_dot_protobuf_dot_empty__pb2.Empty.SerializeToString,
            ),
            'SaveMediaItemPlace': grpc.unary_unary_rpc_method_handler(
                    servicer.SaveMediaItemPlace,
                    request_deserializer=src_dot_protos_dot_api__pb2.MediaItemPlaceRequest.FromString,
                    response_serializer=google_dot_protobuf_dot_empty__pb2.Empty.SerializeToString,
            ),
            'SaveMediaItemThing': grpc.unary_unary_rpc_method_handler(
                    servicer.SaveMediaItemThing,
                    request_deserializer=src_dot_protos_dot_api__pb2.MediaItemThingRequest.FromString,
                    response_serializer=google_dot_protobuf_dot_empty__pb2.Empty.SerializeToString,
            ),
            'SaveMediaItemFaces': grpc.unary_unary_rpc_method_handler(
                    servicer.SaveMediaItemFaces,
                    request_deserializer=src_dot_protos_dot_api__pb2.MediaItemFacesRequest.FromString,
                    response_serializer=google_dot_protobuf_dot_empty__pb2.Empty.SerializeToString,
            ),
            'SaveMediaItemPeople': grpc.unary_unary_rpc_method_handler(
                    servicer.SaveMediaItemPeople,
                    request_deserializer=src_dot_protos_dot_api__pb2.MediaItemPeopleRequest.FromString,
                    response_serializer=google_dot_protobuf_dot_empty__pb2.Empty.SerializeToString,
            ),
            'SaveMediaItemFinalResult': grpc.unary_unary_rpc_method_handler(
                    servicer.SaveMediaItemFinalResult,
                    request_deserializer=src_dot_protos_dot_api__pb2.MediaItemFinalResultRequest.FromString,
                    response_serializer=google_dot_protobuf_dot_empty__pb2.Empty.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'API', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))
    server.add_registered_method_handlers('API', rpc_method_handlers)


 # This class is part of an EXPERIMENTAL API.
class API(object):
    """Missing associated documentation comment in .proto file."""

    @staticmethod
    def GetWorkerConfig(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(
            request,
            target,
            '/API/GetWorkerConfig',
            google_dot_protobuf_dot_empty__pb2.Empty.SerializeToString,
            src_dot_protos_dot_api__pb2.ConfigResponse.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)

    @staticmethod
    def GetMediaItemFaceEmbeddings(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(
            request,
            target,
            '/API/GetMediaItemFaceEmbeddings',
            src_dot_protos_dot_api__pb2.MediaItemFaceEmbeddingsRequest.SerializeToString,
            src_dot_protos_dot_api__pb2.MediaItemFaceEmbeddingsResponse.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)

    @staticmethod
    def GetUsers(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(
            request,
            target,
            '/API/GetUsers',
            google_dot_protobuf_dot_empty__pb2.Empty.SerializeToString,
            src_dot_protos_dot_api__pb2.GetUsersResponse.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)

    @staticmethod
    def SaveMediaItemMetadata(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(
            request,
            target,
            '/API/SaveMediaItemMetadata',
            src_dot_protos_dot_api__pb2.MediaItemMetadataRequest.SerializeToString,
            google_dot_protobuf_dot_empty__pb2.Empty.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)

    @staticmethod
    def SaveMediaItemPreviewThumbnail(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(
            request,
            target,
            '/API/SaveMediaItemPreviewThumbnail',
            src_dot_protos_dot_api__pb2.MediaItemPreviewThumbnailRequest.SerializeToString,
            google_dot_protobuf_dot_empty__pb2.Empty.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)

    @staticmethod
    def SaveMediaItemPlace(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(
            request,
            target,
            '/API/SaveMediaItemPlace',
            src_dot_protos_dot_api__pb2.MediaItemPlaceRequest.SerializeToString,
            google_dot_protobuf_dot_empty__pb2.Empty.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)

    @staticmethod
    def SaveMediaItemThing(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(
            request,
            target,
            '/API/SaveMediaItemThing',
            src_dot_protos_dot_api__pb2.MediaItemThingRequest.SerializeToString,
            google_dot_protobuf_dot_empty__pb2.Empty.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)

    @staticmethod
    def SaveMediaItemFaces(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(
            request,
            target,
            '/API/SaveMediaItemFaces',
            src_dot_protos_dot_api__pb2.MediaItemFacesRequest.SerializeToString,
            google_dot_protobuf_dot_empty__pb2.Empty.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)

    @staticmethod
    def SaveMediaItemPeople(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(
            request,
            target,
            '/API/SaveMediaItemPeople',
            src_dot_protos_dot_api__pb2.MediaItemPeopleRequest.SerializeToString,
            google_dot_protobuf_dot_empty__pb2.Empty.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)

    @staticmethod
    def SaveMediaItemFinalResult(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(
            request,
            target,
            '/API/SaveMediaItemFinalResult',
            src_dot_protos_dot_api__pb2.MediaItemFinalResultRequest.SerializeToString,
            google_dot_protobuf_dot_empty__pb2.Empty.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)
