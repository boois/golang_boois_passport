package boois_passport

type BooisPassportAdapter interface {
	GetPassportInfoByAccount(account string) (PassportInfo,error) // 用来连接数据库来获取用户资料,用来返回给login组件来进行登陆判断
	BeforeLogin(account string,psw string)  error // 登录前的预处理动作,用来给accout和psw做一些处理,如:字母转小写,转义等
	LoginOk(user PassportInfo)  error// 登陆成功
	LoginFail(account string, err error, errCount int, errTimeSpan int64) // 登录失败时的处理动作,可以用来记录日志锁定用户等,errCount是错误的次数,errTimeSpan是锁定解除的剩余时间
	BeforeLogout(account string) error // 登出之前的动作,如果返回的err不为nil可以阻止登出
	AfterLogout(account string) // 登出之后的动作,比如清理关联缓存或者记录下登出时间
	SetSession(account string, user PassportInfo) // 设置session
	GetSession(account string) (PassportInfo,error) // 获取session
	DelSession(account string) // 删除session
	EncryptPsw(psw string) string // 将原始密码加密的方法,请全局保持一致
	ChkAccountExists(account string) bool // 检查一个账号是否已经存在
	ChkNickNameExists(nickname string) bool // 检查一个昵称是否已经存在
	BeforeReg(user PassportInfo) error // 注册前的预处理操作
	Reg(user PassportInfo) error // 开始在后端进行注册操作
	RegFail(user PassportInfo, err error) // 注册失败的处理动作,可以记录失败日志等情况
}


