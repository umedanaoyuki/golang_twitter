package messages

const (
	ErrPasswordLength         = "パスワードは8文字以上15文字以下である必要があります"
	ErrPasswordInvalidChars   = "パスワードは半角英数字と記号(!?-_)のみ使用できます"
	ErrPasswordRequireLower   = "パスワードには小文字を含める必要があります"
	ErrPasswordRequireUpper   = "パスワードには大文字を含める必要があります"
	ErrPasswordRequireNumber  = "パスワードには数字を含める必要があります"
	ErrPasswordRequireSymbol  = "パスワードには記号(!?-_)を1文字以上含める必要があります"
	
	ErrPasswordHashFailed     = "パスワードのハッシュ化に失敗しました"
	ErrEmailAlreadyExists     = "このメールアドレスは既に登録されています"
	
	MsgUserRegistered         = "メールを送信しました。\nアカウントを有効化してください。"
	MsgUserActivated          = "アカウントが有効化されました"
	MsgLoginSuccess           = "ログインに成功しました"
	MsgLogoutSuccess          = "ログアウトしました"

	ErrSessionSaveFailed      = "セッションの保存に失敗しました"

	ErrInvalidEmailOrPassword = "メールアドレスまたはパスワードが正しくありません"
	ErrAccountNotActivated    = "アカウントが有効化されていません"

	ErrTransactionFailed      = "トランザクション開始に失敗しました"
	ErrInvalidTokenOrExpired  = "無効なトークンまたは期限切れです"
	ErrUserActivationFailed   = "ユーザーのアクティブ化に失敗しました"
	ErrTokenDeletionFailed    = "トークンの削除に失敗しました"
)