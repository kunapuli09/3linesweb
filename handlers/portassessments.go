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

func Assessments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	db := r.Context().Value("db").(*sqlx.DB)

	Investment_ID, err := strconv.ParseInt(r.URL.Query().Get("Investment_ID"), 10, 64)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	//fmt.Printf("ApplicationID%v", Investment_ID)
	Assessment_ID, err := strconv.ParseInt(r.URL.Query().Get("Assessment_ID"), 10, 64)
	if err != nil {
		Assessment_ID = 0
	}
	sessionStore := r.Context().Value("sessionStore").(sessions.Store)
	session, _ := sessionStore.Get(r, "3linesweb-session")
	currentUser, ok := session.Values["user"].(*models.UserRow)
	if !ok || !(currentUser.InvestorRelations || currentUser.Admin ) {
		http.Redirect(w, r, "/logout", 302)
		return
	}
	investment, err := models.NewInvestment(db).GetById(nil, Investment_ID)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	// fmt.Printf("Investment%v", investment)
	// fmt.Printf("Assessment_ID%v", Assessment_ID)
	//TODO
	assessment, err := models.NewAssessment(db).GetByInvestmentId(nil, Assessment_ID, Investment_ID)
	// if err != nil {
	// 	fmt.Printf("DB Error %v", err)
	// }
	//fmt.Printf("Assessments%v", Assessments)
	//create session date for page rendering
	data := struct {
		CurrentUser *models.UserRow
		Count       int
		Investment  *models.InvestmentRow
		Assessment  *models.AssessmentRow
	}{
		currentUser,
		getCount(w, r, currentUser.Email),
		investment,
		assessment,
	}
	funcMap := template.FuncMap{
		"safeHTML": func(b string) template.HTML {
			return template.HTML(b)
		},
	}
	tmpl, err := template.New("main").Funcs(funcMap).ParseFiles("templates/portfolio/basic.html.tmpl", "templates/portfolio/editAssessments.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	tmpl.ExecuteTemplate(w, "layout", data)
}

//db call to update
func UpdateAssessment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	db := r.Context().Value("db").(*sqlx.DB)
	var i models.AssessmentRow

	err := r.ParseForm()
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	Assessment_ID, e := strconv.ParseInt(r.FormValue("id"), 10, 64)
	if e != nil {
		libhttp.HandleErrorJson(w, e)
		return
	}
	Investment_ID, e := strconv.ParseInt(r.FormValue("Investment_ID"), 10, 64)
	if e != nil {
		libhttp.HandleErrorJson(w, e)
		return
	}
	investment, err := models.NewInvestment(db).GetById(nil, Investment_ID)
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
	m := structs.Map(i)
	m["StartupName"] = investment.StartupName
	yr3rev := m["YearThreeForecastedRevenue"].(decimal.Decimal)
	marketmultiple := m["MarketMultiple"].(decimal.Decimal)
	hundred, _ := decimal.NewFromString("100")
	ownership := investment.FundOwnershipPercentage.Div(hundred)
	threelinesValueAtExit := yr3rev.Mul(marketmultiple).Mul(ownership)
	m["ThreelinesValueAtExit"] = threelinesValueAtExit
	m["YearThreeExitMultiple"] = threelinesValueAtExit.Div(investment.InvestedCapital).Ceil()
	//fmt.Printf("map %v", m)
	if Assessment_ID == 0 {
		//fmt.Printf("Creating New Assessment %v", m)
		assessment, err2 := models.NewAssessment(db).Create(nil, m)
		if err2 != nil {
			libhttp.HandleErrorJson(w, err2)
			return
		}
		Assessment_ID = assessment.ID

	} else {
		//fmt.Printf("Updating Notes with ApplicationID%v, Assessments_ID%v", Investment_ID, Assessments_ID)

		_, err3 := models.NewAssessment(db).UpdateById(nil, Assessment_ID, m)
		if err3 != nil {
			libhttp.HandleErrorJson(w, err3)
			return
		}
	}
	address := fmt.Sprintf("/viewinvestment?id=%v", Investment_ID)
	http.Redirect(w, r, address, 302)
}

//db call to update
func RemoveAssessment(w http.ResponseWriter, r *http.Request) {
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
	_, err2 := models.NewAssessment(db).DeleteByID(nil, ID)
	if err2 != nil {
		libhttp.HandleErrorJson(w, err2)
		return
	}
	address := fmt.Sprintf("/viewinvestment?Investment_ID=%v", Investment_ID)
	http.Redirect(w, r, address, 302)
}
