"""Faces Component"""
import os
import logging

from grpc import RpcError
import schedule
from google.protobuf.empty_pb2 import Empty   # pylint: disable=no-name-in-module
from annoy import AnnoyIndex

from src.protos.api_pb2_grpc import APIStub
from src.protos.api_pb2 import (  # pylint: disable=no-name-in-module
    MediaItemFacesRequest, MediaItemEmbedding, MediaItemFacePeople,
    MediaItemPeopleRequest, MediaItemFaceEmbeddingsRequest )
from src.components.component import Component
from src.providers.faces.utils import init_faces


class Faces(Component):
    """Faces Component"""
    def __init__(self, api_stub: APIStub, source: str, params: list[str]) -> None:
        super().__init__('faces', api_stub)
        self.source = init_faces(source, params)
        schedule.every(int(os.getenv('SMRITI_WORKER_FACES_TIME', '60'))).minutes.do(self.cluster)

    async def process(self, mediaitem_user_id: str, mediaitem_id: str, _: str, metadata: dict) -> None:
        """Process faces detection for mediaitem"""
        if metadata is None or 'previewPath' not in metadata or ('type' in metadata and metadata['type'] == 'video'):
            return metadata
        try:
            result = self.source.detect(mediaitem_user_id, mediaitem_id, metadata['previewPath'])
            logging.debug(f'extracted faces for user {mediaitem_user_id} mediaitem {mediaitem_id}: {result}')
            if result is not None:
                self._grpc_save_mediaitem_faces(result)
        except Exception as exp:
            logging.error(f'error getting faces response for user {mediaitem_user_id} '+
                          f'mediaitem {mediaitem_id}: {str(exp)}')
        logging.info(f'processed faces for user {mediaitem_user_id} mediaitem {mediaitem_id}')
        return metadata

    def cluster(self) -> None:
        """Get users, get all mediaitems for each user and cluster them"""
        logging.debug('getting all users for clustering mediaitem face embeddings')
        res = self.api_stub.GetUsers(Empty())
        users = res.users
        for user in users:
            try:
                logging.debug(f'getting mediaitem face embeddings for user {user}')
                res = self.api_stub.GetMediaItemFaceEmbeddings(MediaItemFaceEmbeddingsRequest(userId=user))
                mediaitem_face_embeddings = res.mediaItemFaceEmbeddings
                logging.debug(f'clustering total {len(mediaitem_face_embeddings)} faces for user {user}')
                if len(mediaitem_face_embeddings) == 0:
                    continue
                # build tree
                tree = AnnoyIndex(len(mediaitem_face_embeddings[0].embedding.embedding), 'euclidean')
                mediaitems_to_cluster = {}
                new_cluster_idx = -1
                for i, mediaitem_face_embedding in enumerate(mediaitem_face_embeddings):
                    if mediaitem_face_embedding.peopleId is not None and len(mediaitem_face_embedding.peopleId) == 0:
                        mediaitems_to_cluster[i] = None
                    tree.add_item(i, mediaitem_face_embedding.embedding.embedding)
                tree.build(len(mediaitem_face_embeddings))
                # assign and get people for new mediaitems
                for idx in mediaitems_to_cluster:
                    nn_idx, nn_dist = tree.get_nns_by_item(idx, n=2, include_distances=True)
                    nn_idx, nn_dist = nn_idx[1], nn_dist[1]
                    if nn_dist > 0.9:
                        new_cluster_idx += 1
                        mediaitems_to_cluster[idx] = f'{new_cluster_idx}'
                    elif nn_idx in mediaitems_to_cluster:
                        if mediaitems_to_cluster[nn_idx] is None:
                            new_cluster_idx += 1
                            mediaitems_to_cluster[nn_idx] = f'{new_cluster_idx}'  # pylint: disable=modified-iterating-dict
                            mediaitems_to_cluster[idx] = f'{new_cluster_idx}'
                        else:
                            mediaitems_to_cluster[idx] = mediaitems_to_cluster[nn_idx]
                    else:
                        mediaitems_to_cluster[idx] = mediaitem_face_embeddings[nn_idx].peopleId
                mediaItemFacePeople = {}
                for idx in mediaitems_to_cluster:
                    if mediaitem_face_embeddings[idx].mediaItemId not in mediaItemFacePeople:
                        mediaItemFacePeople[mediaitem_face_embeddings[idx].mediaItemId] = {}
                    mediaItemFacePeople[mediaitem_face_embeddings[idx].mediaItemId][mediaitem_face_embeddings[idx].id] = mediaitems_to_cluster[idx]
                logging.debug(f'clustered faces for user {user}: {mediaItemFacePeople}')
                self._grpc_save_mediaitem_people(dict({
                    'userId': user,
                    'mediaItemFacePeople': mediaItemFacePeople,
                }))
            except Exception as exp:
                logging.error(f'error clustering mediaitem face embeddings for user {user}: {str(exp)}')
        logging.info('finished clustering mediaitem face embeddings')

    def _grpc_save_mediaitem_faces(self, result: dict):
        """gRPC call for saving mediaitem faces"""
        try:
            request = MediaItemFacesRequest(
                userId=result['userId'],
                id=result['id'],
                embeddings=[MediaItemEmbedding(embedding=embedding) for embedding in result['embeddings']]
            )
            _ = self.api_stub.SaveMediaItemFaces(request)
        except RpcError as rpc_exp:
            logging.error(
                f'error sending faces for mediaitem {request.id}: {str(rpc_exp)}')

    def _grpc_save_mediaitem_people(self, result: dict):
        """gRPC call for saving mediaitem people"""
        try:
            request = MediaItemPeopleRequest(
                userId=result['userId'],
                mediaItemFacePeople={
                    key: MediaItemFacePeople(facePeople=result['mediaItemFacePeople'][key]) for key in result['mediaItemFacePeople']},
            )
            _ = self.api_stub.SaveMediaItemPeople(request)
        except RpcError as rpc_exp:
            logging.error(
                f'error sending people for mediaitem {request.id}: {str(rpc_exp)}')
