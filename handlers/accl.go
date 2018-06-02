package handlers

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/kunapuli09/3linesweb/libhttp"
	"github.com/kunapuli09/3linesweb/models"
	"net/http"
	"strconv"
)

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
	m["State"] = r.FormValue("State")
	m["Industries"] = r.FormValue("Industries")
	m["Locations"] = r.FormValue("Locations")
	m["Comments"] = r.FormValue("Comments")
	m["CapitalRaised"] = r.FormValue("CapitalRaised")
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
