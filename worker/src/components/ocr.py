"""OCR Component"""
import logging

from src.protos.api_pb2_grpc import APIStub
from src.protos.worker_pb2 import MediaItemComponent, OCR as OCRComponent  # pylint: disable=no-name-in-module
from src.components.component import Component
from src.providers.ocr.utils import init_ocr


class OCR(Component):
    """OCR Component"""
    def __init__(self, api_stub: APIStub, source: str, params: dict) -> None:
        super().__init__(MediaItemComponent.Name(OCRComponent), api_stub)
        self.model = init_ocr(source, params)

    async def process(self, mediaitem_user_id: str, mediaitem_id: str, _: str, metadata: dict) -> None:
        """Process ocr from mediaitem"""
        if metadata is None or 'previewPath' not in metadata:
            return metadata
        try:
            result = self.model.extract(mediaitem_user_id, mediaitem_id,
                                        metadata['type'] if 'type' in metadata else 'photo',
                                        metadata['previewPath'])
            logging.debug(f'extracted ocr for user {mediaitem_user_id} mediaitem {mediaitem_id}: {result}')
            if result:
                ocr_keywords = ' '.join(result['words'] if 'words' in result and result['words'] else '').lower()
                if 'keywords' not in metadata or metadata['keywords'] == '':
                    metadata['keywords'] = ' '.join(list(set(ocr_keywords.split(' '))))
                else:
                    for keyword in ocr_keywords.split(' '):
                        if keyword not in metadata['keywords']:
                            metadata['keywords'] += (' ' + keyword)
                metadata['keywords'] = metadata['keywords'].strip()
        except Exception as exp:
            logging.error('error getting ocr response for user %s mediaitem %s: %s',
                          mediaitem_user_id, mediaitem_id, exp)
        logging.info(f'processed ocr for user {mediaitem_user_id} mediaitem {mediaitem_id}')
        return metadata
