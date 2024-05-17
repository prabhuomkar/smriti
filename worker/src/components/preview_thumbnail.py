"""PreviewThumbnail Component"""
import logging
import random
import base64

import rawpy
from moviepy.editor import VideoFileClip
from PIL import Image as PILImage
from wand.image import Image as WandImage
from grpc import RpcError

from src.protos.api_pb2_grpc import APIStub
from src.protos.api_pb2 import MediaItemPreviewThumbnailRequest  # pylint: disable=no-name-in-module
from src.protos.worker_pb2 import MediaItemComponent, PREVIEW_THUMBNAIL  # pylint: disable=no-name-in-module
from src.components.component import Component


class PreviewThumbnail(Component):
    """PreviewThumbnail Component"""

    PREVIEWABLE_PHOTO_MIME_TYPES = [
        'image/avif', 'image/bmp', 'image/gif', 'image/vnd.microsoft.icon', 'image/x-icon',
        'image/icon', 'image/jpeg', 'image/x-citrix-jpeg', 'image/pjpeg',
        'image/apng', 'image/x-png', 'image/x-citrix-png', 'image/png',
        'image/svg+xml', 'image/webp', 'image/heic', 'image/heif', 'image/tiff'
    ]

    RAW_FILE_TYPES = ['CR2', 'NEF', 'ARW', 'ORF', 'RW2', 'DNG', 'RAF', 'SRW', 'PEF',
                      'MOS', 'CRW', 'MRW', 'RWL', 'SRF', 'RAW']

    PREVIEWABLE_VIDEO_MIME_TYPES = [
        'video/x-msvideo', 'video/mp4', 'video/mpeg', 'video/ogg', 'video/mp2t',
        'video/webm', 'video/3gpp', 'video/3gpp2',
    ]

    def __init__(self, api_stub: APIStub, params: dict) -> None:
        super().__init__(MediaItemComponent.Name(PREVIEW_THUMBNAIL), api_stub)
        self.thumbnail_size = int(params['thumbnail_size'])

    async def process(self, mediaitem_user_id: str, mediaitem_id: str, _: str, metadata: dict) -> dict:
        """Process to generate preview and thumbnail for mediaitem"""
        logging.debug(f'generating preview and thumbnail of mediaitem for \
                      user {mediaitem_user_id} mediaitem {mediaitem_id}')
        result = metadata
        result['userId'] = mediaitem_user_id
        result['id'] = mediaitem_id
        file_path = metadata['sourcePath']
        if metadata['type'] == 'photo' or metadata['type'] == 'unknown':
            # generate preview and thumbnail for a photo
            try:
                result['previewPath'], result['thumbnailPath'], \
                    result['placeholder'] = self._generate_photo_preview_and_thumbnail_and_placeholder(
                    file_path, result['mimeType'], metadata)
                logging.debug(f'extracted preview and thumbnail and placeholder for \
                            user {mediaitem_user_id} photo mediaitem {mediaitem_id}')
            except Exception as exp:
                logging.error(f'error generating preview and thumbnail and placeholder for \
                              user {mediaitem_user_id} photo mediaitem {mediaitem_id}: {str(exp)}')
                result['status'] = 'FAILED'
                self._grpc_save_mediaitem_preview_thumbnail(result)
                return None
        elif metadata['type'] == 'video':
            # generate preview and thumbnail for a video
            try:
                result['previewPath'], result['thumbnailPath'], \
                    result['placeholder'] = self._generate_video_preview_and_thumbnail_and_placeholder(file_path)
                logging.debug(f'extracted preview and thumbnail and placeholder for \
                            user {mediaitem_user_id} video mediaitem {mediaitem_id}')
            except Exception as exp:
                logging.error(f'error generating preview and thumbnail and placeholder for \
                              user {mediaitem_user_id} video mediaitem {mediaitem_id}: {str(exp)}')
                result['status'] = 'FAILED'
                self._grpc_save_mediaitem_preview_thumbnail(result)
                return None

        result['status'] = 'READY'
        self._grpc_save_mediaitem_preview_thumbnail(result)
        logging.info(f'generated preview and thumbnail of mediaitem for \
                      user {mediaitem_user_id} mediaitem {mediaitem_id}')
        return result

    def _grpc_save_mediaitem_preview_thumbnail(self, result: dict):
        """gRPC call for saving mediaitem metadata"""
        try:
            request = MediaItemPreviewThumbnailRequest(
                userId=result['userId'],
                id=result['id'],
                status=result['status'],
                placeholder=result['placeholder'] if 'placeholder' in result else None,
                sourcePath=result['sourcePath'] if 'sourcePath' in result else None,
                previewPath=result['previewPath'] if 'previewPath' in result else None,
                thumbnailPath=result['thumbnailPath'] if 'thumbnailPath' in result else None,
            )
            _ = self.api_stub.SaveMediaItemPreviewThumbnail(request)
        except RpcError as rpc_exp:
            logging.error(
                f'error sending preview and thumbnail for mediaitem {request.id}: {str(rpc_exp)}')

    def _generate_photo_thumbnail_and_placeholder(self, original_file_path: str, preview_file_path: str):
        """Generate thumbnail image from photo"""
        try:
            thumbnail_path = f'{original_file_path}-thumbnail'
            placeholder = ''
            with WandImage(filename=preview_file_path) as img:
                lidx = 0 if img.size[0] > img.size[1] else 1
                sidx = 1 if lidx == 0 else 0
                percent = self.thumbnail_size/float(img.size[lidx])
                size = int((float(img.size[sidx])*float(percent)))
                img.resize(self.thumbnail_size, size)
                with img.convert('jpeg') as converted:
                    converted.save(filename=thumbnail_path)
                lidx = 0 if img.size[0] > img.size[1] else 1
                sidx = 1 if lidx == 0 else 0
                percent = 8/float(img.size[lidx])
                size = int((float(img.size[sidx])*float(percent)))
                img.resize(8, size)
                placeholder_bytes = img.make_blob()
                placeholder = base64.b64encode(placeholder_bytes).decode('utf-8')
            return thumbnail_path, str(placeholder)
        except Exception as exp:
            logging.error(f'error generating photo thumbnail and \
                          placeholder for {original_file_path} {preview_file_path}: {str(exp)}')
            return None, None

    # pylint: disable=too-many-locals
    def _generate_photo_preview_and_thumbnail_and_placeholder(self, original_file_path: str, mime_type: str,
                                                              metadata: dict):
        """Generate preview and thumbnail image for a photo"""
        preview_path = f'{original_file_path}-preview'
        if mime_type in self.PREVIEWABLE_PHOTO_MIME_TYPES and not self._is_raw(metadata):
            try:
                with open(original_file_path, 'rb') as file_reader:
                    with WandImage(file=file_reader) as original:
                        with original.convert('jpeg') as converted:
                            converted.save(filename=preview_path)
                thumbnail_path, placeholder = self._generate_photo_thumbnail_and_placeholder(
                    original_file_path, preview_path)
                return preview_path, thumbnail_path, placeholder
            except Exception as exp:
                logging.warning(f'error generating preview for default \
                                photo mediaitem: {original_file_path}: {str(exp)}')
                try:
                    with open(original_file_path, 'rb') as file_reader:
                        if 'File:FileType'in metadata:
                            with WandImage(file=file_reader, format=metadata['File:FileType']) as original:
                                with original.convert('jpeg') as converted:
                                    converted.save(filename=preview_path)
                    thumbnail_path, placeholder = self._generate_photo_thumbnail_and_placeholder(
                        original_file_path, preview_path)
                    return preview_path, thumbnail_path, placeholder
                except Exception as ft_exp:
                    logging.warning(f'error generating preview for default \
                                    photo mediaitem: {original_file_path}: {str(ft_exp)}')
                    with rawpy.imread(original_file_path) as raw:
                        rgb = raw.postprocess(use_camera_wb=True)
                        img = PILImage.fromarray(rgb)
                        img.save(preview_path, format='JPEG')
                    thumbnail_path, placeholder = self._generate_photo_thumbnail_and_placeholder(
                        original_file_path, preview_path)
                    return preview_path, thumbnail_path, placeholder
        else:
            try:
                with rawpy.imread(original_file_path) as raw:
                    rgb = raw.postprocess(use_camera_wb=True)
                    img = PILImage.fromarray(rgb)
                    img.save(preview_path, format='JPEG')
                thumbnail_path, placeholder = self._generate_photo_thumbnail_and_placeholder(
                    original_file_path, preview_path)
                return preview_path, thumbnail_path, placeholder
            except Exception as exp:
                logging.warning(f'error generating preview for raw photo mediaitem: {original_file_path}: {str(exp)}')
                try:
                    with open(original_file_path, 'rb') as file_reader:
                        with WandImage(file=file_reader) as original:
                            with original.convert('jpeg') as converted:
                                converted.save(filename=preview_path)
                            thumbnail_path, placeholder = self._generate_photo_thumbnail_and_placeholder(
                                original_file_path, preview_path)
                except Exception as nm_exp:
                    logging.warning(f'error generating preview for raw \
                                photo mediaitem: {original_file_path}: {str(nm_exp)}')
                    with open(original_file_path, 'rb') as file_reader:
                        if 'File:FileType'in metadata:
                            with WandImage(file=file_reader, format=metadata['File:FileType']) as original:
                                with original.convert('jpeg') as converted:
                                    converted.save(filename=preview_path)
                    thumbnail_path, placeholder = self._generate_photo_thumbnail_and_placeholder(
                        original_file_path, preview_path)
                    return preview_path, thumbnail_path, placeholder
        return preview_path, thumbnail_path, placeholder

    def _generate_video_thumbnail_and_placeholder(self, preview_video_path: str):
        """Generate thumbnail image from video"""
        try:
            clip = VideoFileClip(preview_video_path)
            video_thumbnail_path = f'{preview_video_path}_thumbnail.jpeg'
            clip.save_frame(video_thumbnail_path, t=random.uniform(0.1, clip.duration))
            placeholder = ''
            with WandImage(filename=video_thumbnail_path) as img:
                lidx = 0 if img.size[0] > img.size[1] else 1
                sidx = 1 if lidx == 0 else 0
                percent = self.thumbnail_size/float(img.size[lidx])
                size = int((float(img.size[sidx])*float(percent)))
                img.resize(self.thumbnail_size, size)
                with img.convert('jpeg') as converted:
                    converted.save(filename=video_thumbnail_path)
                lidx = 0 if img.size[0] > img.size[1] else 1
                sidx = 1 if lidx == 0 else 0
                percent = 8/float(img.size[lidx])
                size = int((float(img.size[sidx])*float(percent)))
                img.resize(8, size)
                placeholder_bytes = img.make_blob()
                placeholder = base64.b64encode(placeholder_bytes).decode('utf-8')
            return video_thumbnail_path, str(placeholder)
        except Exception as exp:
            logging.error(f'error generating video thumbnail and \
                          placeholder for {preview_video_path}: {str(exp)}')
            return None, None

    def _generate_video_preview_and_thumbnail_and_placeholder(self, original_file_path: str):
        """Generate preview and thumbnail image for a video"""
        video = VideoFileClip(original_file_path)
        video_preview_path = f'{original_file_path}-preview.mp4'
        video.write_videofile(video_preview_path, codec='libx264', logger=None, verbose=False)
        video_thumbnail_path, placeholder = self._generate_video_thumbnail_and_placeholder(video_preview_path)
        return video_preview_path, video_thumbnail_path, placeholder

    def _is_raw(self, metadata: dict) -> bool:
        """Detect if the image is RAW irrespective of the image mimetype"""
        if 'EXIF:JpgFromRaw' in metadata or 'EXIF:OriginalRawFileName' in metadata or\
            'EXIF:NewRawImageDigest' in metadata or 'EXIF:RawDataUniqueID' in metadata:
            return True
        if 'File:FileType' in metadata and metadata['File:FileType'] in self.RAW_FILE_TYPES:
            return True
        return False
