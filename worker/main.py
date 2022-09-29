"""Worker"""
from concurrent import futures
import os

import grpc
import worker_pb2_grpc


class WorkerServicer(worker_pb2_grpc.WorkerServicer):
    """Worker gRPC Service"""
    def __init__(self) -> None:
        pass

    def MediaItemProcess(self, request_iterator, context):
        """Process MediaItem"""
        raise NotImplementedError

if __name__ == '__main__':
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    worker_pb2_grpc.add_WorkerServicer_to_server(WorkerServicer(), server)
    server.add_insecure_port(f'[::]:{int(os.getenv("PENSIEVE_WORKER_PORT", "15002"))}')
    server.start()
    server.wait_for_termination()
