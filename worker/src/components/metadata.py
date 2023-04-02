"""Metadata Component"""
import io
import logging
import datetime

import exiftool
import rawpy
from PIL import Image as PILImage
from wand.image import Image as WandImage

from src.components.utils import grpc_save_mediaitem_metadata, getval_from_dict
from src.protos.api_pb2_grpc import APIStub
from src.protos.api_pb2 import MediaItemMetadataRequest  # pylint: disable=no-name-in-module


PREVIEWABLE_PHOTO_MIME_TYPES = [
    'image/avif', 'image/bmp', 'image/gif', 'image/vnd.microsoft.icon', 'image/x-icon',
    'image/icon', 'image/jpeg', 'image/x-citrix-jpeg', 'image/pjpeg',
    'image/apng', 'image/x-png', 'image/x-citrix-png', 'image/png', 'image/tiff',
    'image/svg+xml', 'image/webp', 'image/heic', 'image/heif',
]

async def process_metadata(storage, api_stub: APIStub, mediaitem_id: str) -> None:
    """Process required metadata and generate thumbnail from EXIF data"""
    file_path = storage.get(mediaitem_id)

    # extract metadata
    result = {}
    result['id'] = mediaitem_id
    result['status'] = 'UNSPECIFIED'
    result['sourceUrl'] = file_path
    try:
        with exiftool.ExifToolHelper() as et:
            metadata = et.get_metadata(file_path)[0]
            logging.debug(f'metadata for mediaitem {mediaitem_id}: {metadata}')
            # work(omkar): do this is in cleaner way
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
            logging.debug(f'extracted metadata for mediaitem {mediaitem_id}: {result}')
    except Exception as e:
        logging.error(
            f'error extracting exif metadata for mediaitem {mediaitem_id}: {str(e)}')
        result['status'] = 'FAILED'
        mediaItemMetadataRequest = MediaItemMetadataRequest(
            id=result['id'],
            status=result['status'],
            sourceUrl=result['sourceUrl'],
        )
        grpc_save_mediaitem_metadata(api_stub, mediaItemMetadataRequest)
        return None

    # generate and upload preview and thumbnail
    try:
        preview_bytes, thumbnail_bytes = generate_preview_and_thumbnail(
            file_path, result['mimeType'])
        preview_url = storage.upload(
            mediaitem_id, preview_bytes, 'previews')
        thumbnail_url = storage.upload(
            mediaitem_id, thumbnail_bytes, 'thumbnails')
        result['previewUrl'] = preview_url
        result['thumbnailUrl'] = thumbnail_url
        logging.debug(f'extracted preview and thumbnail for mediaitem {mediaitem_id}: {result}')
    except Exception as e:
        logging.error(
            f'error generating and uploading preview and thumbnail for mediaitem {mediaitem_id}: {str(e)}')
        result['status'] = 'FAILED'
        mediaItemMetadataRequest = MediaItemMetadataRequest(
            id=result['id'],
            status=result['status'],
            mimeType=result['mimeType'] if 'mimeType' in result else None,
            sourceUrl=result['sourceUrl'] if 'sourceUrl' in result else None,
            type=result['type'] if 'type' in result else None,
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
        grpc_save_mediaitem_metadata(api_stub, mediaItemMetadataRequest)
        return None

    result['status'] = 'READY'
    mediaItemMetadataRequest = MediaItemMetadataRequest(
        id=result['id'],
        status=result['status'],
        mimeType=result['mimeType'] if 'mimeType' in result else None,
        sourceUrl=result['sourceUrl'] if 'sourceUrl' in result else None,
        previewUrl=result['previewUrl'] if 'previewUrl' in result else None,
        thumbnailUrl=result['thumbnailUrl'] if 'thumbnailUrl' in result else None,
        type=result['type'] if 'type' in result else None,
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
    grpc_save_mediaitem_metadata(api_stub, mediaItemMetadataRequest)
    logging.info(f'processed metadata for mediaitem {mediaitem_id}')


def generate_thumbnail(preview_bytes: bytes):
    """Generate thumbnail image"""
    # work(omkar): thumbnail size should be configurable through UI
    with WandImage(blob=preview_bytes) as img:
        if img.size[0] > img.size[1]:
            wpercent = 512/float(img.size[0])
            hsize = int((float(img.size[1])*float(wpercent)))
            img.resize(512, hsize)
        else:
            hpercent = 512/float(img.size[1])
            wsize = int((float(img.size[0])*float(hpercent)))
            img.resize(wsize, 512)
        with img.convert('jpeg') as converted:
            return converted.make_blob('jpeg')


def generate_preview_and_thumbnail(original_file_path: str, mime_type: str):
    """Generate preview and thumbnail image"""
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
    return img_bytes, generate_thumbnail(img_bytes)
