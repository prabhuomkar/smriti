import os
import json
import sys
import time
import requests
import re
import multiprocessing as mp


def download_mediaitem(mediaitem):
    source = re.findall(r"href='([^']*)'", mediaitem[7])[0]
    file_name = f'/tmp/{source.split("/")[-3]}-{source.split("/")[-1]}'.lower()
    try:
        response = requests.get(source)
        with open(file_name, 'wb') as file:
            file.write(response.content)
        print(f'downloaded: {file_name}')
    except:
        return None

if __name__ == '__main__':
    items = []
    file_name = '/tmp/raw-list.json'
    if os.path.exists(file_name):
        print('raw list exists')
        with open(file_name, 'r') as f:
            items = json.load(f)
    else:
        print('raw list does not exist')
        res = requests.get(f'https://raw.pixls.us/json/getrepository.php?set=all&_={int(time.time()*1000)}')
        items = res.json()['data']
        with open(file_name, 'w') as f:
            json.dump(items, f)
    cameras = sys.argv[1].split(',') if len(sys.argv) > 0 else []
    data = []
    for mediaitem in items:
        if mediaitem[0].lower() in cameras:
            data.append(mediaitem)
    with mp.Pool(mp.cpu_count()) as pool:
        result = pool.map_async(download_mediaitem, data)
        result.get()
