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
	"sort"
	"strconv"
	"time"
)

func News(w http.ResponseWriter, r *http.Request) {
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
	if !ok || !currentUser.Admin {
		http.Redirect(w, r, "/logout", 302)
		return
	}
	investment, err := models.NewInvestment(db).GetById(nil, Investment_ID)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	allnews, err := models.NewNews(db).GetPendingByInvestmentId(nil, Investment_ID)
	//create empty investmentstructure
	news := models.NewsRow{}
	//create session date for page rendering
	data := struct {
		CurrentUser *models.UserRow
		Count       int
		Investment  *models.InvestmentRow
		News        models.NewsRow
		Existing    []*models.NewsRow
	}{
		currentUser,
		getCount(w, r, currentUser.Email),
		investment,
		news,
		allnews,
	}
	tmpl, err := template.ParseFiles("templates/portfolio/basic.html.tmpl", "templates/portfolio/newnews.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	tmpl.ExecuteTemplate(w, "layout", data)
}

//database call to add new
func AddNews(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	var i models.NewsRow
	db := r.Context().Value("db").(*sqlx.DB)
	sessionStore := r.Context().Value("sessionStore").(sessions.Store)
	session, _ := sessionStore.Get(r, "3linesweb-session")
	currentUser, ok := session.Values["user"].(*models.UserRow)
	if !ok || !currentUser.Admin{
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
	_, err2 := models.NewNews(db).Create(nil, m)
	if err2 != nil {
		fmt.Println("database error")
		libhttp.HandleErrorJson(w, err2)
		return
	}
	address := fmt.Sprintf("/viewinvestment?id=%v", m["Investment_ID"])
	http.Redirect(w, r, address, 302)
}

//db call to update
func RemoveNews(w http.ResponseWriter, r *http.Request) {
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
	_, err2 := models.NewNews(db).DeleteByID(nil, ID)
	if err2 != nil {
		libhttp.HandleErrorJson(w, err2)
		return
	}
	address := fmt.Sprintf("/news?Investment_ID=%v", Investment_ID)
	http.Redirect(w, r, address, 302)
}

//presentation edit view
func EditNews(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	db := r.Context().Value("db").(*sqlx.DB)
	News_ID, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
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
	if !ok || !currentUser.Admin{
		http.Redirect(w, r, "/logout", 302)
		return
	}
	investment, err := models.NewInvestment(db).GetById(nil, Investment_ID)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	news, err := models.NewNews(db).GetById(nil, News_ID)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	//create session data for page rendering
	data := struct {
		CurrentUser *models.UserRow
		Count       int
		Investment  *models.InvestmentRow
		News        *models.NewsRow
	}{
		currentUser,
		getCount(w, r, currentUser.Email),
		investment,
		news,
	}
	funcMap := template.FuncMap{
		"safeHTML": func(b string) template.HTML {
			return template.HTML(b)
		},
	}
	tmpl, err := template.New("main").Funcs(funcMap).ParseFiles("templates/portfolio/basic.html.tmpl", "templates/portfolio/editnews.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	tmpl.ExecuteTemplate(w, "layout", data)
}

//db call to update
func UpdateNews(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	db := r.Context().Value("db").(*sqlx.DB)
	var i models.NewsRow
	err := r.ParseForm()
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	News_ID, e := strconv.ParseInt(r.FormValue("id"), 10, 64)
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
	_, err2 := models.NewNews(db).UpdateById(nil, News_ID, m)
	if err2 != nil {
		libhttp.HandleErrorJson(w, err2)
		return
	}
	address := fmt.Sprintf("/news?Investment_ID=%v", Investment_ID)
	http.Redirect(w, r, address, 302)
}

func Notifications(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	db := r.Context().Value("db").(*sqlx.DB)
	funcMap := template.FuncMap{
		"safeHTML": func(b string) template.HTML {
			return template.HTML(b)
		},
	}
	tmpl, e := template.New("main").Funcs(funcMap).ParseFiles("templates/portfolio/basic.html.tmpl", "templates/portfolio/notifications.html.tmpl")
	if e != nil {
		libhttp.HandleErrorJson(w, e)
		return
	}
	sessionStore := r.Context().Value("sessionStore").(sessions.Store)
	session, _ := sessionStore.Get(r, "3linesweb-session")
	currentUser, ok := session.Values["user"].(*models.UserRow)
	if !ok || !(currentUser.Admin || currentUser.FundOne || currentUser.FundTwo || currentUser.Dsc){
		http.Redirect(w, r, "/logout", 302)
		return
	}
	allnotifications, err := models.NewNotification(db).AllNotifications(nil, currentUser.Email)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	//sort by date descending
	sort.Slice(allnotifications, func(i, j int) bool {
		return allnotifications[i].NewsDate.After(allnotifications[j].NewsDate)
	})

	//create session date for page rendering
	data := struct {
		CurrentUser *models.UserRow
		Count       int
		Existing    []*models.NotificationRow
	}{
		currentUser,
		getCount(w, r, currentUser.Email),
		allnotifications,
	}
	tmpl.ExecuteTemplate(w, "layout", data)
}

//db call to publish
//create a notifiction record for each user
func PublishNotification(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	db := r.Context().Value("db").(*sqlx.DB)
	err := r.ParseForm()
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	News_ID, e := strconv.ParseInt(r.FormValue("id"), 10, 64)
	if e != nil {
		libhttp.HandleErrorJson(w, e)
		return
	}
	Investment_ID, e := strconv.ParseInt(r.FormValue("Investment_ID"), 10, 64)
	if e != nil {
		libhttp.HandleErrorJson(w, e)
		return
	}
	news, err2 := models.NewNews(db).UpdateStatusById(nil, News_ID)
	if err2 != nil {
		libhttp.HandleErrorJson(w, err2)
		return
	}
	investment, err3 := models.NewInvestment(db).GetById(nil, Investment_ID)
	if err3 != nil {
		libhttp.HandleErrorJson(w, err3)
		return
	}
	//get all emails from database
	emails, err3 := models.NewUser(db).AllEmails(nil)
	if err3 != nil {
		libhttp.HandleErrorJson(w, err3)
		return
	}

	//option 1: send  an email
	//NotifyUsers(w, r, emails, news.News, news.Title)

	//option 2: create a notification record for each user with UNREAD status
	_, err4 := models.NewNotification(db).BatchPublish(nil, emails, investment.StartupName, news)
	if err4 != nil {
		libhttp.HandleErrorJson(w, err3)
		return
	}
	address := fmt.Sprintf("/news?Investment_ID=%v", Investment_ID)
	http.Redirect(w, r, address, 302)
}

//db call to update
//update status of notification to READ
func UpdateNotification(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	db := r.Context().Value("db").(*sqlx.DB)
	sessionStore := r.Context().Value("sessionStore").(sessions.Store)
	session, _ := sessionStore.Get(r, "3linesweb-session")
	currentUser, ok := session.Values["user"].(*models.UserRow)
	if !ok || !currentUser.Admin{
		http.Redirect(w, r, "/logout", 302)
		return
	}
	NotificationId, e := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if e != nil {
		libhttp.HandleErrorJson(w, e)
		return
	}
	notification, err2 := models.NewNotification(db).UpdateStatusById(nil, NotificationId)
	if err2 != nil {
		fmt.Println("database error")
		libhttp.HandleErrorJson(w, err2)
		return
	}
	news, err2 := models.NewNews(db).GetById(nil, notification.News_ID)
	if err2 != nil {
		libhttp.HandleErrorJson(w, err2)
		return
	}

	//create session data for page rendering
	data := struct {
		CurrentUser *models.UserRow
		Count       int
		News        *models.NewsRow
	}{
		currentUser,
		getCount(w, r, currentUser.Email),
		news,
	}
	funcMap := template.FuncMap{
		"safeHTML": func(b string) template.HTML {
			return template.HTML(b)
		},
	}
	tmpl, e := template.New("main").Funcs(funcMap).ParseFiles("templates/portfolio/basic.html.tmpl", "templates/portfolio/displaynews.html.tmpl")
	if e != nil {
		libhttp.HandleErrorJson(w, e)
		return
	}
	tmpl.ExecuteTemplate(w, "layout", data)

}
