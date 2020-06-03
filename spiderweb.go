package main

import (
	"encoding/gob"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"net/http"

	"log"
	"strings"
)

type User struct {
	Username      string
	Authenticated bool
}

type userData struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
	Hd            string `json:"hd"`
}

var store *sessions.CookieStore
var googUser *userData

const port string = ":8080"
const cookieName string = "spiderweb-app"

func init() {
	authKeyOne := securecookie.GenerateRandomKey(64)
	encryptionKeyOne := securecookie.GenerateRandomKey(32)

	store = sessions.NewCookieStore(
		authKeyOne,
		encryptionKeyOne,
	)

	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 15,
		HttpOnly: true,
	}

	gob.Register(User{})
}

func main() {
	log.Print("SpIDErweb Startup")

	router := mux.NewRouter()
	router.HandleFunc("/auth/google/login", oauthGoogleLogin)
	router.HandleFunc("/auth/google/callback", oauthGoogleCallback)
	router.HandleFunc("/", index)
	router.HandleFunc("/session/{rest:.*}", session)
	router.HandleFunc("/logout", logout)
	log.Print("Listening on port ", port)
	http.ListenAndServe(port, router)
	//kubernetes.Deploy()
}

func getBackendURL() string {
	backendUrl := "https://spider.ssp.immersion.dev"

	//log.Println(backendUrl)
	name := cleanName(googUser.Name)
	//log.Println(name)

	newBeUrl := "http://" + name + ":3000"
	log.Println(newBeUrl)

	return backendUrl
}

func cleanName(givenName string) string {
	name := strings.ToLower(strings.Replace(givenName, " ", "", -1))
	return name
}

func getUser(s *sessions.Session) User {
	val := s.Values["user"]
	var user = User{}
	user, ok := val.(User)
	if !ok {
		return User{Authenticated: false}
	}
	return user
}
