from behave import *
import requests

from common import API_URL


@when('get all {operation} mediaitems {condition} auth')
def step_impl(context, operation, condition):
    headers = None
    if condition == 'with':
        headers = {'Authorization': f'Bearer {context.access_token}'}
    res = requests.get(API_URL+'/v1/'+get_route(operation), headers=headers)
    context.response = res
    context.mediaitems = res.json()


@then('mediaitem is present in {operation} list')
def step_impl(context, operation):
    assert len(context.mediaitems) == 1
    assert context.mediaitems[0]['filename'] == context.match_mediaitem['filename']
    assert context.mediaitems[0]['mimeType'] == context.match_mediaitem['mimeType']


@then('mediaitem is not present in {operation} list')
def step_impl(context, operation):
    if len(context.mediaitems) > 0:
        assert context.mediaitems[0]['filename'] != context.match_mediaitem['filename']
        assert context.mediaitems[0]['mimeType'] != context.match_mediaitem['mimeType']


@then('mediaitem is present with marked as {operation}')
def step_impl(context, operation):
    assert context.mediaitem[get_field(operation)] == True
    context.match_mediaitem = context.mediaitem


@then('mediaitem is present with unmarked as {operation}')
def step_impl(context, operation):
    assert context.mediaitem[get_field(operation)] == False
    context.match_mediaitem = context.mediaitem


@when('mark {operation} mediaitem {condition} auth')
def step_impl(context, operation, condition):
    headers = None
    if condition == 'with':
        headers = {'Authorization': f'Bearer {context.access_token}'}
    res = requests.post(API_URL+'/v1/'+get_route(operation), json={'mediaItems': [context.mediaitem_id]}, headers=headers)
    context.response = res


@when('unmark {operation} mediaitem {condition} auth')
def step_impl(context, operation, condition):
    headers = None
    if condition == 'with':
        headers = {'Authorization': f'Bearer {context.access_token}'}
    res = requests.delete(API_URL+'/v1/'+get_route(operation), json={'mediaItems': [context.mediaitem_id]}, headers=headers)
    context.response = res


@then('mediaitem is marked as {operation}')
def step_impl(context, operation):
    assert context.response.status_code == 204


@then('mediaitem is unmarked as {operation}')
def step_impl(context, operation):
    assert context.response.status_code == 204


def get_route(operation: str) -> str:
    if operation == 'favourite':
        return 'favourites'
    if operation == 'hide' or operation == 'hidden':
        return 'hidden'
    if operation == 'delete' or operation == 'deleted':
        return 'trash'

def get_field(operation: str) -> str:
    if operation == 'favourite':
        return 'favourite'
    if operation == 'hidden':
        return 'hidden'
    if operation == 'deleted':
        return 'deleted'
