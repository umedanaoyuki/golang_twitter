CREATE TABLE IF NOT EXISTS comments (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tweet_id INTEGER NOT NULL REFERENCES tweets(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT comments_user_id_tweet_id_unique UNIQUE (user_id, tweet_id)
);

CREATE INDEX idx_comments_user_id ON comments(user_id);
CREATE INDEX idx_comments_tweet_id ON comments(tweet_id);
CREATE INDEX idx_comments_created_at ON comments(created_at DESC);