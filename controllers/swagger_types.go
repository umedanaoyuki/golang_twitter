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
	ID           int32   `json:"id" example:"1"`
	UserID       int32   `json:"user_id" example:"1"`
	Content      string  `json:"content" example:"Hello, world!"`
	ImageURL     *string `json:"image_url,omitempty" example:"http://localhost:8080/uploads/1_abc.jpg"`
	CreatedAt    string  `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt    string  `json:"updated_at" example:"2024-01-01T00:00:00Z"`
	LikeCount    int64   `json:"like_count" example:"0"`
	RetweetCount int64   `json:"retweet_count" example:"0"`
}

type CreateTweetResponse struct {
	Tweet CreateTweetItem `json:"tweet"`
}

type CreateTweetItem struct {
	ID       int32   `json:"id" example:"1"`
	UserID   int32   `json:"user_id" example:"1"`
	Content  string  `json:"content" example:"Hello, world!"`
	ImageURL *string `json:"image_url,omitempty" example:"http://localhost:8080/uploads/1_abc.jpg"`
}

type GetTweetResponse struct {
	Tweet SwaggerTweet `json:"tweet"`
}

type GetUserTweetsResponse struct {
	Tweets     []SwaggerTweet `json:"tweets"`
	NextCursor *int32         `json:"next_cursor"`
	HasMore    bool           `json:"has_more" example:"false"`
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
