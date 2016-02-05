package container551

import (
	"github.com/go51/auth551"
	"github.com/go51/cookie551"
	"github.com/go51/log551"
	"github.com/go51/memcache551"
	"github.com/go51/model551"
	"github.com/go51/mysql551"
	"github.com/go51/secure551"
	"github.com/go51/string551"
	"net/http"
	"strconv"
)

type Container struct {
	sid     string
	w       http.ResponseWriter
	r       *http.Request
	logger  *log551.Log551
	cookie  *cookie551.Cookie
	db      *mysql551.Mysql
	session *memcache551.Memcache
	model   *model551.Model
	auth    *auth551.Auth
	user    *auth551.UserModel
}

func New() *Container {
	return &Container{}
}

func (c *Container) SetSID(sid string) {
	c.sid = sid
}

func (c *Container) SID() string {
	return c.sid
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

	// Load user
	c.session.GetModel("reminder_user", &c.user)

	return

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

func (c *Container) Logout() {
	c.cookie.Delete(c.auth.CookieKeyName())
	c.session.Delete("user")
}

func (c *Container) IsLogin() bool {
	return c.user != nil
}

func (c *Container) User() *auth551.UserModel {
	return c.user
}
