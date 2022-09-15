# Design Document

## Requirements
MediaItems & Library:
- As a user I want to upload photos or videos and view them
- As a user I want to view metadata of a mediaitem
- As a user I want to favourite mediaitems and list liked mediaitems
- As a user I want to hide mediaitems and list hidden mediaitems
- As a user I want to delete mediaitems and list deleted mediaitems
- As a user I want to permanently delete mediaitems

Albums:
- As a user I want to create an album of mediaitems
- As a user I want to update name, description or cover of an album
- As a user I want to permanently delete an album
- As a user I want to hide an album

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
- [ExifTool](https://www.exiftool.org/)

#### Extracting Thumbnail
- TBD

#### Supported File Formats
| Type | Extension | Mime Type | Support |
| ---- | --------- | --------- | ------- |
| Photo | .JPG, .JPEG | image/jpg, image/jpeg | ✅ |
| Video | .MP4 | video/mp4 | ✅ |
TODO(omkar): Add more supported file formats

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
