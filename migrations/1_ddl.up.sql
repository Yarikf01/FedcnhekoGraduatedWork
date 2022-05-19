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
COMMIT;