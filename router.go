package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nocubicles/skillbase.io/src/middleware"
	"github.com/nocubicles/skillbase.io/src/routes"
)

func router() *mux.Router {

	router := mux.NewRouter()
	router.Use(middleware.CORS)
	router.Use(middleware.LoggingMiddleware)
	router.HandleFunc("/", middleware.CheckCookie(routes.RenderHome)).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/auth/github/signin", routes.GithubOauthLogin).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
	router.HandleFunc("/auth/github/callback", routes.GithubOauthCallback).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)

	router.HandleFunc("/app", middleware.CheckIsUsedLoggedIn(routes.RenderApp)).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/organizations", middleware.CheckIsUsedLoggedIn(routes.Organization)).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)

	router.HandleFunc("/repositories", middleware.CheckIsUsedLoggedIn(routes.RepoHandler)).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
	router.HandleFunc("/repositories/tracking", middleware.CheckIsUsedLoggedIn(routes.RepoHandler)).Methods(http.MethodPut, http.MethodOptions)

	router.HandleFunc("/labels", middleware.CheckIsUsedLoggedIn(routes.LabelHandler)).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
	router.HandleFunc("/labels/tracking", middleware.CheckIsUsedLoggedIn(routes.LabelHandler)).Methods(http.MethodPut, http.MethodOptions)

	router.HandleFunc("/sync", middleware.CheckIsUsedLoggedIn(routes.Sync)).Methods(http.MethodPost, http.MethodOptions)
	router.Use(mux.CORSMethodMiddleware(router))
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/"))))
	return router
}
