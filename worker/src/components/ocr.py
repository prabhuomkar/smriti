"""OCR Component"""
import logging

from src.protos.api_pb2_grpc import APIStub
from src.components.component import Component
from src.providers.ocr.utils import init_ocr


class OCR(Component):
    """OCR Component"""
    def __init__(self, api_stub: APIStub, source: str, params: list[str]) -> None:
        super().__init__('ocr', api_stub)
        self.model = init_ocr(source, params)

    async def process(self, mediaitem_user_id: str, mediaitem_id: str, _: str, metadata: dict) -> None:
        """Process ocr from mediaitem"""
        if metadata is None or 'previewPath' not in metadata or ('type' in metadata and metadata['type'] == 'video'):
            return metadata
        try:
            result = self.model.extract(mediaitem_user_id, mediaitem_id, metadata['previewPath'])
            logging.debug(f'extracted ocr for user {mediaitem_user_id} mediaitem {mediaitem_id}: {result}')
            if result:
                ocr_keywords = ' '.join(result['words'] if 'words' in result and result['words'] else '').lower()
                if 'keywords' not in metadata or metadata['keywords'] == '':
                    metadata['keywords'] = ocr_keywords
                else:
                    metadata['keywords'] += (' ' + ocr_keywords)
                metadata['keywords'] = metadata['keywords'].strip()
        except Exception as exp:
            logging.error(f'error getting ocr response for user {mediaitem_user_id} '+
                          f'mediaitem {mediaitem_id}: {str(exp)}')
        logging.info(f'processed ocr for user {mediaitem_user_id} mediaitem {mediaitem_id}')
        return metadata
