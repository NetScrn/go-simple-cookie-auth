package web

import (
	"database/sql"
	"github.com/netscrn/gocookieauth/internal/data/sessions"
	"github.com/netscrn/gocookieauth/internal/data/users"
	"github.com/netscrn/gocookieauth/internal/web/controllers"
	"github.com/netscrn/gocookieauth/internal/web/middleware"
	auth "github.com/netscrn/gocookieauth/internal/web/middleware/authentication"
	"net/http"
)

func SetUpMainHandler(db *sql.DB, env string) http.Handler {
	ur := users.NewUserRepo(db)
	tr := sessions.NewTokensRepo(db)

	uc := controllers.NewUsersController(ur)
	sc := controllers.NewSessionController(ur, tr)

	m := http.NewServeMux()
	m.HandleFunc("/user", uc.CreateUser)
	m.HandleFunc("/login", sc.Login)
	m.HandleFunc("/logout", sc.Logout)
	m.HandleFunc("/auth-only", controllers.AuthenticatedOnly)

	h := middleware.CORS(m, env)
	h = middleware.CommonHeaders(h)
	h = auth.Authenticate(h, ur, tr)
	return h
}
