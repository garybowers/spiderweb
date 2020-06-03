package main

import (
	"log"
	"net/http"
	"spiderweb/kubernetes"
)

func index(res http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, cookieName)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		log.Println(req.RemoteAddr, err.Error())
		return
	}
	user := getUser(session)
	log.Println(req.RemoteAddr, user)
	if auth := user.Authenticated; !auth {
		log.Println(req.RemoteAddr, "Session not authenticated, redirecting to sign on")
		err = session.Save(req, res)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			log.Println(req.RemoteAddr, "ERROR", err.Error())
			return
		}
		http.Redirect(res, req, "/auth/google/login", http.StatusFound)
		return
	}
}

func logout(res http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, cookieName)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		log.Println(req.RemoteAddr, err.Error())
		return
	}

	session.Values["user"] = User{}
	session.Options.MaxAge = -1

	err = session.Save(req, res)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		log.Println(req.RemoteAddr, err.Error())
		return
	}
	http.Redirect(res, req, "/", http.StatusFound)
}

func session(res http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, cookieName)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		log.Println(req.RemoteAddr, err.Error())
		return
	}

	user := getUser(session)
	log.Println(req.RemoteAddr, "Session for", user.Username)

	if auth := user.Authenticated; !auth {
		session.AddFlash("You don't have access!")
		err = session.Save(req, res)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			log.Println(req.RemoteAddr, user.Username, err.Error())
			return
		}
		http.Redirect(res, req, "/", http.StatusFound)
		return
	}

	kubernetes.GetDeployment("spider", "dev1")

	serveReverseProxy(getBackendURL(), res, req)
}
