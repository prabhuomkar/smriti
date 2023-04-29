from behave import *
import requests

from common import API_URL, CREATED_PLACE


@when('get all places for mediaitem {condition} auth')
def step_impl(context, condition):
    headers = None
    if condition == 'with':
        headers = {'Authorization': f'Bearer {context.access_token}'}
    mediaitem_id = context.mediaitem_id
    res = requests.get(API_URL+'/v1/mediaItems/'+mediaitem_id+'/places', headers=headers)
    context.response = res
    context.places = res.json()

@when('get all places {condition} auth')
def step_impl(context, condition):
    headers = None
    if condition == 'with':
        headers = {'Authorization': f'Bearer {context.access_token}'}
    res = requests.get(API_URL+'/v1/explore/places', headers=headers)
    context.response = res
    context.places = res.json()

@when('get place {condition} auth')
def step_impl(context, condition):
    headers = None
    if condition == 'with':
        headers = {'Authorization': f'Bearer {context.access_token}'}
    place_id = context.place_id
    res = requests.get(API_URL+'/v1/explore/places/'+place_id, headers=headers)
    context.response = res
    context.place = res.json()

@then('place is present in list')
def step_impl(context):
    assert len(context.places) == 1
    assert context.places[0]['name'] == context.match_place['name']
    assert context.places[0]['city'] == context.match_place['city']
    assert context.places[0]['state'] == context.match_place['state']
    assert context.places[0]['country'] == context.match_place['country']
    assert context.places[0]['postcode'] == context.match_place['postcode']

@then('place is present')
def step_impl(context):
    assert context.place['name'] == context.match_place['name']
    assert context.place['city'] == context.match_place['city']
    assert context.place['state'] == context.match_place['state']
    assert context.place['country'] == context.match_place['country']
    assert context.place['postcode'] == context.match_place['postcode']

@given('a mediaitem exists with place')
def step_impl(context):
    res = requests.get(API_URL+'/v1/mediaItems',
                       headers={'Authorization': f'Bearer {context.access_token}'})
    mediaitems = res.json()
    assert len(mediaitems) == 1
    context.mediaitem_id = mediaitems[0]['id']
    res = requests.get(API_URL+'/v1/explore/places',
                       headers={'Authorization': f'Bearer {context.access_token}'})
    places = res.json()
    assert len(places) == 1
    context.place_id = places[0]['id']
    context.match_place = CREATED_PLACE
