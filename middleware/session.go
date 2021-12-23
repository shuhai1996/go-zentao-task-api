package middleware

import (
	"go-zentao-task/core"
	"go-zentao-task/pkg/session"
)

func Session(c *core.Context) {
	c.Session = &session.Session{
		Name:    "session",
		Session: nil,
		R:       c.Request,
		W:       c.Writer,
		Store:   session.Store,
	}
	c.Next()
}
