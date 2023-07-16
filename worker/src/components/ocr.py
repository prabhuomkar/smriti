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
            return None
        try:
            result = self.model.extract(mediaitem_user_id, mediaitem_id, metadata['previewPath'])
            logging.debug(f'extracted ocr for user {mediaitem_user_id} mediaitem {mediaitem_id}: {result}')

            if 'keywords' not in metadata or metadata['keywords'] == '':
                metadata['keywords'] = ' '.join(result['words']).lower()
            else:
                metadata['keywords'] += (' ' + ' '.join(result['words']).lower())
            metadata['keywords'] = metadata['keywords'].strip()
        except Exception as exp:
            logging.error(f'error getting thing response for user {mediaitem_user_id} '+
                          f'mediaitem {mediaitem_id}: {str(exp)}')
        logging.info(f'processed ocr for user {mediaitem_user_id} mediaitem {mediaitem_id}')
        return metadata
