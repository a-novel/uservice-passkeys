CREATE TABLE passkeys (
    id UUID PRIMARY KEY,

    namespace TEXT NOT NULL,
    encrypted_key TEXT NOT NULL,
    reward json,

    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ
);

--bun:split

CREATE VIEW active_passkeys AS
SELECT * FROM passkeys
WHERE passkeys.expires_at IS NULL OR passkeys.expires_at >= now();
