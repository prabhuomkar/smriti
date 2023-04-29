"""Metadata Component"""
import io
import os
import logging
import datetime
import random

import exiftool
import rawpy
from moviepy.editor import VideoFileClip
from PIL import Image as PILImage
from wand.image import Image as WandImage
from grpc import RpcError

from src.protos.api_pb2_grpc import APIStub
from src.protos.api_pb2 import MediaItemMetadataRequest  # pylint: disable=no-name-in-module
from src.components.component import Component
from src.utils import getval_from_dict


class Metadata(Component):
    """Metadata Component"""

    PREVIEWABLE_PHOTO_MIME_TYPES = [
        'image/avif', 'image/bmp', 'image/gif', 'image/vnd.microsoft.icon', 'image/x-icon',
        'image/icon', 'image/jpeg', 'image/x-citrix-jpeg', 'image/pjpeg',
        'image/apng', 'image/x-png', 'image/x-citrix-png', 'image/png', 'image/tiff',
        'image/svg+xml', 'image/webp', 'image/heic', 'image/heif',
    ]

    PREVIEWABLE_VIDEO_MIME_TYPES = [
        'video/x-msvideo', 'video/mp4', 'video/mpeg', 'video/ogg', 'video/mp2t',
        'video/webm', 'video/3gpp', 'video/3gpp2',
    ]

    def __init__(self, storage, api_stub: APIStub) -> None:
        super().__init__('metadata', storage, api_stub)

    # pylint: disable=too-many-statements
    async def process(self, mediaitem_user_id: str, mediaitem_id: str, _: dict) -> dict:
        """Process required metadata and generate thumbnail from EXIF data"""
        file_path = self.storage.get(mediaitem_id)

        # extract metadata
        result = {}
        result['userId'] = mediaitem_user_id
        result['id'] = mediaitem_id
        result['status'] = 'UNSPECIFIED'
        result['sourceUrl'] = file_path
        result['type'] = 'unknown'
        result['category'] = 'default'
        try:
            with exiftool.ExifToolHelper() as ethelper:
                metadata = ethelper.get_metadata(file_path)[0]
                logging.debug(f'metadata for user {mediaitem_user_id} mediaitem {mediaitem_id}: {metadata}')
                result['mimeType'] = getval_from_dict(metadata, ['File:MIMEType'])
                result['type'] = 'photo' if 'image' in metadata['File:MIMEType'] else \
                    'video' if 'video' in metadata['File:MIMEType'] else 'unknown'
                result['width'] = getval_from_dict(metadata, ['File:ImageWidth', 'EXIF:ExifImageWidth',
                                                              'PNG:ImageWidth', 'QuickTime:ImageWidth',
                                                              'QuickTime:SourceImageWidth'], return_type='int')
                result['height'] = getval_from_dict(metadata, ['File:ImageHeight', 'EXIF:ExifImageHeight',
                                                               'PNG:ImageHeight', 'QuickTime:ImageHeight',
                                                               'QuickTime:SourceImageHeight'], return_type='int')
                if result['height'] is None or result['width'] is None and 'Composite:ImageSize' in metadata:
                    composite_dims = metadata['Composite:ImageSize'].split(' ')
                    if len(composite_dims) == 2:
                        result['width'] = int(composite_dims[0])
                        result['height'] = int(composite_dims[1])
                result['creationTime'] = getval_from_dict(metadata, ['EXIF:DateTimeOriginal', 'EXIF:CreateDate', \
                                            'XMP:CreateDate', 'XMP:DateCreated', 'Composite:SubSecCreateDate', \
                                            'Composite:SubSecDateTimeOriginal', 'QuickTime:CreateDate', \
                                            'QuickTime:TrackModifyDate', 'QuickTime:MediaCreateDate', \
                                            'QuickTime:CreationDate', 'File:FileModifyDate', \
                                            'File:FileAccessDate', 'File:FileInodeChangeDate'])
                # work(omkar): handle timezone when "its time" :P
                creation_time = result['creationTime'].split("+")[0] if result['creationTime'] else None
                result['creationTime'] = datetime.datetime.strptime(creation_time, '%Y:%m:%d %H:%M:%S').replace(
                    tzinfo=datetime.timezone.utc).strftime('%Y-%m-%d %H:%M:%S') if creation_time else None
                result['cameraMake'] = getval_from_dict(metadata, ['EXIF:Make', 'QuickTime:Make'])
                result['cameraModel'] = getval_from_dict(metadata, ['EXIF:Model', 'QuickTime:Model'])
                result['focalLength'] = getval_from_dict(metadata, ['EXIF:FocalLength'])
                result['apertureFNumber'] = getval_from_dict(metadata, ['EXIF:FNumber'])
                result['isoEquivalent'] = getval_from_dict(metadata, ['EXIF:ISO'])
                result['exposureTime'] = getval_from_dict(metadata, ['EXIF:ExposureTime'])
                result['fps'] = getval_from_dict(metadata, ['QuickTime:VideoFrameRate'])
                result['latitude'] = getval_from_dict(metadata, ['EXIF:GPSLatitude', \
                                                                'Composite:GPSLatitude'], return_type='float')
                result['longitude'] = getval_from_dict(metadata, ['EXIF:GPSLongitude', \
                                                                'Composite:GPSLongitude'], return_type='float')
                result['category'] = self._get_mediaitem_category(metadata, result)
                logging.debug(f'extracted metadata for user {mediaitem_user_id} mediaitem {mediaitem_id}: {result}')
        except Exception as exp:
            logging.error(
                f'error extracting exif metadata for user {mediaitem_user_id} mediaitem {mediaitem_id}: {str(exp)}')
            result['status'] = 'FAILED'
            self._grpc_save_mediaitem_metadata(result)
            return None

        if result['type'] == 'photo':
            # generate and upload preview and thumbnail for a photo
            try:
                preview_bytes, thumbnail_bytes = self._generate_photo_preview_and_thumbnail(
                    file_path, result['mimeType'])
                preview_url = self.storage.upload(
                    mediaitem_id, preview_bytes, 'previews')
                thumbnail_url = self.storage.upload(
                    mediaitem_id, thumbnail_bytes, 'thumbnails')
                result['previewUrl'] = preview_url
                result['thumbnailUrl'] = thumbnail_url
                logging.debug(f'extracted preview and thumbnail for \
                            user {mediaitem_user_id} photo mediaitem {mediaitem_id}: {result}')
            except Exception as exp:
                logging.error(
                    f'error generating and uploading preview and thumbnail for \
                        user {mediaitem_user_id} photo mediaitem {mediaitem_id}: {str(exp)}')
                result['status'] = 'FAILED'
                self._grpc_save_mediaitem_metadata(result)
                return None
        elif result['type'] == 'video':
            # generate and upload preview and thumbnail for a video
            try:
                preview_bytes, thumbnail_bytes = self._generate_video_preview_and_thumbnail(file_path)
                preview_url = self.storage.upload(
                    mediaitem_id, preview_bytes, 'previews')
                thumbnail_url = self.storage.upload(
                    mediaitem_id, thumbnail_bytes, 'thumbnails')
                result['previewUrl'] = preview_url
                result['thumbnailUrl'] = thumbnail_url
                logging.debug(f'extracted preview and thumbnail for \
                            user {mediaitem_user_id} video mediaitem {mediaitem_id}: {result}')
            except Exception as exp:
                logging.error(
                    f'error generating and uploading preview and thumbnail for \
                        user {mediaitem_user_id} video mediaitem {mediaitem_id}: {str(exp)}')
                result['status'] = 'FAILED'
                self._grpc_save_mediaitem_metadata(result)
                return None

        # metadata, preview and thumbnail extracted, so its READY
        result['status'] = 'READY'
        self._grpc_save_mediaitem_metadata(result)
        logging.info(f'processed metadata for user {mediaitem_user_id} mediaitem {mediaitem_id}')

        return result

    def _grpc_save_mediaitem_metadata(self, result: dict):
        """gRPC call for saving mediaitem metadata"""
        try:
            request = MediaItemMetadataRequest(
                userId=result['userId'],
                id=result['id'],
                status=result['status'],
                mimeType=result['mimeType'] if 'mimeType' in result else None,
                sourceUrl=result['sourceUrl'] if 'sourceUrl' in result else None,
                previewUrl=result['previewUrl'] if 'previewUrl' in result else None,
                thumbnailUrl=result['thumbnailUrl'] if 'thumbnailUrl' in result else None,
                type=result['type'] if 'type' in result else None,
                category=result['category'] if 'category' in result else None,
                width=result['width'] if 'width' in result else None,
                height=result['height'] if 'height' in result else None,
                creationTime=result['creationTime'] if 'creationTime' in result else None,
                cameraMake=result['cameraMake'] if 'cameraMake' in result else None,
                cameraModel=result['cameraModel'] if 'cameraModel' in result else None,
                focalLength=result['focalLength'] if 'focalLength' in result else None,
                apertureFNumber=result['apertureFNumber'] if 'apertureFNumber' in result else None,
                isoEquivalent=result['isoEquivalent'] if 'isoEquivalent' in result else None,
                exposureTime=result['exposureTime'] if 'exposureTime' in result else None,
                fps=result['fps'] if 'fps' in result else None,
                latitude=result['latitude'] if 'latitude' in result else None,
                longitude=result['longitude'] if 'longitude' in result else None,
            )
            _ = self.api_stub.SaveMediaItemMetadata(request)
        except RpcError as rpc_exp:
            logging.error(
                f'error sending metadata for mediaitem {request.id}: {str(rpc_exp)}')

    def _generate_photo_thumbnail(self, preview_bytes: bytes):
        """Generate thumbnail image from photo"""
        # work(omkar): thumbnail size should be configurable through UI
        with WandImage(blob=preview_bytes) as img:
            lidx = 0 if img.size[0] > img.size[1] else 1
            sidx = 1 if lidx == 0 else 0
            percent = 512/float(img.size[lidx])
            size = int((float(img.size[sidx])*float(percent)))
            img.resize(512, size)
            with img.convert('jpeg') as converted:
                return converted.make_blob('jpeg')

    def _generate_photo_preview_and_thumbnail(self, original_file_path: str, mime_type: str):
        """Generate preview and thumbnail image for a photo"""
        if mime_type in self.PREVIEWABLE_PHOTO_MIME_TYPES:
            with WandImage(filename=original_file_path) as original:
                with original.convert('jpeg') as converted:
                    img_bytes = converted.make_blob('jpeg')
        else:
            with rawpy.imread(original_file_path) as raw:
                rgb = raw.postprocess(use_camera_wb=True)
                img = PILImage.fromarray(rgb)
                img_bytes = io.BytesIO()
                img.save(img_bytes, format='JPEG')
        return img_bytes, self._generate_photo_thumbnail(img_bytes)

    def _generate_video_thumbnail(self, preview_video_path: str):
        """Generate thumbnail image from video"""
        # work(omkar): thumbnail size should be configurable through UI
        clip = VideoFileClip(preview_video_path)
        video_thumbnail_path = f'{preview_video_path}_thumbnail.jpeg'
        clip.save_frame(video_thumbnail_path, t=random.uniform(0.1, clip.duration))
        with open(video_thumbnail_path, 'rb') as thumbnail_file:
            thumbnail_bytes = thumbnail_file.read()
            thumbnail_file.close()
        os.remove(video_thumbnail_path)
        return thumbnail_bytes

    def _generate_video_preview_and_thumbnail(self, original_file_path: str):
        """Generate preview and thumbnail image for a video"""
        video = VideoFileClip(original_file_path)
        video_preview_path = f'{original_file_path}.mp4'
        video.write_videofile(video_preview_path, codec='libx264', logger=None, verbose=False)
        with open(video_preview_path, 'rb') as video_file:
            video_bytes = video_file.read()
            video_file.close()
        video_thumbnail_bytes = self._generate_video_thumbnail(video_preview_path)
        os.remove(video_preview_path)
        return video_bytes, video_thumbnail_bytes

    def _get_mediaitem_category(self, metadata: dict, result: dict) -> str:
        """Get mediaitem category from metadata"""
        if result['type'] == 'photo':
            if 'EXIF:UserComment' in metadata and metadata['EXIF:UserComment'].lower() == 'screenshot' \
                or 'XMP:UserComment' in metadata and metadata['XMP:UserComment'].lower() == 'screenshot':
                return 'screenshot'
            if result['width'] > 10000 and result['height']*4 <= result['width']:
                return 'panorama'
        if result['type'] == 'video':
            if int(result['fps']) > 150:
                return 'slow'
            if 'QuickTime:LivePhotoAuto' in metadata:
                return 'live'
        return 'default'
