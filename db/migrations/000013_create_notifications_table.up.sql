CREATE TABLE IF NOT EXISTS notifications (
    id SERIAL PRIMARY KEY,
    -- 通知を受け取るユーザー
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    -- 通知のきっかけとなった行動をしたユーザー
    actor_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    -- 通知種別: like / follow / comment
    type VARCHAR(20) NOT NULL,
    -- like / comment の対象ツイート（follow の場合は NULL）
    tweet_id INTEGER REFERENCES tweets(id) ON DELETE CASCADE,
    -- comment の場合のコメントID
    comment_id INTEGER REFERENCES comments(id) ON DELETE CASCADE,
    is_read BOOLEAN NOT NULL DEFAULT false,
    read_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT notifications_type_check CHECK (type IN ('like', 'follow', 'comment'))
);

CREATE INDEX idx_notifications_user_id_id ON notifications(user_id, id DESC);
CREATE INDEX idx_notifications_user_id_is_read ON notifications(user_id, is_read);
