package container551

import (
	"github.com/go51/auth551"
	"github.com/go51/cookie551"
	"github.com/go51/log551"
	"github.com/go51/memcache551"
	"github.com/go51/model551"
	"github.com/go51/mysql551"
	"github.com/go51/repository551"
	"github.com/go51/secure551"
	"github.com/go51/string551"
	xoauth2 "golang.org/x/oauth2"
	"net/http"
	"strconv"
)

type Container struct {
	sid     string
	ssid    string
	w       http.ResponseWriter
	r       *http.Request
	logger  *log551.Log551
	cookie  *cookie551.Cookie
	db      *mysql551.Mysql
	session *memcache551.Memcache
	model   *model551.Model
	auth    *auth551.Auth
	user    *auth551.UserModel
	options map[string]string
}

func New() *Container {
	return &Container{}
}

func (c *Container) SetSID(sid string) {
	c.sid = sid
	c.ssid = sid[:10]
}

func (c *Container) SID() string {
	return c.sid
}
func (c *Container) SSID() string {
	return c.ssid
}

func (c *Container) SetResponseWriter(w http.ResponseWriter) {
	c.w = w
}

func (c *Container) ResponseWriter() http.ResponseWriter {
	return c.w
}

func (c *Container) SetRequest(r *http.Request) {
	c.r = r
}

func (c *Container) Request() *http.Request {
	return c.r
}

func (c *Container) SetLogger(logger *log551.Log551) {
	c.logger = logger
}

func (c *Container) Logger() *log551.Log551 {
	return c.logger
}

func (c *Container) SetCookie(cookie *cookie551.Cookie) {
	c.cookie = cookie
}

func (c *Container) Cookie() *cookie551.Cookie {
	return c.cookie
}

func (c *Container) SetDb(db *mysql551.Mysql) {
	c.db = db
}

func (c *Container) Db() *mysql551.Mysql {
	return c.db
}

func (c *Container) SetSession(session *memcache551.Memcache) {
	c.session = session
}

func (c *Container) Session() *memcache551.Memcache {
	return c.session
}

func (c *Container) SetModel(modelManager *model551.Model) {
	c.model = modelManager
}

func (c *Container) ModelManager() *model551.Model {
	return c.model
}

func (c *Container) SetAuth(auth *auth551.Auth) {
	c.auth = auth

	if c.user != nil {
		return
	}

	// Load user from session
	c.session.GetModel("reminder_user", &c.user)

	// Get user id from cookie
	id, err := c.getRemindId()
	if err != nil {
		return
	}

	// Get user model from database
	c.user = c.getUser(id)

	// Set user model to session
	c.session.Set("reminder_user", c.user)

	return

}

func (c *Container) getRemindId() (int64, error) {
	cookieId, err := c.cookie.Get(c.auth.CookieKeyName())
	if err != nil {
		return 0, err
	}

	sid := secure551.Decrypted(cookieId, c.auth.MasterKey())
	id, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		return 0, err
	}

	return id, nil

}

func (c *Container) getUser(id int64) *auth551.UserModel {
	repo := repository551.Load()
	miUser := c.ModelManager().Get("UserModel")
	mUser := repo.Find(c.db, miUser, id)
	user, ok := mUser.(*auth551.UserModel)
	if !ok {
		return nil
	}

	return user
}

func (c *Container) Auth() *auth551.Auth {
	return c.auth
}

func (c *Container) SignIn(user *auth551.UserModel) {
	// Set remind id to cookie
	id := string551.Right("0000000000000000"+strconv.FormatInt(user.Id, 10), 16)
	secureId := secure551.Encrypted(id, c.auth.MasterKey())
	c.cookie.Set(c.auth.CookieKeyName(), secureId, 60*60*24*365)

	// Set user model to session
	c.session.Set("reminder_user", user)

}

func (c *Container) SignOut() {
	c.cookie.Delete(c.auth.CookieKeyName())
	c.session.Delete("reminder_user")
}

func (c *Container) IsSignIn() bool {
	return c.user != nil
}

func (c *Container) User() *auth551.UserModel {
	return c.user
}

func (c *Container) ApiClient(vendor auth551.AuthVendor) *http.Client {
	if c.user == nil {
		return nil
	}
	repo := repository551.Load()
	miUserToken := c.model.Get("UserTokenModel")
	param := map[string]interface{}{
		"user_id": c.user.Id,
	}
	mUserToken := repo.FindOneBy(c.db, miUserToken, param)
	userToken, ok := mUserToken.(*auth551.UserTokenModel)
	if !ok {
		return nil
	}

	token := &xoauth2.Token{
		AccessToken:  userToken.AccessToken,
		TokenType:    userToken.TokenType,
		RefreshToken: userToken.RefreshToken,
		Expiry:       userToken.Expiry,
	}

	return c.auth.Client(vendor, token)

}

func (c *Container) SetCommandOptions(options map[string]string) {
	c.options = options
}

func (c *Container) CommandOption(name string) string {
	return c.options[name]
}
