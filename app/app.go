package app

import (
	"crypto/sha256"
	"crypto/subtle"
	httpServer "net/http"

	"testWB/pkg/http"
	"testWB/pkg/user"
)

type App struct {
	handler *http.Handler
}

func New() *App {
	storage := &user.Storage{}
	return &App{
		handler: http.NewHandler(storage),
	}
}

func (a *App) RunGet(listenedAddr string) error {
	mux := httpServer.NewServeMux()
	mux.HandleFunc("/get", a.basicAuth(a.handler.GetUserGrade))

	return httpServer.ListenAndServe(listenedAddr, mux)
}

func (a *App) RunSet(listenedAddr string) error {
	mux := httpServer.NewServeMux()
	mux.HandleFunc("/set", a.handler.SetUserGrade)

	return httpServer.ListenAndServe(listenedAddr, mux)
}

func (a *App) basicAuth(next httpServer.HandlerFunc) httpServer.HandlerFunc {
	return func(w httpServer.ResponseWriter, r *httpServer.Request) {
		username, password, ok := r.BasicAuth()
		if ok {
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))
			expectedUsernameHash := sha256.Sum256([]byte("username"))
			expectedPasswordHash := sha256.Sum256([]byte("password"))

			usernameMatch := subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1
			passwordMatch := subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1

			if usernameMatch && passwordMatch {
				next.ServeHTTP(w, r)
				return
			}
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		httpServer.Error(w, "Unauthorized", httpServer.StatusUnauthorized)
	}
}
