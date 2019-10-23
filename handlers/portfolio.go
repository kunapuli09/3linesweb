package handlers

import (
	"fmt"
	"github.com/fatih/structs"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/kunapuli09/3linesweb/libhttp"
	"github.com/kunapuli09/3linesweb/models"
	"github.com/leekchan/accounting"
	"github.com/shopspring/decimal"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

const FUNDI = "3Lines 2016 Discretionary Fund, LLC"
const FUNDII = "3Lines Rocket Fund, L.P"
const PERFORMANCE_BELOW_TARGET = "PERFORMANCE_BELOW_TARGET"
const PERFORMANCE_ON_TARGET = "PERFORMANCE_ON_TARGET"
const PERFORMANCE_ABOVE_TARGET = "PERFORMANCE_ABOVE_TARGET"
const POOR_PERFORMANCE = "POOR_PERFORMANCE"
const ACQUIRED_OR_SOLD = "ACQUIRED_OR_SOLD"

var ac = accounting.Accounting{Symbol: "$", Precision: 2}

func EntryAccess(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	sessionStore := r.Context().Value("sessionStore").(sessions.Store)
	session, _ := sessionStore.Get(r, "3linesweb-session")
	currentUser, ok := session.Values["user"].(*models.UserRow)

	if !ok {
		http.Redirect(w, r, "/logout", 302)
		return
	}
	fmt.Printf("%s. Inside EntryAccess", currentUser.Email)
	
    switch defaultView := true; defaultView {
		case currentUser.Admin:
			GetAdminDashboard(w, r)
		case currentUser.FundTwo:
			//fmt.Printf("%s. privileges are FundTwo", currentUser.Email)
			Fund2Dashboard(w,r)
		case currentUser.FundOne:
			//fmt.Printf("%s. privileges are FundOne", currentUser.Email)
			Fund1Dashboard(w,r)
		default:
			//fmt.Printf("%s. privileges are unknown", currentUser.Email)
			NoEntry(w,r)
		}

}

func GetAdminDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	db := r.Context().Value("db").(*sqlx.DB)
	sessionStore := r.Context().Value("sessionStore").(sessions.Store)
	session, _ := sessionStore.Get(r, "3linesweb-session")
	currentUser, ok := session.Values["user"].(*models.UserRow)
	if !ok || !currentUser.Admin{
		http.Redirect(w, r, "/logout", 302)
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
	investments, err := models.NewInvestment(db).GetStartupNames(nil)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	startupnames, amounts := InvestmentSpreadData(investments)

	//data for fund1dashboard.html.tmpl
	data := struct {
		CurrentUser          *models.UserRow
		Count                int
		Investments     	[]*models.InvestmentRow
		StartupNames         []string
		Amounts              []decimal.Decimal
	}{
		currentUser,
		getCount(w, r, currentUser.Email),
		investments,
		startupnames,
		amounts,
	}

	tmpl, err := template.New("main").Funcs(funcMap).ParseFiles("templates/portfolio/basic.html.tmpl", "templates/portfolio/admindashboard.html.tmpl")

	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	tmpl.ExecuteTemplate(w, "layout", data)
}
func Fund1Dashboard(w http.ResponseWriter, r *http.Request) {
	var isFundII bool
	w.Header().Set("Content-Type", "text/html")
	db := r.Context().Value("db").(*sqlx.DB)
	sessionStore := r.Context().Value("sessionStore").(sessions.Store)
	session, _ := sessionStore.Get(r, "3linesweb-session")
	currentUser, ok := session.Values["user"].(*models.UserRow)
	if !ok {
		http.Redirect(w, r, "/logout", 302)
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
	contributions, err := models.NewContribution(db).GetAllByFundNameAndUserId(nil, "", currentUser.ID)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	if len(contributions) > 1 {
		isFundII = true
	}
	fundonecontribution, _ := SplitContributionsByFund(contributions)
	investments, err := models.NewInvestment(db).GetStartupNames(nil)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	fundone, _ := SplitByFund(investments)
	fund1startupnames, fund1amounts := InvestmentSpreadData(fundone)

	//data for fund1dashboard.html.tmpl
	data := struct {
		CurrentUser          *models.UserRow
		Count                int
		FundIInvestments     []*models.InvestmentRow
		FundIContribution    *models.ContributionRow
		StartupNames         []string
		Amounts              []decimal.Decimal
		FundIIInvestorAswell bool
	}{
		currentUser,
		getCount(w, r, currentUser.Email),
		fundone,
		fundonecontribution[0],
		fund1startupnames,
		fund1amounts,
		isFundII,
	}

	tmpl, err := template.New("main").Funcs(funcMap).ParseFiles("templates/portfolio/basic.html.tmpl", "templates/portfolio/fund1dashboard.html.tmpl")

	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	tmpl.ExecuteTemplate(w, "layout", data)
}

func Fund2Dashboard(w http.ResponseWriter, r *http.Request) {
	var isFundI bool
	w.Header().Set("Content-Type", "text/html")
	db := r.Context().Value("db").(*sqlx.DB)
	sessionStore := r.Context().Value("sessionStore").(sessions.Store)
	session, _ := sessionStore.Get(r, "3linesweb-session")
	currentUser, ok := session.Values["user"].(*models.UserRow)
	if !ok {
		http.Redirect(w, r, "/logout", 302)
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
	contributions, err := models.NewContribution(db).GetAllByFundNameAndUserId(nil, "", currentUser.ID)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	if len(contributions) > 1 {
		isFundI = true
	}
	_, fundtwocontribution := SplitContributionsByFund(contributions)
	investments, err := models.NewInvestment(db).GetStartupNames(nil)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	_, fundtwo := SplitByFund(investments)
	fund2startupnames, fund2amounts := InvestmentSpreadData(fundtwo)

	//data for fund2dashboard.html.tmpl
	data := struct {
		CurrentUser         *models.UserRow
		Count               int
		FundIIInvestments   []*models.InvestmentRow
		FundIIContribution  *models.ContributionRow
		StartupNames        []string
		Amounts             []decimal.Decimal
		FundIInvestorAsWell bool
	}{
		currentUser,
		getCount(w, r, currentUser.Email),
		fundtwo,
		fundtwocontribution[0],
		fund2startupnames,
		fund2amounts,
		isFundI,
	}

	tmpl, err := template.New("main").Funcs(funcMap).ParseFiles("templates/portfolio/basic.html.tmpl", "templates/portfolio/fund2dashboard.html.tmpl")

	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	tmpl.ExecuteTemplate(w, "layout", data)
}

func NoEntry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	sessionStore := r.Context().Value("sessionStore").(sessions.Store)
	session, _ := sessionStore.Get(r, "3linesweb-session")
	currentUser, ok := session.Values["user"].(*models.UserRow)
	if !ok {
		http.Redirect(w, r, "/logout", 302)
		return
	}
	data := struct {
		CurrentUser         *models.UserRow
		Count               int
	}{
		currentUser,
		getCount(w, r, currentUser.Email),
	}
	tmpl, err := template.New("main").ParseFiles("templates/portfolio/basic.html.tmpl", "templates/portfolio/noentry.html.tmpl")

	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	tmpl.ExecuteTemplate(w, "layout", data)
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
		"currencyFormat": func(currency decimal.Decimal) string {
			f, _ := currency.Float64()
			return ac.FormatMoney(f)
		},
	}
	tmpl, e := template.New("main").Funcs(funcMap).ParseFiles("templates/portfolio/viewinvestment.html.tmpl", "templates/portfolio/internal.html.tmpl")
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
	if ID == 0 {
		//fmt.Printf("Creating New Notes with ApplicationID%v, ScreenerEmail%v", Application_ID, ScreenerEmail)
		investment, err2 := models.NewInvestment(db).Create(nil, m)
		if err2 != nil {
			libhttp.HandleErrorJson(w, err2)
			return
		}
		ID = investment.ID

	} else {
		//fmt.Printf("Updating Notes with ApplicationID%v, ScreeningNotes_ID%v", Application_ID, ScreeningNotes_ID)

		_, err3 := models.NewInvestment(db).UpdateById(nil, ID, m)
		if err3 != nil {
			libhttp.HandleErrorJson(w, err3)
			return
		}
	}
	address := fmt.Sprintf("/viewinvestment?id=%v", ID)
	http.Redirect(w, r, address, 302)
}

func SplitByFund(investments []*models.InvestmentRow) ([]*models.InvestmentRow, []*models.InvestmentRow) {
	var fundone []*models.InvestmentRow
	var fundtwo []*models.InvestmentRow

	for _, investment := range investments {
		switch fundName := investment.FundLegalName; fundName {
		case FUNDI:
			fundone = append(fundone, investment)
		case FUNDII:
			fundtwo = append(fundtwo, investment)

		default:
			fmt.Printf("%s. is unknown investor type", fundName)
		}
	}
	return fundone, fundtwo
}

func InvestmentSpreadData(investments []*models.InvestmentRow) ([]string, []decimal.Decimal) {
	var startupnames []string
	var amounts []decimal.Decimal

	for _, investment := range investments {
		startupnames = append(startupnames, investment.StartupName)
		amounts = append(amounts, investment.InvestedCapital)
	}
	return startupnames, amounts
}
