package handlers

import (
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

func GetContributions(w http.ResponseWriter, r *http.Request) {
	var i models.SearchContribution
	w.Header().Set("Content-Type", "text/html")
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
	// r.PostForm is a map of our POST form values
	decoder := schema.NewDecoder()
	decoder.RegisterConverter(time.Time{}, ConvertFormDate)
	err1 := decoder.Decode(&i, r.PostForm)
	if err1 != nil {
		libhttp.HandleErrorJson(w, err1)
		return
	}
	contributions, err := models.NewContribution(db).SearchContributions(nil, i)
	//fundone, fundtwo := SplitContributionsByStatus(contributions)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	//create session date for page rendering
	data := struct {
		CurrentUser   *models.UserRow
		Count         int
		Contributions []*models.ContributionRow
	}{
		currentUser,
		getCount(w, r, currentUser.Email),
		contributions,
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
	tmpl, err := template.New("main").Funcs(funcMap).ParseFiles("templates/portfolio/basic.html.tmpl", "templates/portfolio/contributions.html.tmpl")

	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	tmpl.ExecuteTemplate(w, "layout", data)
}

//presentation edit view
func EditContribution(w http.ResponseWriter, r *http.Request) {
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
	if !ok || !currentUser.Admin{
		http.Redirect(w, r, "/logout", 302)
		return
	}

	Contribution, err := models.NewContribution(db).GetById(nil, ID)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	u := models.NewUser(db)

	users, err := u.AllUsers(nil)
	//create session data for page rendering
	data := struct {
		CurrentUser  *models.UserRow
		Count        int
		Contribution *models.ContributionRow
		Users 		[]*models.UserRow
	}{
		currentUser,
		getCount(w, r, currentUser.Email),
		Contribution,
		users,
	}
	funcMap := template.FuncMap{
		"safeHTML": func(b string) template.HTML {
			return template.HTML(b)
		},
	}
	tmpl, err := template.New("main").Funcs(funcMap).ParseFiles("templates/portfolio/basic.html.tmpl", "templates/portfolio/editContribution.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	tmpl.ExecuteTemplate(w, "layout", data)

}

//db call to update
func UpdateContribution(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	db := r.Context().Value("db").(*sqlx.DB)
	var i models.ContributionRow
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
	sessionStore := r.Context().Value("sessionStore").(sessions.Store)
	session, _ := sessionStore.Get(r, "3linesweb-session")
	_, ok := session.Values["user"].(*models.UserRow)
	if !ok {
		http.Redirect(w, r, "/logout", 302)
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
		contribution, err2 := models.NewContribution(db).Create(nil, m)
		if err2 != nil {
			libhttp.HandleErrorJson(w, err2)
			return
		}
		ID = contribution.ID

	} else {
		//fmt.Printf("Updating Contribution with ApplicationID%v, ScreeningNotes_ID%v", Application_ID, ScreeningNotes_ID)

		_, err3 := models.NewContribution(db).UpdateById(nil, ID, m)
		if err3 != nil {
			libhttp.HandleErrorJson(w, err3)
			return
		}
	}
	address := fmt.Sprintf("/contributions")
	http.Redirect(w, r, address, 302)
}
func SplitContributionsByFund(Contributions []*models.ContributionRow) ([]*models.ContributionRow, []*models.ContributionRow) {
	var fundone []*models.ContributionRow
	var fundtwo []*models.ContributionRow

	for _, contribution := range Contributions {
		switch fundName := contribution.FundLegalName; fundName {
		case FUNDI:
			fundone = append(fundone, contribution)
		case FUNDII:
			fundtwo = append(fundtwo, contribution)

		default:
			fmt.Printf("%s. is unknown investor type", fundName)
		}
	}
	return fundone, fundtwo
}

