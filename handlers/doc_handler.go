package handlers

import (
	"crypto/md5"
	"fmt"
	//"github.com/fatih/structs"
	//"github.com/gorilla/schema"
	"errors"
	"github.com/gorilla/sessions"
	"github.com/haisum/recaptcha"
	"github.com/jmoiron/sqlx"
	"github.com/kunapuli09/3linesweb/libhttp"
	"github.com/kunapuli09/3linesweb/models"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func GetInvestmentDocs(w http.ResponseWriter, r *http.Request) {
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
	if !ok || !currentUser.Admin {
		http.Redirect(w, r, "/logout", 302)
		return
	}
	investment, err := models.NewInvestment(db).GetById(nil, Investment_ID)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	alldocs, err := models.NewInvestmentDoc(db).GetAllByInvestmentId(nil, Investment_ID)
	//create empty investmentstructure
	doc := models.InvestmentDocRow{Hash: hash}
	//create session date for page rendering
	data := struct {
		CurrentUser *models.UserRow
		Count       int
		Investment  *models.InvestmentRow
		Doc         models.InvestmentDocRow
		Existing    []*models.InvestmentDocRow
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

func GetUserDocs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	//prepare page
	db := r.Context().Value("db").(*sqlx.DB)
	sessionStore := r.Context().Value("sessionStore").(sessions.Store)
	session, _ := sessionStore.Get(r, "3linesweb-session")
	currentUser, ok := session.Values["user"].(*models.UserRow)
	if !ok || !currentUser.Admin {
		http.Redirect(w, r, "/logout", 302)
		return
	}
	alldocs, err := models.NewUserDoc(db).AllDocs(nil)
	u := models.NewUser(db)
	users, err := u.AllUsers(nil)
	//create session date for page rendering
	data := struct {
		CurrentUser *models.UserRow
		Count       int
		Existing    []*models.UserDocRow
		Users       []*models.UserRow
	}{
		currentUser,
		getCount(w, r, currentUser.Email),
		alldocs,
		users,
	}
	tmpl, err := template.ParseFiles("templates/portfolio/basic.html.tmpl", "templates/portfolio/userdocs.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	tmpl.ExecuteTemplate(w, "layout", data)
}

func GetCareers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	//prepare page
	db := r.Context().Value("db").(*sqlx.DB)
	i := models.NewInvestment(db)
	investments, err := i.GetAllInvestmentsWithoutSyndicates(nil)
	//create session date for page rendering
	data := struct {
		Investments    []*models.InvestmentRow
		SuccessMessage string
	}{
		investments,
		"",
	}
	tmpl, err := template.ParseFiles("templates/careers/careers.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	tmpl.ExecuteTemplate(w, "layout", data)
}

//database call to add new
func AddProposalDoc(w http.ResponseWriter, r *http.Request) {
	var docs []*models.ProposalDocRow
	db := r.Context().Value("db").(*sqlx.DB)
	i := models.NewInvestment(db)
	w.Header().Set("Content-Type", "text/html")
	re := recaptcha.R{
		Secret: os.Getenv("CAPTCHA_SITE_SECRET"),
	}
	token := r.FormValue("g-recaptcha-response")
	log.Println("Verifying Captcha token", token)
	isValid := re.VerifyResponse(token)
	if !isValid {
		log.Printf("Invalid Captcha! These errors ocurred: %v", re.LastError())
		libhttp.HandleErrorJson(w, errors.New("Invalid Captcha!"))
		return
	}
	investment_ID, e1 := strconv.ParseInt(r.FormValue("Investment_ID"), 10, 64)
	if e1 != nil {
		libhttp.HandleErrorJson(w, e1)
		return
	}
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//get a ref to the parsed multipart form
	m := r.MultipartForm

	//get the *fileheaders
	files := m.File["Doc"]
	//---- parse uploaded file------
	//copy each part to destination.
	for i, _ := range files {
		//for each fileheader, get a handle to the actual file
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//create destination file making sure the path is writeable.
		dst, err := os.Create("./docs/" + files[i].Filename)
		defer dst.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//copy the uploaded file to the destination file
		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//Generate MD5 Hash for a new file
		curtime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(curtime, 10))
		hash := fmt.Sprintf("%x", h.Sum(nil))
		doc := models.ProposalDocRow{
			Investment_ID: investment_ID,
			UploadDate:    time.Now(),
			DocPath:       "/files/" + files[i].Filename,
			Hash:          hash,
			DocName:       files[i].Filename,
			Email:         r.FormValue("Email"),
			Phone:         r.FormValue("Phone"),
			CompanyName:   r.FormValue("CompanyName"),
			FullName:      r.FormValue("FullName"),
		}
		fmt.Printf("doc info %v", doc)
		docs = append(docs, &doc)
	}
	//------files uploaded ------
	_, err4 := models.NewProposalDoc(db).BatchInsert(nil, docs)
	if err4 != nil {
		fmt.Println("database error")
		libhttp.HandleErrorJson(w, err4)
		return
	}
	tmpl, err := template.ParseFiles("templates/careers/careers.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	investments, err := i.GetAllInvestmentsWithoutSyndicates(nil)
	//create session date for page rendering
	data := struct {
		Investments    []*models.InvestmentRow
		SuccessMessage string
	}{
		investments,
		"Thank you for your interest to collaborate and 3Lines team will reach out to you soon.",
	}

	tmpl.ExecuteTemplate(w, "layout", data)
}

//database call to add new
func AddInvestmentDocs(w http.ResponseWriter, r *http.Request) {
	var docs []*models.InvestmentDocRow
	db := r.Context().Value("db").(*sqlx.DB)
	w.Header().Set("Content-Type", "text/html")
	sessionStore := r.Context().Value("sessionStore").(sessions.Store)
	session, _ := sessionStore.Get(r, "3linesweb-session")
	currentUser, ok := session.Values["user"].(*models.UserRow)
	if !ok || !currentUser.Admin {
		http.Redirect(w, r, "/logout", 302)
		return
	}
	investment_ID, e1 := strconv.ParseInt(r.FormValue("Investment_ID"), 10, 64)
	if e1 != nil {
		libhttp.HandleErrorJson(w, e1)
		return
	}
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//get a ref to the parsed multipart form
	m := r.MultipartForm

	//get the *fileheaders
	files := m.File["Doc"]
	//---- parse uploaded file------
	//copy each part to destination.
	for i, _ := range files {
		//for each fileheader, get a handle to the actual file
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//create destination file making sure the path is writeable.
		dst, err := os.Create("./docs/" + files[i].Filename)
		defer dst.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//copy the uploaded file to the destination file
		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//Generate MD5 Hash for a new file
		curtime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(curtime, 10))
		hash := fmt.Sprintf("%x", h.Sum(nil))
		doc := models.InvestmentDocRow{
			Investment_ID: investment_ID,
			UploadDate:    time.Now(),
			DocPath:       "/files/" + files[i].Filename,
			Hash:          hash,
			DocName:       files[i].Filename,
		}
		fmt.Printf("doc info %v", doc)
		docs = append(docs, &doc)
	}
	//------files uploaded ------
	_, err4 := models.NewInvestmentDoc(db).BatchInsert(nil, docs)
	if err4 != nil {
		fmt.Println("database error")
		libhttp.HandleErrorJson(w, err4)
		return
	}
	address := fmt.Sprintf("/investmentDocs?Investment_ID=%v", investment_ID)
	http.Redirect(w, r, address, 302)
}

//database call to add new
func AddUserDocs(w http.ResponseWriter, r *http.Request) {
	var docs []*models.UserDocRow
	db := r.Context().Value("db").(*sqlx.DB)
	w.Header().Set("Content-Type", "text/html")
	sessionStore := r.Context().Value("sessionStore").(sessions.Store)
	session, _ := sessionStore.Get(r, "3linesweb-session")
	currentUser, ok := session.Values["user"].(*models.UserRow)
	if !ok || !currentUser.Admin {
		http.Redirect(w, r, "/logout", 302)
		return
	}
	user_ID, e1 := strconv.ParseInt(r.FormValue("User_ID"), 10, 64)
	if e1 != nil {
		libhttp.HandleErrorJson(w, e1)
		return
	}
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//get a ref to the parsed multipart form
	m := r.MultipartForm

	//get the *fileheaders
	files := m.File["Doc"]
	//---- parse uploaded file------
	//copy each part to destination.
	for i, _ := range files {
		//for each fileheader, get a handle to the actual file
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//create destination file making sure the path is writeable.
		dst, err := os.Create("./docs/" + files[i].Filename)
		defer dst.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//copy the uploaded file to the destination file
		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//Generate MD5 Hash for a new file
		curtime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(curtime, 10))
		hash := fmt.Sprintf("%x", h.Sum(nil))
		doc := models.UserDocRow{User_ID: user_ID, UploadDate: time.Now(), DocPath: "/files/" + files[i].Filename, Hash: hash, DocName: files[i].Filename}
		fmt.Printf("doc info %v", doc)
		docs = append(docs, &doc)
	}
	//------files uploaded ------
	_, err4 := models.NewUserDoc(db).BatchInsert(nil, docs)
	if err4 != nil {
		fmt.Println("database error")
		libhttp.HandleErrorJson(w, err4)
		return
	}
	//TODO send back to user profile
	address := fmt.Sprintf("/userdocs")
	http.Redirect(w, r, address, 302)
}

//db call to update
func RemoveInvestmentDoc(w http.ResponseWriter, r *http.Request) {
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
	_, err2 := models.NewInvestmentDoc(db).DeleteByID(nil, ID)
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

	address := fmt.Sprintf("/investmentDocs?Investment_ID=%v", Investment_ID)
	http.Redirect(w, r, address, 302)
}

//db call to update
func RemoveUserDoc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	db := r.Context().Value("db").(*sqlx.DB)
	ID, e := strconv.ParseInt(r.FormValue("id"), 10, 64)
	if e != nil {
		libhttp.HandleErrorJson(w, e)
		return
	}
	_, err2 := models.NewUserDoc(db).DeleteByID(nil, ID)
	if err2 != nil {
		libhttp.HandleErrorJson(w, err2)
		return
	}
	docName := r.FormValue("DocName")
	docPath := "./docs/" + docName
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
	address := fmt.Sprintf("/userdocs")
	http.Redirect(w, r, address, 302)
}
