CREATE schema horgh;
REVOKE CREATE ON schema horgh FROM public;

CREATE TABLE IF NOT EXISTS horgh.%s (
    id SERIAL,
    schema_name text NOT NULL,
    TABLE_NAME text NOT NULL,
    user_name text,
    action_tstamp TIMESTAMP WITH TIME zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    action TEXT NOT NULL CHECK (action IN ('I','D','U')),
    original_data text,
    new_data text,
    query text
) WITH (fillfactor=100);

REVOKE ALL ON horgh.%s FROM public;

GRANT SELECT ON horgh.%s TO public;

CREATE INDEX %s_action_idx
ON horgh.logged_actions(action);