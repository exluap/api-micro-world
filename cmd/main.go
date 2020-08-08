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
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
	"log"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},                            // All origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"}, // Allowing only get, just an example
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
	})

	r.Route("/api/user", func(r chi.Router) {
		r.Post("/register", user.RegisterUser)
		r.Post("/login", user.AuthUser)
	})

	log.Fatal(http.ListenAndServe(":8080", c.Handler(r)))
}
