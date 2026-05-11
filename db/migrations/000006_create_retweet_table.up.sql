CREATE TABLE IF NOT EXISTS retweets (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tweet_id INTEGER NOT NULL REFERENCES tweets(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT retweets_user_id_tweet_id_unique UNIQUE (user_id, tweet_id)
);

CREATE INDEX idx_retweets_user_id ON retweets(user_id);
CREATE INDEX idx_retweets_tweet_id ON retweets(tweet_id);
CREATE INDEX idx_retweets_created_at ON retweets(created_at DESC);