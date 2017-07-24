package backend

import (
	"app/handler"
	"fmt"
	"net/http"
	"time"
)

type AccessData struct {
	ReturnURL string
}
type AuthRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	RememberMe string `json:"rememberme"`
}
type AuthResponse struct {
	Access bool   `json:"access"`
	Role   string `json:"role"`
	Name   string `json:"name"`
}
type Csrf struct {
	Token string
}

type AuthorizationResponse interface {
	CreateAuthorizationResponse() AuthorizationResponse
}

const (
	ACCESS_PATH   = "/access/"
	ACCESS_COOKIE = "ABE_LOGGED_IN"
	CSRF_COOKIE   = "X-CSRF-Token"
)

func HandlePingback(w http.ResponseWriter, r *http.Request) {
	fmt.Println("PONG")

	token := handler.NewToken()

	handler.HandleJsonResponse(w, r, &Csrf{
		Token: token,
	})
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {

	err := handler.VerifyCSRF(w, r)

	if err != nil {
		handler.HandleJsonResponse(w, r, &AuthResponse{
			Access: false,
			Role:   "",
			Name:   "",
		})
		return
	}
	handler.HandleJsonResponse(w, r, &AuthResponse{
		Access: true,
		Role:   "subscriber",
		Name:   "Joshua",
	})
	//log.Println(r.Body)

	// decoder := json.NewDecoder(r.Body)
	// var u AuthRequest
	// err := decoder.Decode(&u)
	// if err != nil {
	// 	panic(err)
	// }
	// defer r.Body.Close()
	// log.Println(u.Email)

	//returnURL := r.URL.Query().Get("return")
	//redirect after login
}

func handleAuthorization(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie(ACCESS_COOKIE)
	if err != nil {
		handler.HandleJsonResponse(w, r, &AuthResponse{
			Access: false,
			Role:   "",
			Name:   "",
		})
		return
	}
	//TODO: Get SessionID from DB, verify with UserID in Cookie, respond with role, access, name
	handler.HandleJsonResponse(w, r, &AuthResponse{
		Access: true,
		Role:   "subscriber",
		Name:   "Joshua",
	})
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	//handler.EnableCSRF(w, r)
	//delete the cookie
	cookie := &http.Cookie{
		Name:   ACCESS_COOKIE,
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
	returnURL := r.URL.Query().Get("return")
	http.Redirect(w, r, fmt.Sprintf("%s#success=true", returnURL), http.StatusSeeOther)
}

func handleSubmit(w http.ResponseWriter, r *http.Request) {
	//handler.EnableCSRF(w, r)

	//TODO: Confirm username and password, then set SessionID along with UserID
	expireInOneDay := time.Now().AddDate(0, 0, 1)
	cookie := &http.Cookie{
		Name:    ACCESS_COOKIE,
		Expires: expireInOneDay,
		Value:   "true",
	}
	http.SetCookie(w, cookie)

	handleAuthorization(w, r)
	returnURL := r.FormValue("returnurl")
	http.Redirect(w, r, fmt.Sprintf("%s#success=true", returnURL), http.StatusSeeOther)
}
