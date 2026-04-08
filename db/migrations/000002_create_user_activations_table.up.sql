CREATE TABLE IF NOT EXISTS user_activations (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) UNIQUE NOT NULL,
    expired_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_user_activations_user_id ON user_activations(user_id);
CREATE INDEX idx_user_activations_token ON user_activations(token);
CREATE INDEX idx_user_activations_expired_at ON user_activations(expired_at);
