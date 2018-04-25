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

func NewInvestmentStructure(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	db := r.Context().Value("db").(*sqlx.DB)
	Investment_ID, err := strconv.ParseInt(r.URL.Query().Get("Investment_ID"), 10, 64)
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
	investment, err := models.NewInvestment(db).GetById(nil, Investment_ID)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	Investmentstructures, err := models.NewInvestmentStructure(db).GetAllByInvestmentId(nil, Investment_ID)
	//create empty investmentstructure
	Investmentstructure := models.InvestmentStructureRow{}
	//create session date for page rendering
	data := struct {
		CurrentUser         *models.UserRow
		Investment          *models.InvestmentRow
		InvestmentStructure models.InvestmentStructureRow
		Existing            []*models.InvestmentStructureRow
	}{
		currentUser,
		investment,
		Investmentstructure,
		Investmentstructures,
	}
	tmpl, err := template.ParseFiles("templates/portfolio/basic.html.tmpl", "templates/portfolio/newinvestmentstructure.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	tmpl.ExecuteTemplate(w, "layout", data)
}

//database call to add new
func AddInvestmentStructure(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	var i models.InvestmentStructureRow
	db := r.Context().Value("db").(*sqlx.DB)
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
	_, err2 := models.NewInvestmentStructure(db).Create(nil, m)
	if err2 != nil {
		fmt.Println("database error")
		libhttp.HandleErrorJson(w, err2)
		return
	}
	address := fmt.Sprintf("/viewinvestment?id=%v", m["Investment_ID"])
	http.Redirect(w, r, address, 302)
}

//db call to update
func RemoveInvestmentStructure(w http.ResponseWriter, r *http.Request) {
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
	_, err2 := models.NewInvestmentStructure(db).DeleteByID(nil, ID)
	if err2 != nil {
		libhttp.HandleErrorJson(w, err2)
		return
	}
	address := fmt.Sprintf("/newinvestmentstructure?Investment_ID=%v", Investment_ID)
	http.Redirect(w, r, address, 302)
}
