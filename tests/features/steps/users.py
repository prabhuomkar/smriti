from behave import *
import requests
from requests.auth import HTTPBasicAuth

from environment import API_URL, ADMIN_USERNAME, ADMIN_PASSWORD


created_user = {'name': 'John Doe', 'username': 'johndoe', 'password': 'johndoeT3st!'}
updated_user = {'name': 'UpdatedJohn Doe', 'username': 'updatedjohndoe', 'password': 'updatedjohndoeT3st!'}

@given('there are no users')
def step_impl(context):
    res = requests.get(API_URL+'/v1/users', auth=HTTPBasicAuth(ADMIN_USERNAME, ADMIN_PASSWORD))
    users = res.json()
    assert len(users) == 0

@given('there is user')
def step_impl(context):
    res = requests.get(API_URL+'/v1/users', auth=HTTPBasicAuth(ADMIN_USERNAME, ADMIN_PASSWORD))
    users = res.json()
    assert len(users) == 1
    context.user_id = users[0]['id']

@when('get all users')
def step_impl(context):
    res = requests.get(API_URL+'/v1/users', auth=HTTPBasicAuth(ADMIN_USERNAME, ADMIN_PASSWORD))
    context.users = res.json()

@when('get user')
def step_impl(context):
    user_id = context.user_id
    res = requests.get(API_URL+'/v1/users/'+user_id, auth=HTTPBasicAuth(ADMIN_USERNAME, ADMIN_PASSWORD))
    context.user = res.json()

@then('user is present in list')
def step_impl(context):
    assert len(context.users) == 1
    assert context.users[0]['name'] == context.match_user['name']
    assert context.users[0]['username'] == context.match_user['username']
    assert context.users[0]['password'] == context.match_user['password']

@then('user is present')
def step_impl(context):
    assert context.user['name'] == context.match_user['name']
    assert context.user['username'] == context.match_user['username']
    assert context.user['password'] == context.match_user['password']

@then('user is not present in list')
def step_impl(context):
    if len(context.users) > 0:
        assert context.users[0]['name'] != context.match_user['name']
        assert context.users[0]['username'] != context.match_user['username']
        assert context.users[0]['password'] != context.match_user['password']

@then('user is not present')
def step_impl(context):
    assert 'name' not in context.user
    assert 'username' not in context.user
    assert 'password' not in context.user
    assert context.user['message'] == 'user not found'

@when('create user is requested')
def step_impl(context):
    res = requests.post(API_URL+'/v1/users', json=created_user, auth=HTTPBasicAuth(ADMIN_USERNAME, ADMIN_PASSWORD))
    context.response = res
    context.user_id = res.json()['id']
    context.match_user = created_user

@when('update user is requested')
def step_impl(context):
    user_id = context.user_id
    res = requests.put(API_URL+'/v1/users/'+user_id, json=updated_user, auth=HTTPBasicAuth(ADMIN_USERNAME, ADMIN_PASSWORD))
    context.response = res
    context.match_user = updated_user

@when('delete user is requested')
def step_impl(context):
    user_id = context.user_id
    res = requests.delete(API_URL+'/v1/users/'+user_id, auth=HTTPBasicAuth(ADMIN_USERNAME, ADMIN_PASSWORD))
    context.response = res
    context.match_user = updated_user

@then('user is created')
def step_impl(context):
    assert context.response.status_code == 201

@then('user is updated')
def step_impl(context):
    assert context.response.status_code == 204

@then('user is deleted')
def step_impl(context):
    assert context.response.status_code == 204