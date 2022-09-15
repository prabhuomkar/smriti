# Notes

- ML models should be agnostic and configurable
- Each action should be configurable from the environment config e.g. Face grouping, Place Identification, File Size Limits
- API-first techniques and then build Web UI & Mobile Apps
- Demos can include: deployment on Fly, Heroku, other SaaS, hosting on Raspberry Pis and dockerized deployment on Cloud
- Logging to make sure consistency across several components

## Database Design

[Postgres Schema](schema.sql) is available.

### MediaItems
- Get all mediaitems for home screen
```sql
SELECT * FROM mediaitems 
WHERE is_hidden=false OR is_deleted=false;
```
- Get one mediaitem
```sql
SELECT * FROM mediaitems 
    INNER JOIN mediaitem_metadata 
        ON mediaitems.id = mediaitem_metadata.mediaitem_id
WHERE mediaitems.id=?;
```

### Library 
- Get all favourite mediaitems
```sql
SELECT * FROM mediaitems 
WHERE is_favourite=true;
```
- Get all hidden mediaitems
```sql
SELECT * FROM mediaitems 
WHERE is_hidden=true;
```
- Get all deleted mediaitems
```sql
SELECT * FROM mediaitems 
WHERE is_deleted=true;
```

### Explore

#### Places
- Get all places
```sql
SELECT * FROM places
WHERE is_hidden=false;
```
- Get one place
```sql
SELECT * FROM places
WHERE id=?;
```
- Get mediaitems for one place
```sql
SELECT * FROM place_mediaitems 
    INNER JOIN mediaitems 
        ON place_mediaitems.mediaitem_id = mediaitems.id
WHERE place_mediaitems.place_id=? AND mediaitems.is_hidden=false OR mediaitems.is_deleted=false;
```

#### Things
- Get all things
```sql
SELECT * FROM things
WHERE is_hidden=false;
```
- Get one thing
```sql
SELECT * FROM things
WHERE id=?;
```
- Get mediaitems for one thing
```sql
SELECT * FROM thing_mediaitems 
    INNER JOIN mediaitems 
        ON thing_mediaitems.mediaitem_id = mediaitems.id
WHERE thing_mediaitems.thing_id=? AND mediaitems.is_hidden=false OR mediaitems.is_deleted=false;
```

#### People
- Get all people
```sql
SELECT * FROM people
WHERE is_hidden=false;
```
- Get one people
```sql
SELECT * FROM people
WHERE id=?;
```
- Get mediaitems for one people
```sql
SELECT * FROM people_mediaitems 
    INNER JOIN mediaitems 
        ON people_mediaitems.mediaitem_id = mediaitems.id
WHERE people_mediaitems.people_id=? AND mediaitems.is_hidden=false OR mediaitems.is_deleted=false;
```

### Albums
- Get all albums
```sql
SELECT * FROM albums
WHERE is_hidden=false;
```
- Get one album
```sql
SELECT * FROM albums
WHERE id=?;
```
- Get mediaitems for one album
```sql
SELECT * FROM album_mediaitems 
    INNER JOIN mediaitems 
        ON album_mediaitems.mediaitem_id = mediaitems.id
WHERE album_mediaitems.album_id=?;
```

### Sharing
- Get all shared albums
```sql
SELECT * FROM shared_albums;
```
- Get one shared album
```sql
SELECT * FROM shared_albums
WHERE id=?;
```
- Get mediaitems for one album
```sql
SELECT * FROM shared_album_mediaitems 
    INNER JOIN mediaitems 
        ON shared_album_mediaitems.mediaitem_id = mediaitems.id
WHERE shared_album_mediaitems.shared_album_id=?;
```

## System Configuration
```bash
# Database
DATABASE_HOST
DATABASE_PORT
DATABASE_USERNAME
DATABASE_PASSWORD
# Cache
CACHE_ENABLED
CACHE_HOST
CACHE_PORT
CACHE_USERNAME
CACHE_PASSWORD
# Core
PIXEL_LOG_LEVEL
# Feature
PIXEL_FEATURE_MEDIAITEM_UPLOAD
PIXEL_FEATURE_FAVOURITES
PIXEL_FEATURE_HIDDEN
PIXEL_FEATURE_TRASH
PIXEL_FEATURE_ALBUMS
PIXEL_FEATURE_EXPLORE
PIXEL_FEATURE_EXPLORE_PLACES
PIXEL_FEATURE_EXPLORE_THINGS
PIXEL_FEATURE_EXPLORE_PEOPLE
PIXEL_FEATURE_SHARING
# Email
SHARING_EMAIL_USERNAME
SHARING_EMAIL_PASSWORD
```