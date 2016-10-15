package service

import (
    "testing"
    "time"

    "golang.org/x/crypto/bcrypt"
)

func TestInvalidToken(t *testing.T) {
    now := time.Now().Unix()
    token := Token{
        Key:        "TEST",
        UserID:     1,
        ExpiresAt:  now,
    }
    if token.IsValid() == true{
        t.Errorf("Token should be invalid")
    }
}

func TestValidToken(t *testing.T) {
    time := time.Now().Add(time.Hour * 24 * 7 * time.Duration(8)).Unix()
    token := Token{
        Key:        "TEST",
        UserID:     1,
        ExpiresAt:  time,
    }
    if token.IsValid() == false{
        t.Errorf("Token should be valid")
    }
}

func TestGenToken(t *testing.T) {
    userID := uint(1)
    token, error := genToken(userID)
    if error != nil {
        t.Errorf(error.Error())
    }
    if token == "" {
        t.Errorf("Failed to generate token")
    }
}

func TestComparePasswordSuccess(t *testing.T) {
    password := "testpassword"
    user := User{
        Username: "testname",
        Password: password,
    }
    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    user.Password = string(hashedPassword)
    valid := user.comparePassword(password)
    if valid == false {
        t.Errorf("Passwords did not match up")
    }
}

func TestComparePasswordDifferentPassword(t *testing.T) {
    password := "testpassword2"
    user := User{
        Username: "testname",
        Password: "testpassword",
    }
    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    user.Password = string(hashedPassword)
    valid := user.comparePassword(password)
    if valid == true {
        t.Errorf("Passwords should not match up")
    }
}

func TestUserBeforeSave(t *testing.T) {
    password := "testpassword"
    user := User {
        Username: "testname",
        Password: password,
    }

    err := user.BeforeSave()
    if err != nil {
        t.Error("User before saved")
    }

    if user.Password == password {
        t.Errorf("Failed to hash password")
    }

    valid := user.comparePassword(password)
    if valid == false {
        t.Errorf("Passwords should match up")
    }
}

func TestTokenBeforeSave(t *testing.T) {
    token := Token{UserID: 1}
    err := token.BeforeSave()
    if err != nil {
        t.Error("Token before save failed")
    }

    if token.Key == "" || token.ExpiresAt == 0 {
        t.Error("Failed to generate token key and expiration time")
    }
}
