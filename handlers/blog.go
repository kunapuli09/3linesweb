package handlers

import (
	"errors"
	"github.com/kunapuli09/3linesweb/libhttp"
	"html/template"
	"net/http"
	"strconv"
)

var m = map[int]string{
	1:  "templates/blog/blog1.html.tmpl",
	2:  "templates/blog/blog2.html.tmpl",
	3:  "templates/blog/blog3.html.tmpl",
	4:  "templates/blog/blog4.html.tmpl",
	5:  "templates/blog/blog5.html.tmpl",
	6:  "templates/blog/blog6.html.tmpl",
	7:  "templates/blog/blog7.html.tmpl",
	8:  "templates/blog/blog8.html.tmpl",
	9:  "templates/blog/blog9.html.tmpl",
	10: "templates/blog/blog10.html.tmpl",
	11: "templates/blog/blog11.html.tmpl",
	12: "templates/blog/blog12.html.tmpl",
	13: "templates/blog/blog13.html.tmpl",
	14: "templates/blog/blog14.html.tmpl",
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
	tmpl, err := template.ParseFiles("templates/blog/blog.html.tmpl", name)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	tmpl.ExecuteTemplate(w, "layout", nil)
}
