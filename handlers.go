package main

import (
	"fmt"
	"github.com/gorilla/mux"
	appsv1 "k8s.io/api/apps/v1"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
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
	http.Redirect(res, req, "/session/", http.StatusFound)
}

func logout(res http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, cookieName)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		log.Println(req.RemoteAddr, err.Error())
		return
	}

	user := getUser(session)
	destroyEnvironment(user)

	session.Options.MaxAge = -1
	err = session.Save(req, res)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		log.Println(req.RemoteAddr, err.Error())
		return
	}

	http.Redirect(res, req, "/dashboard/", http.StatusFound)
}

func dashboard(res http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, cookieName)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		log.Println(req.RemoteAddr, err.Error())
		return
	}

	user := getUser(session)
	log.Println(req.RemoteAddr, "Session for", user.Email)

	if auth := user.Authenticated; !auth {
		session.AddFlash("You don't have access!")
		err = session.Save(req, res)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			log.Println(req.RemoteAddr, user.Email, err.Error())
			return
		}
		http.Redirect(res, req, "/", http.StatusFound)
		return
	}
	var list *appsv1.DeploymentList
	list = kubernetes.GetDeployments(namespace)
	fmt.Fprintf(res, "<html><head><title>SpIDEr</title></head><body>")
	fmt.Fprintf(res, "Email: %s", user.Email)
	fmt.Fprintf(res, "Forename: %s", user.Forename)

	fmt.Fprintf(res, "<table>")
	fmt.Fprintf(res, "<td>Deployment</td><td>Replicas</td><td>Connect</td>")
	for _, d := range list.Items {
		fmt.Fprintf(res, "<tr><td>%s</td><td>%d</td><td><input type='button' value='Connect' /></tr>", d.Name, *d.Spec.Replicas)
	}
	fmt.Fprintf(res, "</table>")
	fmt.Fprintf(res, "</html>")
}

func session(res http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, cookieName)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		log.Println(req.RemoteAddr, err.Error())
		return
	}

	user := getUser(session)
	log.Println(req.RemoteAddr, "Session for", user.Email)

	if auth := user.Authenticated; !auth {
		session.AddFlash("You don't have access!")
		err = session.Save(req, res)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			log.Println(req.RemoteAddr, user.Email, err.Error())
			return
		}
		http.Redirect(res, req, "/", http.StatusFound)
		return
	}

	var i int = 0
	for _, d := range kubernetes.GetDeployments(namespace).Items {
		if cleanEmail(user.Email) == d.Name {
			i = 1
			fmt.Printf("Found existing deployment: * %s (%d replicas)\n", d.Name, *d.Spec.Replicas)
		}
	}

	if i == 0 {
		createEnvironment(user)
		log.Println(req.RemoteAddr, user.Email, "Creating deployment for user %s", cleanEmail(user.Email))

	}
	serveReverseProxy(getBackendURL(cleanEmail(user.Email)), res, req)
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/favicon.ico")
}

func serveReverseProxy(target string, res http.ResponseWriter, req *http.Request) {
	url, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(url)

	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = url.Host

	req.URL.Path = mux.Vars(req)["rest"]
	log.Println("Proxying", target)

	proxy.ServeHTTP(res, req)
}

func int32Ptr(i int32) *int32 { return &i }
func int64Ptr(i int64) *int64 { return &i }
