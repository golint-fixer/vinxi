package manager

import (
	"encoding/json"
	"errors"
	"net/http"

	"gopkg.in/vinxi/vinxi.v0/plugin"
	"gopkg.in/vinxi/vinxi.v0/rule"
)

// Context is used to share request context entities
// across controllers.
type Context struct {
	Manager      *Manager
	Scope        *Scope
	Instance     *Instance
	AdminPlugins *plugin.Layer
	Request      *http.Request
	Response     http.ResponseWriter
	Rule         rule.Rule
	Plugin       plugin.Plugin
}

// ParseBody parses the body.
func (c *Context) ParseBody(bind interface{}) error {
	if c.Request.Header.Get("Content-Type") != "application/json" {
		c.SendError(http.StatusUnsupportedMediaType, "Invalid content type. Must be application/json")
		return errors.New("Invalid type")
	}

	decoder := json.NewDecoder(c.Request.Body)
	return decoder.Decode(&bind)
}

// SendOk is used to serialize and write the response as JSON.
func (c *Context) SendOk(data interface{}) {
	c.Send(http.StatusOK, data)
}

// Send is used to serialize and write the response as JSON with custom status code.
func (c *Context) Send(status int, data interface{}) {
	buf, err := json.Marshal(data)
	if err != nil {
		c.SendError(http.StatusInternalServerError, err.Error())
		return
	}
	c.Response.Header().Set("Content-Type", "application/json")
	c.Response.WriteHeader(status)
	c.Response.Write(buf)
}

// SendError replies with an custom error message and 500 as status code.
func (c *Context) SendError(status int, message string) {
	c.Response.Header().Set("Content-Type", "application/json")

	buf, err := json.Marshal(struct {
		Code    int    `json:"code"`
		Message string `json:"message,omitempty"`
	}{status, message})

	if err != nil {
		c.Response.WriteHeader(500)
		c.Response.Write([]byte(err.Error()))
		return
	}

	c.Response.WriteHeader(status)
	c.Response.Write(buf)
}

// SendNoContent replies with 204 status code.
func (c *Context) SendNoContent() {
	c.Response.WriteHeader(http.StatusNoContent)
}

// SendNotFound replies with 404 status code and custom message.
func (c *Context) SendNotFound(message string) {
	c.SendError(http.StatusNotFound, message)
}
