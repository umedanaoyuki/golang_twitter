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
	
	MsgUserRegistered         = "ユーザー登録が完了しました"
)