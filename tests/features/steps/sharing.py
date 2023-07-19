from behave import *
import requests

from common import API_URL, CREATED_SHARED_ALBUM

@given('there are no shared albums')
def step_impl(context):
    res = requests.get(API_URL+'/v1/albums?shared=true',
                       headers={'Authorization': f'Bearer {context.access_token}'})
    albums = res.json()
    assert len(albums) == 0

@given('a shared album exists')
def step_impl(context):
    res = requests.get(API_URL+'/v1/albums?shared=true',
                       headers={'Authorization': f'Bearer {context.access_token}'})
    albums = res.json()
    assert len(albums) == 1
    context.album_id = albums[0]['id']

@when('get all shared albums {condition} auth')
def step_impl(context, condition):
    headers = None
    if condition == 'with':
        headers = {'Authorization': f'Bearer {context.access_token}'}
    res = requests.get(API_URL+'/v1/albums?shared=true', headers=headers)
    context.response = res
    context.albums = res.json()

@when('get shared album mediaitems {condition} auth')
def step_impl(context, condition):
    headers = None
    if condition == 'with':
        headers = {'Authorization': f'Bearer {context.access_token}'}
    album_id = context.album_id
    res = requests.get(API_URL+'/v1/sharing/'+album_id+'/mediaItems', headers=headers)
    context.response = res
    context.album_mediaitems = res.json()

@when('get shared album {condition} auth')
def step_impl(context, condition):
    headers = None
    if condition == 'with':
        headers = {'Authorization': f'Bearer {context.access_token}'}
    album_id = context.album_id
    res = requests.get(API_URL+'/v1/sharing/'+album_id, headers=headers)
    context.response = res
    context.album = res.json()

@then('shared album is present in list')
def step_impl(context):
    assert len(context.albums) == 1
    for field in context.match_album:
        assert context.albums[0][field] == context.match_album[field]

@then('shared album is present')
def step_impl(context):
    for field in context.match_album:
        assert context.album[field] == context.match_album[field]

@then('shared album is not present in list')
def step_impl(context):
    if len(context.albums) > 0:
        for field in context.match_album:
            assert context.albums[0][field] != context.match_album[field]

@then('shared album is not present')
def step_impl(context):
    for field in context.match_album:
        assert field not in context.album
    assert context.album['message'] == 'shared link not found'

@when('create shared album {condition} auth')
def step_impl(context, condition):
    headers = None
    if condition == 'with':
        headers = {'Authorization': f'Bearer {context.access_token}'}
    res = requests.post(API_URL+'/v1/albums', json=CREATED_SHARED_ALBUM, headers=headers)
    context.response = res

@when('delete shared album {condition} auth')
def step_impl(context, condition):
    headers = None
    if condition == 'with':
        headers = {'Authorization': f'Bearer {context.access_token}'}
    album_id = context.album_id
    res = requests.delete(API_URL+'/v1/albums/'+album_id, headers=headers)
    context.response = res

@when('add shared album mediaitems {condition} auth')
def step_impl(context, condition):
    headers = None
    if condition == 'with':
        headers = {'Authorization': f'Bearer {context.access_token}'}
    album_id = context.album_id
    res = requests.post(API_URL+'/v1/albums/'+album_id+'/mediaItems', json={'mediaItems': [context.mediaitem_id]}, headers=headers)
    context.response = res

@when('remove shared album mediaitems {condition} auth')
def step_impl(context, condition):
    headers = None
    if condition == 'with':
        headers = {'Authorization': f'Bearer {context.access_token}'}
    album_id = context.album_id
    res = requests.delete(API_URL+'/v1/albums/'+album_id+'/mediaItems', json={'mediaItems': [context.mediaitem_id]}, headers=headers)
    context.response = res

@then('shared album mediaitems are {type}')
def step_impl(context, type):
    if type == 'added' or type == 'removed':
        assert context.response.status_code == 204
    elif type == 'present':
        assert len(context.album_mediaitems) == 1
        assert context.album_mediaitems[0]['id'] == context.mediaitem_id
    elif type == 'absent':
        assert len(context.album_mediaitems) == 0

@then('shared album is updated after {type} album mediaitems')
def step_impl(context, type):
    if type == 'add':
        assert context.album['mediaItemsCount'] == 1
        assert context.album['coverMediaItem']['id'] == context.mediaitem_id
    elif type =='remove':
        print(context.album)
        assert context.album['mediaItemsCount'] == 0
        assert context.album['coverMediaItem'] == None

@then('shared album is created')
def step_impl(context):
    assert context.response.status_code == 201
    context.album_id = context.response.json()['id']
    context.match_album = CREATED_SHARED_ALBUM

@then('shared album is deleted')
def step_impl(context):
    assert context.response.status_code == 204
    context.match_album = CREATED_SHARED_ALBUM
