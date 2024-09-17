package session

import (
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
)

var SessionManager *scs.SessionManager

func InitSession() {
	SessionManager = scs.New()

	SessionManager.Lifetime = 24 * time.Hour
	SessionManager.IdleTimeout = 20 * time.Minute
	SessionManager.Cookie.Name = "pixel_vista"
	// SessionManager.Cookie.HttpOnly = true
	// SessionManager.Cookie.Path = "/"
	// SessionManager.Cookie.Persist = true
	SessionManager.Cookie.SameSite = http.SameSiteLaxMode
	// SessionManager.Cookie.Secure = true
}
