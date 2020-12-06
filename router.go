package main

import (
	"github.com/gorilla/mux"
	"github.com/nocubicles/skillbase.io/src/middleware"
	"github.com/nocubicles/skillbase.io/src/routes"
)

func router() *mux.Router {
	router := mux.NewRouter()
	router.Handle("/", middleware.AppHandler(routes.RenderHome)).Methods("GET")

	router.Handle("/auth/github/signin", middleware.AppHandler(routes.GithubOauthLogin)).Methods("GET", "POST")
	router.Handle("/auth/github/callback", middleware.AppHandler(routes.GithubOauthCallback)).Methods("GET", "POST")

	return router
}
