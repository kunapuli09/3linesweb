package handlers

import (
	"crypto/md5"
	"fmt"
	"github.com/fatih/structs"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/kunapuli09/3linesweb/libhttp"
	"github.com/kunapuli09/3linesweb/models"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

func Docs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	//Generate MD5 Hash for a new file
	curtime := time.Now().Unix()
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(curtime, 10))
	hash := fmt.Sprintf("%x", h.Sum(nil))
	//prepare page
	db := r.Context().Value("db").(*sqlx.DB)
	Investment_ID, err := strconv.ParseInt(r.URL.Query().Get("Investment_ID"), 10, 64)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	sessionStore := r.Context().Value("sessionStore").(sessions.Store)
	session, _ := sessionStore.Get(r, "3linesweb-session")
	currentUser, ok := session.Values["user"].(*models.UserRow)
	if !ok || !currentUser.Admin || !currentUser.FundOne || !currentUser.FundTwo{
		http.Redirect(w, r, "/logout", 302)
		return
	}
	investment, err := models.NewInvestment(db).GetById(nil, Investment_ID)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	alldocs, err := models.NewDoc(db).GetAllByInvestmentId(nil, Investment_ID)
	//create empty investmentstructure
	doc := models.DocRow{Hash: hash}
	//create session date for page rendering
	data := struct {
		CurrentUser *models.UserRow
		Count       int
		Investment  *models.InvestmentRow
		Doc         models.DocRow
		Existing    []*models.DocRow
	}{
		currentUser,
		getCount(w, r, currentUser.Email),
		investment,
		doc,
		alldocs,
	}
	tmpl, err := template.ParseFiles("templates/portfolio/basic.html.tmpl", "templates/portfolio/newdoc.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	tmpl.ExecuteTemplate(w, "layout", data)
}

//database call to add new
func AddDoc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	//---- parse uploaded file------
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("Doc")
	if err != nil {
		fmt.Println("file parse error")
		fmt.Println(err)
		libhttp.HandleErrorJson(w, err)
		return
	}
	defer file.Close()
	fmt.Fprintf(w, "%v", handler.Header)
	docPath := "./docs/" + handler.Filename
	f, err1 := os.OpenFile(docPath, os.O_WRONLY|os.O_CREATE, 0666)
	if err1 != nil {
		fmt.Println("file open error")
		fmt.Println(err1)
		return
	}
	defer f.Close()
	io.Copy(f, file)
	//------file upload complete ------
	var i models.DocRow
	db := r.Context().Value("db").(*sqlx.DB)
	err2 := r.ParseForm()
	if err2 != nil {
		libhttp.HandleErrorJson(w, err2)
		return
	}
	decoder := schema.NewDecoder()
	decoder.RegisterConverter(time.Time{}, ConvertFormDate)
	err3 := decoder.Decode(&i, r.PostForm)
	if err3 != nil {
		fmt.Println("decoding error")
		libhttp.HandleErrorJson(w, err3)
		return
	}
	m := structs.Map(i)
	m["DocPath"] = "/files/" + handler.Filename
	m["DocName"] = handler.Filename
	fmt.Printf("map %v", m)
	_, err4 := models.NewDoc(db).Create(nil, m)
	if err4 != nil {
		fmt.Println("database error")
		libhttp.HandleErrorJson(w, err4)
		return
	}
	address := fmt.Sprintf("/viewinvestment?id=%v", m["Investment_ID"])
	http.Redirect(w, r, address, 302)
}

//db call to update
func RemoveDoc(w http.ResponseWriter, r *http.Request) {
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
	_, err2 := models.NewDoc(db).DeleteByID(nil, ID)
	if err2 != nil {
		libhttp.HandleErrorJson(w, err2)
		return
	}
	docPath := r.FormValue("DocPath")
	if _, err := os.Stat(docPath); os.IsNotExist(err) {
		fmt.Printf("file does not exist")
	} else {
		fmt.Printf("file exists. removing it")
		err1 := os.Remove(docPath)
		if err1 != nil {
			fmt.Println("file remove error")
			fmt.Println(err1)
			return
		}

	}

	address := fmt.Sprintf("/docs?Investment_ID=%v", Investment_ID)
	http.Redirect(w, r, address, 302)
}
