CREATE TABLE IF NOT EXISTS image_tweets (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    image_url TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_image_tweets_user_id ON image_tweets(user_id);
CREATE INDEX idx_image_tweets_created_at ON image_tweets(created_at DESC);
