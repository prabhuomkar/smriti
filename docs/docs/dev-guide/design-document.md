# Software Design Document

## Requirements
- [EPICs](https://github.com/users/prabhuomkar/projects/5/views/7)


## Architecture
TODO(omkar): Add architecture diagram

## High Level Design

### Components

#### API
- TODO(omkar): Add Swagger Open API Spec
- Service written in Golang
    - REST API: [echo](https://echo.labstack.com/)
    - RPC: [gRPC + protobuf](https://grpc.io/)
    - Postgres: [gorm](https://gorm.io/)
    - Linting: [golangci-lint](https://golangci-lint.run/)
- Will read/write to Database
- Will exchange protobuf with Worker: TODO(omkar): add _api.proto_

#### Database
TODO(omkar): Add entity relationship diagram
TODO(omkar): Add database schema sql
- Total number of tables: 10
- **Entities**:
    - MediaItem: `mediaitems`
    - Album: `albums`, `album_mediaitems`
    - Explore: `places`, `things`, `people`, `place_mediaitems`, `thing_mediaitems`, `people_mediaitems`
    - User Management: `users`

#### Worker
- Service written in Python
    - RPC: [gRPC + protobuf](https://grpc.io/)
    - Linting: [pylint](https://pypi.org/project/pylint/)
    - Exiftool: [PyExifTool](https://pypi.org/project/PyExifTool/)
    - LibRaw: [rawpy](https://pypi.org/project/rawpy/)
    - CDN: TBD, depends on file storage systems
- Will process images and videos
- Will exchange protobuf with API: TODO(omkar): add _worker.proto_

### Image & Video Processing

#### Parsing Metadata 
- [ExifTool](https://www.exiftool.org/) - Getting EXIF and XMP data

#### Extracting Thumbnail
- [LibRaw](https://www.libraw.org/) - Processing and extracting RAW images

#### Supported File Formats
| Type | Extension | Support |
| ---- | --------- | ------- |
| Photo | .BMP | ✅ |
| Photo | .GIF | ✅ |
| Photo | .HEIC | ❓ |
| Photo | .ICO | ❓ |
| Photo | .JPG | ✅ |
| Photo | .PNG | ✅ |
| Photo | .TIFF | ✅ |
| Photo | .WEBP | ✅ |
| Photo | [RAW Formats](https://www.libraw.org/supported-cameras) | ❓ |

**Post Stable Release, Scope for Video**:

| Type | Extension | Support |
| ---- | --------- | ------- |
| Video | 3GP | ❓ |
| Video | 3G2 | ❓ |
| Video | ASF | ❓ |
| Video | AVI | ❓ |
| Video | DIVX | ❓ |
| Video | M2T | ❓ |
| Video | M2TS | ❓ |
| Video | M4V | ❓ |
| Video | MKV | ❓ |
| Video | MMV | ❓ |
| Video | MOD | ❓ |
| Video | MOV | ❓ |
| Video | MP4 | ❓ |
| Video | MPG | ❓ |
| Video | MTS | ❓ |
| Video | TOD | ❓ |
| Video | WMV | ❓ |
- TODO(omkar): Investiage support for Android Motion Photos, iOS Live Photos

### File Storage & Retrieval
- Support for several file storage systems behind a common interface:
```
interface {
    connect() // initialize connection
    reconnect() // re-establish connection
    upload() // upload the file in chunks
    delete() // delete the file
}
```
- Out of the box incremental support for storage:
    - [Amazon S3](https://aws.amazon.com/s3/)
- Best practices for security and other similar aspects for connecting to storage will be decided later

### Machine Learning Inference
TODO(omkar): Add more information

## Performance
- Benchmarking with several parallel uploads and system configuration
- Results of benchmarks and some graphs

## Testing
- E2E Automation Testing using `behave`

## Security
- Authentication mechanisms; Basic Auth & JWT
- Accessing CDN files using hash keys

## Deployment
- Docker Deployment 
- Volumes and Other Concerns
- Working with HTTPS
