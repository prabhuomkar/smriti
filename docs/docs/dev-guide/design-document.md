# Software Design Document

## High Level Design

### Components

#### API
- Service written in Golang
  - REST API: [echo](https://echo.labstack.com/)
  - RPC: [gRPC + protobuf](https://grpc.io/)
  - Postgres: [gorm](https://gorm.io/)
  - Linting: [golangci-lint](https://golangci-lint.run/)
- Will read/write to Database
- Will read/write to Disk
- Will exchange protobuf with Worker: [api.proto](https://github.com/prabhuomkar/smriti/blob/master/protos/api.proto)

#### Database
- [Schema](https://github.com/prabhuomkar/smriti/blob/master/infra/database/schema.sql)
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
- Will process mediaitem files for several types of tasks
- Will exchange protobuf with API: [worker.proto](https://github.com/prabhuomkar/smriti/blob/master/protos/worker.proto)

### Image & Video Processing

#### Parsing Metadata 
- [ExifTool](https://www.exiftool.org/) - Getting EXIF and XMP data

#### Extracting Thumbnail
- [LibRaw](https://www.libraw.org/) - Processing and extracting RAW images
- [ImageMagick](https://imagemagick.org/index.php) - General purpose extraction

#### Supported File Formats
| Type | Extension | Support |
| ---- | --------- | ------- |
| Photo | .BMP | ✅ |
| Photo | .GIF | ✅ |
| Photo | .HEIC | ✅ |
| Photo | .ICO | ✅ |
| Photo | .JPG | ✅ |
| Photo | .PNG | ✅ |
| Photo | .TIFF | ✅ |
| Photo | .WEBP | ✅ |
| Photo | [RAW Formats](https://raw.pixls.us/) | ✅ |

**Post Stable Release, Scope for Video**:

| Type | Extension | Support |
| ---- | --------- | ------- |
| Video | 3GP | ❓ |
| Video | 3G2 | ❓ |
| Video | ASF | ❓ |
| Video | AVI | ✅ |
| Video | DIVX | ❓ |
| Video | M2T | ❓ |
| Video | M2TS | ❓ |
| Video | M4V | ❓ |
| Video | MKV | ❓ |
| Video | MMV | ❓ |
| Video | MOD | ❓ |
| Video | MOV | ✅ |
| Video | MP4 | ✅ |
| Video | MPG | ❓ |
| Video | MTS | ❓ |
| Video | TOD | ❓ |
| Video | WMV | ❓ |

### File Storage & Retrieval
- Support for several file storage systems behind a common interface:
```
interface {
  upload() // upload the file (in chunks if required)
  delete() // delete the file
  get() // get the file
}
```
- Out of the box incremental support for storage:
  - Disk
  - [MinIO](https://min.io/)
  - [Amazon S3](https://aws.amazon.com/s3/)
- Best practices for security and other similar aspects for connecting to storage will be decided later

### Machine Learning Inference
- Toggles for using CPU and GPU for inference
- Model registration from API on startup:
```yaml
tasks:
- name: places
  service: openstreetmap # google-maps, geojson, paid services, etc.
- name: classification
  download: 
    - s3://model-download-link
  model: convnext
- name: detection
  download: 
    - s3://model-download-link
    - s3://model-labels-link
  model: yolo
- name: ocr
  download: 
    - s3://model-download-link
    - s3://model-labels-link
  model: tbd
- name: faces
  download: 
    - s3://model-download-link
    - s3://model-labels-link
  model: tbd
  params:
    - face_size: 224
    - similarity_threshold: 90
- name: speech
  download: 
    - s3://model-download-link
    - s3://model-labels-link
  model: tbd
  params:
    - token_size: 10
```
- Default Models:
  - Classification - [EfficientNet](https://github.com)
  - Detection - [YOLOv8](https://github.com)
  - Face Detection - [TBD](https://github.com)
  - OCR - [TBD](https://github.com)
  - Speech - [TBD](https://github.com)
  - Search Embeddings - [TBD](https://github.com)

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
