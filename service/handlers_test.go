package service

import (
    "bytes"
    "encoding/json"
    "errors"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"

    "golang.org/x/crypto/bcrypt"
	"github.com/unrolled/render"
    "github.com/urfave/negroni"
	"github.com/gorilla/mux"
)

var (
    formatter = render.New(render.Options{
        IndentJSON: true,
    })
)

type testRepo struct {
    users  []User
    tokens []Token
}

func (t *testRepo) addToken(token Token) error {
    t.tokens = append(t.tokens, token)
    return nil
}

func (t *testRepo) addUser(user User) error {
    t.users = append(t.users, user)
    return nil
}

func (t *testRepo) getUserByUsername(username string) (User, error) {
    for _, user := range t.users {
        if user.Username == username {
            return user, nil
        }
    }
    return User{}, errors.New("User not found")
}

func (t *testRepo) getTokenByKey(key string) (Token, error) {
    for _, token := range t.tokens {
        if token.Key == key {
            return token, nil
        }
    }
    return Token{}, errors.New("Token not found")
}

func TestPostUserHandlerInvalidJSON(t *testing.T) {
    repo := &testRepo{}
    client := &http.Client{}

    server := httptest.NewServer(http.HandlerFunc(postUserHandler(formatter, repo)))
    defer server.Close()

    body := []byte("this is not valid json")
    req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(body))
    if err != nil {
        t.Errorf("Error in creating POST request for createMatchHandler: %v", err)
    }
    res, err := client.Do(req)
    if err != nil {
        t.Errorf("Error in POST to createMatchHandler: %v", err)
    }
    defer res.Body.Close()
    if res.StatusCode != http.StatusBadRequest {
        t.Error("Sending invalid JSON should result in a bad request from server.")
    }
}

func TestPostUserHandlerNotUser(t *testing.T) {
    repo := &testRepo{}
    client := &http.Client{}

    server := httptest.NewServer(http.HandlerFunc(postUserHandler(formatter, repo)))
    defer server.Close()

    body := []byte("{\"test\":\"Not user.\"}")
    req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(body))
    if err != nil {
        t.Errorf("Error in creating second POST request for invalid data on create match: %v", err)
    }
    req.Header.Add("Content-Type", "application/json")
    res, _ := client.Do(req)
    defer res.Body.Close()
    if res.StatusCode != http.StatusBadRequest {
        t.Error("Sending valid JSON but with incorrect or missing fields should result in a bad request and didn't.")
    }
}

func TestPostUserHandlerValidUser(t *testing.T) {
    repo := &testRepo{}
    client := &http.Client{}

    server := httptest.NewServer(http.HandlerFunc(postUserHandler(formatter, repo)))
    defer server.Close()
    body := []byte("{\"username\":\"testname\",\n\"password\":\"testpassword\"}")
    req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(body))

    if err != nil {
        t.Errorf("Error in creating second POST request for invalid data on create match: %v", err)
    }
    req.Header.Add("Content-Type", "application/json")
    res, _ := client.Do(req)
    defer res.Body.Close()
    if res.StatusCode != http.StatusCreated {
        t.Error("Sending valid JSON but with incorrect or missing fields should result in a bad request and didn't.")
    }

    if len(repo.users) != 1 {
        t.Error("Failed to add user.")
    }
    user := repo.users[0]
    if user.Username != "testname" && user.Password != "password"{
        t.Errorf("Failed to add info to user")
    }
}

func TestPostLoginHandlerInvalidJSON(t *testing.T) {
    repo := &testRepo{}
    client := &http.Client{}

    server := httptest.NewServer(http.HandlerFunc(postLoginHandler(formatter, repo)))
    defer server.Close()

    body := []byte("this is not valid json")
    req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(body))
    if err != nil {
        t.Errorf("Error in creating POST request for createMatchHandler: %v", err)
    }
    res, err := client.Do(req)
    if err != nil {
        t.Errorf("Error in POST to createMatchHandler: %v", err)
    }
    defer res.Body.Close()
    if res.StatusCode != http.StatusBadRequest {
        t.Error("Sending invalid JSON should result in a bad request from server.")
    }
}

func TestPostLoginHandlerNotValidLoginInfo(t *testing.T) {
    repo := &testRepo{}
    client := &http.Client{}
    password, _ := bcrypt.GenerateFromPassword([]byte("testpassword"), bcrypt.DefaultCost)
    user := User{Username:"testname", Password: string(password)}
    repo.addUser(user)

    server := httptest.NewServer(http.HandlerFunc(postLoginHandler(formatter, repo)))
    defer server.Close()
    body := []byte("{\"username\":\"testname\",\n\"password\":\"testpassword2\"}")
    req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(body))

    if err != nil {
        t.Errorf("Error in creating second POST request for invalid data on create match: %v", err)
    }
    req.Header.Add("Content-Type", "application/json")
    res, _ := client.Do(req)
    defer res.Body.Close()
    if res.StatusCode != http.StatusBadRequest {
        t.Error("Failed to invalidate login issue.")
    }
}

func TestPostLoginHandlerValidLoginInfo(t *testing.T) {
    repo := &testRepo{}
    client := &http.Client{}
    password, _ := bcrypt.GenerateFromPassword([]byte("testpassword"), bcrypt.DefaultCost)
    user := User{Username:"testname", Password: string(password)}
    repo.addUser(user)

    server := httptest.NewServer(http.HandlerFunc(postLoginHandler(formatter, repo)))
    defer server.Close()
    body := []byte("{\"username\":\"testname\",\n\"password\":\"testpassword\"}")
    req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(body))

    if err != nil {
        t.Errorf("Error in creating second POST request for invalid data on create match: %v", err)
    }
    req.Header.Add("Content-Type", "application/json")
    res, _ := client.Do(req)
    defer res.Body.Close()
    if res.StatusCode != http.StatusOK {
        t.Error("Failed to invalidate login issue.")
    }
}

func TestGetTokenValidationInvalidToken(t *testing.T) {
    var (
        request  *http.Request
        recorder *httptest.ResponseRecorder
    )
    repo  := &testRepo{}

    server := MakeTestServer(repo)
    token := Token{Key: "test", ExpiresAt: time.Now().Unix()}
    repo.addToken(token)

    recorder = httptest.NewRecorder()
    request, _ = http.NewRequest("GET", "/auth/token/test2", nil)
    server.ServeHTTP(recorder, request)

    if recorder.Code != http.StatusNotFound {
        t.Errorf("Expected %v; received %v", http.StatusNotFound, recorder.Code)
    }
}

func TestGetTokenValidationValidToken(t *testing.T) {
    var (
        request  *http.Request
        recorder *httptest.ResponseRecorder
    )
    repo  := &testRepo{}

    server := MakeTestServer(repo)
    token := Token{Key: "test", ExpiresAt: time.Now().Add(time.Hour * 24 * 7 * time.Duration(8)).Unix()}
    repo.addToken(token)

    recorder = httptest.NewRecorder()
    request, _ = http.NewRequest("GET", "/auth/token/"+token.Key, nil)
    server.ServeHTTP(recorder, request)

    if recorder.Code != http.StatusOK {
        t.Errorf("Expected %v; received %v", http.StatusOK, recorder.Code)
    }

    var tokenResponse Token
    err := json.Unmarshal(recorder.Body.Bytes(), &tokenResponse)
    if err != nil {
        t.Errorf("Error unmarshaling token: %s", err)
    }
    if tokenResponse.Key != "test" {
        t.Errorf("Expected token key to be test; received %d", token.Key)
    }
}

func TestGetTokenValidationExpiredToken(t *testing.T) {
    var (
        request  *http.Request
        recorder *httptest.ResponseRecorder
    )
    repo  := &testRepo{}

    server := MakeTestServer(repo)
    token := Token{Key: "test", ExpiresAt: time.Now().Unix()}
    repo.addToken(token)

    recorder = httptest.NewRecorder()
    request, _ = http.NewRequest("GET", "/auth/token/"+token.Key, nil)
    server.ServeHTTP(recorder, request)

    if recorder.Code != http.StatusUnauthorized {
        t.Errorf("Token Expired: Expected %v; received %v", http.StatusUnauthorized, recorder.Code)
    }
}

func MakeTestServer(repository *testRepo) *negroni.Negroni {
	server := negroni.New()
	mx := mux.NewRouter()
	initRoutes(mx, formatter, repository)
	server.UseHandler(mx)
	return server
}
