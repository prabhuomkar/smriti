from behave import *
import requests
from requests.auth import HTTPBasicAuth

from common import API_URL, ADMIN_USERNAME, ADMIN_PASSWORD, CREATED_USER, CREATED_ALBUM, UPDATED_ALBUM


@given('a user is created if does not exist')
def step_impl(context):
    res = requests.get(API_URL+'/v1/users',
                        auth=HTTPBasicAuth(ADMIN_USERNAME, ADMIN_PASSWORD))
    assert res.status_code == 200
    users = res.json()
    if len(users) != 0:
        context.album_id = users[0]['id']
    else:
        res = requests.post(API_URL+'/v1/users', json=CREATED_USER,
                            auth=HTTPBasicAuth(ADMIN_USERNAME, ADMIN_PASSWORD))
        assert res.status_code == 201
        context.user_id = res.json()['id']


@given('there are no albums')
def step_impl(context):
    res = requests.get(API_URL+'/v1/albums',
                       headers={'Authorization': f'Bearer {context.access_token}'})
    albums = res.json()
    assert len(albums) == 0


@given('an album exists')
def step_impl(context):
    res = requests.get(API_URL+'/v1/albums',
                       headers={'Authorization': f'Bearer {context.access_token}'})
    albums = res.json()
    assert len(albums) == 1
    context.album_id = albums[0]['id']


@when('get all albums {condition} auth')
def step_impl(context, condition):
    headers = None
    if condition == 'with':
        headers = {'Authorization': f'Bearer {context.access_token}'}
    res = requests.get(API_URL+'/v1/albums', headers=headers)
    context.response = res
    context.albums = res.json()


@when('get album {condition} auth')
def step_impl(context, condition):
    headers = None
    if condition == 'with':
        headers = {'Authorization': f'Bearer {context.access_token}'}
    album_id = context.album_id
    res = requests.get(API_URL+'/v1/albums/'+album_id, headers=headers)
    context.response = res
    context.album = res.json()


@then('album is present in list')
def step_impl(context):
    assert len(context.albums) == 1
    assert context.albums[0]['name'] == context.match_album['name']
    assert context.albums[0]['description'] == context.match_album['description']


@then('album is present')
def step_impl(context):
    assert context.album['name'] == context.match_album['name']
    assert context.album['description'] == context.match_album['description']


@then('album is not present in list')
def step_impl(context):
    if len(context.albums) > 0:
        assert context.albums[0]['name'] != context.match_album['name']
        assert context.albums[0]['description'] != context.match_album['description']


@then('album is not present')
def step_impl(context):
    assert 'name' not in context.album
    assert 'description' not in context.album
    assert context.album['message'] == 'album not found'


@when('create album {condition} auth')
def step_impl(context, condition):
    headers = None
    if condition == 'with':
        headers = {'Authorization': f'Bearer {context.access_token}'}
    res = requests.post(API_URL+'/v1/albums', json=CREATED_ALBUM, headers=headers)
    context.response = res


@when('update album {condition} auth')
def step_impl(context, condition):
    headers = None
    if condition == 'with':
        headers = {'Authorization': f'Bearer {context.access_token}'}
    album_id = context.album_id
    res = requests.put(API_URL+'/v1/albums/'+album_id,
                       json=UPDATED_ALBUM, headers=headers)
    context.response = res


@when('delete album {condition} auth')
def step_impl(context, condition):
    headers = None
    if condition == 'with':
        headers = {'Authorization': f'Bearer {context.access_token}'}
    album_id = context.album_id
    res = requests.delete(API_URL+'/v1/albums/'+album_id, headers=headers)
    context.response = res


@then('album is created')
def step_impl(context):
    assert context.response.status_code == 201
    context.album_id = context.response.json()['id']
    context.match_album = CREATED_ALBUM


@then('album is updated')
def step_impl(context):
    assert context.response.status_code == 204
    context.match_album = UPDATED_ALBUM


@then('album is deleted')
def step_impl(context):
    assert context.response.status_code == 204
    context.match_album = UPDATED_ALBUM

