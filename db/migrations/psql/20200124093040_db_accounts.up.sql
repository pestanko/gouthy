-- Your SQL goes here

--- DB FUNCTIONS FUNCTIONS
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE OR REPLACE FUNCTION trigger_set_updated()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

--- USERS RELATED
CREATE TABLE IF NOT EXISTS Users
(
    id         uuid      NOT NULL DEFAULT uuid_generate_v4(), -- will be same as entity_id
    username   VARCHAR   NOT NULL UNIQUE,
    password   VARCHAR   NOT NULL,
    name       VARCHAR   NULL,
    email      VARCHAR   NULL,
    user_type  VARCHAR   NULL     DEFAULT 'public',           -- PUBLIC, INTERNAL
    state      VARCHAR   NOT NULL DEFAULT 'created',          -- CREATED, ACTIVE, LOCKED, DELETED
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL,
    PRIMARY KEY (id)
);

CREATE UNIQUE INDEX Users_Username_Idx ON users (username);

CREATE TRIGGER SetUpdated_Users
    BEFORE UPDATE
    ON Users
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_updated();

CREATE TABLE IF NOT EXISTS UserSecrets
(
    id         uuid      NOT NULL DEFAULT uuid_generate_v4(),
    name       VARCHAR   NOT NULL,
    value      VARCHAR   NOT NULL,
    user_id    uuid      NOT NULL REFERENCES Users (id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP NULL,
    PRIMARY KEY (id, user_id)
);

CREATE TRIGGER SetUpdated_UserSecrets
    BEFORE UPDATE
    ON UserSecrets
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_updated();

-- Applications

CREATE TABLE IF NOT EXISTS Applications
(
    id               uuid      NOT NULL DEFAULT uuid_generate_v4(),
    client_id        VARCHAR   NOT NULL UNIQUE,
    codename         VARCHAR   NOT NULL UNIQUE,
    description      VARCHAR   NULL,
    name             VARCHAR   NOT NULL,
    state            VARCHAR   NOT NULL DEFAULT 'created', -- CREATED, ACTIVE, LOCKED, DELETED
    type             VARCHAR   NOT NULL DEFAULT 'public',  -- INTERNAL, PUBLIC, CONFIDENTIAL
    redirect_uris    TEXT      NULL,
    available_scopes TEXT      NULL,

    created_at       TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at       TIMESTAMP NULL,
    PRIMARY KEY (id)
);

CREATE UNIQUE INDEX Applications_codename_Idx ON Applications (codename);
CREATE UNIQUE INDEX Applications_client_id_Idx ON Applications (client_id);


CREATE TRIGGER SetUpdated_Applications
    BEFORE UPDATE
    ON Applications
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_updated();

CREATE TABLE IF NOT EXISTS ApplicationSecrets
(
    id             uuid      NOT NULL DEFAULT uuid_generate_v4(),
    value          VARCHAR   NOT NULL,
    application_id uuid      NOT NULL REFERENCES Applications (id),

    created_at     TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at     TIMESTAMP NULL,
    PRIMARY KEY (id, application_id)
);

CREATE TRIGGER SetUpdated_ApplicationSecrets
    BEFORE UPDATE
    ON ApplicationSecrets
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_updated();