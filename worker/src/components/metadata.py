"""Metadata Component"""
import logging
import datetime
import random
import re
import base64

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
        super().__init__('metadata', api_stub)
        self.thumbnail_size = int(params['thumbnail_size'])

    # pylint: disable=too-many-statements,too-many-branches
    async def process(self, mediaitem_user_id: str, mediaitem_id: str, mediaitem_file_path: str, _: dict) -> dict:
        """Process required metadata and generate thumbnail from EXIF data"""
        file_path = f'{mediaitem_file_path}/{mediaitem_id}'

        # extract metadata
        result = {}
        result['userId'] = mediaitem_user_id
        result['id'] = mediaitem_id
        result['status'] = 'UNSPECIFIED'
        result['sourcePath'] = file_path
        result['type'] = 'unknown'
        result['category'] = 'default'
        try:
            with exiftool.ExifToolHelper() as ethelper:
                ethelper.check_execute = False
                metadata = ethelper.get_metadata(file_path)[0]
                logging.debug(f'metadata for user {mediaitem_user_id} mediaitem {mediaitem_id}: {metadata}')
                result['mimeType'] = getval_from_dict(metadata, ['File:MIMEType'])
                result['type'] = 'photo' if result['mimeType'] and 'image' in result['mimeType'] else \
                    'video' if result['mimeType'] and 'video' in result['mimeType'] else 'unknown'
                result['width'] = getval_from_dict(metadata, ['EXIF:SensorWidth', 'EXIF:ImageWidth',
                                                              'EXIF:ExifImageWidth', 'File:ImageWidth',
                                                              'PNG:ImageWidth', 'XMP:ExifImageWidth',
                                                              'QuickTime:ImageWidth', 'QuickTime:SourceImageWidth'],
                                                              return_type='int')
                result['height'] = getval_from_dict(metadata, ['EXIF:SensorHeight', 'EXIF:ImageHeight',
                                                               'EXIF:ExifImageHeight', 'File:ImageHeight',
                                                               'PNG:ImageHeight', 'XMP:ExifImageHeight',
                                                               'QuickTime:ImageHeight', 'QuickTime:SourceImageHeight'],
                                                               return_type='int')
                if (result['height'] is None or result['width'] is None) and 'Composite:ImageSize' in metadata:
                    composite_dims = metadata['Composite:ImageSize'].split(' ')
                    if len(composite_dims) == 2:
                        result['width'] = int(composite_dims[0])
                        result['height'] = int(composite_dims[1])
                creation_time = getval_from_dict(metadata, ['EXIF:DateTimeOriginal', 'EXIF:CreateDate', \
                                            'XMP:CreateDate', 'XMP:DateCreated', 'XMP:DateTimeOriginal', \
                                            'Composite:SubSecCreateDate', 'Composite:SubSecDateTimeOriginal', \
                                            'QuickTime:CreateDate', 'QuickTime:TrackModifyDate', \
                                            'QuickTime:MediaCreateDate', 'QuickTime:CreationDate', \
                                            'EXIF:ModifyDate', 'XMP:ModifyDate', \
                                            'File:FileModifyDate', 'File:FileAccessDate', 'File:FileInodeChangeDate'])
                # work(omkar): handle timezone when "its time" :P
                if creation_time and re.search(r'[\+]\d{2}:\d{2}', creation_time):
                    creation_time = creation_time.rsplit("+", maxsplit=1)[0] if creation_time else None
                elif creation_time and re.search(r'[\-]\d{2}:\d{2}', creation_time):
                    creation_time = creation_time.rsplit("-", maxsplit=1)[0] if creation_time else None
                if creation_time:
                    creation_time = creation_time.replace('- ', '-')
                    creation_time = creation_time.replace(' -', '-')
                    creation_time = creation_time.replace('T', ' ') if 'T' in creation_time else creation_time
                    creation_time = creation_time.replace('Z', '') if 'Z' in creation_time else creation_time
                if creation_time and '-' not in creation_time:
                    result['creationTime'] = datetime.datetime.strptime(creation_time, '%Y:%m:%d %H:%M:%S').replace(
                        tzinfo=datetime.timezone.utc).strftime('%Y-%m-%d %H:%M:%S') if creation_time else None
                elif creation_time:
                    result['creationTime'] = datetime.datetime.strptime(creation_time, '%Y-%m-%d %H:%M:%S').replace(
                        tzinfo=datetime.timezone.utc).strftime('%Y-%m-%d %H:%M:%S') if creation_time else None
                camera_make = getval_from_dict(metadata, ['EXIF:Make', 'XMP:Make', 'QuickTime:Make'])
                result['cameraMake'] = camera_make.strip() if camera_make else None
                camera_model = getval_from_dict(metadata, ['EXIF:Model', 'XMP:Model', 'QuickTime:Model'])
                result['cameraModel'] = camera_model.strip() if camera_model else None
                result['focalLength'] = getval_from_dict(metadata, ['EXIF:FocalLength'])
                result['apertureFNumber'] = getval_from_dict(metadata, ['EXIF:FNumber'])
                result['isoEquivalent'] = getval_from_dict(metadata, ['EXIF:ISO', 'XMP:ISO'])
                result['exposureTime'] = getval_from_dict(metadata, ['EXIF:ExposureTime'])
                result['fps'] = getval_from_dict(metadata, ['QuickTime:VideoFrameRate'], return_type='float')
                result['fps'] = str(round(result['fps'])) if result['fps'] is not None else None
                result['latitude'] = getval_from_dict(metadata, ['Composite:GPSLatitude'], return_type='float')
                result['longitude'] = getval_from_dict(metadata, ['Composite:GPSLongitude'], return_type='float')
                if result['latitude'] is None or result['longitude'] is None:
                    if 'EXIF:GPSLatitudeRef' in metadata and 'EXIF:GPSLatitude' in metadata and \
                        'EXIF:GPSLongitudeRef' in metadata and 'EXIF:GPSLongitude' in metadata and \
                         metadata['EXIF:GPSLatitudeRef'] != '' and metadata['EXIF:GPSLatitude'] != '' and \
                         metadata['EXIF:GPSLongitudeRef'] != '' and metadata['EXIF:GPSLongitude'] != '':
                        result['latitude'] = float(metadata['EXIF:GPSLatitude']) * (1 \
                            if metadata['EXIF:GPSLatitudeRef'] == 'N' else -1)
                        result['longitude'] = float(metadata['EXIF:GPSLongitude']) * (1 \
                            if metadata['EXIF:GPSLongitudeRef'] == 'E' else -1)
                    elif 'Composite:GPSPosition' in metadata:
                        splits = metadata['Composite:GPSPosition'].split()
                        if len(splits) == 2:
                            result['latitude'] = float(splits[0])
                            result['longitude'] = float(splits[1])
                result['category'] = self._get_mediaitem_category(metadata, result)
                logging.debug(f'extracted metadata for user {mediaitem_user_id} mediaitem {mediaitem_id}: {result}')
        except Exception as exp:
            logging.error(
                f'error extracting exif metadata for user {mediaitem_user_id} mediaitem {mediaitem_id}: {str(exp)}')
            result['status'] = 'FAILED'
            self._grpc_save_mediaitem_metadata(result)
            return None

        if result['type'] == 'photo' or result['type'] == 'unknown':
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
                self._grpc_save_mediaitem_metadata(result)
                return None
        elif result['type'] == 'video':
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
                sourcePath=result['sourcePath'] if 'sourcePath' in result else None,
                previewPath=result['previewPath'] if 'previewPath' in result else None,
                thumbnailPath=result['thumbnailPath'] if 'thumbnailPath' in result else None,
                placeholder=result['placeholder'] if 'placeholder' in result else None,
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
            if 'QuickTime:CaptureMode' in metadata and metadata['QuickTime:CaptureMode'].lower() == 'time-lapse':
                return 'timelapse'
        return 'default'

    def _is_raw(self, metadata: dict) -> bool:
        """Detect if the image is RAW irrespective of the image mimetype"""
        if 'EXIF:JpgFromRaw' in metadata or 'EXIF:OriginalRawFileName' in metadata or\
            'EXIF:NewRawImageDigest' in metadata or 'EXIF:RawDataUniqueID' in metadata:
            return True
        if 'File:FileType' in metadata and metadata['File:FileType'] in self.RAW_FILE_TYPES:
            return True
        return False
