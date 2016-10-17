package service

import (
    "github.com/urfave/negroni"
    "github.com/gorilla/mux"
    "github.com/unrolled/render"
)

// NewServer configures and returns a server.
func NewServer() *negroni.Negroni {
    formatter := render.New(render.Options{
        IndentJSON: true,
    })

    n := negroni.Classic()
    mx := mux.NewRouter()
    repo := &repoHandler{}
    initRoutes(mx, formatter, repo)
    n.UseHandler(mx)
    return n
}

func initRoutes(mx *mux.Router, formatter *render.Render, repo repository) {
    mx.HandleFunc("/auth/login", postLoginHandler(formatter, repo)).Methods("POST")
    mx.HandleFunc("/auth/register", postUserHandler(formatter, repo)).Methods("POST")
    mx.HandleFunc("/auth/token/{token}", getTokenValidate(formatter, repo)).Methods("GET")
    mx.HandleFunc("/ping", getPingHandler(formatter)).Methods("GET")
}
