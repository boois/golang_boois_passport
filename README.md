#golang_boois_passport是什么?
一个定义了回调标准的登陆与注册组件，可以5分钟快速建立一个登陆账号管理系统
设计了回调接口,可以自由使用任何数据源来进行用户数据输出以及session管理



#有问题反馈
在使用中有任何问题，欢迎反馈给我，可以用以下联系方式跟我交流

* 周骁鸣
* 邮件(boois@qq.com)
* 微信 booisming

#使用方法:
##1.建立一个回调对象,以便于处理登陆和注册中所有需要的action


```golang
/*
这个文件配置了登陆所需的内容,依赖第三方组件
go get github.com/boois/golang_boois_passport
 */
package front

import (
	"github.com/boois/golang_boois_passport"
	"errors"
)

//↓下面这个数据源是用一个map来做范例的,实际使用中可以替换成自己的数据源
var userDataSource map[string] boois_passport.PassportInfo
func getUserDataSource() map[string] boois_passport.PassportInfo{
	if userDataSource == nil{
		userDataSource = make(map[string] boois_passport.PassportInfo)
		userDataSource["test"] = boois_passport.PassportInfo{
			Key:"1", // 标识符
			Nickname:"haha", // 昵称
			Account:"test", // 账号
			Psw:"123123",// 密码
			Locked:false, // 是否锁定
			OtherData:make(map[string] string),
		}
		userDataSource["test1"] = boois_passport.PassportInfo{
			Key:"1", // 标识符
			Nickname:"haha", // 昵称
			Account:"nono", // 账号
			Psw:"123123",// 密码
			Locked:true, // 是否锁定
			OtherData:make(map[string] string),
		}
	}
	return userDataSource
}
//↑上面这个数据源是用一个map来做范例的,实际使用中可以替换成自己的数据源
//↓下面这个map是用来存放用户的登录信息的,实际使用中可以替换为memcached之类的缓存
var sessionMap map[string] boois_passport.PassportInfo
func getSessionMap() map[string] boois_passport.PassportInfo{
	if sessionMap == nil{
		sessionMap = make(map[string] boois_passport.PassportInfo)
	}
	return sessionMap
}
//↑上面这个map是用来存放用户的登录信息的,实际使用中可以替换为memcached之类的缓存

// 下面是自定义错误列表,可以自己定义相关的错误文字
var errorMsgs map[int] string

func myMsgsMap() map[int] string{
	if errorMsgs == nil {
		errorMsgs = make(map[int] string)
		errorMsgs[boois_passport.ERR_ACC_EMPTY] = "账号不能为空"
		errorMsgs[boois_passport.ERR_ACC_LEN_FAIL] = "2账号的长度只能为%d-%d"
		errorMsgs[boois_passport.ERR_ACC_FMT_FAIL] = "3账号格式错误"
		errorMsgs[boois_passport.ERR_PSW_EMPTY] = "密码不能为空"
		errorMsgs[boois_passport.ERR_PSW_LEN_FAIL] = "密码的长度只能为%d-%d"
		errorMsgs[boois_passport.ERR_PSW_FMT_FAIL] = "密码格式错误"
		errorMsgs[boois_passport.ERR_TIME_LOCKED] = "超过了错误次数,请稍后在%d秒后再试"
		errorMsgs[boois_passport.ERR_PSW_FAIL] = "密码错误,还有%d次机会"
		errorMsgs[boois_passport.ERR_USER_LOCKED] = "您已被锁定登录,请联系管理员解锁"
		errorMsgs[boois_passport.ERR_USER_NONE] = "没有获取到用户资料"
		errorMsgs[boois_passport.ERR_CK_FAIL] = "cookies记录读取失败"
		errorMsgs[boois_passport.ERR_CK_SIGN_FAIL] = "cookies签名验证失败,可能cookies被篡改"
		errorMsgs[boois_passport.ERR_KICK_USER] = "服务器设置了用户互踢,同一个账号同一时间只能登陆一个用户"
		errorMsgs[boois_passport.ERR_ACC_EXISTS] = "账号已经存在"
		errorMsgs[boois_passport.ERR_NICKNAME_EXISTS] = "昵称已经存在"
	}
	return errorMsgs
}


var login *boois_passport.BooisPassport = nil
func BooisPassport() *boois_passport.BooisPassport {
	if login == nil {
		login = &boois_passport.BooisPassport{ // 登陆配置
			BooisPassportAdapter:&MyPassportAdapter{},
			ErrorMsgs:myMsgsMap(),
			AccountRegex:"^[a-zA-Z][a-zA-Z0-9_]*$",
			AccountMinLen:4,
			AccountMaxLen:24,
			PswRegex:"^[a-zA-Z0-9_]*$",
			PswMinLen:4,
			PswMaxLen:24,
			ErrSecondSpan:3,
			ErrCount:5,
			AllowNicknameRepeat:false,
			CookiesMaxAge:604800,
			LoginUserUnique:true,
		}
	}
	return login
}

type MyPassportAdapter struct {}

// 用来连接数据库来获取用户资料,用来返回给login组件来进行登陆判断
func (this MyPassportAdapter) GetPassportInfoByAccount(account string) (boois_passport.PassportInfo,error){
	if v,ok := getUserDataSource()[account];ok{
		return v,nil
	}
	return boois_passport.PassportInfo{},errors.New("账号错误!")
}
// 登录前的预处理动作,用来给accout和psw做一些处理,如:字母转小写,转义等
func (this MyPassportAdapter) BeforeLogin(account string,psw string)  error{
	return nil
}
// 登录前的预处理动作,用来给accout和psw做一些处理,如:字母转小写,转义等
func (this MyPassportAdapter) LoginOk(user boois_passport.PassportInfo)  error{
	return nil
}
// 登录失败时的处理动作,可以用来记录日志锁定用户等,errCount是错误的次数,errTimeSpan是锁定解除的剩余时间
func (this MyPassportAdapter) LoginFail(account string, err error, errCount int, errTimeSpan int64){
	println("登陆失败")
}
// 登出之前的动作,如果返回的err不为nil可以阻止登出
func (this MyPassportAdapter) BeforeLogout(key string) error{
	return nil
}
// 登出之后的动作,比如清理关联缓存或者记录下登出时间
func (this MyPassportAdapter) AfterLogout(account string){

}
// 设置session
func (this MyPassportAdapter) SetSession(account string, user boois_passport.PassportInfo){
	getSessionMap()[account] = user
}
// 获取session
func (this MyPassportAdapter) GetSession(account string) (boois_passport.PassportInfo,error){
	if v,ok := getSessionMap()[account];ok{
		return v,nil
	}
	return boois_passport.PassportInfo{},errors.New("没有找到用户")
}
// 删除session
func (this MyPassportAdapter) DelSession(account string){
	delete(getSessionMap(),account)
}
// 密码加密方法
func (this MyPassportAdapter) EncryptPsw(psw string) string{
	return psw
}
// 向数据源中添加一个用户账号
func (this MyPassportAdapter) AddAccount(user boois_passport.PassportInfo) error  {
	return nil
}
// 检查一个账号是否已经存在
func (this MyPassportAdapter) ChkAccountExists(account string) bool  {
	if _,ok := getUserDataSource()[account];ok{
		return true
	}
	return false
}
// 检查一个昵称是否已经存在
func (this MyPassportAdapter) ChkNickNameExists(nickname string) bool  {
	return false
}
// 注册前的预处理操作
func (this MyPassportAdapter) BeforeReg(user boois_passport.PassportInfo) error{
	return nil
}
// 开始在后端进行注册操作
func (this MyPassportAdapter) Reg(user boois_passport.PassportInfo) error{
	// 这里向数据源写入数据,密码如果需要加密请记得使用和当前一致的密码加密方式 this.EncryptPsw
	getUserDataSource()[user.Account] = user
	println("成功添加一个用户")
	return nil
}
// 注册失败的处理动作,可以记录失败日志等情况
func (this MyPassportAdapter) RegFail(user boois_passport.PassportInfo, err error){

}
```
## 2. 在登陆页面中调用它
```html
<html>
<form method="post">
    <br/>
    {{.Err}}<br/><br/>
    账号 <input type="text" name="account"/><br/>
    密码 <input type="text" name="psw"/><br/>
    <button type="submit">登陆</button> <br/>
    <a href="/reg">注册</a><br/>
</form>
</html>
```
```golang
package front

import (
    "net/http"
	"github.com/boois/golang_get_tmp_from_oss_or_local"
	"html/template"
	"cfg"
)

func Index(w http.ResponseWriter, r *http.Request,url_args ...string) { // url_args 是路由解析后传入的url中的相关值
	r.ParseForm()       // 解析参数,默认并不会解析
	var html = boois_temp_utils.GetTemp("index.html",cfg.TMP_CACHE,cfg.OSS_MODE,cfg.OSS_URL)
	t,_ := template.New("page").Parse(html)

	err_str := "请登录"

	if r.Method == "POST"{
		post_acc := r.Form.Get("account")
		post_psw := r.Form.Get("psw")
		println("post_acc:",post_acc,"  post_psw:",post_psw)
		logininfo,err := BooisPassport().Login(w,post_acc,post_psw) // 登陆
		if err != nil{
			err_str = err.Error()
		}else{
			println("成功登录:",logininfo.Account)
			// 跳转
			http.Redirect(w,r,"/uc",http.StatusFound)
		}
	}
	data := struct {
		Err string
		Title string // 注意大写开头才能被调用
		Items []string
	}{
		Err:err_str,
		Title:"my page",
		Items:	[]string {"1","2"},
	}
	t.Execute(w,data)
}
```
## 3. 在登出页面中调用它
```golang
package front

import (
    "net/http"
)

func Logout(w http.ResponseWriter, r *http.Request,url_args ...string) { // url_args 是路由解析后传入的url中的相关值
	r.ParseForm()
	BooisPassport().Logout(r.Form.Get("account"))
	http.Redirect(w,r,"/",http.StatusFound)
}
```
## 4. 在注册页面中调用它
```html
<html>
<form method="post">
{{.Err}} <br/>
昵称 <input type="text" name="nickname"/>
账号 <input type="text" name="account"/>
密码 <input type="text" name="psw"/>
<button type="submit">注册</button>
</form>
</html>
```
```golang
package front

import (
    "net/http"
	"github.com/boois/golang_get_tmp_from_oss_or_local"
	"html/template"
	"cfg"
	"github.com/boois/golang_boois_passport"
)

func Reg(w http.ResponseWriter, r *http.Request,url_args ...string) { // url_args 是路由解析后传入的url中的相关值
	r.ParseForm()       // 解析参数,默认并不会解析
	var html = boois_temp_utils.GetTemp("reg.html",cfg.TMP_CACHE,cfg.OSS_MODE,cfg.OSS_URL)
	t,_ := template.New("page").Parse(html)

	post_nickname := r.Form.Get("nickname")
	post_acc := r.Form.Get("account")
	post_psw := r.Form.Get("psw")
	println("post_acc:",post_acc,"  post_psw:",post_psw,"  nickname:",post_nickname)
	post_user := boois_passport.PassportInfo{
			Key:post_acc, // 标识符
			Nickname:post_nickname, // 昵称
			Account:post_acc, // 账号
			Psw:post_psw,// 密码
			Locked:false, // 是否锁定
			OtherData:make(map[string] string),
	}
	user,err := BooisPassport().Register(post_user) // 登陆

	err_str := "注册成功"
	if err != nil{
		err_str = err.Error()
	}else{
		println(user.Account)
		// 跳转
		http.Redirect(w,r,"/",http.StatusFound)
	}
	data := struct {
		Err string
		Title string // 注意大写开头才能被调用
		Items []string
	}{
		Err:err_str,
		Title:"my page",
		Items:	[]string {"1","2"},
	}
	t.Execute(w,data)
}
```
## 5. 在登陆后的页面中调用它
```html
<html>
<form method="post">
<br/>
欢迎您 {{.User.Nickname}}({{.User.Account}})  <a href="/logout?account={{.User.Account}}">退出</a>
</form>
</html>
```
```golang
package front

import (
    "net/http"
	"github.com/boois/golang_get_tmp_from_oss_or_local"
	"html/template"
	"cfg"
	"github.com/boois/golang_boois_passport"
)

func Uc(w http.ResponseWriter, r *http.Request,url_args ...string) { // url_args 是路由解析后传入的url中的相关值
	user,err := BooisPassport().GetCookiesUser(r)
	if err != nil {
		println(err.Error())
		http.Redirect(w,r,"/",http.StatusFound)
		return
	}
	var html = boois_temp_utils.GetTemp("uc.html",cfg.TMP_CACHE,cfg.OSS_MODE,cfg.OSS_URL)
	t,_ := template.New("page").Parse(html)

	data := struct {
		User boois_passport.PassportInfo // 注意大写开头才能被调用
	}{
		User:user,
	}
	t.Execute(w,data)
}
```

