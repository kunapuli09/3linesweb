package handlers

import (
	"errors"
	"github.com/kunapuli09/3linesweb/libhttp"
	"html/template"
	"net/http"
	"strconv"
)

var m = map[int]string{
	1: "templates/blog/blog1.html.tmpl",
	2: "templates/blog/blog2.html.tmpl",
	3: "templates/blog/blog3.html.tmpl",
	4: "templates/blog/blog4.html.tmpl",
}

func GetBlog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	blogNumber := r.FormValue("blogNumber")
	i, err := strconv.Atoi(blogNumber)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	name, ok := m[i]
	if !ok {
		libhttp.HandleErrorJson(w, errors.New("no blog"))
		return
	}
	tmpl, err := template.ParseFiles(name)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	tmpl.Execute(w, r)

}
