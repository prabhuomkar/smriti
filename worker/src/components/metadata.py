"""Metadata Component"""
import io
import os
import asyncio
import logging
import datetime
import random

import exiftool
import rawpy
from moviepy.editor import VideoFileClip
from PIL import Image as PILImage
from wand.image import Image as WandImage

from src.components.place import process_place
from src.components.utils import grpc_save_mediaitem_metadata, getval_from_dict
from src.protos.api_pb2_grpc import APIStub


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


# pylint: disable=too-many-statements
async def process_metadata(storage, api_stub: APIStub, mediaitem_user_id: str, mediaitem_id: str) -> None:
    """Process required metadata and generate thumbnail from EXIF data"""
    file_path = storage.get(mediaitem_id)

    # extract metadata
    result = {}
    result['userId'] = mediaitem_user_id
    result['id'] = mediaitem_id
    result['status'] = 'UNSPECIFIED'
    result['sourceUrl'] = file_path
    try:
        with exiftool.ExifToolHelper() as et:
            metadata = et.get_metadata(file_path)[0]
            logging.debug(f'metadata for user {mediaitem_user_id} mediaitem {mediaitem_id}: {metadata}')
            result['mimeType'] = getval_from_dict(metadata, ['File:MIMEType'])
            result['type'] = 'photo' if 'image' in metadata['File:MIMEType'] else \
                'video' if 'video' in metadata['File:MIMEType'] else None
            result['width'] = getval_from_dict(metadata, ['File:ImageWidth', 'EXIF:ExifImageWidth', \
                                        'QuickTime:ImageWidth', 'QuickTime:SourceImageWidth'], return_type='int')
            result['height'] = getval_from_dict(metadata, ['File:ImageHeight', 'EXIF:ExifImageHeight', \
                                        'QuickTime:ImageHeight', 'QuickTime:SourceImageHeight'], return_type='int')
            result['creationTime'] = getval_from_dict(metadata, ['EXIF:DateTimeOriginal', 'EXIF:CreateDate', \
                                        'XMP:CreateDate', 'XMP:DateCreated', 'Composite:SubSecCreateDate', \
                                        'Composite:SubSecDateTimeOriginal', 'QuickTime:CreateDate', \
                                        'QuickTime:TrackModifyDate', 'QuickTime:MediaCreateDate', \
                                        'QuickTime:CreationDate', 'File:FileModifyDate', \
                                        'File:FileAccessDate', 'File:FileInodeChangeDate'])
            # work(omkar): handle timezone when "its time" :P
            creationTime = result['creationTime'].split("+")[0] if result['creationTime'] else None
            result['creationTime'] = datetime.datetime.strptime(creationTime, '%Y:%m:%d %H:%M:%S').replace(
                tzinfo=datetime.timezone.utc).strftime('%Y-%m-%d %H:%M:%S') if creationTime else None
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
            logging.debug(f'extracted metadata for user {mediaitem_user_id} mediaitem {mediaitem_id}: {result}')
    except Exception as e:
        logging.error(
            f'error extracting exif metadata for user {mediaitem_user_id} mediaitem {mediaitem_id}: {str(e)}')
        result['status'] = 'FAILED'
        grpc_save_mediaitem_metadata(api_stub, result)
        return None

    if result['type'] == 'photo':
        # generate and upload preview and thumbnail for a photo
        try:
            preview_bytes, thumbnail_bytes = generate_photo_preview_and_thumbnail(
                file_path, result['mimeType'])
            preview_url = storage.upload(
                mediaitem_id, preview_bytes, 'previews')
            thumbnail_url = storage.upload(
                mediaitem_id, thumbnail_bytes, 'thumbnails')
            result['previewUrl'] = preview_url
            result['thumbnailUrl'] = thumbnail_url
            logging.debug(f'extracted preview and thumbnail for \
                        user {mediaitem_user_id} photo mediaitem {mediaitem_id}: {result}')
        except Exception as e:
            logging.error(
                f'error generating and uploading preview and thumbnail for \
                    user {mediaitem_user_id} photo mediaitem {mediaitem_id}: {str(e)}')
            result['status'] = 'FAILED'
            grpc_save_mediaitem_metadata(api_stub, result)
            return None
    else:
        # generate and upload preview and thumbnail for a video
        try:
            preview_bytes, thumbnail_bytes = generate_video_preview_and_thumbnail(file_path)
            preview_url = storage.upload(
                mediaitem_id, preview_bytes, 'previews')
            thumbnail_url = storage.upload(
                mediaitem_id, thumbnail_bytes, 'thumbnails')
            result['previewUrl'] = preview_url
            result['thumbnailUrl'] = thumbnail_url
            logging.debug(f'extracted preview and thumbnail for \
                        user {mediaitem_user_id} video mediaitem {mediaitem_id}: {result}')
        except Exception as e:
            logging.error(
                f'error generating and uploading preview and thumbnail for \
                    user {mediaitem_user_id} video mediaitem {mediaitem_id}: {str(e)}')
            result['status'] = 'FAILED'
            grpc_save_mediaitem_metadata(api_stub, result)
            return None

    # metadata, preview and thumbnail extracted, so its READY
    result['status'] = 'READY'
    grpc_save_mediaitem_metadata(api_stub, result)
    logging.info(f'processed metadata for user {mediaitem_user_id} mediaitem {mediaitem_id}')

    # send mediaitem for further component processing
    # work(omkar): find correct way to do for ML: things, people
    if 'latitude' in result and 'longitude' in result:
        asyncio.ensure_future(process_place(api_stub, mediaitem_user_id, mediaitem_id,
                                            [result['latitude'], result['longitude']]))
    return None


def generate_photo_thumbnail(preview_bytes: bytes):
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


def generate_photo_preview_and_thumbnail(original_file_path: str, mime_type: str):
    """Generate preview and thumbnail image for a photo"""
    if mime_type in PREVIEWABLE_PHOTO_MIME_TYPES:
        with WandImage(filename=original_file_path) as original:
            with original.convert('jpeg') as converted:
                img_bytes = converted.make_blob('jpeg')
    else:
        with rawpy.imread(original_file_path) as raw:
            rgb = raw.postprocess(use_camera_wb=True)
            img = PILImage.fromarray(rgb)
            img_bytes = io.BytesIO()
            img.save(img_bytes, format='JPEG')
    return img_bytes, generate_photo_thumbnail(img_bytes)

def generate_video_thumbnail(preview_video_path: str):
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

def generate_video_preview_and_thumbnail(original_file_path: str):
    """Generate preview and thumbnail image for a video"""
    video = VideoFileClip(original_file_path)
    video_preview_path = f'{original_file_path}.mp4'
    video.write_videofile(video_preview_path, codec='libx264', logger=None, verbose=False)
    with open(video_preview_path, 'rb') as video_file:
        video_bytes = video_file.read()
        video_file.close()
    video_thumbnail_bytes = generate_video_thumbnail(video_preview_path)
    os.remove(video_preview_path)
    return video_bytes, video_thumbnail_bytes
