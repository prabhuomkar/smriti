from behave import *
import requests
from requests.auth import HTTPBasicAuth

from environment import API_URL, ADMIN_USERNAME, ADMIN_PASSWORD, CREATED_USER


@given('a user is created')
def step_impl(context):
    res = requests.post(API_URL+'/v1/users', json=CREATED_USER,
                        auth=HTTPBasicAuth(ADMIN_USERNAME, ADMIN_PASSWORD))
    assert res.status_code == 201
    context.user_id = res.json()['id']


@then('a user is deleted')
def step_impl(context):
    res = requests.delete(API_URL+'/v1/users/'+context.user_id, json=CREATED_USER,
                          auth=HTTPBasicAuth(ADMIN_USERNAME, ADMIN_PASSWORD))
    assert res.status_code == 204


@when('user logs in')
def step_impl(context):
    res = requests.post(API_URL+'/v1/auth/login', json={
                        'username': CREATED_USER['username'], 'password': CREATED_USER['password']})
    context.response = res


@when('user refreshes token')
def step_impl(context):
    res = requests.post(API_URL+'/v1/auth/refresh',
                        headers={'Authorization': 'Bearer '+context.refresh_token})
    context.response = res


@when('user logs out')
def step_impl(context):
    res = requests.post(API_URL+'/v1/auth/logout',
                        headers={'Authorization': 'Bearer '+context.access_token})
    context.response = res


@then('token is generated')
def step_impl(context):
    assert context.response.status_code == 200
    body = context.response.json()
    assert 'accessToken' in body
    assert 'refreshToken' in body
    context.access_token = body['accessToken']
    context.refresh_token = body['refreshToken']


@then('token is refreshed')
def step_impl(context):
    assert context.response.status_code == 200
    body = context.response.json()
    assert 'accessToken' in body
    assert 'refreshToken' in body
    context.access_token = body['accessToken']
    context.refresh_token = body['refreshToken']


@then('token is deleted')
def step_impl(context):
    assert context.response.status_code == 204


@then('auth error is found')
def step_impl(context):
    assert context.response.status_code == 401
    assert context.response.json()['message'] == 'Unauthorized'
