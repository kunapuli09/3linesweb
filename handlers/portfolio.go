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
	"sort"
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
	//fmt.Printf("%s. Inside EntryAccess", currentUser.Email)

	switch defaultView := true; defaultView {
	case currentUser.Admin:
		GetAdminDashboard(w, r)
	case currentUser.Dsc:
		InvestorDashboard(w, r)
	case currentUser.Investor:
		InvestorDashboard(w, r)
	case currentUser.BlogReader:
		BlogDashboard(w, r)
	default:
		fmt.Println("reached default as well")
	}

}

func GetAdminDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	db := r.Context().Value("db").(*sqlx.DB)
	sessionStore := r.Context().Value("sessionStore").(sessions.Store)
	session, _ := sessionStore.Get(r, "3linesweb-session")
	currentUser, ok := session.Values["user"].(*models.UserRow)
	if !ok || !currentUser.Admin {
		http.Redirect(w, r, "/logout", 302)
		return
	}
	investments := []*models.InvestmentRow{}
	funcMap := template.FuncMap{
		"safeHTML": func(b string) template.HTML {
			return template.HTML(b)
		},
		"currencyFormat": func(currency decimal.Decimal) string {
			f, _ := currency.Float64()
			return ac.FormatMoney(f)
		},
	}
	contributions, err := models.NewContribution(db).AllContributions(nil)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	participatedFundNamesForBackend, participatedFundNamesForWeb, capitalcontributions := CapitalContributionDataForPieChart(contributions)

	if len(participatedFundNamesForBackend) > 0 {
		i, err := models.NewInvestment(db).GetUserInvestments(nil, unique(participatedFundNamesForBackend))
		if err != nil {
			libhttp.HandleErrorJson(w, err)
			return
		}
		investments = append(investments, i...)
	}
	startupnames, amounts := CapitalSpreadDataForBarChart(investments)

	//data for entryaccess.html.tmpl
	data := struct {
		CurrentUser          *models.UserRow
		Count                int
		Contributions        []*models.ContributionRow
		Investments          []*models.InvestmentRow
		StartupNames         []string
		Amounts              []decimal.Decimal
		FundNames            []string
		CapitalContributions []decimal.Decimal
	}{
		currentUser,
		getCount(w, r, currentUser.Email),
		contributions,
		investments,
		startupnames,
		amounts,
		participatedFundNamesForWeb,
		capitalcontributions,
	}

	tmpl, err := template.New("main").Funcs(funcMap).ParseFiles("templates/portfolio/basic.html.tmpl", "templates/portfolio/admindashboard.html.tmpl")

	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	tmpl.ExecuteTemplate(w, "layout", data)
}

func GetRevenueSummaryDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	db := r.Context().Value("db").(*sqlx.DB)
	sessionStore := r.Context().Value("sessionStore").(sessions.Store)
	session, _ := sessionStore.Get(r, "3linesweb-session")
	currentUser, ok := session.Values["user"].(*models.UserRow)
	if !ok || !currentUser.Admin {
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
	revenues, err := models.NewInvestment(db).GetRevenueSummary(nil)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	table := BuildRevenueSummaryDisplayTable(revenues)

	sort.Slice(table, func(i, j int) bool {
		return table[i].InvestmentMultiple.GreaterThan(table[j].InvestmentMultiple)
	})

	//data for entryaccess.html.tmpl
	data := struct {
		CurrentUser *models.UserRow
		Count       int
		Table       []*models.RevenueDisplay
	}{
		currentUser,
		getCount(w, r, currentUser.Email),
		table,
	}

	tmpl, err := template.New("main").Funcs(funcMap).ParseFiles("templates/portfolio/basic.html.tmpl", "templates/portfolio/adminrevenuesummary.html.tmpl")

	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	tmpl.ExecuteTemplate(w, "layout", data)
}

func InvestorDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	db := r.Context().Value("db").(*sqlx.DB)
	sessionStore := r.Context().Value("sessionStore").(sessions.Store)
	session, _ := sessionStore.Get(r, "3linesweb-session")
	currentUser, ok := session.Values["user"].(*models.UserRow)
	if !ok {
		http.Redirect(w, r, "/logout", 302)
		return
	}
	investments := []*models.InvestmentRow{}
	var displayPieChart bool
	funcMap := template.FuncMap{
		"safeHTML": func(b string) template.HTML {
			return template.HTML(b)
		},
		"currencyFormat": func(currency decimal.Decimal) string {
			f, _ := currency.Float64()
			return ac.FormatMoney(f)
		},
	}
	investmentdocs, err := models.NewUserDoc(db).GetAllByUserId(nil, currentUser.ID)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	contributions, err := models.NewContribution(db).GetAllByFundNameAndUserId(nil, "", currentUser.ID)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	participatedFundNamesForBackend, participatedFundNamesForWeb, capitalcontributions := CapitalContributionDataForPieChart(contributions)

	if len(participatedFundNamesForBackend) > 0 {
		i, err := models.NewInvestment(db).GetUserInvestments(nil, participatedFundNamesForBackend)
		if err != nil {
			libhttp.HandleErrorJson(w, err)
			return
		}
		investments = append(investments, i...)
	}
	startupnames, amounts := CapitalSpreadDataForBarChart(investments)
	if len(contributions) > 1 {
		displayPieChart = true
	}
	//Display Fund I /Fund II Specific Data
	_, f1Investor := Find(participatedFundNamesForBackend, FUNDI)
	_, f2Investor := Find(participatedFundNamesForBackend, FUNDII)
	//data for entryaccess.html.tmpl
	data := struct {
		CurrentUser          *models.UserRow
		Count                int
		Contributions        []*models.ContributionRow
		InvestmentDocs       []*models.UserDocRow
		Investments          []*models.InvestmentRow
		StartupNames         []string
		Amounts              []decimal.Decimal
		FundNames            []string
		CapitalContributions []decimal.Decimal
		DisplayPieChart      bool
		FundIInvestor        bool
		FundIIInvestor       bool
	}{
		currentUser,
		getCount(w, r, currentUser.Email),
		contributions,
		investmentdocs,
		investments,
		startupnames,
		amounts,
		participatedFundNamesForWeb,
		capitalcontributions,
		displayPieChart,
		f1Investor,
		f2Investor,
	}

	tmpl, err := template.New("main").Funcs(funcMap).ParseFiles("templates/portfolio/basic.html.tmpl", "templates/portfolio/investordashboard.html.tmpl")

	if err != nil {
		fmt.Println(err)
		libhttp.HandleErrorJson(w, err)
		return
	}
	tmpl.ExecuteTemplate(w, "layout", data)
}

func BlogDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	sessionStore := r.Context().Value("sessionStore").(sessions.Store)
	session, _ := sessionStore.Get(r, "3linesweb-session")
	currentUser, ok := session.Values["user"].(*models.UserRow)

	if !ok {
		http.Redirect(w, r, "/logout", 302)
		return
	}
	//data for entryaccess.html.tmpl
	data := struct {
		CurrentUser *models.UserRow
	}{
		currentUser,
	}

	tmpl, err := template.ParseFiles("templates/blog/blogdashboard.html.tmpl", "templates/blog/securedblogs.html.tmpl")

	if err != nil {
		fmt.Println(err)
		libhttp.HandleErrorJson(w, err)
		return
	}
	tmpl.ExecuteTemplate(w, "bloglayout", data)
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
		CurrentUser *models.UserRow
		Count       int
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
func ViewAdminDashboard(w http.ResponseWriter, r *http.Request) {
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
	AllDocs, err := models.NewInvestmentDoc(db).GetAllByInvestmentId(nil, ID)

	//create session date for page rendering
	data := struct {
		CurrentUser                  *models.UserRow
		Count                        int
		Investment                   *models.InvestmentRow
		Existing                     []*models.FinancialResultsRow
		ExistingNews                 []*models.NewsRow
		ExistingCapitalStructures    []*models.CapitalizationStructure
		ExistingInvestmentStructures []*models.InvestmentStructureRow
		ExistingDocs                 []*models.InvestmentDocRow
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
	AllDocs, err := models.NewInvestmentDoc(db).GetAllByInvestmentId(nil, ID)

	//create session date for page rendering
	data := struct {
		CurrentUser                  *models.UserRow
		Count                        int
		Investment                   *models.InvestmentRow
		Existing                     []*models.FinancialResultsRow
		ExistingNews                 []*models.NewsRow
		ExistingCapitalStructures    []*models.CapitalizationStructure
		ExistingInvestmentStructures []*models.InvestmentStructureRow
		ExistingDocs                 []*models.InvestmentDocRow
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
	if !ok || !currentUser.Admin {
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

func CapitalSpreadDataForBarChart(investments []*models.InvestmentRow) ([]string, []decimal.Decimal) {
	var startupnames []string
	var amounts []decimal.Decimal

	for _, investment := range investments {
		startupnames = append(startupnames, investment.StartupName)
		amounts = append(amounts, investment.InvestedCapital)
	}
	return startupnames, amounts
}

func CapitalContributionDataForPieChart(contributions []*models.ContributionRow) ([]string, []string, []decimal.Decimal) {
	var participatedFundNamesForBackend []string
	var participatedFundNamesForWeb []string
	var amounts []decimal.Decimal
	for _, contribution := range contributions {
		contributedCapitalWithOwnership := fmt.Sprintf("%s", contribution.FundLegalName)
		participatedFundNamesForBackend = append(participatedFundNamesForBackend, contribution.FundLegalName)
		participatedFundNamesForWeb = append(participatedFundNamesForWeb, contributedCapitalWithOwnership)
		amounts = append(amounts, contribution.InvestmentAmount)
	}
	//fmt.Printf("%v \n %v",participatedFundNamesForWeb, amounts)
	return participatedFundNamesForBackend, participatedFundNamesForWeb, amounts
}

func unique(slice []string) []string {
	keys := make(map[string]bool)
	var participatedFundNamesForBackend []string
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			participatedFundNamesForBackend = append(participatedFundNamesForBackend, entry)
		}
	}
	return participatedFundNamesForBackend
}

// Find takes a slice and looks for an element in it. If found it will
// return it's key, otherwise it will return -1 and a bool of false.
func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func BuildRevenueSummaryDisplayTable(revenues []*models.RevenueSummary) []*models.RevenueDisplay {
	data := make(map[string]*models.RevenueDisplay)
	zero, _ := decimal.NewFromString("0")
	percentage, _ := decimal.NewFromString("100")
	for _, revenue := range revenues {
		decimal.DivisionPrecision = 4
		RTC := revenue.Revenue.Div(revenue.TotalCapitalRaised).Mul(percentage)
		if _, ok := data[revenue.StartupName]; !ok {
			display := &models.RevenueDisplay{
				ID:                 revenue.ID,
				StartupName:        revenue.StartupName,
				InvestedCapital:    revenue.InvestedCapital,
				TotalCapitalRaised: revenue.TotalCapitalRaised,
				ReportedValue:      revenue.ReportedValue,
				InvestmentMultiple: revenue.InvestmentMultiple,
			}

			if revenue.ReportingDate.Year() == time.Now().Year() {
				display.ForecastedRevenue = revenue.Revenue
				display.ForecastedEBIDTA = revenue.EBIDTA
				display.ForecastedRevenueToCapital = RTC
			}
			if revenue.ReportingDate.Year() == (time.Now().Year() - 1) {
				display.LastYearRevenue = revenue.Revenue
				display.LastYearEBIDTA = revenue.EBIDTA
				display.LastYearRevenueToCapital = RTC
			}
			data[revenue.StartupName] = display
		} else {
			if revenue.ReportingDate.Year() == time.Now().Year() {
				data[revenue.StartupName].ForecastedRevenue = revenue.Revenue
				data[revenue.StartupName].ForecastedEBIDTA = revenue.EBIDTA
				data[revenue.StartupName].ForecastedRevenueToCapital = RTC //Undefined
			}
			if revenue.ReportingDate.Year() == (time.Now().Year() - 1) {
				data[revenue.StartupName].LastYearRevenue = revenue.Revenue
				data[revenue.StartupName].LastYearEBIDTA = revenue.EBIDTA
				data[revenue.StartupName].LastYearRevenueToCapital = RTC //Undefined
			}

		}

	}

	for _, display := range data {
		if display.LastYearEBIDTA.LessThan(zero) {
			display.IsLastYearEBIDTANegative = true
		}
		if display.LastYearRevenue.LessThan(zero) {
			display.IsLastYearRevenueNegative = true
		}
		if display.ForecastedEBIDTA.LessThan(zero) {
			display.IsForecastedEBIDTANegative = true
		}
		if display.ForecastedRevenue.LessThan(zero) {
			display.IsForecastedRevenueNegative = true
		}
		if display.LastYearRevenueToCapital.LessThan(percentage) {
			display.IsLastYearRevenueToCapitalNegative = true
		}
		if display.ForecastedRevenueToCapital.LessThan(percentage) {
			display.IsForecastedRevenueToCapitalNegative = true
		}
	}

	v := make([]*models.RevenueDisplay, 0, len(data))
	for _, value := range data {
		v = append(v, value)
	}
	return v
}
