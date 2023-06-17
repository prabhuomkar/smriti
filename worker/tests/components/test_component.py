"""Components Component Test"""
import pytest

from src.components import Component


@pytest.mark.asyncio
async def test_component_process_failed_exception():
    with pytest.raises(NotImplementedError):
        result = await Component('name', None).process('mediaitem_user_id', 'mediaitem_id', 'mediaitem_file_path', {})
        assert result == None
