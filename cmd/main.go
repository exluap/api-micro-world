/**
 * Project api-microworld created by exluap
 * Date: 02.08.2020 00:24
 * twitter: https://twitter.com/exluap
 * keybase: https://keybase.io/exluap
 * website: https://exluap.com
 */

package main

import (
	"github.com/exluap/api-microworld/internal/user"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/user/register", user.RegisterUser).Methods("GET")
	api.HandleFunc("/user/login", user.AuthUser).Methods("GET")

	err := http.ListenAndServe("/", r)
	if err != nil {
		log.Fatal("Error with starting server: ", err)
	}
}
