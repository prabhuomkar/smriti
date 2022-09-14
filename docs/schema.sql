CREATE TYPE "mediaitem_type" AS ENUM (
  'photo',
  'video'
);

CREATE TYPE "mediaitem_status" AS ENUM (
  'UNSPECIFIED',
  'PROCESSING',
  'READY',
  'FAILED'
);

CREATE TABLE "mediaitems" (
  "id" uuid NOT NULL,
  "filename" varchar NOT NULL,
  "description" varchar,
  "mime_type" varchar NOT NULL,
  "source_url" varchar NOT NULL,
  "thumbnail_url" varchar NOT NULL,
  "is_favourite" boolean,
  "is_archived" boolean,
  "is_deleted" boolean,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "mediaitem_metadata" (
  "mediaitem_id" uuid PRIMARY KEY NOT NULL,
  "mediaitem_type" mediaitem_type,
  "width" int,
  "height" int,
  "creation_time" timestamp,
  "camera_make" varchar,
  "camera_model" varchar,
  "focal_length" varchar,
  "aperture_fnumber" varchar,
  "iso_equivalent" varchar,
  "exposure_time" varchar,
  "fps" varchar,
  "status" mediaitem_status
);

CREATE TABLE "places" (
  "id" uuid PRIMARY KEY NOT NULL,
  "postcode" varchar UNIQUE,
  "suburb" varchar,
  "road" varchar,
  "town" varchar,
  "city" varchar,
  "county" varchar,
  "district" varchar,
  "state" varchar,
  "country" varchar,
  "cover_mediaitem_id" uuid,
  "cover_mediaitem_thumbnail_url" varchar
);

CREATE TABLE "place_mediaitems" (
  "place_id" uuid NOT NULL,
  "mediaitem_id" uuid PRIMARY KEY NOT NULL,
  "latitude" varchar,
  "longitude" varchar,
  PRIMARY KEY ("place_id", "mediaitem_id")
);

CREATE TABLE "things" (
  "id" uuid PRIMARY KEY NOT NULL,
  "name" varchar UNIQUE,
  "cover_mediaitem_id" uuid,
  "cover_mediaitem_thumbnail_url" varchar
);

CREATE TABLE "thing_mediaitems" (
  "thing_id" uuid NOT NULL,
  "mediaitem_id" uuid PRIMARY KEY NOT NULL,
  PRIMARY KEY ("thing_id", "mediaitem_id")
);

CREATE TABLE "people" (
  "id" uuid PRIMARY KEY NOT NULL,
  "name" varchar UNIQUE,
  "cover_mediaitem_id" uuid,
  "cover_mediaitem_thumbnail_url" varchar
);

CREATE TABLE "people_mediaitems" (
  "people_id" uuid NOT NULL,
  "mediaitem_id" uuid PRIMARY KEY NOT NULL,
  PRIMARY KEY ("people_id", "mediaitem_id")
);

CREATE TABLE "albums" (
  "id" uuid PRIMARY KEY NOT NULL,
  "name" varchar NOT NULL,
  "description" varchar,
  "cover_mediaitem_id" uuid,
  "cover_mediaitem_thumbnail_url" varchar,
  "mediaitems_count" int,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "album_mediaitems" (
  "album_id" uuid NOT NULL,
  "mediaitem_id" uuid NOT NULL,
  PRIMARY KEY ("album_id", "mediaitem_id")
);

CREATE TABLE "shared_albums" (
  "id" uuid PRIMARY KEY NOT NULL,
  "name" varchar NOT NULL,
  "cover_mediaitem_id" uuid,
  "cover_mediaitem_thumbnail_url" varchar,
  "mediaitems_count" int,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "shared_album_mediaitems" (
  "shared_album_id" uuid NOT NULL,
  "mediaitem_id" uuid NOT NULL,
  PRIMARY KEY ("shared_album_id", "mediaitem_id")
);

ALTER TABLE "mediaitem_metadata" ADD FOREIGN KEY ("mediaitem_id") REFERENCES "mediaitems" ("id");

ALTER TABLE "place_mediaitems" ADD FOREIGN KEY ("place_id") REFERENCES "places" ("id");

ALTER TABLE "place_mediaitems" ADD FOREIGN KEY ("mediaitem_id") REFERENCES "mediaitems" ("id");

ALTER TABLE "thing_mediaitems" ADD FOREIGN KEY ("thing_id") REFERENCES "things" ("id");

ALTER TABLE "thing_mediaitems" ADD FOREIGN KEY ("mediaitem_id") REFERENCES "mediaitems" ("id");

ALTER TABLE "people_mediaitems" ADD FOREIGN KEY ("people_id") REFERENCES "people" ("id");

ALTER TABLE "people_mediaitems" ADD FOREIGN KEY ("mediaitem_id") REFERENCES "mediaitems" ("id");

ALTER TABLE "album_mediaitems" ADD FOREIGN KEY ("album_id") REFERENCES "albums" ("id");

ALTER TABLE "album_mediaitems" ADD FOREIGN KEY ("mediaitem_id") REFERENCES "mediaitems" ("id");

ALTER TABLE "shared_album_mediaitems" ADD FOREIGN KEY ("shared_album_id") REFERENCES "shared_albums" ("id");

ALTER TABLE "shared_album_mediaitems" ADD FOREIGN KEY ("mediaitem_id") REFERENCES "mediaitems" ("id");
