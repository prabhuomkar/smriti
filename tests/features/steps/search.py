from behave import *
import requests

from common import API_URL


@when('search for mediaitems {condition} auth')
def step_impl(context, condition):
    headers = None
    if condition == 'with':
        headers = {'Authorization': f'Bearer {context.access_token}'}
    res = requests.get(API_URL+'/v1/search?q=pizza', headers=headers)
    context.response = res
    context.mediaitems = res.json()

@then('searched mediaitem is present in list')
def step_impl(context):
    assert len(context.mediaitems) == 1
    assert context.mediaitem_id == context.mediaitems[0]['id']
