"""Utils Tests"""
from src.utils import getval_from_dict


def test_getval_from_dict():
    data = dict({'str_key': 'value', 'int_key': 42, 'float_key': 3.142})
    assert getval_from_dict(data, ['str_key']) == 'value'
    assert getval_from_dict(data, ['int_key'], 'int') == 42
    assert getval_from_dict(data, ['float_key'], 'float') == 3.142
    assert getval_from_dict(data, ['invalid']) == None
