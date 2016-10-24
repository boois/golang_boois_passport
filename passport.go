package boois_passport

import (
	"errors"
	"time"
	"fmt"
	"regexp"
	"net/http"
	"strings"
	"crypto/md5"
	"encoding/hex"
)


type BooisPassport struct {
	BooisPassportAdapter BooisPassportAdapter // 适配器
	ErrorMsgs 			 map[int] string 	  // 适配器
	AccountRegex         string               // 账号规则
	AccountMinLen        int                  // 账号最小长度
	AccountMaxLen        int                  // 账号最大长度
	PswRegex             string               // 密码规则
	PswMinLen            int                  // 密码最小长度
	PswMaxLen            int                  // 密码最大长度
	ErrSecondSpan        int64                // 出错后能够再次尝试的秒数
	ErrCount             int                  // 允许出错的总次数
	ErrCountMap			 map[string] int 	  // 缓存错误次数
	ErrCountDateMap		 map[string] int64    // 缓存错误次数超出时间
	AllowNicknameRepeat  bool 				  // 缓存错误次数超出时间
	CookiesMaxAge		 int				  // cookies缓存的时间,如果不使用cookies可设为-1
	LoginUserUnique		 bool				  // 是否保持唯一登录,如果开启此选项,则同一个账号多点登录会互踢

}

func (this *BooisPassport) getErrMsgsMap() map[int] string {
	if this.ErrorMsgs == nil {
		this.ErrorMsgs = make(map[int] string)
	}
	return this.ErrorMsgs
}
func (this *BooisPassport) getErrMsg(key int) string{
	if v,ok := this.getErrMsgsMap()[key];ok{
		return v
	}
	if v,ok := GetDefaultErrorMsgsMap()[key];ok {
		return v
	}
	return "未定义的错误"
}

// 获取错误次数Map
func (this *BooisPassport) getErrCountMap() map[string] int{
	if this.ErrCountMap == nil{
		this.ErrCountMap = make(map[string] int)
	}
	return this.ErrCountMap
}
// 获取错误次数DataMap
func (this *BooisPassport) getErrCountDateMap() map[string] int64{
	if this.ErrCountDateMap == nil{
		this.ErrCountDateMap = make(map[string] int64)
	}
	return this.ErrCountDateMap
}

// 获取错误次数
func (this *BooisPassport) GetErrCount(account string) int{
	if v,ok := this.getErrCountMap()[account];ok{
		return v
	}
	return 0
}

// 增加错误次数
func (this *BooisPassport) PlusErrCount(account string) int{
	currentCount := this.GetErrCount(account)
	this.getErrCountMap()[account] =  currentCount + 1
	return currentCount
}

// 清理错误次数
func (this *BooisPassport) clearErrCount(account string)  {
	this.getErrCountMap()[account] = 0
	delete(this.getErrCountDateMap(),account)
}

// 处理登陆失败回调
func (this *BooisPassport) loginFail(account string, err error) {
	errCount := 0
	if v,ok := this.getErrCountMap()[account];ok {
		errCount = v
	}

	this.BooisPassportAdapter.LoginFail(account, err,errCount,this.chkErrDate(account))
}
// 检查账号格式
func (this *BooisPassport) chkAccountFormat(account string) error{
	// 先检查长度
	acc_arr := []rune(account)
	if len(acc_arr) == 0 {
		return errors.New(this.getErrMsg(ERR_ACC_EMPTY))//"账号不能为空"
	}
	if len(acc_arr) < this.AccountMinLen || len(acc_arr) > this.AccountMaxLen {
		return errors.New(fmt.Sprintf(this.getErrMsg(ERR_ACC_LEN_FAIL),this.AccountMinLen,this.AccountMaxLen))//"账号的长度只能为%d-%d"
	}
	// 再检查格式
	if this.AccountRegex == "" {this.AccountRegex = "^[\\s\\S]*$"}
	if reg,err := regexp.Compile(this.AccountRegex); err == nil {
		if reg.FindAllString(account,-1) == nil{
			return errors.New(this.getErrMsg(ERR_ACC_FMT_FAIL))//"账号格式错误"
		}
	}
	return nil
}
// 检查密码格式
func (this *BooisPassport) chkPswFormat(psw string) error{
	// 先检查长度
	psw_arr := []rune(psw)
	if len(psw_arr) == 0 {
		return errors.New(this.getErrMsg(ERR_PSW_EMPTY))//"密码不能为空"
	}
	if len(psw_arr) < this.PswMinLen || len(psw_arr) > this.PswMaxLen {
		return errors.New(fmt.Sprintf(this.getErrMsg(ERR_PSW_LEN_FAIL),this.PswMinLen,this.PswMaxLen))//"密码的长度只能为%d-%d"
	}
	// 再检查格式
	if this.PswRegex == "" {this.PswRegex = "^[\\s\\S]*$"}
	if reg,err := regexp.Compile(this.PswRegex); err == nil {
		if reg.FindAllString(psw,-1) == nil{
			return errors.New(this.getErrMsg(ERR_PSW_FMT_FAIL))//"密码格式错误"
		}
	}
	return nil
}
// cookies签名
func (this *BooisPassport) sign(v string) string{
	newstr := "__boois_passport__"+this.BooisPassportAdapter.EncryptPsw(v)+"__boois_passport__"
	m := md5.New()
	m.Write([]byte(newstr))
	cipherStr := m.Sum(nil)
	return hex.EncodeToString(cipherStr)
}
// 登陆
func (this *BooisPassport) Login(w http.ResponseWriter,account string, psw string) (PassportInfo,error) {
	// 1.预处理登陆前的事件
	if err := this.BooisPassportAdapter.BeforeLogin(&account,&psw);err != nil {
		this.loginFail(account,err) // 失败回调
		return PassportInfo{},err
	}
	//  2.先检查是否超过了错误次数,如果超过了,则直接返回错误
	timespan := this.chkErrDate(account);
	if timespan > 0 {
		err := errors.New(fmt.Sprintf(this.getErrMsg(ERR_TIME_LOCKED),timespan))//"超过了错误次数,请稍后在%d秒后再试"
		this.loginFail(account,err) // 失败回调
		return PassportInfo{},err
	}
	if _,ok := this.getErrCountDateMap()[account];ok{ // 如果已经过了错误时间间隔,则清理错误记录信息
		this.clearErrCount(account)
	}
	// 3. 先检查account,psw的输入格式
	if err := this.chkAccountFormat(account);err != nil {
		this.loginFail(account,err) // 失败回调
		return PassportInfo{},err
	}
	if err := this.chkPswFormat(psw);err != nil {
		this.loginFail(account,err) // 失败回调
		return PassportInfo{},err
	}
	// 4. 通过account 来获取用户信息,并
	login_info,err := this.BooisPassportAdapter.GetPassportInfoByAccount(account)
	if err != nil { // 如果无法获取用户信息,直接返回错误
		this.loginFail(account,err) // 失败回调
		return login_info,err
	}
	// 5. 比对密码
	if login_info.Psw != this.BooisPassportAdapter.EncryptPsw(psw) {
		//如果密码不一致,先累计一次登录错误记录
		count := this.plusErrCount(account)
		err := errors.New(fmt.Sprintf(this.getErrMsg(ERR_PSW_FAIL),this.ErrCount-count))//"密码错误,还有%d次机会"
		this.loginFail(account,err) // 失败回调
		return login_info,err
	}
	// 6. 用户是否被锁定
	if login_info.Locked {
		return login_info,errors.New(this.getErrMsg(ERR_USER_LOCKED))//"您已被锁定登录,请联系管理员解锁"
	}
	// 7. 成功,清理错误信息记录
	this.clearErrCount(account)
	// 8. 写入cookies
	login_info.LoginDate = time.Now().Unix()
	if this.CookiesMaxAge >0 {
		ck := &http.Cookie{
			Name:   "__boois_passport_user_account__",
			Value:   fmt.Sprintf("%s|%d|%s",login_info.Account,login_info.LoginDate,this.sign(login_info.Account)), // 为account做一个签名,防止客户端随机猜测
			Path:     "/",
			HttpOnly: false,
			MaxAge: this.CookiesMaxAge,
		}
		http.SetCookie(w,ck) //写入cookies
	}
	// 9. 将当前的信息序列化到session中去
	this.BooisPassportAdapter.SetSession(login_info.Account,login_info)
	// 10. 处理成功后的回调事件
	if err := this.BooisPassportAdapter.LoginOk(login_info);err!=nil {
		this.loginFail(account,err) // 失败回调
		return login_info,err
	}
	return login_info,err
}

// 增加一次错误次数,达到最大错误次数时,会记录一个时间戳
func (this *BooisPassport) plusErrCount(account string) int{
	currentCount := this.PlusErrCount(account)
	if currentCount >= this.ErrCount-1 {
		this.getErrCountDateMap()[account] = time.Now().Unix() //记录下当前时间戳
	}
	return currentCount
}

// 检查是否在错误限定时间内,int返回值是剩余时间
func (this *BooisPassport) chkErrDate(account string) int64{
	if v,ok := this.getErrCountDateMap()[account]; ok {
		return this.ErrSecondSpan - (time.Now().Unix() - v)
	}
	return 0
}

// 登出
func (this *BooisPassport) Logout(account string) error{
	if err:=this.BooisPassportAdapter.BeforeLogout(account);err!=nil{
		return err
	}
	this.BooisPassportAdapter.DelSession(account)

	this.BooisPassportAdapter.AfterLogout(account)
	return nil
}

// 从session中获取用户
func (this *BooisPassport) GetSessionUser(account string) (PassportInfo,error) {
	return this.BooisPassportAdapter.GetSession(account)
}
// 从cookies中获取session用户
func (this *BooisPassport) GetCookiesUser(r *http.Request) (PassportInfo,error) {
	ck,err := r.Cookie("__boois_passport_user_account__");//从cookies中获取用户资料
	if ck==nil || err != nil{
		return PassportInfo{},errors.New(this.getErrMsg(ERR_USER_NONE))//"没有获取到用户资料"
	}
	ckstr := ck.Value
	args := strings.Split(ckstr,"|")
	if len(args) != 3 {
		return PassportInfo{},errors.New(this.getErrMsg(ERR_CK_FAIL))//"cookies记录读取失败"
	}
	user, err := this.GetSessionUser(args[0])
	if err != nil {
		return user,err
	}
	//验证cookies签名
	sign := this.sign(user.Account)
	if args[2] != sign {
		return user,errors.New(this.getErrMsg(ERR_CK_SIGN_FAIL))//"cookies签名验证失败,可能cookies被篡改"
	}
	// 如果需要互踢的话,就要验证account的时间戳
	if this.LoginUserUnique {
		if args[1] != fmt.Sprint(user.LoginDate){
			return user,errors.New(this.getErrMsg(ERR_KICK_USER))//"服务器设置了用户互踢,同一个账号同一时间只能登陆一个用户"
		}
	}
	return user,nil
}

// 注册一个用户
func (this *BooisPassport) Register(user PassportInfo) (PassportInfo,error) {
	// 1.预处理登陆前的事件
	if err := this.BooisPassportAdapter.BeforeReg(&user);err != nil {
		this.BooisPassportAdapter.RegFail(user,err) // 失败回调
		return user,err
	}
	// 2. 检查account,psw的输入格式
	if err := this.chkAccountFormat(user.Account);err != nil {
		this.loginFail(user.Account,err) // 失败回调
		return user,err
	}
	if err := this.chkPswFormat(user.Psw);err != nil {
		this.loginFail(user.Psw,err) // 失败回调
		return user,err
	}
	// 3. 检查account和nickname是否已经存在
	if this.BooisPassportAdapter.ChkAccountExists(user.Account) {
		return user,errors.New(this.getErrMsg(ERR_ACC_EXISTS))//"账号已经存在"
	}
	if !this.AllowNicknameRepeat && this.BooisPassportAdapter.ChkNickNameExists(user.Nickname){
		return user,errors.New(this.getErrMsg(ERR_NICKNAME_EXISTS))//"昵称已经存在"
	}
	// 4. 检查通过后,将注册对象交给Reg
	if err:= this.BooisPassportAdapter.Reg(user);err != nil {
		return user,err
	}
	return user,nil
}





