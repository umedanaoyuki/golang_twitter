CREATE TYPE group_member_role AS ENUM ('admin', 'member');
CREATE TYPE group_member_status AS ENUM ('invited', 'accepted');

CREATE TABLE IF NOT EXISTS group_members (
    id SERIAL PRIMARY KEY,
    group_id INTEGER NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role group_member_role NOT NULL DEFAULT 'member',
    status group_member_status NOT NULL DEFAULT 'invited',
    invited_by INTEGER REFERENCES users(id) ON DELETE SET NULL,
    accepted_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT group_members_group_id_user_id_unique UNIQUE (group_id, user_id)
);

CREATE INDEX idx_group_members_group_id ON group_members(group_id);
CREATE INDEX idx_group_members_user_id ON group_members(user_id);
CREATE INDEX idx_group_members_group_id_status ON group_members(group_id, status);
CREATE INDEX idx_group_members_created_at ON group_members(created_at DESC);
