import time
import os
import re
from behave import *
import requests
import datetime

from common import API_URL


def get_exif(url: str) -> dict:
    res = requests.get(url)
    data = res.text
    width, height, make, model, creation_time = None, None, None, None, None
    subtype = ''
    for item in data.split('\n'):
        if '.SubImage1.NewSubfileType' in item and item.split(maxsplit=1)[1] == 'Primary image':
            subtype = '.SubImage1'
        if '.SubImage2.NewSubfileType' in item and item.split(maxsplit=1)[1] == 'Primary image':
            subtype = '.SubImage2'
        if ('.ImageWidth' in item and subtype != '' and f'{subtype}.ImageWidth' in item):
            width = int(item.split()[1])
        elif (('.ImageHeight' in item and subtype != '' and f'{subtype}.ImageHeight' in item) or ('.ImageLength' in item and subtype != '' and f'{subtype}.ImageLength' in item)):
            height = int(item.split()[1])
        elif '.Make' in item and not 'Maker' in item:
            make = item.split(maxsplit=1)[1].strip()
        elif '.Model' in item:
            model = item.split(maxsplit=1)[1].strip()
        elif '.DateTime' in item:
            creation_time = (item.split()[1]+' '+item.split()[2])
            creation_time = datetime.datetime.strptime(creation_time, '%Y:%m:%d %H:%M:%S').replace(
                    tzinfo=datetime.timezone.utc).strftime('%Y-%m-%dT%H:%M:%SZ')
    if width == None or height == None:
        for item in data.split('\n'):
            if '.ImageWidth' in item:
                width = int(item.split()[1])
            elif '.ImageHeight' in item or '.ImageLength' in item:
                height = int(item.split()[1])
            if width is not None and height is not None:
                break
    return {'width': width, 'height': height, 'cameraMake': make, 'cameraModel': model,
            'creationTime': creation_time}

def download_upload_remove(mediaitem):
    file_name = f'/tmp/{mediaitem["source"].split("/")[-1]}'.lower()
    # upload
    headers = {'Authorization': f'Bearer {mediaitem["access_token"]}'}
    files = {'file': open(file_name, 'rb')}
    res = requests.post(API_URL+'/v1/mediaItems', files=files, headers=headers)
    assert res.status_code == 201
    res = res.json()
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
    for mediaitem_id, upload_mediaitem in zip(context.upload_mediaitem_ids, context.upload_mediaitems):
        while True:
            headers = {'Authorization': f'Bearer {context.access_token}'}
            res = requests.get(API_URL+'/v1/mediaItems/'+mediaitem_id, headers=headers)
            res = res.json()
            if res['status'] == 'READY':
                mediaitem_exif = upload_mediaitem['exif']
                print(res, mediaitem_exif)
                assert res['id'] == mediaitem_id
                assert res['width'] == mediaitem_exif['width']
                assert res['height'] == mediaitem_exif['height']
                assert res['cameraMake'] == mediaitem_exif['cameraMake']
                assert res['cameraModel'] == mediaitem_exif['cameraModel']
                assert res['creationTime'] == mediaitem_exif['creationTime']
                break
            if res['status'] == 'FAILED':
                raise Exception(f'failed to process mediaitem: {mediaitem_id}, response: {res}')
            time.sleep(2)

