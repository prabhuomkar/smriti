import time
import os
import re
from behave import *
import requests
from multiprocessing import Pool

from common import API_URL


def get_exif(url: str) -> dict:
    res = requests.get(url)
    data = res.text
    width, height, make = None, None, None
    for item in data.split('\n'):
        if '.ImageWidth' in item:
            width = item.split()[1]
        elif '.ImageHeight' in item:
            height = item.split()[1]
        elif '.Make' in item and not 'Maker' in item:
            make = item.split()[1]
        if width is not None and height is not None and make is not None:
            break
    return {'width': width, 'height': height, 'make': make}

def download_upload_remove(mediaitem):
    # download
    file_name = f'/tmp/{mediaitem["source"].split("/")[-1]}'
    response = requests.get(mediaitem["source"])
    with open(file_name, 'wb') as file:
        file.write(response.content)
    # upload
    headers = {'Authorization': f'Bearer {mediaitem["access_token"]}'}
    files = {'file': open(file_name, 'rb')}
    res = requests.post(API_URL+'/v1/mediaItems', files=files, headers=headers)
    assert res.status_code == 201
    res = res.json()
    # remove
    os.remove(file_name)
    return res['id']

@given('get list of raw mediaitems to upload')
def step_impl(context):
    res = requests.get(f'https://raw.pixls.us/json/getrepository.php?set=all&_={int(time.time()*1000)}')
    data = res.json()['data']
    cameras = [row['camera'] for row in context.table]
    context.upload_mediaitems = []
    for mediaitem in data:
        if mediaitem[0].lower() in cameras:
            context.upload_mediaitems.append({'access_token': context.access_token,
                                              'source': re.findall(r"href='([^']*)'", mediaitem[7])[0],
                                              'exif': get_exif(re.findall(r"href='([^']*)'", mediaitem[8])[0])})

@when('upload raw mediaitems')
def step_impl(context):
    context.upload_mediaitem_ids = []
    for mediaitem in context.upload_mediaitems:
        mediaitem_id = download_upload_remove(mediaitem)
        context.upload_mediaitem_ids.append(mediaitem_id)

@then('get raw mediaitems with auth and validate it is present')
def step_impl(context):
    raise Exception(context.result)
