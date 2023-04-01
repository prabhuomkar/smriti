"""Metadata Component"""
import io
import logging
import datetime

import exiftool
import rawpy
from PIL import Image as PILImage
from wand.image import Image as WandImage

from src.protos.api_pb2 import MediaItemMetadataRequest  # pylint: disable=no-name-in-module


PREVIEWABLE_PHOTO_MIME_TYPES = [
    'image/bmp', 'image/gif', 'image/vnd.microsoft.icon', 'image/x-icon', 'image/ico',
    'image/icon', 'text/ico', 'application/ico', 'image/jpeg', 'image/apng',
    'image/png', 'image/tiff', 'image/svg+xml', 'image/webp', 'image/avif',
]


async def process_metadata(storage, api_stub, mediaitem_id: str) -> None:
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
            result['mimeType'] = metadata['File:MIMEType'] if 'File:MIMEType' in metadata else None
            result['type'] = 'photo' if 'image' in metadata['File:MIMEType'] else \
                'video' if 'video' in metadata['File:MIMEType'] else None
            result['width'] = metadata['File:ImageWidth'] if 'File:ImageWidth' in metadata else \
                metadata['EXIF:ExifImageWidth'] if 'EXIF:ExifImageWidth' in metadata else \
                metadata['QuickTime:ImageWidth'] if 'QuickTime:ImageWidth' in metadata else \
                metadata['QuickTime:SourceImageWidth'] if 'QuickTime:SourceImageWidth' in metadata else None
            result['height'] = metadata['File:ImageHeight'] if 'File:ImageHeight' in metadata else \
                metadata['EXIF:ExifImageHeight'] if 'EXIF:ExifImageHeight' in metadata else \
                metadata['QuickTime:ImageHeight'] if 'QuickTime:ImageHeight' in metadata else \
                metadata['QuickTime:SourceImageHeight'] if 'QuickTime:SourceImageHeight' in metadata else None
            result['creationTime'] = metadata['EXIF:DateTimeOriginal'] if 'EXIF:DateTimeOriginal' in metadata else \
                metadata['EXIF:CreateDate'] if 'EXIF:CreateDate' in metadata else \
                metadata['XMP:CreateDate'] if 'XMP:CreateDate' in metadata else \
                metadata['XMP:DateCreated'] if 'XMP:DateCreated' in metadata else \
                metadata['Composite:SubSecCreateDate'] if 'Composite:SubSecCreateDate' in metadata else \
                metadata['Composite:SubSecDateTimeOriginal'] if 'Composite:SubSecDateTimeOriginal' in metadata else \
                metadata['QuickTime:CreateDate'] if 'QuickTime:CreateDate' in metadata else \
                metadata['QuickTime:TrackModifyDate'] if 'QuickTime:TrackModifyDate' in metadata else \
                metadata['QuickTime:MediaCreateDate'] if 'QuickTime:MediaCreateDate' in metadata else \
                metadata['QuickTime:CreationDate'] if 'QuickTime:CreationDate' in metadata else \
                metadata['File:FileModifyDate'] if 'File:FileModifyDate' in metadata else \
                metadata['File:FileAccessDate'] if 'File:FileAccessDate' in metadata else \
                metadata['File:FileInodeChangeDate'] if 'File:FileInodeChangeDate' in metadata else None
            # work(omkar): handle timezone when "its time" :P
            creationTime = result['creationTime'].split("+")[0] if result['creationTime'] else None
            result['creationTime'] = datetime.datetime.strptime(creationTime, '%Y:%m:%d %H:%M:%S').replace(
                tzinfo=datetime.timezone.utc).strftime('%Y-%m-%d %H:%M:%S') if creationTime else None
            result['cameraMake'] = metadata['EXIF:Make'] if 'EXIF:Make' in metadata else \
                metadata['QuickTime:Make'] if 'QuickTime:Make' in metadata else None
            result['cameraModel'] = metadata['EXIF:Model'] if 'EXIF:Model' in metadata else \
                metadata['QuickTime:Model'] if 'QuickTime:Model' in metadata else None
            result['focalLength'] = str(
                metadata['EXIF:FocalLength']) if 'EXIF:FocalLength' in metadata else None
            result['apertureFNumber'] = str(
                metadata['EXIF:FNumber']) if 'EXIF:FNumber' in metadata else None
            result['isoEquivalent'] = str(
                metadata['EXIF:ISO']) if 'EXIF:ISO' in metadata else None
            result['exposureTime'] = str(
                metadata['EXIF:ExposureTime']) if 'EXIF:ExposureTime' in metadata else None
            result['fps'] = str(
                metadata['QuickTime:VideoFrameRate']) if 'QuickTime:VideoFrameRate' in metadata else None
            result['latitude'] = float(metadata['EXIF:GPSLatitude']) if 'EXIF:GPSLatitude' in metadata else \
                float(metadata['Composite:GPSLatitude']
                      ) if 'Composite:GPSLatitude' in metadata else None
            result['longitude'] = float(metadata['EXIF:GPSLongitude']) if 'EXIF:GPSLongitude' in metadata else \
                float(metadata['Composite:GPSLongitude']
                      ) if 'Composite:GPSLongitude' in metadata else None
            logging.debug(f'extracted metadata for mediaitem {mediaitem_id}: {result}')
    except Exception as e:
        logging.error(
            f'error extracting exif metadata for mediaitem {mediaitem_id}: {str(e)}')
        result['status'] = 'FAILED'
        _ = api_stub.SaveMediaItemMetadata(MediaItemMetadataRequest(
            id=result['id'],
            status=result['status']
        ))
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
        _ = api_stub.SaveMediaItemMetadata(MediaItemMetadataRequest(
            id=result['id'],
            status=result['status']
        ))
        return None

    result['status'] = 'READY'
    _ = api_stub.SaveMediaItemMetadata(MediaItemMetadataRequest(
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
    ))
    logging.info(f'processed metadata for mediaitem {mediaitem_id}')


def generate_thumbnail(preview_bytes: bytes):
    """Generate thumbnail image"""
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
