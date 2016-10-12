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

    initRoutes(mx, formatter)
    n.UseHandler(mx)
    return n
}

func initRoutes(mx *mux.Router, formatter *render.Render) {
    mx.HandleFunc("/auth/login", postLoginHandler(formatter)).Methods("POST")
    mx.HandleFunc("/auth/register", postUserHandler(formatter)).Methods("POST")
    mx.HandleFunc("/auth/token/{token}", getTokenValidate(formatter)).Methods("GET")
}
