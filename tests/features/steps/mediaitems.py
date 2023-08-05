import time
from behave import *
import requests

from common import API_URL, CREATED_MEDIAITEM, UPDATED_MEDIAITEM


@given('there are no mediaitems')
def step_impl(context):
    res = requests.get(API_URL+'/v1/mediaItems',
                       headers={'Authorization': f'Bearer {context.access_token}'})
    mediaitems = res.json()
    assert len(mediaitems) == 0

@given('a mediaitem exists')
def step_impl(context):
    res = requests.get(API_URL+'/v1/mediaItems',
                       headers={'Authorization': f'Bearer {context.access_token}'})
    mediaitems = res.json()
    if len(mediaitems) == 1:
        context.mediaitem_id = mediaitems[0]['id']
        context.mediaitem_type = mediaitems[0]['mediaItemType']
    else:
        mediaitem_ids = [mediaitem['id'] for mediaitem in mediaitems]
        assert context.mediaitem_id in mediaitem_ids
        for mediaitem in mediaitems:
            if mediaitem['id'] == context.mediaitem_id:
                context.mediaitem_type = mediaitem['mediaItemType']

@when('get all mediaitems {condition} auth')
def step_impl(context, condition):
    headers = None
    if condition == 'with':
        headers = {'Authorization': f'Bearer {context.access_token}'}
    res = requests.get(API_URL+'/v1/mediaItems', headers=headers)
    context.response = res
    context.mediaitems = res.json()

@when('get mediaitem {condition} auth')
def step_impl(context, condition):
    headers = None
    if condition == 'with':
        headers = {'Authorization': f'Bearer {context.access_token}'}
    mediaitem_id = context.mediaitem_id
    res = requests.get(API_URL+'/v1/mediaItems/'+mediaitem_id, headers=headers)
    context.response = res
    context.mediaitem = res.json()

@then('mediaitem is present in list')
def step_impl(context):
    assert len(context.mediaitems) == 1
    for field in context.match_mediaitem:
        if field == 'description':
            if 'description' in context.mediaitems[0]:
                assert context.mediaitems[0]['description'] == context.match_mediaitem['description']
        else:
            assert context.mediaitems[0][field] == context.match_mediaitem[field]

@then('mediaitem is present')
def step_impl(context):
    for field in context.match_mediaitem:
        if field == 'description':
            if 'description' in context.mediaitem:
                assert context.mediaitem['description'] == context.match_mediaitem['description']
        else:
            assert context.mediaitem[field] == context.match_mediaitem[field]
    context.mediaitem_id = context.mediaitem['id']

@then('mediaitem is not present in list')
def step_impl(context):
    print(context.mediaitems)
    assert len(context.mediaitems) == 0

@then('mediaitem is not present')
def step_impl(context):
    for field in context.match_mediaitem:
        assert field not in context.mediaitem
    assert context.mediaitem['message'] == 'mediaitem not found'

@when('upload {name} {type} mediaitem {condition} auth')
def step_impl(context, name, type, condition):
    headers = None
    if condition == 'with':
        headers = {'Authorization': f'Bearer {context.access_token}'}
    files = {'file': open(f'data/{"IMG_0543.HEIC" if name == "default" and type == "photo" else "IMG_6470.MOV" if name == "default" and type =="video" else name}','rb')}
    res = requests.post(API_URL+'/v1/mediaItems', files=files, headers=headers)
    context.response = res
    context.mediaitem_type = type

@when('upload {name} {type} mediaitem with auth if does not exist and wait {seconds} seconds')
def step_impl(context, name, type, seconds):
    headers = {'Authorization': f'Bearer {context.access_token}'}
    res = requests.get(API_URL+'/v1/mediaItems',
                       headers={'Authorization': f'Bearer {context.access_token}'})
    mediaitems = res.json()
    if len(mediaitems) == 0:
        files = {'file': open(f'data/{"IMG_0543.HEIC" if name == "default" and type == "photo" else "IMG_6470.MOV" if name == "default" and type =="video" else name}','rb')}
        res = requests.post(API_URL+'/v1/mediaItems', files=files, headers=headers)
    context.response = res
    context.mediaitem_type = type
    time.sleep(int(seconds))

@when('update mediaitem {condition} auth')
def step_impl(context, condition):
    headers = None
    if condition == 'with':
        headers = {'Authorization': f'Bearer {context.access_token}'}
    mediaitem_id = context.mediaitem_id
    res = requests.put(API_URL+'/v1/mediaItems/'+mediaitem_id,
                       json=UPDATED_MEDIAITEM[context.mediaitem_type], headers=headers)
    context.response = res

@when('delete mediaitem {condition} auth')
def step_impl(context, condition):
    headers = None
    if condition == 'with':
        headers = {'Authorization': f'Bearer {context.access_token}'}
    mediaitem_id = context.mediaitem_id
    res = requests.delete(API_URL+'/v1/mediaItems/'+mediaitem_id, headers=headers)
    context.response = res

@then('mediaitem is uploaded')
def step_impl(context):
    assert context.response.status_code == 201
    context.mediaitem_id = context.response.json()['id']
    context.match_mediaitem = CREATED_MEDIAITEM[context.mediaitem_type]
    while True:
        headers = {'Authorization': f'Bearer {context.access_token}'}
        mediaitem_id = context.mediaitem_id
        res = requests.get(API_URL+'/v1/mediaItems/'+mediaitem_id, headers=headers)
        res = res.json()
        if res['status'] == 'READY':
            break
        if res['status'] == 'FAILED':
            raise Exception('failed to upload mediaitem')
        time.sleep(2)

@then('mediaitem is uploaded or exists')
def step_impl(context):
    assert context.response.status_code in [200, 201]
    if context.response.status_code == 201:
        context.mediaitem_id = context.response.json()['id']
    else:
        context.mediaitem_id = context.response.json()[0]['id']
    while True:
        headers = {'Authorization': f'Bearer {context.access_token}'}
        mediaitem_id = context.mediaitem_id
        res = requests.get(API_URL+'/v1/mediaItems/'+mediaitem_id, headers=headers)
        res = res.json()
        if res['status'] == 'READY':
            break
        if res['status'] == 'FAILED':
            raise Exception('failed to upload mediaitem')
        time.sleep(2)
    context.match_mediaitem = CREATED_MEDIAITEM[context.mediaitem_type]

@then('mediaitem is updated')
def step_impl(context):
    assert context.response.status_code == 204
    context.match_mediaitem = UPDATED_MEDIAITEM[context.mediaitem_type]

@then('mediaitem is deleted')
def step_impl(context):
    assert context.response.status_code == 204
    context.match_mediaitem = UPDATED_MEDIAITEM[context.mediaitem_type]

