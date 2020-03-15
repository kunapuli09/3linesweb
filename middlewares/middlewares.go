// Package middlewares provides common middleware handlers.
package middlewares

import (
	"net/http"

	"context"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/kunapuli09/3linesweb/models"
	"fmt"
)

func SetDB(db *sqlx.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			req = req.WithContext(context.WithValue(req.Context(), "db", db))

			next.ServeHTTP(res, req)
		})
	}
}

func SetSessionStore(sessionStore sessions.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			req = req.WithContext(context.WithValue(req.Context(), "sessionStore", sessionStore))

			next.ServeHTTP(res, req)
		})
	}
}

// MustLogin is a middleware that checks existence of current user.
func MustLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		sessionStore := req.Context().Value("sessionStore").(sessions.Store)
		session, _ := sessionStore.Get(req, "3linesweb-session")
		userRowInterface := session.Values["user"]

		if userRowInterface == nil {
			http.Redirect(res, req, "/login", 302)
			return
		}
		next.ServeHTTP(res, req)
	})
}

// MustSecure is a middleware that checks existence of current user.
func MustSecure(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "text/html")
		path := req.URL.EscapedPath()
		//prepare page
		db := req.Context().Value("db").(*sqlx.DB)
		sessionStore := req.Context().Value("sessionStore").(sessions.Store)
		session, _ := sessionStore.Get(req, "3linesweb-session")
		currentUser, ok := session.Values["user"].(*models.UserRow)
		//fmt.Println(req)
		if !ok {
			fmt.Printf("Not logged in but tried to access %v", path)
			http.Redirect(res, req, "/logout", 302)
			return
		}
		alldocs, _ := models.NewUserDoc(db).GetAllByUserId(nil, currentUser.ID)
		_, exists := Find(alldocs, path)
		if !exists {
			fmt.Printf("No privileges but tried to access %v", path)
			http.Redirect(res, req, "/logout", 302)
			return
		}
		next.ServeHTTP(res, req)
	})
}

func Find(slice []*models.UserDocRow, val string) (int, bool) {
    for i, item := range slice {
        if item.DocName == val {
            return i, true
        }
    }
    return -1, false
}
