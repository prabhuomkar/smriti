"""Utils"""

def getval_from_dict(data, keys: list[str], return_type: str = 'str') -> str | int | float | None:
    """Get possible value from dictionary"""
    for key in keys:
        if key in data and data[key] != '':
            if return_type == 'int':
                return int(data[key])
            if return_type == 'float':
                return float(data[key])
            return str(data[key])
    return None
