"""Metadata Component"""
import io
import logging

import exiftool
import rawpy
from PIL import Image as PILImage
from wand.image import Image as WandImage

PREVIEWABLE_PHOTO_MIME_TYPES = [
    'image/bmp', 'image/gif', 'image/vnd.microsoft.icon', 'image/x-icon', 'image/ico',
    'image/icon', 'text/ico', 'application/ico', 'image/jpeg', 'image/apng',
    'image/png', 'image/tiff', 'image/svg+xml', 'image/webp', 'image/avif',
]


def process_metadata(storage, id: str) -> dict:  # pylint: disable=redefined-builtin
    """Process required metadata and generate thumbnail from EXIF data"""
    file_path = storage.get(id)

    # extract metadata
    result = {}
    try:
        with exiftool.ExifToolHelper() as et:
            metadata = et.get_metadata(file_path)[0]
            logging.debug(f'metadata for mediaitem {id}: {metadata}')
            result['mimeType'] = metadata['File:MIMEType']
            result['type'] = 'image' if 'image' in metadata['File:MIMEType'] else \
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
                metadata['QuickTime:CreationDate'] if 'QuickTime:CreationDate' in metadata else None
            result['cameraMake'] = metadata['EXIF:Make'] if 'EXIF:Make' in metadata else \
                metadata['QuickTime:Make'] if 'QuickTime:Make' in metadata else None
            result['cameraModel'] = metadata['EXIF:Model'] if 'EXIF:Model' in metadata else \
                metadata['QuickTime:Model'] if 'QuickTime:Model' in metadata else None
            result['focalLength'] = metadata['EXIF:FocalLength'] if 'EXIF:FocalLength' in metadata else None
            result['apertureFNumber'] = metadata['EXIF:FNumber'] if 'EXIF:FNumber' in metadata else None
            result['isoEquivalent'] = metadata['EXIF:ISO'] if 'EXIF:ISO' in metadata else None
            result['exposureTime'] = metadata['EXIF:ExposureTime'] if 'EXIF:ExposureTime' in metadata else None
            result['fps'] = metadata['QuickTime:VideoFrameRate'] if 'QuickTime:VideoFrameRate' in metadata else None
            result['latitude'] = metadata['EXIF:GPSLatitude'] if 'EXIF:GPSLatitude' in metadata else \
                metadata['Composite:GPSLatitude'] if 'Composite:GPSLatitude' in metadata else None
            result['longitude'] = metadata['EXIF:GPSLongitude'] if 'EXIF:GPSLongitude' in metadata else \
                metadata['Composite:GPSLongitude'] if 'Composite:GPSLongitude' in metadata else None
    except Exception as e:
        logging.error(
            f'error extracting exif metadata for mediaitem {id}: {str(e)}')

    # generate preview and thumbnail
    try:
        preview_bytes, thumbnail_bytes = generate_preview_and_thumbnail(
            file_path, result['mimeType'])
    except Exception as e:
        logging.error(
            f'error generating preview and thumbnail for mediaitem {id}: {str(e)}')

    # upload preview url and thumbnail url
    preview_url = storage.upload(id, 0, preview_bytes, 'previews')
    thumbnail_url = storage.upload(id, 0, thumbnail_bytes, 'thumbnails')
    result['previewUrl'] = preview_url
    result['thumbnailUrl'] = thumbnail_url

    return result


def generate_thumbnail(preview_bytes: bytes):
    """Generate thumbnail image"""
    with WandImage(blob=preview_bytes) as img:
        if img.size[0] > img.size[1]:
            wpercent = (512/float(img.size[0]))
            hsize = int((float(img.size[1])*float(wpercent)))
            img.resize(512, hsize)
        else:
            hpercent = (512/float(img.size[1]))
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
