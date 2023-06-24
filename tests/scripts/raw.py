import sys
import time
import requests
import re
import multiprocessing as mp


def download_mediaitem(mediaitem):
    source = re.findall(r"href='([^']*)'", mediaitem[7])[0]
    file_name = f'/tmp/{source.split("/")[-1]}'.lower()
    try:
        response = requests.get(source)
        with open(file_name, 'wb') as file:
            file.write(response.content)
        print(f'downloaded: {file_name}')
    except:
        return None

if __name__ == '__main__':
    res = requests.get(f'https://raw.pixls.us/json/getrepository.php?set=all&_={int(time.time()*1000)}')
    items = res.json()['data']
    cameras = sys.argv[1].split(',') if len(sys.argv) > 0 else []
    data = []
    for mediaitem in items:
        if mediaitem[0].lower() in cameras:
            data.append(mediaitem)
    with mp.Pool(mp.cpu_count()) as pool:
        result = pool.map_async(download_mediaitem, data)
        result.get()
