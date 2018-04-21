package handlers

import (
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/gorilla/schema"
	"github.com/jmoiron/sqlx"
	"github.com/kunapuli09/3linesweb/libhttp"
	"github.com/kunapuli09/3linesweb/models"
	"html/template"
	"net/http"
	"strconv"
	"github.com/fatih/structs"
	"reflect"
  	"time"
)
var decoder = schema.NewDecoder()

func ConvertFormDate(value string) reflect.Value {
	s, _ := time.Parse("2006-Jan-02", value)
	return reflect.ValueOf(s)
}

func GetPortfolio(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	sessionStore := r.Context().Value("sessionStore").(sessions.Store)
	session, _ := sessionStore.Get(r, "3linesweb-session")
	currentUser, ok := session.Values["user"].(*models.UserRow)
	if !ok {
		http.Redirect(w, r, "/logout", 302)
		return
	}

	data := struct {
		CurrentUser *models.UserRow
	}{
		currentUser,
	}

	tmpl, err := template.ParseFiles("templates/portfolio/portfolio.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	tmpl.Execute(w, data)

	//tmpl.ExecuteTemplate(w, "layout", data)
}

//presentation view for new investment
func ViewInvestment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	db := r.Context().Value("db").(*sqlx.DB)
	ID, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	tmpl, e := template.ParseFiles("templates/portfolio/viewinvestment.html.tmpl","templates/portfolio/basic.html.tmpl" )
	if e != nil {
		libhttp.HandleErrorJson(w, e)
		return
	}
	investment, err := models.NewInvestment(db).GetById(nil, ID)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	fmt.Printf(" Retrieved id %v name %s industry %s", investment.ID, investment.StartupName, investment.Industry)
	tmpl.ExecuteTemplate(w, "View", investment)
}
//presentation view for new investment
func NewInvestment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	tmpl, err := template.ParseFiles("templates/portfolio/newinvestment.html.tmpl", "templates/portfolio/basic.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	tmpl.ExecuteTemplate(w, "layout", nil)
}
//database call to add new
func Add(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	var i models.InvestmentRow
	db := r.Context().Value("db").(*sqlx.DB)
	err := r.ParseForm()
    if err != nil {
        // Handle error
    }
    // r.PostForm is a map of our POST form values
    decoder.RegisterConverter(time.Time{}, ConvertFormDate)
    err1 := decoder.Decode(&i, r.PostForm)
    if err1 != nil {
        libhttp.HandleErrorJson(w, err1)
		return
    }
    m := structs.Map(i)
    fmt.Printf("map %v", m)
    investment, err2 := models.NewInvestment(db).Create(nil, m)
	if err2 != nil {
		libhttp.HandleErrorJson(w, err2)
		return
	}

	fmt.Printf(" id %v name %s industry %s", investment.ID, investment.StartupName, investment.Industry)
	GetPortfolio(w, r)
}
//presentation edit view
func EditInvestment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	db := r.Context().Value("db").(*sqlx.DB)
	ID, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	tmpl, e := template.ParseFiles("templates/portfolio/editinvestment.html.tmpl","templates/portfolio/basic.html.tmpl" )
	if e != nil {
		libhttp.HandleErrorJson(w, e)
		return
	}
	investment, err := models.NewInvestment(db).GetById(nil, ID)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	fmt.Printf(" Retrieved id %v name %s industry %s", investment.ID, investment.StartupName, investment.Industry)
	tmpl.ExecuteTemplate(w, "Edit", investment)
}
//db call to update
func Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	db := r.Context().Value("db").(*sqlx.DB)
	var investment *models.InvestmentRow
	var err error
	name := r.FormValue("StartupName")
	industry := r.FormValue("Industry")
	ID, err := strconv.ParseInt(r.FormValue("id"), 10, 64)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	fmt.Printf(" id %v" , ID)
	investment, err = models.NewInvestment(db).UpdateById(nil, ID, name, industry)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	fmt.Printf(" Updated id %v name %s industry %s", investment.ID, investment.StartupName, investment.Industry)
	GetPortfolio(w, r)
}
