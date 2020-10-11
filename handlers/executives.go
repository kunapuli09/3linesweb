package handlers

import (
	"fmt"
	"errors"
	"github.com/haisum/recaptcha"
	"github.com/kunapuli09/3linesweb/libhttp"
	"github.com/kunapuli09/3linesweb/models"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
	"github.com/gorilla/schema"
	"github.com/fatih/structs"
	"github.com/jmoiron/sqlx"
	"database/sql"
)



func GetExecutives(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	//prepare page
	tmpl, err := template.ParseFiles("templates/executives/executives.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	tmpl.ExecuteTemplate(w, "layout", nil)
}

//database call to add new
func AddExecutive(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	var i models.ExecutiveRow
	db := r.Context().Value("db").(*sqlx.DB)
	err := r.ParseForm()
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	decoder := schema.NewDecoder()
	decoder.RegisterConverter(sql.NullString{}, ConvertSQLNullString)
	err1 := decoder.Decode(&i, r.PostForm)
	fmt.Printf("Form %v", r.PostForm)
	if err1 != nil {
		fmt.Println("decoding error")
		libhttp.HandleErrorJson(w, err1)
		return
	}
	re := recaptcha.R{
		Secret: os.Getenv("CAPTCHA_SITE_SECRET"),
	}
	token := r.FormValue("rcres")
	log.Println("Verifying Captcha token", token)
	isValid := re.VerifyResponse(token)
	if !isValid {
		log.Printf("Invalid Captcha! These errors ocurred: %v", re.LastError())
		libhttp.HandleErrorJson(w, errors.New("Invalid Captcha!"))
		return
	}
	m := structs.Map(i)
	m["ApplicationDate"] = time.Now()
	fmt.Printf("map %v", m)
	//check for duplicate entry
	email, ok2 := m["Email"].(string)
	name, ok3 := m["Name"].(string)
	socialhandle, ok4 := m["SocialMediaHandle"].(string)

	if ok2 && ok3 && ok4 {
		exists := models.NewExecutive(db).GetExisting(nil, email, name, socialhandle)
		if exists == true {
			err3 := errors.New("Duplicate Entry. An executive with the same Social Media Handle or Email Already Exists")
			//fmt.Println(err2)
			libhttp.HandleErrorJson(w, err3)
			return
		}
	}

	_, err4 := models.NewExecutive(db).Create(nil, m)
	if err4 != nil {
		fmt.Println("Executive Information is not Valid", err4)
		libhttp.HandleErrorJson(w, err4)
		return
	}
	http.Redirect(w, r, "/", 302)
}

