// Copyright Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"app/backend"
	"app/config"
	"app/handler"
	"log"
	"net/http"
	"os"

	"github.com/rs/cors"
	"github.com/zenazn/goji/graceful"

	goji "goji.io"
	"goji.io/pat"
)

// var (
// 	appLog logFunc
// )

func main() {
	router := goji.NewMux()
	router.HandleFunc(pat.Get("/access/pingback"), backend.HandlePingback)
	router.HandleFunc(pat.Post("/access/login"), backend.HandleLogin)
	router.Handle(pat.Get("/*"), http.FileServer(http.Dir(config.DIST_FOLDER)))

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Content-Type", "X-CSRF-Token"},
		Debug:            false,
	})
	os.Setenv("PORT", ":8000")

	// err := graceful.ListenAndServeTLS(os.Getenv("PORT"), "server.crt", "server.key", c.Handler(router))
	err := graceful.ListenAndServeTLS(os.Getenv("PORT"), "server.crt", "server.key", c.Handler(router))

	if err != nil {
		log.Fatal("ListenAndServes: ", err)
	}
}

func HandleNotFound(h http.Handler) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
}

func ServeStaticFiles(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Host == config.OLD_URL || handler.IsInsecureRequest(r) {
			handler.RedirectToSecureVersion(w, r)
			return
		}
		handler.SetDefaultMaxAge(w)
		h.ServeHTTP(w, r)
	})
}

func exists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}
