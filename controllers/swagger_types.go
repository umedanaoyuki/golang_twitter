package controllers

// Swagger 用のレスポンス型（実装の gin.H 等と対応）

type HealthCheckResponse struct {
	Status string `json:"status" example:"ok"`
}

type StatusOKResponse struct {
	Status string `json:"status" example:"ok"`
}

type MessageResponse struct {
	Message string `json:"message" example:"処理が完了しました"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"エラーメッセージ"`
}

type RegisterResponse struct {
	Message string               `json:"message" example:"ユーザー登録が完了しました"`
	User    RegisterUserResponse `json:"user"`
}

type RegisterUserResponse struct {
	ID        int32  `json:"id" example:"1"`
	Email     string `json:"email" example:"user@example.com"`
	CreatedAt string `json:"created_at" example:"2024-01-01T00:00:00Z"`
}

type LoginResponse struct {
	Message string    `json:"message" example:"ログインに成功しました"`
	User    LoginUser `json:"user"`
}

type LoginUser struct {
	ID        int32  `json:"id" example:"1"`
	Email     string `json:"email" example:"user@example.com"`
	IsActive  bool   `json:"is_active" example:"true"`
	CreatedAt string `json:"created_at" example:"2024-01-01T00:00:00Z"`
}

type SwaggerTweet struct {
	ID           int32  `json:"id" example:"1"`
	UserID       int32  `json:"user_id" example:"1"`
	Content      string `json:"content" example:"Hello, world!"`
	CreatedAt    string `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt    string `json:"updated_at" example:"2024-01-01T00:00:00Z"`
	LikeCount    int64  `json:"like_count" example:"0"`
	RetweetCount int64  `json:"retweet_count" example:"0"`
}

type CreateTweetResponse struct {
	Tweet CreateTweetItem `json:"tweet"`
}

type CreateTweetItem struct {
	UserID  int32  `json:"user_id" example:"1"`
	Content string `json:"content" example:"Hello, world!"`
}

type CreateImageTweetResponse struct {
	Tweet CreateImageTweetItem `json:"tweet"`
}

type CreateImageTweetItem struct {
	UserID   int32  `json:"user_id" example:"1"`
	ImageURL string `json:"image_url" example:"https://example.com/image.jpg"`
}

type PresignImageTweetResponse struct {
	Key       string `json:"key" example:"uploads/1_abc123.jpg"`
	UploadURL string `json:"upload_url" example:"https://s3.example.com/bucket/uploads/1_abc123.jpg?X-Amz-Signature=..."`
	PublicURL string `json:"public_url" example:"https://example.com/bucket/uploads/1_abc123.jpg"`
}

type GetTweetResponse struct {
	Tweet SwaggerTweet `json:"tweet"`
}

type GetUserTweetsResponse struct {
	Tweets     []SwaggerTweet `json:"tweets"`
	NextCursor *int32         `json:"next_cursor"`
	HasMore    bool           `json:"has_more" example:"false"`
}

type GetCurrentUserTweetsResponse struct {
	User       SwaggerUserDetail `json:"user"`
	Tweets     []SwaggerTweet    `json:"tweets"`
	NextCursor *int32            `json:"next_cursor"`
	HasMore    bool              `json:"has_more" example:"false"`
}

type SwaggerUserDetail struct {
	ID        int32  `json:"id" example:"1"`
	Email     string `json:"email" example:"user@example.com"`
	IsActive  bool   `json:"is_active" example:"true"`
	CreatedAt string `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt string `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

type GetUserResponse struct {
	User SwaggerUserDetail `json:"user"`
}

type SwaggerBookmark struct {
	ID        int32  `json:"id" example:"1"`
	UserID    int32  `json:"user_id" example:"1"`
	TweetID   int32  `json:"tweet_id" example:"1"`
	CreatedAt string `json:"created_at" example:"2024-01-01T00:00:00Z"`
	Content   string `json:"content" example:"ツイート本文"`
}

type SwaggerGroup struct {
	ID        int32  `json:"id" example:"1"`
	Name      string `json:"name" example:"グループ名"`
	UserID    int32  `json:"user_id" example:"1"`
	CreatedAt string `json:"created_at" example:"2024-01-01T00:00:00Z"`
}

type CreateGroupResponse struct {
	Group SwaggerGroup `json:"group"`
}

type CreateMessageResponse struct {
	Message CreateMessageItem `json:"message"`
}

type CreateMessageItem struct {
	UserID  int32  `json:"user_id" example:"1"`
	GroupID int32  `json:"group_id" example:"1"`
	Content string `json:"content" example:"メッセージ本文"`
}

type SwaggerMessage struct {
	ID        int32  `json:"id" example:"1"`
	UserID    int32  `json:"user_id" example:"1"`
	GroupID   int32  `json:"group_id" example:"1"`
	Content   string `json:"content" example:"メッセージ本文"`
	CreatedAt string `json:"created_at" example:"2024-01-01T00:00:00Z"`
}

type GetMessagesResponse struct {
	Messages []SwaggerMessage `json:"messages"`
}

type SwaggerRetweetItem struct {
	Tweet SwaggerTweet `json:"tweet"`
}

type GetUserRetweetsResponse struct {
	Retweets   []SwaggerRetweetItem `json:"retweets"`
	NextCursor *int32               `json:"next_cursor"`
	HasMore    bool                 `json:"has_more" example:"false"`
}

type SwaggerFollow struct {
	ID             int32  `json:"id" example:"1"`
	UserID         int32  `json:"user_id" example:"1"`
	FollowedUserID int32  `json:"followed_user_id" example:"2"`
	CreatedAt      string `json:"created_at" example:"2024-01-01T00:00:00Z"`
}

type GetFollowersResponse struct {
	Followers  []SwaggerFollow `json:"followers"`
	NextCursor *int32          `json:"next_cursor"`
	HasMore    bool            `json:"has_more" example:"false"`
}

type GetFollowingResponse struct {
	Following  []SwaggerFollow `json:"following"`
	NextCursor *int32          `json:"next_cursor"`
	HasMore    bool            `json:"has_more" example:"false"`
}

type SwaggerComment struct {
	ID        int32  `json:"id" example:"1"`
	UserID    int32  `json:"user_id" example:"1"`
	TweetID   int32  `json:"tweet_id" example:"1"`
	Content   string `json:"content" example:"いいツイートですね！"`
	CreatedAt string `json:"created_at" example:"2024-01-01T00:00:00Z"`
}

type CreateCommentResponse struct {
	Comment SwaggerComment `json:"comment"`
}

type GetCommentsResponse struct {
	Comments   []SwaggerComment `json:"comments"`
	NextCursor *int32           `json:"next_cursor"`
	HasMore    bool             `json:"has_more" example:"false"`
}
type SwaggerUserProfile struct {
	ID        int32  `json:"id" example:"1"`
	UserID    int32  `json:"user_id" example:"1"`
	Name      string `json:"name" example:"Taro007"`
	Bio       string `json:"bio" example:"自己紹介文のサンプル文です"`
	ImageURL  string `json:"image_url" example:"https://example.com/avatar.png"`
	Location  string `json:"location" example:"東京"`
	CreatedAt string `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt string `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

type UserProfileResponse struct {
	Profile SwaggerUserProfile `json:"profile"`
}

type PresignProfileImageResponse struct {
	Key       string `json:"key" example:"uploads/1_abc123.jpg"`
	UploadURL string `json:"upload_url" example:"https://s3.example.com/bucket/uploads/1_abc123.jpg?X-Amz-Signature=..."`
	PublicURL string `json:"public_url" example:"https://example.com/bucket/uploads/1_abc123.jpg"`
}
