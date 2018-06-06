package handlers

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/kunapuli09/3linesweb/libhttp"
	"github.com/kunapuli09/3linesweb/models"
	"github.com/gorilla/sessions"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

func NewApplication(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	//create empty investmentstructure
	tmpl, err := template.ParseFiles("templates/portfolio/basic.html.tmpl", "templates/portfolio/appl.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	tmpl.ExecuteTemplate(w, "layout", nil)
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
		Existing    []*models.ApplRow
	}{
		currentUser,
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
	db := r.Context().Value("db").(*sqlx.DB)
	err := r.ParseForm()
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	m := make(map[string]interface{})
	m["FirstName"] = r.FormValue("FirstName")
	m["LastName"] = r.FormValue("LastName")
	m["Email"] = r.FormValue("Email")
	m["Phone"] = r.FormValue("Phone")
	m["CompanyName"] = r.FormValue("CompanyName")
	m["Website"] = r.FormValue("Website")
	m["Title"] = r.FormValue("Title")
	m["Industries"] = r.FormValue("Industries")
	m["Locations"] = r.FormValue("Locations")
	m["Comments"] = r.FormValue("Comments")
	m["CapitalRaised"] = r.FormValue("CapitalRaised")
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
