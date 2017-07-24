package handler

import (
	"fmt"
	"net/http"
	"strings"

	jose "github.com/dvsekhvalnov/jose2go"
	"github.com/dvsekhvalnov/jose2go/keys/ecc"
)

const NEW_ADDRESS = "https://ampbyexample.com"
const DEFAULT_MAX_AGE = 60

var hmacSampleSecret = []byte("my_secret_key")

func RedirectToSecureVersion(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, NEW_ADDRESS+r.URL.Path, http.StatusMovedPermanently)
}

func IsInsecureRequest(r *http.Request) bool {
	return r.TLS == nil && !strings.HasPrefix(r.Host, "localhost")
}

func isFormPostRequest(method string, w http.ResponseWriter) bool {
	if method != "POST" {
		http.Error(w, "post only", http.StatusMethodNotAllowed)
		return false
	}
	return true
}
func VerifyCSRF(w http.ResponseWriter, r *http.Request) error {

	var token string

	// Get token from the Authorization header

	tokens, ok := r.Header["X-Csrf-Token"]
	if ok && len(tokens) >= 1 {
		token = tokens[0]
	}

	return CheckThenMarkToken(token)
}
func GenerateCSRF(w http.ResponseWriter, r *http.Request) string {
	payload := `{"hello":"world"}`

	privateKey := ecc.NewPrivate([]byte{4, 114, 29, 223, 58, 3, 191, 170, 67, 128, 229, 33, 242, 178, 157, 150, 133, 25, 209, 139, 166, 69, 55, 26, 84, 48, 169, 165, 67, 232, 98, 9},
		[]byte{131, 116, 8, 14, 22, 150, 18, 75, 24, 181, 159, 78, 90, 51, 71, 159, 214, 186, 250, 47, 207, 246, 142, 127, 54, 183, 72, 72, 253, 21, 88, 53},
		[]byte{42, 148, 231, 48, 225, 196, 166, 201, 23, 190, 229, 199, 20, 39, 226, 70, 209, 148, 29, 70, 125, 14, 174, 66, 9, 198, 80, 251, 95, 107, 98, 206})

	token, err := jose.Sign(payload, jose.ES256, privateKey)

	if err == nil {
		//go use token
		fmt.Printf("\ntoken = %v\n", token)
	}
	return token
}

func GetOrigin(r *http.Request) string {
	origin := r.Header.Get("Origin")
	if origin != "" {
		return origin
	}
	return GetHost(r)
}

func GetSourceOrigin(r *http.Request) string {
	// TODO perform checks if source origin is allowed
	return r.URL.Query().Get("__amp_source_origin")
}

func GetHost(r *http.Request) string {
	if r.TLS == nil {
		return "http://" + r.Host
	}
	return "https://" + r.Host
}

func SetContentTypeJson(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func SetDefaultMaxAge(w http.ResponseWriter) {
	SetMaxAge(w, DEFAULT_MAX_AGE)
}

func SetMaxAge(w http.ResponseWriter, age int) {
	w.Header().Set("cache-control", fmt.Sprintf("max-age=%d, public, must-revalidate", age))
}

func handlePost(w http.ResponseWriter, r *http.Request, postHandler func(http.ResponseWriter, *http.Request)) {
	if r.Method != "POST" {
		http.Error(w, "post only", http.StatusMethodNotAllowed)
		return
	}
	postHandler(w, r)
}
