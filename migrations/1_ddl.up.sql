CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS pg_trgm;

BEGIN;

DROP TYPE IF EXISTS account_type;
CREATE TYPE account_type AS ENUM ('private', 'public');

CREATE TABLE IF NOT EXISTS user_profile
(
    id           SERIAL,
    auth_id      varchar(200) COLLATE pg_catalog."default" NOT NULL UNIQUE,
    name         varchar(80) COLLATE pg_catalog."default"  NOT NULL,
    created      timestamptz                               NOT NULL DEFAULT CURRENT_TIMESTAMP,
    account_type account_type                              NOT NULL DEFAULT 'public',
    email        varchar(50) COLLATE pg_catalog."default",
    phone        varchar(30) COLLATE pg_catalog."default",
    bio          text,
    nick         varchar(80) COLLATE pg_catalog."default"  NOT NULL,
    avatar       varchar(100) COLLATE pg_catalog."default",
    updated      timestamptz                               NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT pk_user PRIMARY KEY (id)
);

-- TODO: add avatar placeholder for deleted user
INSERT INTO user_profile (id, auth_id, name, nick, email, phone, bio)
VALUES (-1, 'deleted_auth_id', 'deleted user', '[deleted]', '', '', '');

CREATE INDEX IF NOT EXISTS idx_user_profile_name_lower ON user_profile (LOWER(name));
CREATE UNIQUE INDEX IF NOT EXISTS unq_user_profile_nick ON user_profile (LOWER(nick));
CREATE UNIQUE INDEX IF NOT EXISTS unq_user_profile_phone ON user_profile (phone) WHERE TRIM(phone) <> '';

CREATE TABLE IF NOT EXISTS user_settings
(
    user_id               integer,
    push_notification     boolean NOT NULL DEFAULT true,
    follow_accepted_push  boolean NOT NULL DEFAULT true,
    push_token            varchar(200) COLLATE pg_catalog."default",
    follow_request_push   boolean NOT NULL DEFAULT true,
    new_review_push       boolean NOT NULL DEFAULT true,
    tagged_push           boolean NOT NULL DEFAULT true,
    commented_push        boolean NOT NULL DEFAULT true,
    liked_push            boolean NOT NULL DEFAULT true,
    comment_answered_push boolean NOT NULL DEFAULT true,
    CONSTRAINT pk_user_setting PRIMARY KEY (user_id),
    CONSTRAINT fk_user_setting_user_profile_id
        FOREIGN KEY (user_id)
            REFERENCES user_profile (id) ON
            DELETE CASCADE
);

DROP TYPE IF EXISTS geo_location_type;
CREATE TYPE geo_location_type AS ENUM ('google', 'here', 'foursquare');

CREATE TABLE IF NOT EXISTS place
(
    id              SERIAL,
    created         timestamptz                               NOT NULL DEFAULT CURRENT_TIMESTAMP,
    external_id     varchar(200) COLLATE pg_catalog."default" NOT NULL,
    location_type   geo_location_type                         NOT NULL,
    longitude       numeric(9, 6)                             NOT NULL,
    latitude        numeric(9, 6)                             NOT NULL,
    point           geometry(Point)                           NOT NULL,
    point_3857      geometry(Point, 3857)                     NOT NULL, -- SRID [Google Maps: 3857, Google Earth: 4326]
    title           varchar(256) COLLATE pg_catalog."default" NOT NULL,
    tags            varchar(512) COLLATE pg_catalog."default" NOT NULL,
    country         varchar(256) COLLATE pg_catalog."default" NOT NULL,
    region          varchar(256) COLLATE pg_catalog."default" NOT NULL,
    county          varchar(256) COLLATE pg_catalog."default" NOT NULL,
    city            varchar(256) COLLATE pg_catalog."default" NOT NULL,
    postal          varchar(50) COLLATE pg_catalog."default"  NOT NULL,
    address_line    varchar(256) COLLATE pg_catalog."default" NOT NULL,
    phone_number    varchar(32) COLLATE pg_catalog."default"  NOT NULL,
    website         varchar(350) COLLATE pg_catalog."default" NOT NULL,
    completed       bool                                      NOT NULL DEFAULT false,
    rating          int                                       NOT NULL DEFAULT 0,
    reviews_counter int                                       NOT NULL DEFAULT 0,
    sources         jsonb                                     NOT NULL,
    timezone        varchar(100) COLLATE pg_catalog."default" NOT NULL,
    hours           jsonb                                     NOT NULL,
    CONSTRAINT pk_place PRIMARY KEY (id),
    CONSTRAINT unq_place_external_id_location_type UNIQUE (external_id, location_type)
);

CREATE INDEX IF NOT EXISTS idx_place_point_3857 ON place USING gist (point_3857);
CREATE INDEX IF NOT EXISTS idx_place_title ON place USING gin (title gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_place_location_type_completed ON place (location_type, completed);

CREATE TABLE IF NOT EXISTS recipe
(
    id        SERIAL,
    created   timestamptz                                NOT NULL DEFAULT CURRENT_TIMESTAMP,
    title     varchar(1024) COLLATE pg_catalog."default" NOT NULL,
    user_id   int                                        NOT NULL,
    content   text                                       NOT NULL,
    completed boolean                                    NOT NULL DEFAULT false,
    CONSTRAINT pk_recipe PRIMARY KEY (id),
    CONSTRAINT fk_recipe_user_profile_user_id
        FOREIGN KEY (user_id)
            REFERENCES user_profile (id)
);

CREATE INDEX IF NOT EXISTS idx_recipe_title ON recipe (LOWER(title));

CREATE TABLE IF NOT EXISTS recipe_photo
(
    id          SERIAL,
    recipe_id   int                                       NOT NULL,
    created     timestamptz                               NOT NULL DEFAULT CURRENT_TIMESTAMP,
    photo_index smallint                                  NOT NULL,
    photo_key   varchar(120) COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT pk_recipe_photo PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS review
(
    id              SERIAL,
    user_id         int         NOT NULL,
    place_id        int,
    created         timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    content         text        NOT NULL,
    ready           boolean     NOT NULL DEFAULT false,
    rating          smallint    NOT NULL DEFAULT 0,
    likes_counter   int         NOT NULL DEFAULT 0,
    is_kitchen      bool        NOT NULL DEFAULT false,
    image_crop_mode smallint    NOT NULL DEFAULT 0,
    loop_animation  boolean     NOT NULL DEFAULT false,
    updated         timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    recipe_id       int,
    CONSTRAINT pk_review PRIMARY KEY (id),
    CONSTRAINT fk_review_user_profile_user_id
        FOREIGN KEY (user_id)
            REFERENCES user_profile (id),
    CONSTRAINT fk_review_place_place_id
        FOREIGN KEY (place_id)
            REFERENCES place (id),
    CONSTRAINT fk_review_recipe_recipe_id
        FOREIGN KEY (recipe_id)
            REFERENCES recipe (id)
);

CREATE INDEX IF NOT EXISTS idx_review_user_id ON review (user_id);

CREATE INDEX IF NOT EXISTS idx_review_place_id ON review (place_id);

CREATE TABLE IF NOT EXISTS review_like
(
    review_id int NOT NULL,
    user_id   int NOT NULL,
    CONSTRAINT unq_review_like_review_id_user_id UNIQUE (review_id, user_id),
    CONSTRAINT fk_review_like_user_profile_user_id
        FOREIGN KEY (user_id)
            REFERENCES user_profile (id) ON DELETE CASCADE,
    CONSTRAINT fk_review_like_review_review_id
        FOREIGN KEY (review_id)
            REFERENCES review (id) ON DELETE CASCADE
);

CREATE TYPE media_type AS ENUM ('image', 'video');

CREATE TABLE IF NOT EXISTS review_media
(
    id          SERIAL,
    review_id   int                                       NOT NULL,
    created     timestamptz                               NOT NULL DEFAULT CURRENT_TIMESTAMP,
    media_index smallint                                  NOT NULL,
    media_key   varchar(120) COLLATE pg_catalog."default" NOT NULL,
    datestamp   timestamptz,
    media_type  media_type                                NOT NULL DEFAULT 'image',
    CONSTRAINT pk_review_media PRIMARY KEY (id),
    CONSTRAINT unq_review_media_review_id_media_index UNIQUE (review_id, media_index)
);

CREATE INDEX IF NOT EXISTS idx_review_media_review_id ON review_media (review_id);

CREATE TABLE IF NOT EXISTS review_comment
(
    id            SERIAL,
    review_id     int         NOT NULL,
    author_id     int         NOT NULL,
    parent_id     int         NOT NULL DEFAULT 0,
    published_at  timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    content       text        NOT NULL,
    replies_count int         NOT NULL DEFAULT 0,
    CONSTRAINT pk_review_comment PRIMARY KEY (id),
    CONSTRAINT fk_review_comment_user_profile_author_id
        FOREIGN KEY (author_id)
            REFERENCES user_profile (id) ON DELETE CASCADE,
    CONSTRAINT fk_review_comment_review_review_id
        FOREIGN KEY (review_id)
            REFERENCES review (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS review_preview
(
    review_id   int                                       NOT NULL,
    media_index int,
    created     timestamptz                               NOT NULL DEFAULT CURRENT_TIMESTAMP,
    preview_id  varchar(36) COLLATE pg_catalog."default"  NOT NULL,
    deep_link   varchar(64) COLLATE pg_catalog."default"  NOT NULL,
    photo_key   varchar(120) COLLATE pg_catalog."default" NOT NULL DEFAULT '',
    video_key   varchar(120) COLLATE pg_catalog."default" NOT NULL DEFAULT ''
);

DROP TYPE IF EXISTS follower_status;
CREATE TYPE follower_status AS ENUM ('pending', 'accepted');

CREATE TABLE IF NOT EXISTS user_follower
(
    id          SERIAL,
    user_id     int             NOT NULL,
    follower_id int             NOT NULL,
    updated     timestamptz     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    status      follower_status NOT NULL DEFAULT 'pending',
    CONSTRAINT pk_user_follower PRIMARY KEY (id),
    CONSTRAINT fk_user_follower_user_profile_user_id
        FOREIGN KEY (user_id)
            REFERENCES user_profile (id) ON DELETE CASCADE,
    CONSTRAINT fk_user_follower_user_profile_follower_id
        FOREIGN KEY (follower_id)
            REFERENCES user_profile (id) ON DELETE CASCADE,
    CONSTRAINT unq_user_follower_user_id_follower_id UNIQUE (user_id, follower_id)
);

CREATE INDEX IF NOT EXISTS idx_user_follower_user_id ON user_follower (user_id);

CREATE INDEX IF NOT EXISTS idx_user_follower_follower_id ON user_follower (follower_id);

CREATE INDEX IF NOT EXISTS idx_user_follower_status ON user_follower (status);

DROP TYPE IF EXISTS complaint_types;
CREATE TYPE complaint_types AS ENUM ('review', 'user');

CREATE TABLE IF NOT EXISTS complaint
(
    id             SERIAL,
    created        timestamptz     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    author_id      int             NOT NULL,
    complaint_type complaint_types NOT NULL,
    object_id      int             NOT NULL,
    reason         text            NOT NULL,
    handled        boolean         NOT NULL DEFAULT false,
    sent           boolean         NOT NULL DEFAULT false,
    sent_at        timestamptz,
    CONSTRAINT pk_complaint PRIMARY KEY (id),
    CONSTRAINT fk_complaint_user_profile_author_id
        FOREIGN KEY (author_id)
            REFERENCES user_profile (id) ON DELETE CASCADE
);

CREATE SEQUENCE IF NOT EXISTS id_generator
    INCREMENT BY 1
    MINVALUE 1
    START WITH 1;

CREATE TABLE IF NOT EXISTS geo_hash
(
    id           SERIAL,
    created      timestamptz     NOT NULL                  DEFAULT CURRENT_TIMESTAMP,
    hash         varchar(7)      NOT NULL UNIQUE, --geo precision: 7
    center_point geometry(Point) NOT NULL,
    sync         timestamptz     NOT NULL                  DEFAULT CURRENT_TIMESTAMP,
    agg_id       varchar(128) COLLATE pg_catalog."default" DEFAULT '',
    country      varchar(60) COLLATE pg_catalog."default"  DEFAULT '',
    region       varchar(100) COLLATE pg_catalog."default" DEFAULT '',
    CONSTRAINT pk_geo_hash PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS geo_hash_place
(
    id       SERIAL,
    hash_id  int NOT NULL,
    place_id int NOT NULL,
    CONSTRAINT pk_geo_hash_place PRIMARY KEY (id),
    CONSTRAINT unq_geo_hash_place_hash_id_place_id UNIQUE (hash_id, place_id),
    CONSTRAINT fk_geo_hash_place_hash_id
        FOREIGN KEY (hash_id)
            REFERENCES geo_hash (id),
    CONSTRAINT fk_geo_hash_place_place_id
        FOREIGN KEY (place_id)
            REFERENCES place (id)
);

CREATE INDEX IF NOT EXISTS idx_geo_hash_place_hash_id ON geo_hash_place (hash_id);

CREATE INDEX IF NOT EXISTS idx_geo_hash_place_place_id ON geo_hash_place (place_id);

CREATE TABLE IF NOT EXISTS activation_code
(
    user_id         int        NOT NULL DEFAULT 0,
    code            varchar(6) NOT NULL UNIQUE,
    activated       bool       NOT NULL DEFAULT false,
    activation_time timestamptz,
    sent            bool       NOT NULL DEFAULT false,
    sent_time       timestamptz,
    sent_user_id    int
);

DROP TYPE IF EXISTS subscription_type;
CREATE TYPE subscription_type AS ENUM ('email', 'sms');

CREATE TABLE IF NOT EXISTS subscription
(
    id        SERIAL,
    created   timestamptz                              NOT NULL DEFAULT CURRENT_TIMESTAMP,
    subs_type subscription_type                        NOT NULL DEFAULT 'email',
    target    varchar(50) COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT unq_subscription_subs_type_target UNIQUE (subs_type, target)
);

CREATE TABLE IF NOT EXISTS activity_log
(
    id          SERIAL,
    created     timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_id     int         NOT NULL,
    action_type int         NOT NULL DEFAULT 0,
    doer_id     int         NOT NULL,
    root_id     int         NOT NULL DEFAULT 0,
    object_id   int         NOT NULL,
    parent_id   int         NOT NULL DEFAULT 0,
    CONSTRAINT pk_activity_log PRIMARY KEY (id),
    CONSTRAINT fk_activity_log_user_profile_doer_id
        FOREIGN KEY (doer_id)
            REFERENCES user_profile (id) ON DELETE CASCADE,
    CONSTRAINT fk_activity_log_user_profile_user_id
        FOREIGN KEY (user_id)
            REFERENCES user_profile (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_activity_log_action_type ON activity_log (action_type);

CREATE INDEX IF NOT EXISTS idx_activity_log_doer_id ON activity_log (doer_id);

CREATE INDEX IF NOT EXISTS idx_activity_log_object_id ON activity_log (object_id);

CREATE TABLE IF NOT EXISTS hashtag
(
    id       SERIAL,
    value    varchar(150) NOT NULL UNIQUE,
    featured bool         NOT NULL DEFAULT false,
    CONSTRAINT pk_hashtag PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS hashtag_review
(
    id         SERIAL,
    hashtag_id int NOT NULL,
    review_id  int NOT NULL,
    counter    int NOT NULL DEFAULT 1,
    CONSTRAINT unq_hashtag_review_hashtag_id_review_id UNIQUE (hashtag_id, review_id),
    CONSTRAINT fk_hashtag_review_hashtag_id
        FOREIGN KEY (hashtag_id)
            REFERENCES hashtag (id),
    CONSTRAINT fk_hashtag_review_review_id
        FOREIGN KEY (review_id)
            REFERENCES review (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS media
(
    id          SERIAL,
    created     timestamptz                               NOT NULL DEFAULT CURRENT_TIMESTAMP,
    object_id   int                                       NOT NULL,
    object_type smallint                                  NOT NULL,
    media_key   varchar(120) COLLATE pg_catalog."default" NOT NULL,
    media_type  varchar(32) COLLATE pg_catalog."default"  NOT NULL,
    CONSTRAINT pk_media PRIMARY KEY (id),
    CONSTRAINT unq_object_type_media_key UNIQUE (object_type, media_key)
);

CREATE TABLE IF NOT EXISTS collection
(
    id          SERIAL,
    author_id   int                                       NOT NULL,
    name        varchar(20) COLLATE pg_catalog."default"  NOT NULL,
    cover_image varchar(120) COLLATE pg_catalog."default" NOT NULL,
    created     timestamptz                               NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT pk_collection PRIMARY KEY (id),
    CONSTRAINT fk_collection_user_profile_author_id
        FOREIGN KEY (author_id)
            REFERENCES user_profile (id)
);

CREATE TABLE IF NOT EXISTS collection_review
(
    collection_id int NOT NULL,
    review_id     int NOT NULL,
    CONSTRAINT unq_collection_review_collection_id_review_id UNIQUE (collection_id, review_id),
    CONSTRAINT fk_collection_review_collection_id
        FOREIGN KEY (collection_id)
            REFERENCES collection (id) ON DELETE CASCADE,
    CONSTRAINT fk_collection_review_review_id
        FOREIGN KEY (review_id)
            REFERENCES review (id) ON DELETE CASCADE
);

CREATE OR REPLACE VIEW unique_place_ready_reviews_reviewers_view AS
(
SELECT DISTINCT u.id           AS reviewer_id,
                u.auth_id      AS reviewer_auth_id,
                u.name         AS reviewer_name,
                u.created      AS reviewer_created,
                u.account_type AS reviewer_account_type,
                u.email        AS reviewer_email,
                u.bio          AS reviewer_bio,
                u.nick         AS reviewer_nick,
                u.avatar       AS reviewer_avatar,
                p.id           AS place_id,
                u.account_type AS review_visibility
FROM place p
         JOIN review r ON r.place_id = p.id AND r.ready = true
         JOIN user_profile u ON r.user_id = u.id
    );

CREATE OR REPLACE VIEW ready_reviews_with_users_view AS
(
SELECT r.id              AS review_id,
       r.user_id         AS review_user_id,
       r.place_id        AS place_id,
       r.created         AS review_created,
       r.content         AS review_content,
       r.ready           AS review_ready,
       ru.account_type   AS review_visibility,
       r.rating          AS review_rating,
       r.is_kitchen      AS review_is_kitchen,
       r.likes_counter   AS review_likes_counter,
       ru.name           AS review_user_name,
       ru.nick           AS review_user_nick,
       ru.account_type   AS review_user_account_type,
       ru.bio            AS review_user_bio,
       ru.avatar         AS review_user_avatar,
       uf.follower_id    AS review_follower_id,
       r.image_crop_mode AS review_image_crop_mode
FROM review r
         JOIN user_profile ru ON r.user_id = ru.id
         LEFT JOIN user_follower uf ON ru.id = uf.user_id AND uf.status = 'accepted'
WHERE r.ready = true
    );

CREATE OR REPLACE VIEW tagged_latest_ready_reviews_view AS
(
SELECT h.value           AS review_tag,
       r.id              AS review_id,
       r.user_id         AS review_user_id,
       r.recipe_id       AS recipe_id,
       r.place_id        AS place_id,
       r.created         AS review_created,
       r.updated         AS review_updated,
       r.content         AS review_content,
       r.ready           AS review_ready,
       ru.account_type   AS review_visibility,
       r.rating          AS review_rating,
       r.is_kitchen      AS review_is_kitchen,
       r.likes_counter   AS review_likes_counter,
       r.image_crop_mode AS review_image_crop_mode,
       r.loop_animation  AS review_loop_animation,
       uf.follower_id    AS review_follower_id
FROM hashtag_review hr
         JOIN review r ON r.id = hr.review_id
         JOIN hashtag h ON h.id = hr.hashtag_id
         JOIN user_profile ru ON r.user_id = ru.id
         LEFT JOIN user_follower uf ON ru.id = uf.user_id AND uf.status = 'accepted'
WHERE r.ready
ORDER BY r.id DESC
    );

CREATE OR REPLACE VIEW most_popular_public_hash_tags_view AS
(
SELECT h.id       AS id,
       h.value    AS value,
       h.featured AS featured,
       count(*)   AS reviews_count
FROM hashtag h
         JOIN hashtag_review hr ON h.id = hr.hashtag_id
         JOIN review r on hr.review_id = r.id
         JOIN user_profile up on r.user_id = up.id
WHERE up.account_type = 'public'
GROUP BY h.id
ORDER BY count(*) DESC
    );

CREATE OR REPLACE VIEW tagged_public_reviews_media_view AS
(
SELECT hr.hashtag_id  AS hashtag_id,
       rm.id          AS id,
       rm.review_id   AS review_id,
       rm.created     AS created,
       rm.media_index AS media_index,
       rm.media_key   AS media_key,
       rm.datestamp   AS datestamp,
       rm.media_type  AS media_type
FROM review_media rm
         JOIN hashtag_review hr on rm.review_id = hr.review_id
         JOIN review r on hr.review_id = r.id
         JOIN user_profile up on r.user_id = up.id
WHERE up.account_type = 'public'
ORDER BY media_index, review_id DESC
    );

CREATE TABLE IF NOT EXISTS notification_event
(
    id         SERIAL,
    created    timestamptz                              NOT NULL DEFAULT CURRENT_TIMESTAMP,
    type       varchar(64) COLLATE pg_catalog."default" NOT NULL,
    data       jsonb                                    NOT NULL,
    attributes jsonb                                    NOT NULL
);

COMMIT;