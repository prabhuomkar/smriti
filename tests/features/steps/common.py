import os


API_URL = os.getenv('API_URL', 'http://localhost:5001')
ADMIN_USERNAME = 'carousel'
ADMIN_PASSWORD = 'carouselT3st!'


CREATED_USER = {'name': 'John Doe', 'username': 'johndoe', 'password': 'johndoeT3st!',
                'features':'{"albums":true,"favourites":true,"hidden":true,"trash":true,"explore":true,"places":true}'}
UPDATED_USER = {'name': 'UpdatedJohn Doe', 'username': 'updatedjohndoe', 'password': 'updatedjohndoeT3st!',
                'features':'{"albums":true,"favourites":true,"hidden":true,"trash":true,"explore":true,"places":true}'}

CREATED_ALBUM = {'name': 'Album Name', 'description': 'Album Description'}
UPDATED_ALBUM = {'name': 'Updated Album Name', 'description': 'Updated Album Description'}

CREATED_MEDIAITEM = {'filename': 'IMG_0543.HEIC', 'mimeType': 'image/heic',
                     'mediaItemType': 'photo', 'mediaItemCategory': 'default',
                     'description': None, 'favourite': False, 'hidden': False}
UPDATED_MEDIAITEM = {'filename': 'IMG_0543.HEIC', 'mimeType': 'image/heic',
                     'mediaItemType': 'photo', 'mediaItemCategory': 'default',
                     'description': 'Updated MediaItem Description', 'favourite': True, 'hidden': True}

CREATED_PLACE = {'name': 'Mumbai', 'city': 'Mumbai', 'state': 'Maharashtra', 'postcode': '400050', 'country': 'India'}