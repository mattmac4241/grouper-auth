package service

import (
    "net/http"
    "io/ioutil"
    "encoding/json"

    "github.com/gorilla/mux"
    "github.com/unrolled/render"
)

//postUserHandler handles user creation
func postUserHandler(formatter *render.Render, repo repository) http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
        var user User
        payload, _ := ioutil.ReadAll(req.Body)
        err := json.Unmarshal(payload, &user)
        if err != nil || (user == User{}) {
            formatter.JSON(w, http.StatusBadRequest, "Failed to parse user.")
            return
        }
        err = repo.addUser(user)
        if err != nil {
            formatter.JSON(w, http.StatusInternalServerError, "Failed to create user.")
            return
        }
        formatter.JSON(w, http.StatusCreated, "User succesfully created.")
    }
}

//Handles login by creating and returing a new token
func postLoginHandler(formatter *render.Render, repo repository) http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
        var user User
        var token Token
        payload, _ := ioutil.ReadAll(req.Body)
        err := json.Unmarshal(payload, &user)
        if err != nil {
            formatter.JSON(w, http.StatusBadRequest, "Failed to parse user.")
            return
        }
        oldPassword := user.Password
        user, err = repo.getUserByUsername(user.Username)
        if err != nil {
            formatter.JSON(w, http.StatusBadRequest, "Username/password not valid")
            return
        }
        compare := user.comparePassword(oldPassword)
        if compare == false {
            formatter.JSON(w, http.StatusBadRequest, "Username/password not valid")
            return
        }
        token.UserID = user.ID
        err = repo.addToken(token)
        if err != nil {
            formatter.JSON(w, http.StatusInternalServerError, "Failed to create token.")
            return
        }
        formatter.JSON(w, http.StatusOK, token)
    }
}

//getTokenHandler takes a token
func getTokenValidate(formatter *render.Render, repo repository) http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
        vars := mux.Vars(req)
        key := vars["token"]
        token, err := repo.getTokenByKey(key)
        if err != nil {
            formatter.JSON(w, http.StatusNotFound, "Failed to find token")
            return
        }
        if token.IsValid() == false {
            formatter.JSON(w, http.StatusUnauthorized, "Token expired")
            return
        }
        formatter.JSON(w, http.StatusOK, token)
    }
}
