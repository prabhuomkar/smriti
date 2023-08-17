"""Download and Setup Models"""
import os
import requests
import zipfile

DOWNLOAD_MODELS_URL='https://www.dropbox.com/scl/fi/czw3bz5w8mcioqu393xkm/models.zip?rlkey=l1olp762n7iu3qmh4ie52ic3w&dl=1'


print('ℹ️ downloading models, hang on...')
if not os.path.exists('models.zip'):
    response = requests.get(DOWNLOAD_MODELS_URL)
    if response.status_code == 200:
        local_file_path = 'models.zip'
        with open(local_file_path, 'wb') as local_file:
            for chunk in response.iter_content(chunk_size=8192):
                local_file.write(chunk)
    else:
        print('❌ error downloading models')
        exit(0)
print('✅ downloaded models')
with zipfile.ZipFile('models.zip', 'r') as zip_ref:
    zip_ref.extractall('.')
