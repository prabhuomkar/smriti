from behave import *
import requests
from requests.auth import HTTPBasicAuth

from common import API_URL, ADMIN_USERNAME, ADMIN_PASSWORD, CREATED_USER, UPDATED_USER


@given('there are no users')
def step_impl(context):
    res = requests.get(API_URL+'/v1/users',
                       auth=HTTPBasicAuth(ADMIN_USERNAME, ADMIN_PASSWORD))
    users = res.json()
    assert len(users) == 0

@given('a user exists')
def step_impl(context):
    res = requests.get(API_URL+'/v1/users',
                       auth=HTTPBasicAuth(ADMIN_USERNAME, ADMIN_PASSWORD))
    users = res.json()
    assert len(users) == 1
    context.user_id = users[0]['id']

@when('get all users {condition} auth')
def step_impl(context, condition):
    auth = None
    if condition == 'with':
        auth = HTTPBasicAuth(ADMIN_USERNAME, ADMIN_PASSWORD)
    res = requests.get(API_URL+'/v1/users', auth=auth)
    context.response = res
    context.users = res.json()

@when('get user {condition} auth')
def step_impl(context, condition):
    auth = None
    if condition == 'with':
        auth = HTTPBasicAuth(ADMIN_USERNAME, ADMIN_PASSWORD)
    user_id = context.user_id
    res = requests.get(API_URL+'/v1/users/'+user_id, auth=auth)
    context.response = res
    context.user = res.json()

@then('user is present in list')
def step_impl(context):
    assert len(context.users) == 1
    assert context.users[0]['name'] == context.match_user['name']
    assert context.users[0]['username'] == context.match_user['username']

@then('user is present')
def step_impl(context):
    assert context.user['name'] == context.match_user['name']
    assert context.user['username'] == context.match_user['username']

@then('user is not present in list')
def step_impl(context):
    if len(context.users) > 0:
        assert context.users[0]['name'] != context.match_user['name']
        assert context.users[0]['username'] != context.match_user['username']

@then('user is not present')
def step_impl(context):
    assert 'name' not in context.user
    assert 'username' not in context.user
    assert context.user['message'] == 'user not found'

@when('create {name} user {condition} auth')
def step_impl(context, name, condition):
    auth = None
    if condition == 'with':
        auth = HTTPBasicAuth(ADMIN_USERNAME, ADMIN_PASSWORD)
    res = requests.post(API_URL+'/v1/users', json=CREATED_USER[name], auth=auth)
    context.response = res

@when('update {name} user {condition} auth')
def step_impl(context, name, condition):
    auth = None
    if condition == 'with':
        auth = HTTPBasicAuth(ADMIN_USERNAME, ADMIN_PASSWORD)
    user_id = context.user_id
    res = requests.put(API_URL+'/v1/users/'+user_id,
                       json=UPDATED_USER[name], auth=auth)
    context.response = res

@when('delete user {condition} auth')
def step_impl(context, condition):
    auth = None
    if condition == 'with':
        auth = HTTPBasicAuth(ADMIN_USERNAME, ADMIN_PASSWORD)
    user_id = context.user_id
    res = requests.delete(API_URL+'/v1/users/'+user_id, auth=auth)
    context.response = res

@then('user {name} is created')
def step_impl(context, name):
    assert context.response.status_code == 201
    context.user_id = context.response.json()['id']
    context.match_user = CREATED_USER[name]

@then('user {name} is updated')
def step_impl(context, name):
    assert context.response.status_code == 204
    context.match_user = UPDATED_USER[name]

@then('user {name} is deleted')
def step_impl(context, name):
    assert context.response.status_code == 204
    context.match_user = UPDATED_USER[name]

