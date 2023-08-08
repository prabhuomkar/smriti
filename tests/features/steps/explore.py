from behave import *
import requests

from common import API_URL, CREATED_PLACE, CREATED_THING


@when('get all explored {type} for mediaitem {condition} auth')
def step_impl(context, type, condition):
    headers = None
    if condition == 'with':
        headers = {'Authorization': f'Bearer {context.access_token}'}
    mediaitem_id = context.mediaitem_id
    res = requests.get(API_URL+'/v1/mediaItems/'+mediaitem_id+'/'+type, headers=headers)
    context.response = res
    if type == 'places':
        context.places = res.json()
    elif type == 'things':
        context.things = res.json()
    elif type == 'people':
        context.people = res.json()

@when('get all explored {type} {condition} auth')
def step_impl(context, type, condition):
    headers = None
    if condition == 'with':
        headers = {'Authorization': f'Bearer {context.access_token}'}
    res = requests.get(API_URL+'/v1/explore/'+type, headers=headers)
    context.response = res
    if type == 'places':
        context.places = res.json()
    elif type == 'things':
        context.things = res.json()
    elif type == 'people':
        context.people = res.json()

@when('get explored {type} {condition} auth')
def step_impl(context, type, condition):
    headers = None
    if condition == 'with':
        headers = {'Authorization': f'Bearer {context.access_token}'}
    type_id = context.place_id if type == 'place' else context.thing_id if type == 'thing' else context.person_id if type == 'person' else None
    res = requests.get(API_URL+'/v1/explore/'+get_plural(type)+'/'+type_id, headers=headers)
    context.response = res
    if type == 'place':
        context.place = res.json()
    elif type == 'thing':
        context.thing = res.json()
    elif type == 'person':
        context.person = res.json()

@then('explored {type} is present in list')
def step_impl(context, type):
    if type == 'place':
        assert len(context.places) == 1
        assert context.places[0]['name'] == context.match_place['name']
        assert context.places[0]['city'] == context.match_place['city']
        assert context.places[0]['state'] == context.match_place['state']
        assert context.places[0]['country'] == context.match_place['country']
        assert context.places[0]['postcode'] == context.match_place['postcode']
    elif type == 'thing':
        assert len(context.things) == 1
        assert context.things[0]['name'] == context.match_thing['name']
    elif type == 'person':
        assert len(context.people) >= 1
        assert context.person_id in [person['id'] for person in context.people]

@then('explored {type} is present')
def step_impl(context, type):
    if type == 'place':
        assert context.place['name'] == context.match_place['name']
        assert context.place['city'] == context.match_place['city']
        assert context.place['state'] == context.match_place['state']
        assert context.place['country'] == context.match_place['country']
        assert context.place['postcode'] == context.match_place['postcode']
    elif type == 'thing':
        assert context.thing['name'] == context.match_thing['name']
    elif type == 'person':
        assert context.person['id'] == context.person_id

@when('get all mediaitems for {type} {condition} auth')
def step_impl(context, type, condition):
    headers = None
    if condition == 'with':
        headers = {'Authorization': f'Bearer {context.access_token}'}
    type_id = context.place_id if type == 'place' else context.thing_id if type == 'thing' else None
    res = requests.get(API_URL+'/v1/explore/'+get_plural(type)+'/'+type_id+'/mediaItems', headers=headers)
    context.response = res
    if type == 'place':
        context.place_mediaitems = res.json()
    elif type == 'thing':
        context.thing_mediaitems = res.json()

@then('mediaitem with {type} is present in list')
def step_impl(context, type):
    if type == 'place':
        assert len(context.place_mediaitems) == 1
        for field in context.match_mediaitem:
            if context.match_mediaitem[field] != None:
                assert context.place_mediaitems[0][field] == context.match_mediaitem[field]
    elif type == 'thing':
        assert len(context.thing_mediaitems) == 1
        for field in context.match_mediaitem:
            if context.match_mediaitem[field] != None:
                assert context.thing_mediaitems[0][field] == context.match_mediaitem[field]

@given('a mediaitem exists with {type}')
def step_impl(context, type):
    res = requests.get(API_URL+'/v1/mediaItems',
                       headers={'Authorization': f'Bearer {context.access_token}'})
    mediaitems = res.json()
    assert len(mediaitems) == 1
    context.mediaitem_id = mediaitems[0]['id']
    res = requests.get(API_URL+'/v1/explore/'+get_plural(type),
                       headers={'Authorization': f'Bearer {context.access_token}'})
    types = res.json()
    assert len(types) >= 1
    if type == 'place':
        context.place_id = types[0]['id']
        context.match_place = CREATED_PLACE
    elif type == 'thing':
        context.thing_id = types[0]['id']
        context.match_thing = CREATED_THING
    elif type == 'person':
        context.person_id = types[0]['id']
        context.match_people = {}

def get_plural(type: str) -> str:
    return type+'s' if type != 'person' else 'people'