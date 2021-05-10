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

func ScreeningNotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	db := r.Context().Value("db").(*sqlx.DB)

	Application_ID, err := strconv.ParseInt(r.URL.Query().Get("Application_ID"), 10, 64)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	//fmt.Printf("ApplicationID%v", Application_ID)
	ScreeningNotes_ID, err := strconv.ParseInt(r.URL.Query().Get("ScreeningNotes_ID"), 10, 64)
	if err != nil {
		ScreeningNotes_ID = 0
	}
	sessionStore := r.Context().Value("sessionStore").(sessions.Store)
	session, _ := sessionStore.Get(r, "3linesweb-session")
	currentUser, ok := session.Values["user"].(*models.UserRow)
	if !ok || !(currentUser.InvestorRelations || currentUser.Admin || currentUser.Dsc) {
		http.Redirect(w, r, "/logout", 302)
		return
	}
	Application, err := models.NewAppl(db).GetById(nil, Application_ID)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	//fmt.Printf("ScreeningNotes_ID%v", ScreeningNotes_ID)
	screeningNotes, err := models.NewScreeningNotes(db).GetByApplicationIdAndScreener(nil, ScreeningNotes_ID, Application_ID, currentUser.Email)
	if err != nil {
		fmt.Printf("DB Error %v", err)
	}
	//fmt.Printf("screeningnotes%v", screeningNotes)
	//create session date for page rendering
	data := struct {
		CurrentUser    *models.UserRow
		Count          int
		Application    *models.ApplRow
		ScreeningNotes *models.ScreeningNotesRow
	}{
		currentUser,
		getCount(w, r, currentUser.Email),
		Application,
		screeningNotes,
	}
	funcMap := template.FuncMap{
		"safeHTML": func(b string) template.HTML {
			return template.HTML(b)
		},
	}
	tmpl, err := template.New("main").Funcs(funcMap).ParseFiles("templates/portfolio/basic.html.tmpl", "templates/portfolio/editscreeningnotes.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	tmpl.ExecuteTemplate(w, "layout", data)
}

//db call to update
func UpdateScreeningNotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	db := r.Context().Value("db").(*sqlx.DB)
	var i models.ScreeningNotesRow

	err := r.ParseForm()
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	ScreeningNotes_ID, e := strconv.ParseInt(r.FormValue("id"), 10, 64)
	if e != nil {
		libhttp.HandleErrorJson(w, e)
		return
	}
	Application_ID, e := strconv.ParseInt(r.FormValue("Application_ID"), 10, 64)
	if e != nil {
		libhttp.HandleErrorJson(w, e)
		return
	}
	// ScreenerEmail := r.FormValue("ScreenerEmail")
	// if e != nil {
	// 	libhttp.HandleErrorJson(w, e)
	// 	return
	// }
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
	if ScreeningNotes_ID == 0 {
		//fmt.Printf("Creating New Notes with ApplicationID%v, ScreenerEmail%v", Application_ID, ScreenerEmail)
		notes, err2 := models.NewScreeningNotes(db).Create(nil, m)
		if err2 != nil {
			libhttp.HandleErrorJson(w, err2)
			return
		}
		ScreeningNotes_ID = notes.ID

	} else {
		//fmt.Printf("Updating Notes with ApplicationID%v, ScreeningNotes_ID%v", Application_ID, ScreeningNotes_ID)

		_, err3 := models.NewScreeningNotes(db).UpdateById(nil, ScreeningNotes_ID, m)
		if err3 != nil {
			libhttp.HandleErrorJson(w, err3)
			return
		}
	}
	address := fmt.Sprintf("/fundingappl?id=%v", Application_ID)
	http.Redirect(w, r, address, 302)
}

//db call to update
func RemoveScreeningNotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	db := r.Context().Value("db").(*sqlx.DB)
	ID, e := strconv.ParseInt(r.FormValue("id"), 10, 64)
	if e != nil {
		libhttp.HandleErrorJson(w, e)
		return
	}
	Application_ID, e1 := strconv.ParseInt(r.FormValue("Application_ID"), 10, 64)
	if e1 != nil {
		libhttp.HandleErrorJson(w, e1)
		return
	}
	_, err2 := models.NewScreeningNotes(db).DeleteByID(nil, ID)
	if err2 != nil {
		libhttp.HandleErrorJson(w, err2)
		return
	}
	address := fmt.Sprintf("/ScreeningNotes?Application_ID=%v", Application_ID)
	http.Redirect(w, r, address, 302)
}
