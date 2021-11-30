package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/netscrn/gocookieauth/data"
	"github.com/netscrn/gocookieauth/security"
)

type UsersManger interface {
	GetUserByID(userId int) (*data.User, error)
	SaveUser(user *data.User) error
}

type UsersController struct {
	um UsersManger
}

func NewUsersController(um UsersManger) UsersController {
	return UsersController{
		um: um,
	}
}

func (uc UsersController) CreateUser(w http.ResponseWriter, r *http.Request)  {
	h := w.Header()
	h.Set("Content-Type", "application/json;charset=utf-8")
	h.Set("Access-Control-Allow-Origin", "*")
	h.Set("Access-Control-Allow-Headers","*")
	h.Set("X-Content-Type-Options", "nosniff")
	h.Set("X-Frame-Options", "DENY")
	h.Set("X-XSS-Protection", "0")
	h.Set("Cache-Control", "no-store")
	h.Set("Content-Security-Policy","default-src 'none'; frame-ancestors 'none'; sandbox")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}
	defer r.Body.Close()

	userDataJson, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
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
		fmt.Println(err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	digest, err := security.CreatePasswordHash(regUserData.Password)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	u := data.User{
		Email: regUserData.Email,
		PasswordDigest: digest,
	}
	err = uc.um.SaveUser(&u)

	if err == data.ErrSuchEmailIsAlreadyExists {
		w.WriteHeader(http.StatusConflict)	
		_, err = w.Write([]byte(fmt.Sprintf(`{"error": %d, "error_desc": "%v"}`, 0, err)))
		return
	} else if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte(fmt.Sprintf(`{"userId": %d}`, u.Id)))
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}