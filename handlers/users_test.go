package handlers

import (
	"context"
	"encoding/gob"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/kunapuli09/3linesweb/models"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func init() {
	if err := os.Chdir("../.."); err != nil {
		panic(err)
	}
}
func TestGetHome(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "http://localhost:8888/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(GetHome)
	handler.ServeHTTP(res, req)
	//fmt.Println(res.Body.String())

	// Check the status code is what we expect.
	if status := res.Code; status != http.StatusOK {
		t.Errorf("Home handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// // Check the response body is what we expect.
	// expected := `{"alive": true}`
	// if res.Body.String() != expected {
	//     t.Errorf("Home handler returned unexpected body: got %v want %v",
	//         res.Body.String(), expected)
	// }
}

func newDBForTest(t *testing.T) *sqlx.DB {

	//defaultDSN := os.Getenv("DB_URL")
	//defaultDSN := strings.Replace("docker:docker@tcp(db:3307)/3lineswebtest", "-", "_", -1)
	defaultDSN := strings.Replace("docker:docker@tcp(localhost:3306)/3lineswebtest", "-", "_", -1)
	db, err := sqlx.Connect("mysql", defaultDSN)
	if err != nil {
		t.Fatalf("Connecting to local MySQL should never fail. Error: %v", err)
	}
	return db
}

func newSessionStoreForTest(t *testing.T) sessions.Store {
	defaultStore := "3lineswebtest"
	store := sessions.NewCookieStore([]byte(defaultStore))
	store.Options = &sessions.Options{
		MaxAge:   60 * 30,
		HttpOnly: true,
	}
	return store
}

func TestPostSignup(t *testing.T) {

	token := `03AHaCkAb61RioFRrk_uW5318qBUThwIF_QDXTsxrmSac6u1o-262B-NKJgBKMf1vFfSFdZEnuL7w80me-4s73Mw-ESb9KN7KiRKKCSP9lFx4_nG_uCS79QcwME3i0qrm3XiwJffSpJ0xswyZz6KKGNFVSPqgBIID7iBLvb3RPLQk9LigvZSqRNtqv0G-sHHeJNFfBD0CkYJAs-x0a1SlBfdJrN9KUIcM5OCtF6M7BEO3lvRbv_YXwwaw9pMqBhmLkV0GI5XO_IjX855U_yn8lTv9KnhrNq0pekalmZGUBSd3qhv8RPUneabdW9C2lndBZHkUdnhSWMsm1wXvfszvWLHEuMLyq15-gWY3Ha2GHlqS9ZVi7oO2H04tZU8TXM5_qxRYYsJlYY-Vi-Be7DhPTlM3YcwcGlo0dI3s3vQ6bVnO_pzNkDJWczrKvpfkpoe3Q5DRBucv128H7`

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	url := fmt.Sprintf("http://localhost:9998/signup?Email=test@test.com&Phone=1234567890&Password=test&PasswordAgain=test&token=%s", token)
	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	gob.Register("github.com/kunapuli09/3linesweb/models/models.UserRow{}")
	ctx := req.Context()
	ctx = context.WithValue(ctx, "db", newDBForTest(t))
	ctx = context.WithValue(ctx, "sessionStore", newSessionStoreForTest(t))
	req = req.WithContext(ctx)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	handler := http.HandlerFunc(PostSignup)
	handler.ServeHTTP(res, req)

	//real test if session saved the user
	sessionStore := req.Context().Value("sessionStore").(sessions.Store)
	session, _ := sessionStore.Get(req, "3linesweb-session")
	currentUser := session.Values["user"].(*models.UserRow)

	if currentUser.Email != req.FormValue("Email") {
		t.Errorf("Signup handler saved session: got %v want %v",
			currentUser.Email, req.FormValue("Email"))
	}

	if currentUser.Email != "Notsame" {
		t.Errorf("Signup handler saved session bad test: got %v want %v",
			currentUser.Email, "Notsame")
	}

	//body, _ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(body))
	// Check the status code is what we expect.
	if status := res.Code; status != 302 {
		t.Errorf("Signup handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
