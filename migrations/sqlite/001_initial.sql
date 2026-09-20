-- +goose Up
CREATE TABLE users (
    id BLOB(16) PRIMARY KEY NOT NULL,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT,
    name TEXT NOT NULL
);

CREATE TABLE organizations (
    id BLOB(16) PRIMARY KEY NOT NULL,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE repositories (
    id BLOB(16) PRIMARY KEY NOT NULL,
    organization_id BLOB(16) NOT NULL,
    name TEXT NOT NULL,
    identifier TEXT NOT NULL,
    FOREIGN KEY (organization_id) REFERENCES organizations (id) ON DELETE CASCADE,
    UNIQUE (organization_id, identifier)
);

CREATE TABLE organization_memberships (
    id BLOB(16) PRIMARY KEY NOT NULL,
    organization_id BLOB(16) NOT NULL,
    user_id BLOB(16) NOT NULL,
    permission TEXT NOT NULL CHECK (permission IN ('user', 'admin')),
    approved BOOLEAN NOT NULL,
    FOREIGN KEY (organization_id) REFERENCES organizations (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    UNIQUE (organization_id, user_id)
);

CREATE TABLE repository_permissions (
    id BLOB(16) PRIMARY KEY NOT NULL,
    repository_id BLOB(16) NOT NULL,
    user_id BLOB(16) NOT NULL,
    permission TEXT NOT NULL CHECK (permission IN ('user', 'admin')),
    FOREIGN KEY (repository_id) REFERENCES repositories (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    UNIQUE (repository_id, user_id)
);

CREATE TABLE vcs_identities (
    id BLOB(16) PRIMARY KEY NOT NULL,
    user_id BLOB(16) NOT NULL,
    identity TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    UNIQUE (user_id, identity)
);

CREATE TABLE personal_access_tokens (
    id BLOB(16) PRIMARY KEY NOT NULL,
    user_id BLOB(16) NOT NULL,
    name TEXT NOT NULL,
    token_hash TEXT NOT NULL UNIQUE,
    created_at DATETIME NOT NULL,
    expires_at DATETIME,
    last_used_at DATETIME,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE locks (
    id BLOB(16) PRIMARY KEY NOT NULL,
    repository_id BLOB(16) NOT NULL,
    user_id BLOB(16) NOT NULL,
    pattern TEXT NOT NULL,
    reason TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    expires_at DATETIME,
    FOREIGN KEY (repository_id) REFERENCES repositories (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
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
