package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/fatih/structs"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/kunapuli09/3linesweb/libhttp"
	"github.com/kunapuli09/3linesweb/models"
	"github.com/shopspring/decimal"
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
	var i models.Search
	w.Header().Set("Content-Type", "text/html")
	db := r.Context().Value("db").(*sqlx.DB)
	sessionStore := r.Context().Value("sessionStore").(sessions.Store)
	session, _ := sessionStore.Get(r, "3linesweb-session")
	currentUser, ok := session.Values["user"].(*models.UserRow)
	if !ok || !(currentUser.Admin || currentUser.Dsc) {
		http.Redirect(w, r, "/logout", 302)
		return
	}
	err := r.ParseForm()
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	// r.PostForm is a map of our POST form values
	decoder := schema.NewDecoder()
	decoder.RegisterConverter(time.Time{}, ConvertFormDate)
	err1 := decoder.Decode(&i, r.PostForm)
	if err1 != nil {
		libhttp.HandleErrorJson(w, err1)
		return
	}
	allreqs, err := models.NewAppl(db).Search(nil, i)
	//create session date for page rendering
	data := struct {
		CurrentUser *models.UserRow
		Count       int
		Existing    []*models.ApplRow
	}{
		currentUser,
		getCount(w, r, currentUser.Email),
		allreqs,
	}
	tmpl, err := template.ParseFiles("templates/portfolio/internal.html.tmpl", "templates/portfolio/fundingreqs.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	tmpl.ExecuteTemplate(w, "layout", data)
}

func FundingAppl(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	db := r.Context().Value("db").(*sqlx.DB)
	sessionStore := r.Context().Value("sessionStore").(sessions.Store)
	session, _ := sessionStore.Get(r, "3linesweb-session")
	currentUser, ok := session.Values["user"].(*models.UserRow)
	if !ok || !(currentUser.Admin || currentUser.Dsc) {
		http.Redirect(w, r, "/logout", 302)
		return
	}
	ID, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	appl, err := models.NewAppl(db).GetById(nil, ID)
	allNotes, err := models.NewScreeningNotes(db).AllScreeningNotesByApplicationId(nil, ID)
	//create session date for page rendering
	data := struct {
		CurrentUser *models.UserRow
		Count       int
		Existing    *models.ApplRow
		AllNotes    []*models.ScreeningNotesRow
	}{
		currentUser,
		getCount(w, r, currentUser.Email),
		appl,
		allNotes,
	}
	funcMap := template.FuncMap{
		"safeHTML": func(b string) template.HTML {
			return template.HTML(b)
		},
		"currencyFormat": func(currency decimal.Decimal) string {
			f, _ := currency.Float64()
			return ac.FormatMoney(f)
		},
		"screeningStatus": func(b string) string {
			if b == "ARCHIVE" {
				return b + "D"
			}
			if b == "FASTTRACK" {
				return "FAST TRACKED"
			}
			return b + "ED"
		},
		"scoreDescription": func(b int8) template.HTML {
			if b < 4 {
				return template.HTML(`<span style="color:#28a745">LOW</span>`)
			}
			if b > 3 && b < 8 {
				return template.HTML(`<span style="color:#ffc107">MEDIUM</span>`)
			}
			if b > 7 && b < 11 {
				return template.HTML(`<span style="color:#dc3545">HIGH</span>`)
			}
			return template.HTML("NOT REVIWED")
		},
	}
	tmpl, err := template.New("main").Funcs(funcMap).ParseFiles("templates/portfolio/internal.html.tmpl", "templates/portfolio/fundingappl.html.tmpl")
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
	decoder.RegisterConverter(sql.NullString{}, ConvertSQLNullString)
	err1 := decoder.Decode(&i, r.PostForm)
	fmt.Printf("Form %v", r.PostForm)
	if err1 != nil {
		fmt.Println("decoding error")
		libhttp.HandleErrorJson(w, err1)
		return
	}
	m := structs.Map(i)
	m["ApplicationDate"] = time.Now()
	m["Title"] = "Removed"
	m["LastUpdatedTime"] = time.Now()
	m["LastUpdatedBy"] = m["Email"].(string)
	m["Status"] = "Initial Review Pending By DSC"
	fmt.Printf("map %v", m)
	//check for duplicate entry
	phone, ok1 := m["Phone"].(string)
	email, ok2 := m["Email"].(string)
	website, ok3 := m["Website"].(string)
	companyname, ok4 := m["CompanyName"].(string)

	if ok1 {
		if len(phone) > 12 {
			err2 := errors.New("Maximum 12 Digits in a Phone Number including Country Code. No + Sign Reqired for International.")
			//fmt.Println(err5)
			libhttp.HandleErrorJson(w, err2)
			return
		}
	}
	if ok2 && ok3 && ok4 {
		exists := models.NewAppl(db).GetExisting(nil, email, website, companyname)
		if exists == true {
			err3 := errors.New("Duplicate Entry. An application with Website, CompanyName or Email Already Exists")
			//fmt.Println(err2)
			libhttp.HandleErrorJson(w, err3)
			return
		}
	}

	_, err4 := models.NewAppl(db).Create(nil, m)
	if err4 != nil {
		fmt.Println("Application Information is not Valid", err4)
		libhttp.HandleErrorJson(w, err4)
		return
	}
	http.Redirect(w, r, "/", 302)
}

//database call to add update
func UpdateApplication(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	var i models.ApplRow
	db := r.Context().Value("db").(*sqlx.DB)
	err := r.ParseForm()
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	sessionStore := r.Context().Value("sessionStore").(sessions.Store)
	session, _ := sessionStore.Get(r, "3linesweb-session")
	currentUser, ok := session.Values["user"].(*models.UserRow)
	if !ok {
		http.Redirect(w, r, "/logout", 302)
		return
	}
	ID, e := strconv.ParseInt(r.FormValue("id"), 10, 64)
	if e != nil {
		libhttp.HandleErrorJson(w, e)
		return
	}
	existing, e := models.NewAppl(db).GetById(nil, ID)
	if e != nil {
		fmt.Println("Failed to Get  Existing Application", e)
		libhttp.HandleErrorJson(w, e)
		return
	}
	decoder := schema.NewDecoder()
	decoder.RegisterConverter(sql.NullString{}, ConvertSQLNullString)
	decoder.RegisterConverter(time.Time{}, ConvertFormDate)
	err1 := decoder.Decode(&i, r.PostForm)
	if err1 != nil {
		fmt.Println("decoding error")
		libhttp.HandleErrorJson(w, err1)
		return
	}
	m := structs.Map(i)
	m["Title"] = existing.Title
	m["ApplicationDate"] = existing.ApplicationDate
	m["LastUpdatedTime"] = time.Now()
	m["LastUpdatedBy"] = currentUser.Email
	fmt.Printf("map %v", m)
	_, err4 := models.NewAppl(db).UpdateById(nil, ID, m)
	if err4 != nil {
		fmt.Println("Failed to Update Information", err4)
		libhttp.HandleErrorJson(w, err4)
		return
	}
	http.Redirect(w, r, "/fundingreqs", 302)
}

//db call to remove
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
