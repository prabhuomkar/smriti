"""Generate test data for local testing"""
import os
import requests
import zipfile
import shutil
import time
import multiprocessing as mp


API_URL ='http://localhost:5001'
DOWNLOAD_SAMPLES_URL = 'https://www.dropbox.com/scl/fo/yyy82163nqh5ii5aqm8vz/h?rlkey=bjlvrz198fu9zntu6tmvgb2ya&dl=1'
MAX_PARALLEL_UPLOADS = 3

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

# upload function
def upload(file_type, file):
    print(f'uploading {file_type} {file}')
    files = {'file': open(f'samples/{file_type}/{file}', 'rb')}
    res = requests.post(f'{API_URL}/v1/mediaItems', files=files, headers=headers)
    if res.status_code != 201:
        res = res.json()
        print(f'❌ error uploading sample mediaitem: {res["message"]}')
        exit(0)
    res = res.json()
    while True:
        res = requests.get(f'{API_URL}/v1/mediaItems/{res["id"]}', files=files, headers=headers)
        assert res.status_code == 200
        res = res.json()
        if res['status'] == 'READY' or res['status'] == 'FAILED':
            print(f'finished {file_type} {file}')
            break
        time.sleep(5)
    return res['id']

# download and upload mediaitems
mediaitems = []
print('ℹ️ downloading sample mediaitems, hang on...')
if not os.path.exists('samples.zip'):
    response = requests.get(DOWNLOAD_SAMPLES_URL)
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
for file_type in os.listdir('samples'):
    files = os.listdir(f'samples/{file_type}')
    for i in range(0, len(files), MAX_PARALLEL_UPLOADS):
        chunk = files[i:i+MAX_PARALLEL_UPLOADS]
        with mp.Pool(processes=mp.cpu_count()-1) as pool:
            mediaitems = pool.starmap(upload, list(zip([file_type for _ in range(len(chunk))], chunk)))
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
