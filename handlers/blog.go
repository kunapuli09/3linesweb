package handlers

import (
	"errors"
	"github.com/gorilla/sessions"
	"github.com/kunapuli09/3linesweb/libhttp"
	"github.com/kunapuli09/3linesweb/models"
	"html/template"
	"net/http"
	"strconv"
	"fmt"
)

type Blog struct {
	Name   string
	Secure bool
}

var m = map[int]*Blog{
	1:  &Blog{"templates/blog/blog1.html.tmpl", false},
	2:  &Blog{"templates/blog/blog2.html.tmpl", false},
	3:  &Blog{"templates/blog/blog3.html.tmpl", false},
	4:  &Blog{"templates/blog/blog4.html.tmpl", false},
	5:  &Blog{"templates/blog/blog5.html.tmpl", false},
	6:  &Blog{"templates/blog/blog6.html.tmpl", false},
	7:  &Blog{"templates/blog/blog7.html.tmpl", false},
	8:  &Blog{"templates/blog/blog8.html.tmpl", false},
	9:  &Blog{"templates/blog/blog9.html.tmpl", false},
	10: &Blog{"templates/blog/blog10.html.tmpl", false},
	11: &Blog{"templates/blog/blog11.html.tmpl", false},
	12: &Blog{"templates/blog/blog12.html.tmpl", false},
	13: &Blog{"templates/blog/blog13.html.tmpl", false},
	14: &Blog{"templates/blog/blog14.html.tmpl", true},
}

func GetBlog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	blogNumber := r.FormValue("blogNumber")
	i, err := strconv.Atoi(blogNumber)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	b, ok1 := m[i]
	if !ok1 {
		libhttp.HandleErrorJson(w, errors.New("no blog"))
		return
	}
	tmpl, err := template.ParseFiles("templates/blog/blogdashboard.html.tmpl", b.Name)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	sessionStore := r.Context().Value("sessionStore").(sessions.Store)
	session, _ := sessionStore.Get(r, "3linesweb-session")
	currentUser, ok := session.Values["user"].(*models.UserRow)
	if ok {
		//data
		if currentUser.BlogReader || currentUser.Investor || currentUser.Dsc || currentUser.Admin {
			data := struct {
				CurrentUser *models.UserRow
			}{
				currentUser,
			}
			fmt.Println("user logged in executing blog template")
			tmpl.ExecuteTemplate(w, "bloglayout", data)
		}

	}else{
		if b.Secure == true {
			http.Redirect(w, r, "/logout", 302)
			return
		}else{
			fmt.Println("not a secure blog, so executing blog template")
			tmpl.ExecuteTemplate(w, "bloglayout", nil)
		}
	}
	
}

