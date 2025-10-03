"""Download and Setup Models"""
import os
import requests
import zipfile


GITHUB_API = "https://api.github.com/repos/prabhuomkar/smriti/releases"
DOWNLOAD_MODELS_URL="https://github.com/prabhuomkar/smriti/releases/download/%s/"
TYPES = ["faces", "ocr", "search"]

print("ℹ️ downloading models, hang on...")
if not os.path.exists("models.zip"):
    response = requests.get(GITHUB_API)
    if response.status_code == 200:
        releases = response.json()
        tag_name = releases[0]["tag_name"]
    else:
        print("❌ error getting release tag name")
        exit(0)
    for type in TYPES:
        download_url = f"{DOWNLOAD_MODELS_URL % tag_name}{type}.zip"
        print(f"Downloading models: {download_url}")
        response = requests.get(download_url, stream=True, allow_redirects=True)
        if response.status_code == 200:
            local_file_path = f"{type}.zip"
            with open(local_file_path, "wb") as local_file:
                for chunk in response.iter_content(chunk_size=1024):
                    local_file.write(chunk)
        else:
            print("❌ error downloading models")
            exit(0)
        print(f"✅ downloaded {type} models")
        with zipfile.ZipFile(f"{type}.zip", "r") as zip_ref:
            zip_ref.extractall("models/")
