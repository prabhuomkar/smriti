"""Finalize Component"""
import logging

from grpc import RpcError

from src.protos.api_pb2_grpc import APIStub
from src.protos.api_pb2 import MediaItemFinalResultRequest  # pylint: disable=no-name-in-module
from src.components.component import Component


class Finalize(Component):
    """Finalize Component"""

    def __init__(self, api_stub: APIStub) -> None:
        super().__init__('finalize', api_stub)

    async def process(self, mediaitem_user_id: str, mediaitem_id: str, _: str, metadata: dict) -> dict:
        """Process finalizing mediaitem"""
        logging.debug(f'finalizing mediaitem for user {mediaitem_user_id} mediaitem {mediaitem_id}')
        result = {}
        result['userId'] = mediaitem_user_id
        result['id'] = mediaitem_id
        result['keywords'] = metadata['keywords']
        result['embedding'] = metadata['embedding']
        self._grpc_final_save_mediaitem(result)
        logging.debug(f'finalized mediaitem for user {mediaitem_user_id} mediaitem {mediaitem_id}')
        return None

    def _grpc_final_save_mediaitem(self, result: dict):
        """gRPC call for final save of mediaitem"""
        try:
            request = MediaItemFinalResultRequest(
                userId=result['userId'],
                id=result['id'],
                keywords=result['keywords'],
                embedding=result['embedding']
            )
            _ = self.api_stub.SaveMediaItemFinalResult(request)
        except RpcError as rpc_exp:
            logging.error(
                f'error finalizing save for mediaitem {request.id}: {str(rpc_exp)}')
