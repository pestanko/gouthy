CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

--- DB FUNCTIONS

CREATE OR REPLACE FUNCTION trigger_set_updated()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
