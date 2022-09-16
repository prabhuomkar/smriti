# Design Document

## Requirements
- [EPICs](https://github.com/users/prabhuomkar/projects/5/views/7)

Explore:
- As a user I want to see all places from my mediaitems as a list
- As a user I want to see all mediaitems for a place as a list
- As a user I want to see all things from my mediaitems as a list
- As a user I want to see all mediaitems for a thing as a list
- As a user I want to see all people from my mediaitems as a list
- As a user I want to see all mediaitems for a person as a list
- As a user I want to update cover mediaitem for a person
- As a user I want to see list of images which are N years ago from today

## Architecture
TODO(omkar): Add architecture diagram

## High Level Design

### Components

#### API
- [Swagger Open API Spec](assets/swagger.yaml)
- Routing
- Common Middlewares
- Packaging and Structure
- Logging

#### Database
- [Postgres DB Schema](assets/schema.sql)
- Total number of tables: 12
- **Entities**:
    - MediaItem: `mediaitems`, `mediaitem_metadata`
    - Album: `albums`, `album_mediaitems`
    - Shared Album: `shared_albums`, `shared_album_mediaitems`
    - Explore: `places`, `things`, `people`, `place_mediaitems`, `things_mediaitems`, `people_mediaitems`

#### Worker
TODO(omkar): Background processing of files & inference

### Image & Video Processing

#### Parsing Metadata 
- [ExifTool](https://www.exiftool.org/) - Getting EXIF and XMP data

#### Extracting Thumbnail
- [LibRaw](https://www.libraw.org/) - Processing and extracting RAW images

#### Supported File Formats
| Type | Extension | Support |
| ---- | --------- | ------- |
| Photo | .BMP | ❓ |
| Photo | .GIF | ❓ |
| Photo | .HEIC | ❓ |
| Photo | .ICO | ❓ |
| Photo | .JPG | ❓ |
| Photo | .PNG | ❓ |
| Photo | .TIFF | ❓ |
| Photo | .WEBP | ❓ |
| Photo | [RAW Formats](https://www.libraw.org/supported-cameras) | ❓ |
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
    Init() // initialize client connection
    Upload() // upload the file in chunks
    Delete() // delete the file
}
```
- Out of the box incremental support for MinIO, Amazon S3, Google Storage

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
