--- Account is the "Accountable entity", can be ether USER or MACHINE
CREATE TABLE IF NOT EXISTS Entities
(
    id           uuid               DEFAULT uuid_generate_v4(),
    entity_type  VARCHAR   NOT NULL DEFAULT 'user',
    entity_state VARCHAR   NOT NULL DEFAULT 'created',

    created_at   TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

CREATE TRIGGER SetUpdated_Entities
    BEFORE UPDATE
    ON Entities
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_updated();

CREATE TABLE IF NOT EXISTS Secrets
(
    id         uuid      NOT NULL DEFAULT uuid_generate_v4(),
    name       VARCHAR   NOT NULL,
    value      VARCHAR   NOT NULL,
    entity_id  uuid      NOT NULL REFERENCES Entities (id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP NULL,
    PRIMARY KEY (id, entity_id)
);

CREATE TRIGGER SetUpdated_Secrets
    BEFORE UPDATE
    ON Secrets
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_updated();

CREATE TABLE IF NOT EXISTS LoginAudits
(
    id           uuid      NOT NULL DEFAULT uuid_generate_v4(),
    login_method VARCHAR   NOT NULL,
    entity_id    uuid      NOT NULL REFERENCES Entities (id),
    created_at   TIMESTAMP NOT NULL DEFAULT NOW(),
    message      TEXT      NULL,
    ip           VARCHAR   NULL,
    ua           VARCHAR   NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS EntityStateAudit
(
    prev_state VARCHAR   NOT NULL,
    curr_state VARCHAR   NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    entity_id  uuid      NOT NULL REFERENCES Entities (id),
    updated_by uuid      NOT NULL REFERENCES Entities (id),
    PRIMARY KEY (entity_id, created_at, prev_state, curr_state)
);

--- SPECIFIC ENTITIES - USERS, MACHINES

CREATE TABLE IF NOT EXISTS Users
(
    id         uuid      NOT NULL, -- will be same as entity_id
    username   VARCHAR   NOT NULL UNIQUE,
    password   VARCHAR   NOT NULL,
    name       VARCHAR   NULL,
    email      VARCHAR   NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

CREATE TRIGGER SetUpdated_Users
    BEFORE UPDATE
    ON Users
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_updated();

CREATE TABLE IF NOT EXISTS Machines
(
    id        uuid    NOT NULL,
    codename  VARCHAR NOT NULL UNIQUE,
    name      VARCHAR NULL,
    PRIMARY KEY (id)
);

CREATE TRIGGER SetUpdated_Machines
    BEFORE UPDATE
    ON Machines
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_updated();

-- Helper tables

CREATE TABLE IF NOT EXISTS AutomaticSecurityCodes
(
    id         uuid               DEFAULT uuid_generate_v4(),
    code       VARCHAR   NOT NULL,
    entity_id  uuid      NOT NULL REFERENCES Entities (id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    used_at    timestamp NULL,
    primary key (entity_id, id)
);


