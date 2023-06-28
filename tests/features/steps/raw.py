import time
import os
import re
import json
from behave import *
import requests
import datetime

from common import API_URL, FILES_TO_SKIP


def get_exif(url: str) -> dict:
    res = requests.get(url)
    data = res.text
    width, height, make, model, creation_time = None, None, None, None, None
    subtype = ''
    for item in data.split('\n'):
        if '.SubImage1.NewSubfileType' in item and item.split(maxsplit=1)[1] == 'Primary image' and subtype == '':
            subtype = '.SubImage1'
        if '.SubImage2.NewSubfileType' in item and item.split(maxsplit=1)[1] == 'Primary image' and subtype == '':
            subtype = '.SubImage2'
        if ('.ImageWidth' in item and subtype != '' and f'{subtype}.ImageWidth' in item):
            width = int(item.split()[1])
        elif (('.ImageHeight' in item and subtype != '' and f'{subtype}.ImageHeight' in item) or ('.ImageLength' in item and subtype != '' and f'{subtype}.ImageLength' in item)):
            height = int(item.split()[1])
        elif '.Make' in item and not 'Maker' in item and make is None:
            make = item.split(maxsplit=1)[1].strip()
        elif '.Model' in item and model is None:
            model = item.split(maxsplit=1)[1].strip()
        elif '.DateTime' in item:
            creation_time = item.split(maxsplit=1)[1].strip()
            if creation_time and re.search(r'[\+]\d{2}:\d{2}', creation_time):
                creation_time = creation_time.rsplit("+", maxsplit=1)[0]
            elif creation_time and re.search(r'[\-]\d{2}:\d{2}', creation_time):
                creation_time = creation_time.rsplit("-", maxsplit=1)[0]
            if 'T' not in creation_time and 'Z' not in creation_time:
                if '-' in creation_time:
                    creation_time = creation_time.replace(' ', '', -1)
                    creation_time = datetime.datetime.strptime(creation_time, '%Y-%m-%d%H:%M:%S').replace(
                        tzinfo=datetime.timezone.utc).strftime('%Y-%m-%dT%H:%M:%SZ')
                else:
                    creation_time = datetime.datetime.strptime(creation_time, '%Y:%m:%d %H:%M:%S').replace(
                    tzinfo=datetime.timezone.utc).strftime('%Y-%m-%dT%H:%M:%SZ')
            elif 'T' in creation_time and 'Z' not in creation_time:
                creation_time = datetime.datetime.strptime(creation_time, '%Y-%m-%dT%H:%M:%S').replace(
                        tzinfo=datetime.timezone.utc).strftime('%Y-%m-%dT%H:%M:%SZ')
    if width == None or height == None:
        for item in data.split('\n'):
            if '.ImageWidth' in item or '.SensorWidth' in item:
                width = int(item.split()[1])
            elif '.ImageHeight' in item or '.ImageLength' in item or '.SensorHeight' in item:
                height = int(item.split()[1])
            if width is not None and height is not None:
                break
    return {'width': width, 'height': height, 'cameraMake': make, 'cameraModel': model,
            'creationTime': creation_time}

def download_upload_remove(mediaitem):
    file_name = f'/tmp/{mediaitem["source"].split("/")[-3]}-{mediaitem["source"].split("/")[-1]}'.lower()
    # upload
    headers = {'Authorization': f'Bearer {mediaitem["access_token"]}'}
    files = {'file': open(file_name, 'rb')}
    res = requests.post(API_URL+'/v1/mediaItems', files=files, headers=headers)
    assert res.status_code == 201
    res = res.json()
    return res['id']

@given('get list of raw mediaitems to upload')
def step_impl(context):
    file_name = '/tmp/raw-list.json'
    data = []
    if os.path.exists(file_name):
        with open(file_name, 'r') as f:
            data = json.load(f)
    else:
        res = requests.get(f'https://raw.pixls.us/json/getrepository.php?set=all&_={int(time.time()*1000)}')
        data = res.json()['data']
        with open(file_name, 'w') as f:
            json.dump(data, f)
    cameras = [row['camera'] for row in context.table]
    context.upload_mediaitems = []
    for mediaitem in data:
        if mediaitem[0].lower() in cameras:
            source = re.findall(r"href='([^']*)'", mediaitem[7])[0]
            if f'{source.split("/")[-3]}-{source.split("/")[-1]}'.lower() not in FILES_TO_SKIP:
                context.upload_mediaitems.append({'access_token': context.access_token, 'source': source,
                                                  'exif': get_exif(re.findall(r"href='([^']*)'", mediaitem[8])[0])})

@then('raw mediaitems are ready to upload')
def step_impl(context):
    for mediaitem in context.upload_mediaitems:
        file_name = f'/tmp/{mediaitem["source"].split("/")[-3]}-{mediaitem["source"].split("/")[-1]}'.lower()
        try:
            assert os.path.exists(file_name)
            assert os.path.getsize(file_name) > 1000
        except Exception as exp:
            raise Exception(f'assertion failed for {file_name}')

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
                try:
                    mediaitem_exif = upload_mediaitem['exif']
                    assert res['id'] == mediaitem_id
                    assert res['width'] == mediaitem_exif['width']
                    assert res['height'] == mediaitem_exif['height']
                    if mediaitem_exif['cameraMake'] is not None:
                        assert 'cameraMake' in res
                        assert res['cameraMake'] == mediaitem_exif['cameraMake']
                    if mediaitem_exif['cameraModel'] is not None:
                        assert 'cameraModel' in res
                        assert res['cameraModel'] == mediaitem_exif['cameraModel']
                    if mediaitem_exif['creationTime'] is not None:
                        assert 'creationTime' in res
                        assert res['creationTime'] == mediaitem_exif['creationTime']
                    break
                except Exception as e:
                    raise Exception(f'failed assertion for mediaitem: {upload_mediaitem["source"]} response: {res} wanted: {upload_mediaitem["exif"]}')
            if res['status'] == 'FAILED':
                raise Exception(f'failed to process mediaitem: {upload_mediaitem["source"]}, response: {res} wanted: {upload_mediaitem["exif"]}')
            time.sleep(2)

