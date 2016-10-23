package boois_passport

type PassportInfo struct {
	Key string // 标识符
	Nickname string // 昵称
	Account string // 账号
	Psw string // 密码
	Locked bool // 账号是否被锁定登陆
	Token string // 鉴权码
	LoginDate int64 //登录的时间戳
	OtherData map[string] string // 额外附加的资料
}
