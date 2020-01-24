
--- Account is the "Accountable entity", can be ether USER or MACHINE
CREATE TABLE IF NOT EXISTS Accounts
(
    id            uuid               DEFAULT uuid_generate_v4(),
    account_type  VARCHAR   NOT NULL DEFAULT 'user',
    account_id    uuid      NOT NULL,
    account_state VARCHAR   NOT NULL DEFAULT 'created',

    created_at    TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

CREATE IF NOT EXISTS TRIGGER SetUpdated_Clients
    BEFORE UPDATE
    ON Accounts
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_updated();


CREATE TABLE IF NOT EXISTS Secrets
(
    id         uuid      NOT NULL DEFAULT uuid_generate_v4(),
    name       VARCHAR   NOT NULL,
    value      VARCHAR   NOT NULL,
    account_id uuid      NOT NULL REFERENCES Accounts (id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP NULL,
    PRIMARY KEY (id, account_id)
);

CREATE TRIGGER IF NOT EXISTS SetUpdated_Secrets
    BEFORE UPDATE
    ON Secrets
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_updated();

CREATE TABLE IF NOT EXISTS LastLoginAudit
(
    id              uuid      NOT NULL DEFAULT uuid_generate_v4(),
    login_method    VARCHAR   NOT NULL,
    client_id       uuid      NOT NULL REFERENCES Accounts (id),
    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    message         TEXT      NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS AccountStateAudit
(
    prev_state VARCHAR   NOT NULL,
    curr_state VARCHAR   NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    account_id uuid      NOT NULL REFERENCES Accounts(id)
    updated_by uuid      NOT NULL REFERENCES Accounts (id)
    PRIMARY KEY (account_id, created_at, prev_state, curr_state)
);


--- SPECIFIC ACCOUNTS

CREATE TABLE IF NOT EXISTS Users
(
    id          uuid DEFAULT uuid_generate_v4(),
    username    VARCHAR NOT NULL UNIQUE,
    password    VARCHAR NOT NULL,
    name        VARCHAR NULL,
    email       VARCHAR NULL,
    account_id  uuid    NOT NULL REFERENCES Accounts (id),
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS Machines
(
    id          uuid DEFAULT uuid_generate_v4(),
    codename    VARCHAR NOT NULL UNIQUE,
    name        VARCHAR NULL,
    account_id  uuid NOT NULL REFERENCES Accounts (id),
);

