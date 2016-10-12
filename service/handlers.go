package service

import (
    "net/http"
    "io/ioutil"
    "encoding/json"

    "github.com/gorilla/mux"
    "github.com/unrolled/render"
)

//postUserHandler handles user creation
func postUserHandler(formatter *render.Render) http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
        var user User
        payload, _ := ioutil.ReadAll(req.Body)
        err := json.Unmarshal(payload, &user)
        if err != nil {
            formatter.JSON(w, http.StatusBadRequest, "Failed to parse user.")
        }
        err = DB.Create(&user).Error
        if err != nil {
            formatter.JSON(w, http.StatusInternalServerError, "Failed to create user.")
            return
        }
        formatter.JSON(w, http.StatusCreated, "User succesfully created.")
    }
}

//Handles login by creating and returing a new token
func postLoginHandler(formatter *render.Render) http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
        var user User
        var token Token
        payload, _ := ioutil.ReadAll(req.Body)
        err := json.Unmarshal(payload, &user)
        if err != nil {
            formatter.JSON(w, http.StatusBadRequest, "Failed to parse user.")
        }
        oldPassword := user.Password
        err = DB.Where("username=?", user.Username).First(&user).Error
        if err != nil {
            formatter.JSON(w, http.StatusNotFound, "Username/password not valid")
            return
        }
        compare := user.comparePassword(oldPassword)
        if compare == false {
            formatter.JSON(w, http.StatusNotFound, "Username/password not valid")
            return
        }
        token.UserID = user.ID
        err = DB.Create(&token).Error
        if err != nil {
            formatter.JSON(w, http.StatusInternalServerError, "Failed to create token.")
            return
        }
        formatter.JSON(w, http.StatusCreated, token)
    }
}

//getTokenHandler takes a token
func getTokenValidate(formatter *render.Render) http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
        token := Token{}
        vars := mux.Vars(req)
        key := vars["token"]
        err := DB.Where("key= ?", key).First(&token).Error
        if err != nil {
            formatter.JSON(w, http.StatusNotFound, "Failed to find token")
            return
        }
        if token.IsValid() == false {
            formatter.JSON(w, http.StatusUnauthorized, "Token expired")
        }
        formatter.JSON(w, http.StatusFound, token)
    }
}
