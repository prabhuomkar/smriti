import os


API_URL = os.getenv('API_URL', 'http://localhost:5001')
ADMIN_USERNAME = 'smriti'
ADMIN_PASSWORD = 'smritiT3st!'

CREATED_USER = {'name': 'John Doe', 'username': 'johndoe', 'password': 'johndoeT3st!','features':'{"albums":true,'+
                '"favourites":true,"hidden":true,"trash":true,"explore":true,"places":true,"things":true,'+
                '"people":true,"sharing":true}'}
UPDATED_USER = {'name': 'UpdatedJohn Doe', 'username': 'updatedjohndoe', 'password': 'updatedjohndoeT3st!','features':'{"albums"'+
                ':true,"favourites":true,"hidden":true,"trash":true,"explore":true,"places":true,"things":true,'+
                '"people":true,"sharing":true}'}

CREATED_ALBUM = {'name': 'Album Name', 'description': 'Album Description'}
CREATED_SHARED_ALBUM = {'name': 'Album Name', 'description': 'Album Description', 'shared': True}
UPDATED_ALBUM = {'name': 'Updated Album Name', 'description': 'Updated Album Description'}

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

CREATED_THING = {'name': 'Pizza'}

FILES_TO_SKIP = ['3839-samsung - sm-g973u - 16bit (2.1132075471698).dng', '1087-leica - leica m monochrom (typ 246) - 12bit (3:2).dng',
                 '672-pentax - pentax optio s4.raw', '778-xiaomi - yi.raw', '3896-phase one - iq4 150mp - unknown (8) (4:3).iiq',
                 '4272-plustek - opticfilm 8200i se - 16bit (3:2).dng', '3901-phase one - iq4 150mp - unknown (8) (4:3).iiq',
                 '4953-phase one - ixm-rs150f - iiq sv2 (4:3).iiq', '3900-phase one - iq4 150mp - iiq l (4:3).iiq']