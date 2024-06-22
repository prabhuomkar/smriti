from behave import *
import requests

from common import API_URL


@when('search query {query} for mediaitems {condition} auth')
def step_impl(context, query, condition):
    headers = None
    if condition == 'with':
        headers = {'Authorization': f'Bearer {context.access_token}'}
    res = requests.get(API_URL+'/v1/search?q='+query, headers=headers)
    context.response = res
    context.mediaitems = res.json()

@then('searched mediaitem is present in list')
def step_impl(context):
    assert len(context.mediaitems) == 1
    assert context.mediaitem_id == context.mediaitems[0]['id']

@then('searched mediaitem is not present in list')
def step_impl(context):
    assert len(context.mediaitems) == 0
    assert context.mediaitem_id not in [mediaitem['id'] for mediaitem in context.mediaitems]
