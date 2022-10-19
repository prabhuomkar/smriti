CREATE TYPE "mediaitem_status" AS ENUM (
  'UNSPECIFIED',
  'PROCESSING',
  'READY',
  'FAILED'
);

CREATE TYPE "mediaitem_type" AS ENUM (
  'photo',
  'video'
);

CREATE TABLE "mediaitems" (
  "id" uuid PRIMARY KEY NOT NULL,
  "filename" varchar,
  "description" varchar,
  "mime_type" varchar,
  "source_url" varchar,
  "preview_url" varchar,
  "thumbnail_url" varchar,
  "is_favourite" boolean DEFAULT FALSE NOT NULL,
  "is_hidden" boolean DEFAULT FALSE NOT NULL,
  "is_deleted" boolean DEFAULT FALSE NOT NULL,
  "status" mediaitem_status,
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
  "latitude" double precision,
  "longitude" double precision,
  "fps" varchar,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "places" (
  "id" uuid PRIMARY KEY NOT NULL,
  "name" varchar,
  "postcode" varchar,
  "town" varchar,
  "city" varchar,
  "state" varchar,
  "country" varchar,
  "cover_mediaitem_id" uuid,
  "is_hidden" boolean DEFAULT FALSE NOT NULL,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "place_mediaitems" (
  "place_id" uuid NOT NULL,
  "mediaitem_id" uuid NOT NULL,
  PRIMARY KEY ("place_id", "mediaitem_id")
);

CREATE TABLE "things" (
  "id" uuid PRIMARY KEY NOT NULL,
  "name" varchar UNIQUE,
  "cover_mediaitem_id" uuid,
  "is_hidden" boolean DEFAULT FALSE NOT NULL,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "thing_mediaitems" (
  "thing_id" uuid NOT NULL,
  "mediaitem_id" uuid NOT NULL,
  PRIMARY KEY ("thing_id", "mediaitem_id")
);

CREATE TABLE "people" (
  "id" uuid PRIMARY KEY NOT NULL,
  "name" varchar UNIQUE,
  "cover_mediaitem_id" uuid,
  "is_hidden" boolean DEFAULT FALSE NOT NULL,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "people_mediaitems" (
  "people_id" uuid NOT NULL,
  "mediaitem_id" uuid NOT NULL,
  PRIMARY KEY ("people_id", "mediaitem_id")
);

CREATE TABLE "albums" (
  "id" uuid PRIMARY KEY NOT NULL,
  "name" varchar NOT NULL,
  "description" varchar,
  "is_shared" boolean DEFAULT FALSE NOT NULL,
  "cover_mediaitem_id" uuid,
  "mediaitems_count" int DEFAULT 0 NOT NULL,
  "is_hidden" boolean DEFAULT FALSE NOT NULL,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "album_mediaitems" (
  "album_id" uuid NOT NULL,
  "mediaitem_id" uuid NOT NULL,
  PRIMARY KEY ("album_id", "mediaitem_id")
);

CREATE TABLE "users" (
  "id" uuid PRIMARY KEY NOT NULL,
  "name" varchar NOT NULL,
  "username" varchar NOT NULL,
  "password" varchar NOT NULL,
  "created_at" timestamp,
  "updated_at" timestamp
);

ALTER TABLE "albums" ADD FOREIGN KEY ("cover_mediaitem_id") REFERENCES "mediaitems" ("id");

ALTER TABLE "places" ADD FOREIGN KEY ("cover_mediaitem_id") REFERENCES "mediaitems" ("id");

ALTER TABLE "things" ADD FOREIGN KEY ("cover_mediaitem_id") REFERENCES "mediaitems" ("id");

ALTER TABLE "people" ADD FOREIGN KEY ("cover_mediaitem_id") REFERENCES "mediaitems" ("id");

ALTER TABLE "place_mediaitems" ADD FOREIGN KEY ("place_id") REFERENCES "places" ("id");

ALTER TABLE "place_mediaitems" ADD FOREIGN KEY ("mediaitem_id") REFERENCES "mediaitems" ("id");

ALTER TABLE "thing_mediaitems" ADD FOREIGN KEY ("thing_id") REFERENCES "things" ("id");

ALTER TABLE "thing_mediaitems" ADD FOREIGN KEY ("mediaitem_id") REFERENCES "mediaitems" ("id");

ALTER TABLE "people_mediaitems" ADD FOREIGN KEY ("people_id") REFERENCES "people" ("id");

ALTER TABLE "people_mediaitems" ADD FOREIGN KEY ("mediaitem_id") REFERENCES "mediaitems" ("id");

ALTER TABLE "album_mediaitems" ADD FOREIGN KEY ("album_id") REFERENCES "albums" ("id");

ALTER TABLE "album_mediaitems" ADD FOREIGN KEY ("mediaitem_id") REFERENCES "mediaitems" ("id");
