package session

import (
	"github.com/boj/redistore"
	"github.com/gorilla/sessions"
	"go-zentao-task/pkg/config"
	"go-zentao-task/pkg/gredis"
	"net/http"
)

var Store sessions.Store

const (
	MaxAge   = 1200
	HttpOnly = true
)

func Setup() {
	var secure bool
	if config.Get("app.scheme") == "https" {
		secure = true
	}

	switch config.Get("app.session.store") {
	case "cookie":
		store := sessions.NewCookieStore([]byte(config.Get("app.session.key")))
		store.MaxAge(MaxAge)
		store.Options.HttpOnly = HttpOnly
		store.Options.Secure = secure
		store.Options.SameSite = http.SameSiteLaxMode
		Store = store
	case "file":
		store := sessions.NewFilesystemStore("runtime/sessions", []byte(config.Get("app.session.key")))
		store.MaxAge(MaxAge)
		store.Options.HttpOnly = HttpOnly
		store.Options.Secure = secure
		store.Options.SameSite = http.SameSiteLaxMode
		Store = store
	default:
		store, _ := redistore.NewRediStoreWithPool(gredis.RedisPool, []byte(config.Get("app.session.key")))
		store.SetMaxAge(MaxAge)
		store.Options.HttpOnly = HttpOnly
		store.Options.Secure = secure
		store.Options.SameSite = http.SameSiteLaxMode
		Store = store
	}
}

type Session struct {
	Name    string
	Session *sessions.Session
	R       *http.Request
	W       http.ResponseWriter
	Store   sessions.Store
}

func (s *Session) session() *sessions.Session {
	if s.Session == nil {
		s.Session, _ = s.Store.Get(s.R, s.Name)
	}
	return s.Session
}

func (s *Session) GetSession(key interface{}) interface{} {
	return s.session().Values[key]
}

func (s *Session) SetSession(key interface{}, val interface{}) {
	s.session().Values[key] = val
}

func (s *Session) DeleteSession(key interface{}) {
	delete(s.session().Values, key)
}

func (s *Session) FlushSession() {
	s.session().Options.MaxAge = -1
}

func (s *Session) AddSessionFlash(value interface{}, vars ...string) {
	s.session().AddFlash(value, vars...)
}

func (s *Session) SessionFlashes(vars ...string) []interface{} {
	return s.session().Flashes(vars...)
}

func (s *Session) SaveSession() error {
	return s.session().Save(s.R, s.W)
}
