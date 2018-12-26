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
	"github.com/shopspring/decimal"
	"github.com/leekchan/accounting"
)

// type NewsHtmlRow struct {
// 	ID            int64
// 	Investment_ID int64
// 	NewsDate      time.Time
// 	News          template.HTML
// }

const PENDING = "PENDING"
const COMPLETE = "COMPLETE"
const ARCHIVE = "ARCHIVE"

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
	investments, err := models.NewInvestment(db).GetStartupNames(nil)
	archive, pending, complete := SplitByStatus(investments)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	//create session date for page rendering
	data := struct {
		CurrentUser *models.UserRow
		Count       int
		Investments []*models.InvestmentRow
		Pending     []*models.InvestmentRow
		Archive     []*models.InvestmentRow
	}{
		currentUser,
		getCount(w, r, currentUser.Email),
		complete,
		pending,
		archive,
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
	ac := accounting.Accounting{Symbol: "$", Precision: 2}
	ID, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	funcMap := template.FuncMap{
		"safeHTML": func(b string) template.HTML {
			return template.HTML(b)
		},
		"currencyFormat": func(currency decimal.Decimal) string {
			f, _ := currency.Float64()
			return ac.FormatMoney(f)
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

	//create session date for page rendering
	data := struct {
		CurrentUser                  *models.UserRow
		Count                        int
		Investment                   *models.InvestmentRow
		Existing                     []*models.FinancialResultsRow
		ExistingNews                 []*models.NewsRow
		ExistingCapitalStructures    []*models.CapitalizationStructure
		ExistingInvestmentStructures []*models.InvestmentStructureRow
		ExistingDocs                 []*models.DocRow
	}{
		currentUser,
		getCount(w, r, currentUser.Email),
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
		Count       int
		Investment  *models.InvestmentRow
	}{
		currentUser,
		getCount(w, r, currentUser.Email),
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
		Count       int
		Investment  *models.InvestmentRow
	}{
		currentUser,
		getCount(w, r, currentUser.Email),
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
func SplitByStatus(investments []*models.InvestmentRow) ([]*models.InvestmentRow, []*models.InvestmentRow, []*models.InvestmentRow) {
	var pending []*models.InvestmentRow
	var complete []*models.InvestmentRow
	var archive []*models.InvestmentRow

	for _, investment := range investments {
		switch status := investment.Status; status {
		case PENDING:
			pending = append(pending, investment)
		case COMPLETE:
			complete = append(complete, investment)
		case ARCHIVE:
			archive = append(archive, investment)
		default:
			fmt.Printf("%s. is unknown status", status)
		}
	}
	return archive, pending, complete
}
