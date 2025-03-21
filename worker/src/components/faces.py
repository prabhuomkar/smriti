"""Faces Component"""
import logging

from grpc import RpcError
import schedule
from google.protobuf.empty_pb2 import Empty   # pylint: disable=no-name-in-module
from annoy import AnnoyIndex

from src.protos.api_pb2_grpc import APIStub
from src.protos.api_pb2 import (  # pylint: disable=no-name-in-module
    MediaItemFacesRequest, MediaItemEmbedding, MediaItemFacePeople,
    MediaItemPeopleRequest, MediaItemFaceEmbeddingsRequest, MediaItemFaceEmbedding)
from src.protos.worker_pb2 import MediaItemComponent, FACES  # pylint: disable=no-name-in-module
from src.components.component import Component
from src.providers.faces.utils import init_faces


class Faces(Component):
    """Faces Component"""
    def __init__(self, api_stub: APIStub, source: str, params: dict) -> None:
        super().__init__(MediaItemComponent.Name(FACES), api_stub)
        self.source = init_faces(source, params)
        self.clustering_framework = params['clustering']
        self.data = None
        schedule.every(int(params['minutes'])).minutes.do(self.cluster)

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
            logging.error('error getting faces response for user %s mediaitem %s: %s',
                          mediaitem_user_id, mediaitem_id, exp)
        logging.info(f'processed faces for user {mediaitem_user_id} mediaitem {mediaitem_id}')
        return metadata

    # pylint: disable=too-many-locals,too-many-branches
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
                # build tree/index
                if self.clustering_framework == 'annoy':
                    self._build_annoy(mediaitem_face_embeddings)
                # select items which need to be clustered
                mediaitems_to_cluster = {}
                new_cluster_idx = -1
                for i, mediaitem_face_embedding in enumerate(mediaitem_face_embeddings):
                    if mediaitem_face_embedding.peopleId is not None and len(mediaitem_face_embedding.peopleId) == 0:
                        mediaitems_to_cluster[i] = None
                # assign and get people for new mediaitems
                for idx in mediaitems_to_cluster:
                    nn_idx, nn_dist = self._get_nn_for_mediaitem(idx)
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
                mediaitem_face_people = {}
                for idx, face_embed in mediaitems_to_cluster.items():
                    mf_embedding = mediaitem_face_embeddings[idx]
                    if mf_embedding.mediaItemId not in mediaitem_face_people:
                        mediaitem_face_people[mf_embedding.mediaItemId] = {}
                    mediaitem_face_people[mf_embedding.mediaItemId][mf_embedding.id] = face_embed
                logging.debug(f'clustered faces for user {user}: {mediaitem_face_people}')
                self._grpc_save_mediaitem_people(dict({
                    'userId': user,
                    'mediaitem_face_people': mediaitem_face_people,
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
                embeddings=[MediaItemEmbedding(embedding=embedding) for embedding in result['embeddings']],
                thumbnails=result['thumbnails']
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
                    key: MediaItemFacePeople(facePeople=result['mediaitem_face_people'][key]) \
                        for key in result['mediaitem_face_people']},
            )
            _ = self.api_stub.SaveMediaItemPeople(request)
        except RpcError as rpc_exp:
            logging.error(
                f'error sending people for mediaitem {request.id}: {str(rpc_exp)}')

    def _build_annoy(self, mediaitem_face_embeddings: list[MediaItemFaceEmbedding]):
        """Cluster faces using Annoy library"""
        # build tree
        self.data = AnnoyIndex(len(mediaitem_face_embeddings[0].embedding.embedding), 'euclidean')
        for i, mediaitem_face_embedding in enumerate(mediaitem_face_embeddings):
            self.data.add_item(i, mediaitem_face_embedding.embedding.embedding)
        self.data.build(len(mediaitem_face_embeddings))

    def _get_nn_for_mediaitem(self, mediaitem_idx: int):
        """Get nearest neighbour index for mediaitem"""
        if self.clustering_framework == 'annoy':
            nn_idx, nn_dist = self.data.get_nns_by_item(mediaitem_idx, n=2, include_distances=True)
            return (nn_idx[1], nn_dist[1])
        return None
