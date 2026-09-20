-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT,
    name TEXT NOT NULL
);

CREATE TABLE organizations (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE repositories (
    id UUID PRIMARY KEY,
    organization_id UUID NOT NULL REFERENCES organizations (id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    identifier TEXT NOT NULL,
    UNIQUE (organization_id, identifier)
);

CREATE TABLE organization_memberships (
    id UUID PRIMARY KEY,
    organization_id UUID NOT NULL REFERENCES organizations (id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    permission TEXT NOT NULL CHECK (permission IN ('user', 'admin')),
    approved BOOLEAN NOT NULL,
    UNIQUE (organization_id, user_id)
);

CREATE TABLE repository_permissions (
    id UUID PRIMARY KEY,
    repository_id UUID NOT NULL REFERENCES repositories (id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    permission TEXT NOT NULL CHECK (permission IN ('user', 'admin')),
    UNIQUE (repository_id, user_id)
);

CREATE TABLE vcs_identities (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    identity TEXT NOT NULL,
    UNIQUE (user_id, identity)
);

CREATE TABLE personal_access_tokens (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    token_hash TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL,
    expires_at TIMESTAMPTZ,
    last_used_at TIMESTAMPTZ
);

CREATE TABLE locks (
    id UUID PRIMARY KEY,
    repository_id UUID NOT NULL REFERENCES repositories (id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    pattern TEXT NOT NULL,
    reason TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    expires_at TIMESTAMPTZ
);

CREATE INDEX idx_repositories_organization_id ON repositories (organization_id);
CREATE INDEX idx_organization_memberships_user_id ON organization_memberships (user_id);
CREATE INDEX idx_repository_permissions_user_id ON repository_permissions (user_id);
CREATE INDEX idx_vcs_identities_user_id ON vcs_identities (user_id);
CREATE INDEX idx_personal_access_tokens_user_id ON personal_access_tokens (user_id);
CREATE INDEX idx_locks_repository_id ON locks (repository_id);
CREATE INDEX idx_locks_user_id ON locks (user_id);

-- +goose Down
DROP TABLE locks;
DROP TABLE personal_access_tokens;
DROP TABLE vcs_identities;
DROP TABLE repository_permissions;
DROP TABLE organization_memberships;
DROP TABLE repositories;
DROP TABLE organizations;
DROP TABLE users;
