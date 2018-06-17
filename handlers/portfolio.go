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

// type NewsHtmlRow struct {
// 	ID            int64
// 	Investment_ID int64
// 	NewsDate      time.Time
// 	News          template.HTML
// }

const PENDING_STATUS = "PENDING"

func GetPortfolio(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	db := r.Context().Value("db").(*sqlx.DB)
	sessionStore := r.Context().Value("sessionStore").(sessions.Store)
	session, _ := sessionStore.Get(r, "3linesweb-session")
	currentUser, ok := session.Values["user"].(*models.UserRow)
	if !ok {
		http.Redirect(w, r, "/logout", 302)
		return
	}
	var pending []*models.InvestmentRow
	var complete []*models.InvestmentRow
	investments, err := models.NewInvestment(db).GetStartupNames(nil)
	//TODO loop through investments for status
	//if not admin, then show complete investments
	//else, show all investments
	for _, v := range investments {
			if v.Status == PENDING_STATUS {
				pending = append(pending, v)
			}else{
				complete = append(complete, v)
			}
	}

	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	//create session date for page rendering
	data := struct {
		CurrentUser *models.UserRow
		Investments []*models.InvestmentRow
		Pending []*models.InvestmentRow
	}{
		currentUser,
		complete,
		pending,
	}
	funcMap := template.FuncMap{
		"safeHTML": func(b string) template.HTML {
			return template.HTML(b)
		},
	}
	tmpl, err := template.New("main").Funcs(funcMap).ParseFiles("templates/portfolio/basic.html.tmpl", "templates/portfolio/portfolio.html.tmpl")

	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	tmpl.ExecuteTemplate(w, "layout", data)

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
	funcMap := template.FuncMap{
		"safeHTML": func(b string) template.HTML {
			return template.HTML(b)
		},
	}
	tmpl, e := template.New("main").Funcs(funcMap).ParseFiles("templates/portfolio/viewinvestment.html.tmpl", "templates/portfolio/basic.html.tmpl")
	// tmpl, e := template.ParseFiles("templates/portfolio/viewinvestment.html.tmpl", "templates/portfolio/basic.html.tmpl")
	if e != nil {
		libhttp.HandleErrorJson(w, e)
		return
	}
	sessionStore := r.Context().Value("sessionStore").(sessions.Store)
	session, _ := sessionStore.Get(r, "3linesweb-session")
	currentUser, ok := session.Values["user"].(*models.UserRow)
	if !ok {
		http.Redirect(w, r, "/logout", 302)
		return
	}
	investment, err := models.NewInvestment(db).GetById(nil, ID)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	//TODO do a big join query
	AllFinancialResults, err := models.NewFinancialResults(db).GetAllByInvestmentId(nil, ID)
	AllNews, err := models.NewNews(db).GetAllByInvestmentId(nil, ID)
	AllCapitalStructures, err := models.NewCapitalStructure(db).GetAllByInvestmentId(nil, ID)
	AllInvestmentStructures, err := models.NewInvestmentStructure(db).GetAllByInvestmentId(nil, ID)
	AllDocs, err := models.NewDoc(db).GetAllByInvestmentId(nil, ID)

	// AllHtmlNews := make([]*NewsHtmlRow, len(AllNews))
	// for i, element := range AllNews {
	// 	AllHtmlNews[i] = &NewsHtmlRow{
	//     		ID: element.ID,
	//     		Investment_ID: element.Investment_ID,
	//     		NewsDate: element.NewsDate,
	//     		News: template.HTML(element.News),
	//     	}
	// }
	// for _, element := range AllHtmlNews {
	// 	fmt.Println(element.News)
	// }

	//create session date for page rendering
	data := struct {
		CurrentUser                  *models.UserRow
		Investment                   *models.InvestmentRow
		Existing                     []*models.FinancialResultsRow
		ExistingNews                 []*models.NewsRow
		ExistingCapitalStructures    []*models.CapitalizationStructure
		ExistingInvestmentStructures []*models.InvestmentStructureRow
		ExistingDocs                 []*models.DocRow
	}{
		currentUser,
		investment,
		AllFinancialResults,
		AllNews,
		AllCapitalStructures,
		AllInvestmentStructures,
		AllDocs,
	}
	tmpl.ExecuteTemplate(w, "layout", data)
}

//presentation view for new investment
func NewInvestment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	sessionStore := r.Context().Value("sessionStore").(sessions.Store)
	session, _ := sessionStore.Get(r, "3linesweb-session")
	currentUser, ok := session.Values["user"].(*models.UserRow)
	if !ok {
		http.Redirect(w, r, "/logout", 302)
		return
	}

	investment := &models.InvestmentRow{}
	//create session data for page rendering
	data := struct {
		CurrentUser *models.UserRow
		Investment  *models.InvestmentRow
	}{
		currentUser,
		investment,
	}
	tmpl, err := template.ParseFiles("templates/portfolio/newinvestment.html.tmpl", "templates/portfolio/basic.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	tmpl.ExecuteTemplate(w, "layout", data)
}

//database call to add new
func Add(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	var i models.InvestmentRow
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
	_, err2 := models.NewInvestment(db).Create(nil, m)
	if err2 != nil {
		fmt.Println("database error")
		libhttp.HandleErrorJson(w, err2)
		return
	}
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
	sessionStore := r.Context().Value("sessionStore").(sessions.Store)
	session, _ := sessionStore.Get(r, "3linesweb-session")
	currentUser, ok := session.Values["user"].(*models.UserRow)
	if !ok {
		http.Redirect(w, r, "/logout", 302)
		return
	}

	investment, err := models.NewInvestment(db).GetById(nil, ID)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	//create session data for page rendering
	data := struct {
		CurrentUser *models.UserRow
		Investment  *models.InvestmentRow
	}{
		currentUser,
		investment,
	}
	funcMap := template.FuncMap{
		"safeHTML": func(b string) template.HTML {
			return template.HTML(b)
		},
	}
	tmpl, err := template.New("main").Funcs(funcMap).ParseFiles("templates/portfolio/basic.html.tmpl", "templates/portfolio/editinvestment.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	tmpl.ExecuteTemplate(w, "layout", data)

}

//db call to update
func Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	db := r.Context().Value("db").(*sqlx.DB)
	var i models.InvestmentRow
	err := r.ParseForm()
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	ID, e := strconv.ParseInt(r.FormValue("id"), 10, 64)
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
	_, err2 := models.NewInvestment(db).UpdateById(nil, ID, m)
	if err2 != nil {
		libhttp.HandleErrorJson(w, err2)
		return
	}
	address := fmt.Sprintf("/viewinvestment?id=%v", ID)
	http.Redirect(w, r, address, 302)
}
