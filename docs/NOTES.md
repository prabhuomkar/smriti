# Notes

- ML models should be agnostic and configurable
- Each action should be configurable from the environment config e.g. Face grouping, Place Identification, File Size Limits
- API-first techniques and then build Web UI & Mobile Apps
- Demos can include: deployment on Fly, Heroku, other SaaS, hosting on Raspberry Pis and dockerized deployment on Cloud
- Logging to make sure consistency across several components
- Give a comparison of supported file formats between Google Photos and Pensieve

## Database Design

[Postgres Schema](assets/schema.sql) is available.

### MediaItems
- Get all mediaitems for home screen
```sql
SELECT * FROM mediaitems 
WHERE (is_hidden=false OR is_deleted=false) AND status=READY;
```
- Get one mediaitem
```sql
SELECT * FROM mediaitems 
WHERE id=? AND status=READY;
```
- Get places for one mediaitem
```sql
SELECT * FROM places 
    INNER JOIN place_mediaitems
        ON places.id = place_mediaitems.place_id
WHERE place_mediaitems.mediaitem_id=?;
```
- Get things for one mediaitem
```sql
SELECT * FROM things 
    INNER JOIN thing_mediaitems
        ON things.id = thing_mediaitems.thing_id
WHERE thing_mediaitems.mediaitem_id=?;
```
- Get places for one mediaitem
```sql
SELECT * FROM people 
    INNER JOIN people_mediaitems
        ON people.id = people_mediaitems.people_id
WHERE people_mediaitems.mediaitem_id=?;
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
- Get all shared albums
```sql
SELECT * FROM albums
WHERE is_shared=true AND is_hidden=false;
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
PENSIEVE_LOG_LEVEL
# Feature
PENSIEVE_FEATURE_MEDIAITEM_UPLOAD
PENSIEVE_FEATURE_FAVOURITES
PENSIEVE_FEATURE_HIDDEN
PENSIEVE_FEATURE_TRASH
PENSIEVE_FEATURE_ALBUMS
PENSIEVE_FEATURE_EXPLORE
PENSIEVE_FEATURE_EXPLORE_PLACES
PENSIEVE_FEATURE_EXPLORE_THINGS
PENSIEVE_FEATURE_EXPLORE_PEOPLE
PENSIEVE_FEATURE_SHARING
# Email
SHARING_EMAIL_USERNAME
SHARING_EMAIL_PASSWORD
```
