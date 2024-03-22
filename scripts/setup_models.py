"""Download and Setup Models"""
import os
import requests
import zipfile
import tqdm

GITHUB_API = 'https://api.github.com/repos/prabhuomkar/smriti/releases'
DOWNLOAD_MODELS_URL='https://github.com/prabhuomkar/smriti/releases/download/%s/models.zip'


print('ℹ️ downloading models, hang on...')
if not os.path.exists('models.zip'):
    response = requests.get(GITHUB_API)
    if response.status_code == 200:
        releases = response.json()
        tag_name = releases[0]['tag_name']
    else:
        print('❌ error getting release tag name')
        exit(0)
    download_url = f'{DOWNLOAD_MODELS_URL % tag_name}'
    print(f'Downloading models: {download_url}')
    response = requests.get(download_url, stream=True, allow_redirects=True)
    if response.status_code == 200:
        local_file_path = 'models.zip'
        with open(local_file_path, 'wb') as local_file:
            for chunk in response.iter_content(chunk_size=1024):
                local_file.write(chunk)
    else:
        print('❌ error downloading models')
        exit(0)
print('✅ downloaded models')
with zipfile.ZipFile('models.zip', 'r') as zip_ref:
    zip_ref.extractall('models/')
