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
	"github.com/rs/cors"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},                  // All origins
		AllowedMethods:   []string{"GET", "POST", "PUT"}, // Allowing only get, just an example
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
	})

	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/user/register", user.RegisterUser).Methods("GET")
	api.HandleFunc("/user/login", user.AuthUser).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", c.Handler(r)))
}
