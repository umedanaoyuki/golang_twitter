package messages

// エラーメッセージ
const (
	// バリデーションエラー
	ErrPasswordLength         = "パスワードは8文字以上15文字以下である必要があります"
	ErrPasswordInvalidChars   = "パスワードは半角英数字と記号(!?-_)のみ使用できます"
	ErrPasswordRequireLower   = "パスワードには小文字を含める必要があります"
	ErrPasswordRequireUpper   = "パスワードには大文字を含める必要があります"
	ErrPasswordRequireNumber  = "パスワードには数字を含める必要があります"
	ErrPasswordRequireSymbol  = "パスワードには記号(!?-_)を1文字以上含める必要があります"
	
	// サービスエラー
	ErrPasswordHashFailed     = "パスワードのハッシュ化に失敗しました"
	ErrEmailAlreadyExists     = "このメールアドレスは既に登録されています"
	
	// 成功メッセージ
	MsgUserRegistered         = "ユーザー登録が完了しました"
	MsgEmailSent              = "メールを送信しました"
)

// HTTPエラーメッセージ
const (
	ErrInvalidRequest         = "入力内容に誤りがあります"
	ErrRegistrationFailed     = "登録に失敗しました"
)
