package middleware

import (
	"go-zentao-task-api/core"
	"go-zentao-task-api/pkg/session"
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
