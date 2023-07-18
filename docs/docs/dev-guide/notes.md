# Developer Notes

- [x] ML models should be agnostic and configurable
- [x] Each action should be configurable from the environment config e.g. Face grouping, Place Identification, File Size Limits
- [x] API-first techniques and then build Web UI & Mobile Apps
- [x] Demos can include: deployment on Fly, Heroku, other SaaS, hosting on Raspberry Pis and dockerized deployment on Cloud
- [x] Logging to make sure consistency across several components
- [x] Give a comparison of supported file formats between Google Photos and Smriti
- [ ] Send email on failures to owner or some metric of daily upload success/failed
- [x] Release Workflow in GitHub Actions, to deploy Docker images to Docker Hub Registry
- [x] Infra Cases: graceful shutdown of api & worker, database reconnects & grpc retry/reconnect
- [ ] i18n support for languages incrementally: English, Chinese, Japanese, German, French, etc.
- [x] Try to achieve best ratings incrementally as on: https://github.com/meichthys/foss_photo_libraries

## Features

| Feature | Support |
| - | - |
| Demo | ✅ |
| Freeness | ✅ |
| Docker Installation | ✅ |
| Albums | ✅ |
| Photo Map | ✅ |
| Multiple User Support | ✅ |
| EXIF Data | ✅  |
| Photo Discovery | 🚧 |
| Object/Face Recognition | 🚧 |
| Photo Search | ✅ |
| LivePhotos Support | ✅ |
| Video Support | ✅ |
| Automatic Mobile Upload | 🚧 |
| Timeline | 🚧 |
| Web App | 🚧 |
| Android App | 🚧 |
| iOS App | 🚧 |
| Desktop App | 🚧 |
| Photo Sharing | 🚧 |
| Basic Editing | 🚧 |
