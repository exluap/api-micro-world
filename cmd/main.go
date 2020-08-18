/**
 * Project api-microworld created by exluap
 * Date: 02.08.2020 00:24
 * twitter: https://twitter.com/exluap
 * keybase: https://keybase.io/exluap
 * website: https://exluap.com
 */

package main

import (
	_ "github.com/exluap/api-microworld/docs"
	"github.com/exluap/api-microworld/internal/user"
	"github.com/exluap/api-microworld/internal/utils"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
)

// @title MicroWorld Swagger API
// @version 1.0
// @description Swagger API for microworld project
// @basepath /api

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},                            // All origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"}, // Allowing only get, just an example
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
	})

	r.Route("/api", func(r chi.Router) {

		//Что касается  пользователей
		r.Route("/user", func(r chi.Router) {
			r.Post("/register", user.RegisterUser)
			r.Post("/login", user.AuthUser)

			r.Route("/{userId}", func(r chi.Router) {
				r.Use(utils.JwtAuthentication)
				r.Get("/info", user.GetUserInfo)
				r.Delete("/", user.DeleteUser)
			})
		})

	})

	r.Get("/swagger/*", httpSwagger.Handler())

	log.Fatal(http.ListenAndServe(":8080", c.Handler(r)))
}
