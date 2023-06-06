"""Component"""
from src.protos.api_pb2_grpc import APIStub


class Component:
    """Component"""

    def __init__(self, name: str, api_stub: APIStub) -> None:
        self.name = name
        self.api_stub = api_stub

    async def process(self, mediaitem_user_id: str, mediaitem_id: str, mediaitem_file_path: str, metadata: dict):
        """Component level processing"""
        raise NotImplementedError
