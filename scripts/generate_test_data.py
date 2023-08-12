"""Generate test data for local testing"""
import os
import requests
import zipfile
import shutil
import time


API_URL ='http://localhost:5001'
DOWNLOAD_PHOTOS_URL = 'https://www.dropbox.com/sh/q7yqg7uufaflqjg/AABAo-QEwNIAjxyZSJD0ICzDa?dl=1'

# create user
res = requests.post(f'{API_URL}/v1/users', auth=('smriti', 'smritiT3st!'), json={'name':'Jeff Dean','username':'jeffdean','password':'jeffT3st!','features':'{"albums":true,'+
                '"favourites":true,"hidden":true,"trash":true,"explore":true,"places":true,"things":true,"people":true,"sharing":true}'})
if res.status_code == 201:
    print('✅ created user')
else:
    res = res.json()
    if 'violates unique constraint' not in res['message']:
        print(f'❌ error creating user: {res["message"]}')
        exit(0)

# create auth token
res = requests.post(f'{API_URL}/v1/auth/login', json={'username':'jeffdean','password':'jeffT3st!'})
assert res.status_code == 200
res = res.json()
access_token = res['accessToken']
print('✅ got auth token')
headers = {'Authorization': f'Bearer {access_token}'}

# download and upload mediaitems
mediaitems = []
print('ℹ️ downloading sample mediaitems, hang on...')
if not os.path.exists('samples.zip'):
    response = requests.get(DOWNLOAD_PHOTOS_URL)
    if response.status_code == 200:
        local_file_path = 'samples.zip'
        with open(local_file_path, 'wb') as local_file:
            for chunk in response.iter_content(chunk_size=8192):
                local_file.write(chunk)
    else:
        print('❌ error downloading sample mediaitems')
        exit(0)
print('✅ downloaded sample mediaitems')
with zipfile.ZipFile('samples.zip', 'r') as zip_ref:
    zip_ref.extractall('samples')
print('ℹ️ uploading sample mediaitems, hang on...')
for file in os.listdir('samples'):
    files = {'file': open(f'samples/{file}', 'rb')}
    res = requests.post(f'{API_URL}/v1/mediaItems', files=files, headers=headers)
    if res.status_code != 201:
        res = res.json()
        print(f'❌ error uploading sample mediaitem: {res["message"]}')
        exit(0)
    res = res.json()
    mediaitems.append(res['id'])
    while True:
        res = requests.get(f'{API_URL}/v1/mediaItems/{res["id"]}', files=files, headers=headers)
        assert res.status_code == 200
        res = res.json()
        if res['status'] == 'READY':
            break
        time.sleep(2)
print('✅ uploaded sample mediaitems')

# create albums
first_album_mediaitems = mediaitems[:int(len(mediaitems)/2)]
second_album_mediaitems = mediaitems[int(len(mediaitems)/2):]
res = requests.post(f'{API_URL}/v1/albums', json={'name':'First Album'}, headers=headers)
if res.status_code == 201:
    res = res.json()
else:
    res = res.json()
    print(f'❌ error creating first album: {res["message"]}')
    exit(0)
res = requests.post(f'{API_URL}/v1/albums/{res["id"]}/mediaItems', json={'mediaItems':first_album_mediaitems}, headers=headers)
if res.status_code != 204:
    res = res.json()
    print(f'❌ error adding mediaitems to first album: {res["message"]}')
    exit(0)
res = requests.post(f'{API_URL}/v1/albums', json={'name':'Second Album'}, headers=headers)
if res.status_code == 201:
    res = res.json()
else:
    res = res.json()
    print(f'❌ error creating first album: {res["message"]}')
    exit(0)
res = requests.post(f'{API_URL}/v1/albums/{res["id"]}/mediaItems', json={'mediaItems':second_album_mediaitems}, headers=headers)
if res.status_code != 204:
    res = res.json()
    print(f'❌ error adding mediaitems to second album: {res["message"]}')
    exit(0)
print('✅ created 2 albums with mediaitems')

# clearing extracted files
shutil.rmtree('samples')
