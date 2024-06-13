import time
from behave import *
import requests

from common import API_URL


@when('create jobs for {components} components {condition} auth')
def step_impl(context, components, condition):
    headers = None
    if condition == 'with':
        headers = {'Authorization': f'Bearer {context.access_token}'}
    res = requests.post(API_URL+'/v1/jobs', json={'components': components}, headers=headers)
    context.response = res
    context.job = res.json()

@then('job is created')
def step_impl(context):
    assert context.response.status_code == 201

@when('get jobs {condition} auth and wait {seconds} seconds')
def step_impl(context, condition, seconds):
    time.sleep(int(seconds))
    headers = None
    if condition == 'with':
        headers = {'Authorization': f'Bearer {context.access_token}'}
    res = requests.get(API_URL+'/v1/jobs', headers=headers)
    context.response = res
    context.jobs = res.json()

@when('get jobs {condition} auth and wait until {status}')
def step_impl(context, condition, status):
    while True:
        headers = None
        if condition == 'with':
            headers = {'Authorization': f'Bearer {context.access_token}'}
        res = requests.get(API_URL+'/v1/jobs', headers=headers)
        context.response = res
        context.jobs = res.json()
        res_job = None
        for job in context.jobs:
            if job['id'] == context.job['id']:
                res_job = job
        if res_job['status'] == status.upper():
            break
        time.sleep(10)

@then('job is {status} and present in list')
def step_impl(context, status):
    assert context.job['id'] in [job['id'] for job in context.jobs]
    for job in context.jobs:
        if job['id'] == context.job['id']:
            assert status.upper() == job['status']
            break

@when('get job {condition} auth')
def step_impl(context, condition):
    headers = None
    if condition == 'with':
        headers = {'Authorization': f'Bearer {context.access_token}'}
    res = requests.get(API_URL+'/v1/jobs/'+context.job['id'], headers=headers)
    context.response = res
    context.job_res = res.json()

@then('job is present')
def step_impl(context):
    assert context.job_res is not None
    assert context.job_res['status'] in ['SCHEDULED', 'RUNNING', 'PAUSED', 'COMPLETED', 'STOPPED']

@when('get job mediaitem {component}')
def step_impl(context, component):
    headers = {'Authorization': f'Bearer {context.access_token}'}
    res = requests.get(API_URL+'/v1/mediaItems/'+context.mediaitem_id+'/'+component, headers=headers)
    context.response = res
    context.mediaitem_things = res.json()

@then('job mediaitem related {component} are {condition} in list')
def step_impl(context, component, condition):
    data = context.mediaitem_things if component == 'things' else \
            context.mediaitem_people if component == 'people' else \
            context.mediaitem_places if component == 'places' else []
    if condition == 'absent':
        assert len(data) == 0
    else:
        assert len(data) != 0
