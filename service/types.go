package service

import (
    "os"
    "time"
    "golang.org/x/crypto/bcrypt"
    "github.com/jinzhu/gorm"
    "github.com/dgrijalva/jwt-go"
)

//User struct handler user registration
type User struct {
    gorm.Model
    Username string `gorm:"not null;unique"json:"username"`
    Password string `gorm:"not null"json:"password"`
}

//Token struct handles authentication
type Token struct {
    gorm.Model
    Key         string   `json:"token"`
    UserID      uint      `json:"user_id"`
    ExpiresAt   int64    `json:"expires_at"`
}

//BeforeSave hash password
func (u *User) BeforeSave() (err error){
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
    u.Password = string(hashedPassword)
    if err != nil {
        return err
    }
    return nil
}

//BeforeSave generate token key
func (t *Token) BeforeSave() (err error){
    token, err := genToken(t.UserID)
    if err != nil {
        return err
    }
    t.ExpiresAt = time.Now().Add(time.Hour * 24 * 7 * time.Duration(8)).Unix()
    t.Key = string(token)
    return nil
}

//IsValid checks if a token has expired
func (t *Token) IsValid() bool {
    return t.ExpiresAt > time.Now().Unix()
}

//generates token string from username
func genToken(userID uint) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user": userID,
    })
    tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
    return tokenString, err
}

func (u *User) comparePassword(password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
    return err == nil
}
