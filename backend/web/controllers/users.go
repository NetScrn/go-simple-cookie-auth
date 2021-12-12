package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/netscrn/gocookieauth/data/users"
	auth "github.com/netscrn/gocookieauth/web/middleware/authentication"
	"github.com/netscrn/gocookieauth/web/security"
)

type UsersController struct {
	um users.Manger
}

func NewUsersController(um users.Manger) UsersController {
	return UsersController{
		um: um,
	}
}

func (uc UsersController) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	if _, ok := auth.TokenFromContext(ctx); ok {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	userDataJson, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("CreateUser - can't read req body: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(userDataJson) == 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	regUserData := struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}{}
	err = json.Unmarshal(userDataJson, &regUserData)
	if err != nil {
		fmt.Printf("CreateUser - can't parse req body json: %v\n", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	digest, err := security.CreatePasswordHash(regUserData.Password)
	if err != nil {
		fmt.Printf("CreateUser - can't create password hash: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	u := users.User{
		Email:          regUserData.Email,
		PasswordDigest: digest,
	}
	err = uc.um.SaveUser(ctx, &u)
	if err == users.ErrSuchEmailIsAlreadyExists {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(fmt.Sprintf(`{"error": %d, "error_desc": "%v"}`, 0, err)))
		return
	} else if err != nil {
		fmt.Printf("CreateUser - can't save user: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf(`{"userId": %d}`, u.Id)))
	}
}
