"""Worker"""
from concurrent import futures
import logging
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
    logging.basicConfig(level=logging.getLevelName(os.getenv('PENSIEVE_LOG_LEVEL', 'INFO')))
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    worker_pb2_grpc.add_WorkerServicer_to_server(WorkerServicer(), server)
    port = int(os.getenv('PENSIEVE_WORKER_PORT', '15002'))
    server.add_insecure_port(f'[::]:{port}')
    logging.info(f'starting grpc server on: {port}')
    server.start()
    server.wait_for_termination()
