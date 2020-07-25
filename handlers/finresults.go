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

func NewFinancialResults(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	db := r.Context().Value("db").(*sqlx.DB)
	sessionStore := r.Context().Value("sessionStore").(sessions.Store)
	session, _ := sessionStore.Get(r, "3linesweb-session")
	currentUser, ok := session.Values["user"].(*models.UserRow)
	if !ok || !currentUser.Admin {
		http.Redirect(w, r, "/logout", 302)
		return
	}
	Investment_ID, err := strconv.ParseInt(r.URL.Query().Get("Investment_ID"), 10, 64)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	investment, err := models.NewInvestment(db).GetById(nil, Investment_ID)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	AllFinancialResults, err := models.NewFinancialResults(db).GetAllByInvestmentId(nil, Investment_ID)

	//create empty investmentstructure
	FinancialResults := models.FinancialResultsRow{}
	//create session date for page rendering
	data := struct {
		CurrentUser      *models.UserRow
		Count            int
		Investment       *models.InvestmentRow
		FinancialResults models.FinancialResultsRow
		Existing         []*models.FinancialResultsRow
	}{
		currentUser,
		getCount(w, r, currentUser.Email),
		investment,
		FinancialResults,
		AllFinancialResults,
	}
	tmpl, err := template.ParseFiles("templates/portfolio/basic.html.tmpl", "templates/portfolio/newfinancials.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	tmpl.ExecuteTemplate(w, "layout", data)
}

//database call to add new
func AddFinancialResults(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	var i models.FinancialResultsRow
	db := r.Context().Value("db").(*sqlx.DB)
	sessionStore := r.Context().Value("sessionStore").(sessions.Store)
	session, _ := sessionStore.Get(r, "3linesweb-session")
	_, ok := session.Values["user"].(*models.UserRow)
	if !ok {
		http.Redirect(w, r, "/logout", 302)
		return
	}
	err := r.ParseForm()
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	decoder := schema.NewDecoder()
	decoder.RegisterConverter(time.Time{}, ConvertFormDate)
	err1 := decoder.Decode(&i, r.PostForm)
	if err1 != nil {
		fmt.Println("decoding error")
		libhttp.HandleErrorJson(w, err1)
		return
	}
	m := structs.Map(i)
	//fmt.Printf("map %v", m)
	_, err2 := models.NewFinancialResults(db).Create(nil, m)
	if err2 != nil {
		fmt.Println("database error")
		libhttp.HandleErrorJson(w, err2)
		return
	}
	address := fmt.Sprintf("/viewinvestment?id=%v", m["Investment_ID"])
	http.Redirect(w, r, address, 302)
}

//presentation edit view
func EditFinancialResults(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	db := r.Context().Value("db").(*sqlx.DB)
	Fin_ID, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	Investment_ID, err := strconv.ParseInt(r.URL.Query().Get("Investment_ID"), 10, 64)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	sessionStore := r.Context().Value("sessionStore").(sessions.Store)
	session, _ := sessionStore.Get(r, "3linesweb-session")
	currentUser, ok := session.Values["user"].(*models.UserRow)
	if !ok || !currentUser.Admin {
		http.Redirect(w, r, "/logout", 302)
		return
	}
	investment, err := models.NewInvestment(db).GetById(nil, Investment_ID)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	finresults, err := models.NewFinancialResults(db).GetById(nil, Fin_ID)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	//create session data for page rendering
	data := struct {
		CurrentUser      *models.UserRow
		Count            int
		Investment       *models.InvestmentRow
		FinancialResults *models.FinancialResultsRow
	}{
		currentUser,
		getCount(w, r, currentUser.Email),
		investment,
		finresults,
	}
	funcMap := template.FuncMap{
		"safeHTML": func(b string) template.HTML {
			return template.HTML(b)
		},
	}
	tmpl, err := template.New("main").Funcs(funcMap).ParseFiles("templates/portfolio/basic.html.tmpl", "templates/portfolio/editfinancialresults.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	tmpl.ExecuteTemplate(w, "layout", data)
}

//db call to update
func UpdateFinancialResults(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	db := r.Context().Value("db").(*sqlx.DB)
	var i models.FinancialResultsRow
	err := r.ParseForm()
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	Fin_ID, e := strconv.ParseInt(r.FormValue("id"), 10, 64)
	if e != nil {
		libhttp.HandleErrorJson(w, e)
		return
	}
	Investment_ID, e := strconv.ParseInt(r.FormValue("Investment_ID"), 10, 64)
	if e != nil {
		libhttp.HandleErrorJson(w, e)
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
	m := structs.Map(i)
	//fmt.Printf("map %v", m)
	_, err2 := models.NewFinancialResults(db).UpdateById(nil, Fin_ID, m)
	if err2 != nil {
		libhttp.HandleErrorJson(w, err2)
		return
	}
	address := fmt.Sprintf("/newfinancials?Investment_ID=%v", Investment_ID)
	http.Redirect(w, r, address, 302)
}

//db call to update
func RemoveFinancialResults(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	db := r.Context().Value("db").(*sqlx.DB)
	ID, e := strconv.ParseInt(r.FormValue("id"), 10, 64)
	if e != nil {
		libhttp.HandleErrorJson(w, e)
		return
	}
	Investment_ID, e1 := strconv.ParseInt(r.FormValue("Investment_ID"), 10, 64)
	if e1 != nil {
		libhttp.HandleErrorJson(w, e1)
		return
	}
	_, err2 := models.NewFinancialResults(db).DeleteByID(nil, ID)
	if err2 != nil {
		libhttp.HandleErrorJson(w, err2)
		return
	}
	address := fmt.Sprintf("/newfinancials?Investment_ID=%v", Investment_ID)
	http.Redirect(w, r, address, 302)
}
