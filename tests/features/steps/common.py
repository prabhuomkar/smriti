import os


API_URL = os.getenv('API_URL', 'http://localhost:5001')
ADMIN_USERNAME = 'smriti'
ADMIN_PASSWORD = 'smritiT3st!'

CREATED_USER = {
    'default': {
        'name': 'John Doe', 'username': 'johndoe', 'password': 'johndoeT3st!','features':'{"albums":true,'+
        '"favourites":true,"hidden":true,"trash":true,"explore":true,"places":true,"things":true,'+
        '"people":true,"sharing":true}'
    },
    'jobs': {
        'name': 'Steve Jobs', 'username': 'stevejobs', 'password': 'johndoeT3st!','features':'{"albums":true,'+
        '"favourites":true,"hidden":true,"trash":true,"sharing":true,"jobs":true}'
    }
}
UPDATED_USER = {
    'default': {
        'name': 'UpdatedJohn Doe', 'username': 'updatedjohndoe', 'password': 'updatedjohndoeT3st!','features':'{"albums"'+
        ':true,"favourites":true,"hidden":true,"trash":true,"explore":true,"places":true,"things":true,'+
        '"people":true,"sharing":true}'
    },
    'jobs': {
        'name': 'Steve Updated Jobs', 'features':'{"albums":true,"favourites":true,"hidden":true,"trash":true,'+
        '"explore":true,"places":true,"things":true,"people":true,"sharing":true,"jobs":true}'
    }
}

CREATED_ALBUM = {'name': 'Album Name', 'description': 'Album Description'}
CREATED_SHARED_ALBUM = {'name': 'Album Name', 'description': 'Album Description', 'shared': True}
UPDATED_ALBUM = {'name': 'Updated Album Name', 'description': 'Updated Album Description'}

CREATED_MEDIAITEM = {'PHOTO':{'filename': 'IMG_0543.HEIC', 'mimeType': 'image/heic', 'status': 'READY', 'cameraMake': 'Apple',
                     'cameraModel': 'iPhone 12 mini', 'focalLength': '4.2 mm', 'apertureFNumber': '1.6', 'isoEquivalent': '640',
                     'exposureTime': '1/25', 'mediaItemType': 'PHOTO', 'mediaItemCategory': 'LIVE',
                     'description': None, 'favourite': False, 'hidden': False},
                     'VIDEO':{'filename': 'IMG_6470.MOV', 'mimeType': 'video/quicktime', 'status': 'READY', 'cameraMake': 'Apple',
                     'cameraModel': 'iPhone 12 mini', 'fps': '30', 'mediaItemType': 'VIDEO', 'mediaItemCategory': 'DEFAULT',
                     'description': None, 'favourite': False, 'hidden': False}}
UPDATED_MEDIAITEM = {'PHOTO':{'filename': 'IMG_0543.HEIC', 'mimeType': 'image/heic', 'status': 'READY', 'cameraMake': 'Apple',
                     'cameraModel': 'iPhone 12 mini', 'focalLength': '4.2 mm', 'apertureFNumber': '1.6', 'isoEquivalent': '640',
                     'exposureTime': '1/25', 'mediaItemType': 'PHOTO', 'mediaItemCategory': 'LIVE',
                     'description': 'Updated MediaItem Description', 'favourite': True, 'hidden': False},
                     'VIDEO':{'filename': 'IMG_6470.MOV', 'mimeType': 'video/quicktime', 'status': 'READY', 'cameraMake': 'Apple',
                     'cameraModel': 'iPhone 12 mini', 'fps': '30', 'mediaItemType': 'VIDEO', 'mediaItemCategory': 'DEFAULT',
                     'description': None, 'favourite': False, 'hidden': False}}

CREATED_PLACE = {'name': 'Mumbai', 'area': 'Zone 3', 'locality': 'Mumbai', 'postcode': '400050', 'country': 'India'}

CREATED_THING = {'name': 'Pizza'}

FILES_TO_SKIP = ['3839-samsung - sm-g973u - 16bit (2.1132075471698).dng', '1087-leica - leica m monochrom (typ 246) - 12bit (3:2).dng',
                 '672-pentax - pentax optio s4.raw', '778-xiaomi - yi.raw', '3896-phase one - iq4 150mp - unknown (8) (4:3).iiq',
                 '4272-plustek - opticfilm 8200i se - 16bit (3:2).dng', '3901-phase one - iq4 150mp - unknown (8) (4:3).iiq',
                 '4953-phase one - ixm-rs150f - iiq sv2 (4:3).iiq', '3900-phase one - iq4 150mp - iiq l (4:3).iiq']