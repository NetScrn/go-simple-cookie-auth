package controllers

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/netscrn/gocookieauth/data/sessions"
	"github.com/netscrn/gocookieauth/data/users"
	auth "github.com/netscrn/gocookieauth/web/middleware/authentication"
	"github.com/netscrn/gocookieauth/web/security"
	"io/ioutil"
	"net/http"
	"time"
)

type SessionsController struct {
	um users.Manger
	tm sessions.TokenManager
}

func NewSessionController(um users.Manger, tm sessions.TokenManager) SessionsController {
	return SessionsController{
		um: um,
		tm: tm,
	}
}

func (sc SessionsController) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if _, ok := auth.TokenFromContext(ctx); ok {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	loginDataJson, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Login - can't read req body: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(loginDataJson) == 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	loginData := struct {
		Email      string `json:"email"`
		Password   string `json:"password"`
		RememberMe bool   `json:"remember_me"`
	}{}

	err = json.Unmarshal(loginDataJson, &loginData)
	if err != nil {
		fmt.Printf("Login - can't parse req body json: %v\n", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	u, err := sc.um.GetUserByEmail(ctx, loginData.Email)
	if err == users.ErrNoUserFound {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ok, err := security.IsPassMatchHash(loginData.Password, u.PasswordDigest)
	if err != nil {
		fmt.Printf("Login - can't compare password hashes: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if !ok {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var expiry time.Time
	var now = time.Now()
	if loginData.RememberMe {
		expiry = now.Add(24 * 30 * 3 * time.Hour)
	} else {
		expiry = now.Add(24 * time.Hour)
	}

	token := sessions.Token{
		UserID:     u.Id,
		Expiry:     expiry,
		Attributes: nil,
	}

	tokenId, err := sc.tm.Create(ctx, token)
	if err != nil {
		fmt.Printf("Login - can't create token: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	c := http.Cookie{
		Name:     "u_session",
		Value:    tokenId,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}
	if loginData.RememberMe {
		c.MaxAge = int(expiry.Sub(now) / time.Second)
	}
	http.SetCookie(w, &c)

	rawCsrfToken := sha256.Sum256([]byte(tokenId))
	csrfToken := base64.StdEncoding.EncodeToString(rawCsrfToken[:])
	resBody := []byte(fmt.Sprintf(`{"token": "%s"}`, csrfToken))

	w.WriteHeader(http.StatusCreated)
	w.Write(resBody)
}

func (sc SessionsController) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != "DELETE" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	_, ok := auth.TokenFromContext(ctx)
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	tokenId, err := r.Cookie("u_session")
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	c := &http.Cookie{
		Name:     "u_session",
		Value:    "",
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}
	http.SetCookie(w, c)

	err = sc.tm.Revoke(ctx, tokenId.Value)
	if err == sessions.ErrNoTokenWasDeleted {
		fmt.Printf("Logout - no token was deleted")
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
