import os


API_URL = os.getenv('API_URL', 'http://localhost:5001')
ADMIN_USERNAME = 'smriti'
ADMIN_PASSWORD = 'smritiT3st!'

CREATED_USER = {'name': 'John Doe', 'username': 'johndoe', 'password': 'johndoeT3st!',
                'features':'{"albums":true,"favourites":true,"hidden":true,"trash":true,"explore":true,"places":true}'}
UPDATED_USER = {'name': 'UpdatedJohn Doe', 'username': 'updatedjohndoe', 'password': 'updatedjohndoeT3st!',
                'features':'{"albums":true,"favourites":true,"hidden":true,"trash":true,"explore":true,"places":true}'}

CREATED_ALBUM = {'name': 'Album Name', 'description': 'Album Description'}
UPDATED_ALBUM = {'name': 'Updated Album Name', 'description': 'Updated Album Description', 'shared': True}

CREATED_MEDIAITEM = {'photo':{'filename': 'IMG_0543.HEIC', 'mimeType': 'image/heic', 'status': 'READY', 'cameraMake': 'Apple',
                     'cameraModel': 'iPhone 12 mini', 'focalLength': '4.2', 'apertureFNumber': '1.6', 'isoEquivalent': '640',
                     'exposureTime': '0.04', 'mediaItemType': 'photo', 'mediaItemCategory': 'default',
                     'description': None, 'favourite': False, 'hidden': False},
                     'video':{'filename': 'IMG_6470.MOV', 'mimeType': 'video/quicktime', 'status': 'READY', 'cameraMake': 'Apple',
                     'cameraModel': 'iPhone 12 mini', 'fps': '30', 'mediaItemType': 'video', 'mediaItemCategory': 'default',
                     'description': None, 'favourite': False, 'hidden': False}}
UPDATED_MEDIAITEM = {'photo':{'filename': 'IMG_0543.HEIC', 'mimeType': 'image/heic', 'status': 'READY', 'cameraMake': 'Apple',
                     'cameraModel': 'iPhone 12 mini', 'focalLength': '4.2', 'apertureFNumber': '1.6', 'isoEquivalent': '640',
                     'exposureTime': '0.04', 'mediaItemType': 'photo', 'mediaItemCategory': 'default',
                     'description': 'Updated MediaItem Description', 'favourite': True, 'hidden': False},
                     'video':{'filename': 'IMG_6470.MOV', 'mimeType': 'video/quicktime', 'status': 'READY', 'cameraMake': 'Apple',
                     'cameraModel': 'iPhone 12 mini', 'fps': '30', 'mediaItemType': 'video', 'mediaItemCategory': 'default',
                     'description': None, 'favourite': False, 'hidden': False}}

CREATED_PLACE = {'name': 'Mumbai', 'city': 'Mumbai', 'state': 'Maharashtra', 'postcode': '400050', 'country': 'India'}