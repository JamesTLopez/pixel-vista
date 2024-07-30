package session

import (
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
)

var SessionManager *scs.SessionManager

func InitSession() {
	SessionManager = scs.New()

	SessionManager.Lifetime = 1 * time.Hour
	// SessionManager.IdleTimeout = 20 * time.Minute
	SessionManager.Cookie.Name = "pixel_vista_id"
	// SessionManager.Cookie.Domain = "e"
	SessionManager.Cookie.HttpOnly = true
	SessionManager.Cookie.Path = "/"
	SessionManager.Cookie.Persist = true
	SessionManager.Cookie.SameSite = http.SameSiteStrictMode
	SessionManager.Cookie.Secure = true
}
