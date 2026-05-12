CREATE TABLE IF NOT EXISTS follows (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    followed_user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT follows_user_id_followed_user_id_unique UNIQUE (user_id, followed_user_id)
);

CREATE INDEX idx_follows_user_id ON follows(user_id);
CREATE INDEX idx_follows_followed_user_id ON follows(followed_user_id);
CREATE INDEX idx_follows_created_at ON follows(created_at DESC);