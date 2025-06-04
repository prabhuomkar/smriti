CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE album_mediaitems (
    mediaitem_id uuid NOT NULL,
    album_id uuid NOT NULL
);

CREATE TABLE albums (
    id uuid NOT NULL,
    user_id text,
    name text,
    description text,
    is_shared boolean DEFAULT false,
    is_hidden boolean DEFAULT false,
    mediaitems_count bigint DEFAULT 0,
    cover_mediaitem_id uuid,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);

CREATE TABLE jobs (
    id uuid NOT NULL,
    user_id text,
    status text,
    components text,
    last_mediaitem_id text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);

CREATE TABLE mediaitem_embeddings (
    mediaitem_id uuid,
    embedding vector
);

CREATE TABLE mediaitem_faces (
    id uuid NOT NULL,
    mediaitem_id uuid,
    people_id uuid,
    embedding vector,
    thumbnail text
);

CREATE TABLE mediaitems (
    id uuid NOT NULL,
    user_id text,
    filename text,
    hash text,
    description text,
    mime_type text,
    source_url text,
    preview_url text,
    thumbnail_url text,
    placeholder text,
    is_favourite boolean DEFAULT false,
    is_hidden boolean DEFAULT false,
    is_deleted boolean DEFAULT false,
    status text,
    mediaitem_type text,
    mediaitem_category text,
    width bigint,
    height bigint,
    creation_time timestamp with time zone,
    camera_make text,
    camera_model text,
    focal_length text,
    aperture_fnumber text,
    iso_equivalent text,
    exposure_time text,
    latitude numeric,
    longitude numeric,
    fps text,
    exif_data text,
    keywords text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);

CREATE TABLE people (
    id uuid NOT NULL,
    user_id text,
    name text,
    is_hidden boolean DEFAULT false,
    cover_mediaitem_id uuid,
    cover_mediaitem_face_id uuid,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);

CREATE TABLE people_mediaitems (
    mediaitem_id uuid NOT NULL,
    people_id uuid NOT NULL
);

CREATE TABLE place_mediaitems (
    mediaitem_id uuid NOT NULL,
    place_id uuid NOT NULL
);

CREATE TABLE places (
    id uuid NOT NULL,
    user_id text,
    name text,
    postcode text,
    country text,
    locality text,
    area text,
    is_hidden boolean DEFAULT false,
    cover_mediaitem_id uuid,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);

CREATE TABLE thing_mediaitems (
    mediaitem_id uuid NOT NULL,
    thing_id uuid NOT NULL
);

CREATE TABLE things (
    id uuid NOT NULL,
    user_id text,
    name text,
    is_hidden boolean DEFAULT false,
    cover_mediaitem_id uuid,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);

CREATE TABLE users (
    id uuid NOT NULL,
    name text,
    username text,
    password text,
    features text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);

ALTER TABLE ONLY album_mediaitems
ADD CONSTRAINT album_mediaitems_pkey PRIMARY KEY (mediaitem_id, album_id);

ALTER TABLE ONLY albums
ADD CONSTRAINT albums_pkey PRIMARY KEY (id);

ALTER TABLE ONLY jobs
ADD CONSTRAINT jobs_pkey PRIMARY KEY (id);

ALTER TABLE ONLY mediaitem_faces
ADD CONSTRAINT mediaitem_faces_pkey PRIMARY KEY (id);

ALTER TABLE ONLY mediaitems
ADD CONSTRAINT mediaitems_pkey PRIMARY KEY (id);

ALTER TABLE ONLY people_mediaitems
ADD CONSTRAINT people_mediaitems_pkey PRIMARY KEY (mediaitem_id, people_id);

ALTER TABLE ONLY people
ADD CONSTRAINT people_pkey PRIMARY KEY (id);

ALTER TABLE ONLY place_mediaitems
ADD CONSTRAINT place_mediaitems_pkey PRIMARY KEY (mediaitem_id, place_id);

ALTER TABLE ONLY places
ADD CONSTRAINT places_pkey PRIMARY KEY (id);

ALTER TABLE ONLY thing_mediaitems
ADD CONSTRAINT thing_mediaitems_pkey PRIMARY KEY (mediaitem_id, thing_id);

ALTER TABLE ONLY things
ADD CONSTRAINT things_pkey PRIMARY KEY (id);

ALTER TABLE ONLY users
ADD CONSTRAINT uni_users_username UNIQUE (username);

ALTER TABLE ONLY users
ADD CONSTRAINT users_pkey PRIMARY KEY (id);

CREATE UNIQUE INDEX idx_albums_id ON albums USING btree (id);

CREATE UNIQUE INDEX idx_jobs_id ON jobs USING btree (id);

CREATE UNIQUE INDEX idx_mediaitem_faces_id ON mediaitem_faces USING btree (id);

CREATE UNIQUE INDEX idx_mediaitems_id ON mediaitems USING btree (id);

CREATE UNIQUE INDEX idx_mediaitems_user_id_hash ON mediaitems USING btree (user_id, hash);

CREATE UNIQUE INDEX idx_people_id ON people USING btree (id);

CREATE UNIQUE INDEX idx_places_id ON places USING btree (id);

CREATE UNIQUE INDEX idx_things_id ON things USING btree (id);

ALTER TABLE ONLY album_mediaitems
ADD CONSTRAINT fk_album_mediaitems_album FOREIGN KEY (album_id) REFERENCES albums(id) ON DELETE CASCADE;

ALTER TABLE ONLY album_mediaitems
ADD CONSTRAINT fk_album_mediaitems_media_item FOREIGN KEY (mediaitem_id) REFERENCES mediaitems(id) ON DELETE CASCADE;

ALTER TABLE ONLY albums
ADD CONSTRAINT fk_albums_cover_media_item FOREIGN KEY (cover_mediaitem_id) REFERENCES mediaitems(id);

ALTER TABLE ONLY mediaitem_embeddings
ADD CONSTRAINT fk_mediaitems_embeddings FOREIGN KEY (mediaitem_id) REFERENCES mediaitems(id) ON DELETE CASCADE;

ALTER TABLE ONLY mediaitem_faces
ADD CONSTRAINT fk_mediaitems_faces FOREIGN KEY (mediaitem_id) REFERENCES mediaitems(id) ON DELETE CASCADE;

ALTER TABLE ONLY people
ADD CONSTRAINT fk_people_cover_media_item FOREIGN KEY (cover_mediaitem_id) REFERENCES mediaitems(id);

ALTER TABLE ONLY people
ADD CONSTRAINT fk_people_cover_media_item_face FOREIGN KEY (cover_mediaitem_face_id) REFERENCES mediaitem_faces(id);

ALTER TABLE ONLY people_mediaitems
ADD CONSTRAINT fk_people_mediaitems_media_item FOREIGN KEY (mediaitem_id) REFERENCES mediaitems(id) ON DELETE CASCADE;

ALTER TABLE ONLY people_mediaitems
ADD CONSTRAINT fk_people_mediaitems_people FOREIGN KEY (people_id) REFERENCES people(id) ON DELETE CASCADE;

ALTER TABLE ONLY place_mediaitems
ADD CONSTRAINT fk_place_mediaitems_media_item FOREIGN KEY (mediaitem_id) REFERENCES mediaitems(id) ON DELETE CASCADE;

ALTER TABLE ONLY place_mediaitems
ADD CONSTRAINT fk_place_mediaitems_place FOREIGN KEY (place_id) REFERENCES places(id) ON DELETE CASCADE;

ALTER TABLE ONLY places
ADD CONSTRAINT fk_places_cover_media_item FOREIGN KEY (cover_mediaitem_id) REFERENCES mediaitems(id);

ALTER TABLE ONLY thing_mediaitems
ADD CONSTRAINT fk_thing_mediaitems_media_item FOREIGN KEY (mediaitem_id) REFERENCES mediaitems(id) ON DELETE CASCADE;

ALTER TABLE ONLY thing_mediaitems
ADD CONSTRAINT fk_thing_mediaitems_thing FOREIGN KEY (thing_id) REFERENCES things(id) ON DELETE CASCADE;

ALTER TABLE ONLY things
ADD CONSTRAINT fk_things_cover_media_item FOREIGN KEY (cover_mediaitem_id) REFERENCES mediaitems(id);