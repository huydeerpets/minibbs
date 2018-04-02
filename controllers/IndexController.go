package controllers

import (
	"fmt"
	"minibbs/filters"
	"minibbs/models"
	"net/http"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/sluu99/uuid"
)

type IndexController struct {
	beego.Controller
}

// Index .
func (c *IndexController) Index() {
	c.Data["PageTitle"] = "首页"
	c.Data["IsLogin"], c.Data["UserInfo"] = filters.IsLogin(c.Controller.Ctx)
	page, _ := strconv.Atoi(c.Ctx.Input.Query("p"))
	if page == 0 {
		page = 1
	}
	size, _ := beego.AppConfig.Int("page.size")
	tagId, _ := strconv.Atoi(c.Ctx.Input.Query("tagId"))
	c.Data["TagId"] = tagId
	tag := models.Tag{Id: tagId}
	c.Data["Page"] = models.TopicManager.PageTopic(page, size, &tag)
	c.Data["Tags"] = models.FindAllTag()
	c.Layout = "layout/layout.tpl"
	c.TplName = "index.tpl"
}

// LoginPage .
func (c *IndexController) LoginPage() {
	IsLogin, _ := filters.IsLogin(c.Ctx)
	if IsLogin {
		c.Redirect("/", 302)
	} else {
		beego.ReadFromRequest(&c.Controller)
		u := models.UserManager.FindPermissionByUser(1)
		beego.Debug(u) // ????????????????????????????????
		c.Data["PageTitle"] = "登录"
		c.Layout = "layout/layout.tpl"
		c.TplName = "login.tpl"
	}
}

// Login .
func (c *IndexController) Login() {
	flash := beego.NewFlash()
	username, password := c.Input().Get("username"), c.Input().Get("password")

	exsit, user, err := models.UserManager.Login(username, password)
	if err != nil {
		flash.Error(err.Error())
		flash.Store(&c.Controller)
		c.Redirect("/login", 302)
	}

	if exsit {
		c.SetSecureCookie(beego.AppConfig.String("cookie.secure"), beego.AppConfig.String("cookie.token"), user.Token, 30*24*60*60, "/", beego.AppConfig.String("cookie.domain"), false, true)
		c.Redirect("/", 302)
	}

	flash.Error("用户名或密码错误")
	flash.Store(&c.Controller)
	c.Redirect("/login", 302)
}

// RegisterPage .
func (c *IndexController) RegisterPage() {
	isLogin, _ := filters.IsLogin(c.Ctx)

	if isLogin {
		c.Redirect("/", http.StatusFound)
		return
	}

	beego.ReadFromRequest(&c.Controller)
	c.Data["PageTitle"] = "用户注册"
	c.Layout = "layout/layout.tpl"
	c.TplName = "register.tpl"
	return

}

// Register .
func (c *IndexController) Register() {
	flash := beego.NewFlash()
	username, password, email := c.Input().Get("username"), c.Input().Get("password"), c.Input().Get("email")
	if len(username) == 0 || len(password) == 0 || len(email) == 0 {
		flash.Error("输入框不能为空")
		flash.Store(&c.Controller)
		c.Redirect("/register", http.StatusFound)
		return
	}

	var token = uuid.Rand().Hex() // token 唯一

	user := models.User{
		Username: username,
		Password: password,
		Email:    email,
		Token:    token,
		Image:    "/static/imgs/default.png",
	}

	if exsit, _ := models.UserManager.FindUserByUserName(username); exsit {
		flash.Error("用户名已被注册")
		flash.Store(&c.Controller)
		c.Redirect("/register", http.StatusFound)
		return
	}

	if exsit, _ := models.UserManager.FindUserByUserEmail(email); exsit {
		flash.Error("邮箱已被注册")
		flash.Store(&c.Controller)
		c.Redirect("/register", http.StatusFound)
		return
	}

	if isEmailRegist, _ := beego.AppConfig.Bool("emailRegist"); isEmailRegist {
		authURL := models.EmailManager.GenerateAuthURL(email)
		fmt.Printf("\n%s\n", authURL)
		models.EmailManager.SetTheme("用户帐号激活") //设置主题
		models.EmailManager.SetEmailContent(authURL)

		err := models.EmailManager.InitSendCfg(email, username)
		if err != nil {
			fmt.Printf("send email init error[%s]", err.Error())
			flash.Error("发送注册邮件初始化时发生错误，请联系管理员")
			flash.Store(&c.Controller)
			c.Redirect("/register", http.StatusFound)
			return
		}

		err = models.EmailManager.SendEmail()
		if err != nil {
			fmt.Println(err.Error())
			flash.Error("发送注册邮件时发生错误，请联系管理员")
			flash.Store(&c.Controller)
			c.Redirect("/register", http.StatusFound)
			return
		}
		flash.Success("注册验证邮件已经发送到您的邮箱，请激活后再登录")
	} else {
		user.Active = true // if email regist is false
	}

	if err := models.UserManager.SaveUser(&user); err != nil {
		flash.Error("注册用户失败:" + err.Error())
		flash.Store(&c.Controller)
		c.Redirect("/register", http.StatusFound)
		return
	}

	flash.Success("注册成功")
	flash.Store(&c.Controller)
	c.Redirect("/register", http.StatusFound)
	return
}

// ActiveAccount activation user account by check email
func (c *IndexController) ActiveAccount() {
	flash := beego.NewFlash()
	token := c.GetString("token")
	fmt.Println("token: " + token)

	if models.EmailManager != nil {
		isAccess, email := models.EmailManager.CheckEmailURL(token)

		if isAccess {
			err := models.UserManager.ActiveAccount(email)
			if err != nil {
				// glog.Errorf("active user by email error[%s]\n", err.Error())
				flash.Error("激活账户时发生错误，请联系管理员 " + err.Error())
				flash.Store(&c.Controller)
				c.Redirect("/login", http.StatusFound)
				return
			}
		}
		//需要设置cookie???????
		// c.SetSecureCookie(beego.AppConfig.String("cookie.secure"), beego.AppConfig.String("cookie.token"), token, 30 * 24 * 60 * 60, "/", beego.AppConfig.String("cookie.domain"), false, true)
		flash.Success("激活账户成功")
		flash.Store(&c.Controller)
		c.Redirect("/login", http.StatusFound)
		return
	}

	flash.Error("发送注册邮件时发生错误，请联系管理员")
	flash.Store(&c.Controller)
	c.Redirect("/login", http.StatusFound)
	return
}

// Logout .
func (c *IndexController) Logout() {
	c.SetSecureCookie(beego.AppConfig.String("cookie.secure"), beego.AppConfig.String("cookie.token"), "", -1, "/", beego.AppConfig.String("cookie.domain"), false, true)
	c.Redirect("/", 302)
}

// About .
func (c *IndexController) About() {
	c.Data["IsLogin"], c.Data["UserInfo"] = filters.IsLogin(c.Controller.Ctx)
	c.Data["PageTitle"] = "公告"
	c.Layout = "layout/layout.tpl"
	c.TplName = "about.tpl"
}
