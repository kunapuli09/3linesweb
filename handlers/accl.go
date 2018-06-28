package handlers

import (
	"fmt"
	"github.com/fatih/structs"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/kunapuli09/3linesweb/libhttp"
	"github.com/kunapuli09/3linesweb/models"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

func NewApplication(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	//create empty investmentstructure
	tmpl, err := template.ParseFiles("templates/portfolio/appl.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	tmpl.ExecuteTemplate(w, "content", nil)
}

func FundingRequests(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	db := r.Context().Value("db").(*sqlx.DB)
	sessionStore := r.Context().Value("sessionStore").(sessions.Store)
	session, _ := sessionStore.Get(r, "3linesweb-session")
	currentUser, ok := session.Values["user"].(*models.UserRow)
	if !ok {
		http.Redirect(w, r, "/logout", 302)
		return
	}
	allreqs, err := models.NewAppl(db).AllAppls(nil)
	//create session date for page rendering
	data := struct {
		CurrentUser *models.UserRow
		Count int
		Existing    []*models.ApplRow
	}{
		currentUser,
		getCount(w,r, currentUser.Email),
		allreqs,
	}
	tmpl, err := template.ParseFiles("templates/portfolio/basic.html.tmpl", "templates/portfolio/fundingreqs.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	tmpl.ExecuteTemplate(w, "layout", data)
}

//database call to add new
func AddApplication(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	var i models.ApplRow
	db := r.Context().Value("db").(*sqlx.DB)
	err := r.ParseForm()
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	decoder := schema.NewDecoder()
	err1 := decoder.Decode(&i, r.PostForm)
	if err1 != nil {
		fmt.Println("decoding error")
		libhttp.HandleErrorJson(w, err1)
		return
	}
	m := structs.Map(i)
	m["ApplicationDate"] = time.Now()
	fmt.Printf("map %v", m)
	_, err2 := models.NewAppl(db).Create(nil, m)
	if err2 != nil {
		fmt.Println("database error", err2)
		libhttp.HandleErrorJson(w, err2)
		return
	}
	http.Redirect(w, r, "/", 302)
}

//db call to update
func RemoveApplication(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	db := r.Context().Value("db").(*sqlx.DB)
	ID, e := strconv.ParseInt(r.FormValue("id"), 10, 64)
	if e != nil {
		libhttp.HandleErrorJson(w, e)
		return
	}
	_, err2 := models.NewAppl(db).DeleteByID(nil, ID)
	if err2 != nil {
		libhttp.HandleErrorJson(w, err2)
		return
	}
	http.Redirect(w, r, "/", 302)
}
