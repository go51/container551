package container551

import (
	"github.com/go51/log551"
	"net/http"
)

type Container struct {
	sid    string
	w      http.ResponseWriter
	r      *http.Request
	logger *log551.Log551
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
