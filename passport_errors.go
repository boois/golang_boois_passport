/*
当前文件用来自定义错误信息
*/
package boois_passport

const (
	ERR_ACC_EMPTY = iota
	ERR_ACC_LEN_FAIL
	ERR_ACC_FMT_FAIL
	ERR_PSW_EMPTY
	ERR_PSW_LEN_FAIL
	ERR_PSW_FMT_FAIL
	ERR_TIME_LOCKED
	ERR_PSW_FAIL
	ERR_USER_LOCKED
	ERR_USER_NONE
	ERR_CK_FAIL
	ERR_CK_SIGN_FAIL
	ERR_KICK_USER
	ERR_ACC_EXISTS
	ERR_NICKNAME_EXISTS

)

var errorMsgs map[int] string

func GetDefaultErrorMsgsMap() map[int] string{
	if errorMsgs == nil {
		errorMsgs = make(map[int] string)
		errorMsgs[ERR_ACC_EMPTY] = "账号不能为空"
		errorMsgs[ERR_ACC_LEN_FAIL] = "账号的长度只能为%d-%d"
		errorMsgs[ERR_ACC_FMT_FAIL] = "账号格式错误"
		errorMsgs[ERR_PSW_EMPTY] = "密码不能为空"
		errorMsgs[ERR_PSW_LEN_FAIL] = "密码的长度只能为%d-%d"
		errorMsgs[ERR_PSW_FMT_FAIL] = "密码格式错误"
		errorMsgs[ERR_TIME_LOCKED] = "超过了错误次数,请稍后在%d秒后再试"
		errorMsgs[ERR_PSW_FAIL] = "密码错误,还有%d次机会"
		errorMsgs[ERR_USER_LOCKED] = "您已被锁定登录,请联系管理员解锁"
		errorMsgs[ERR_USER_NONE] = "没有获取到用户资料"
		errorMsgs[ERR_CK_FAIL] = "cookies记录读取失败"
		errorMsgs[ERR_CK_SIGN_FAIL] = "cookies签名验证失败,可能cookies被篡改"
		errorMsgs[ERR_KICK_USER] = "服务器设置了用户互踢,同一个账号同一时间只能登陆一个用户"
		errorMsgs[ERR_ACC_EXISTS] = "账号已经存在"
		errorMsgs[ERR_NICKNAME_EXISTS] = "昵称已经存在"
	}
	return errorMsgs
}


