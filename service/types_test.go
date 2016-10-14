package service

import (
    "testing"
    "time"
)

func TestInValidToken(t *testing.T) {
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
