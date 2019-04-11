package application

import (
	"github.com/carbocation/interpose"
	_ "github.com/go-sql-driver/mysql"
	gorilla_mux "github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"net/http"

	"github.com/kunapuli09/3linesweb/handlers"
	"github.com/kunapuli09/3linesweb/middlewares"
)

// New is the constructor for Application struct.
func New(config *viper.Viper) (*Application, error) {
	dsn := config.Get("dsn").(string)

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, err
	}

	cookieStoreSecret := config.Get("cookie_secret").(string)

	app := &Application{}
	app.config = config
	app.dsn = dsn
	app.db = db
	app.sessionStore = sessions.NewCookieStore([]byte(cookieStoreSecret))

	return app, nil
}

// Application is the application object that runs HTTP server.
type Application struct {
	config       *viper.Viper
	dsn          string
	db           *sqlx.DB
	sessionStore sessions.Store
}

func (app *Application) MiddlewareStruct() (*interpose.Middleware, error) {
	middle := interpose.New()
	middle.Use(middlewares.SetDB(app.db))
	middle.Use(middlewares.SetSessionStore(app.sessionStore))

	middle.UseHandler(app.mux())

	return middle, nil
}

func (app *Application) mux() *gorilla_mux.Router {
	MustLogin := middlewares.MustLogin

	router := gorilla_mux.NewRouter()

	router.Handle("/", http.HandlerFunc(handlers.GetHome)).Methods("GET")
	router.HandleFunc("/signup", handlers.GetSignup).Methods("GET")
	router.HandleFunc("/signup", handlers.PostSignup).Methods("POST")
	router.HandleFunc("/events", handlers.GetEvents).Methods("GET")
	router.HandleFunc("/performance", handlers.GetPerformance).Methods("GET")
	router.HandleFunc("/login", handlers.GetLogin).Methods("GET")
	router.HandleFunc("/login", handlers.PostLogin).Methods("POST")
	router.HandleFunc("/resetEmail", handlers.PasswordResetEmail).Methods("POST")
	router.HandleFunc("/reset", handlers.GetReset).Methods("GET")
	router.HandleFunc("/reset", handlers.Reset).Methods("POST")
	router.HandleFunc("/logout", handlers.GetLogout).Methods("GET")
	router.HandleFunc("/blog", handlers.GetBlog).Methods("GET")
	router.HandleFunc("/contact", handlers.PostEmail).Methods("POST")
	router.HandleFunc("/appl", handlers.NewApplication).Methods("GET")
	router.HandleFunc("/application", handlers.AddApplication).Methods("POST")
	router.HandleFunc("/fundingreqs", handlers.FundingRequests).Methods("GET")
	router.HandleFunc("/notifications", handlers.Notifications).Methods("GET")
	// router.HandleFunc("/notifyinvestors", handlers.NotifyInvestors).Methods("POST")
	router.HandleFunc("/portfolio", handlers.GetPortfolio).Methods("GET")
	router.HandleFunc("/viewinvestment", handlers.ViewInvestment).Methods("GET")
	router.HandleFunc("/newinvestment", handlers.NewInvestment).Methods("GET")
	router.HandleFunc("/editinvestment", handlers.EditInvestment).Methods("GET")
	router.HandleFunc("/update", handlers.Update).Methods("POST")
	router.HandleFunc("/add", handlers.Add).Methods("POST")

	//investor contribution details
	router.HandleFunc("/contributions", handlers.GetContributions).Methods("GET")
	router.HandleFunc("/newcontribution", handlers.NewContribution).Methods("GET")
	router.HandleFunc("/editcontribution", handlers.EditContribution).Methods("GET")
	router.HandleFunc("/updatecontribution", handlers.UpdateContribution).Methods("POST")
	router.HandleFunc("/addcontribution", handlers.AddContribution).Methods("POST")

	router.HandleFunc("/newfinancials", handlers.NewFinancials).Methods("GET")
	router.HandleFunc("/newinvestmentstructure", handlers.NewInvestmentStructure).Methods("GET")
	router.HandleFunc("/newcapitalstructure", handlers.NewCapitalStructure).Methods("GET")
	router.HandleFunc("/addCapitalStructure", handlers.AddCapitalStructure).Methods("POST")
	router.HandleFunc("/addInvestmentStructure", handlers.AddInvestmentStructure).Methods("POST")
	router.HandleFunc("/addFinancialResults", handlers.AddFinancialResults).Methods("POST")
	router.HandleFunc("/news", handlers.News).Methods("GET")

	router.HandleFunc("/addNews", handlers.AddNews).Methods("POST")
	router.HandleFunc("/editNews", handlers.EditNews).Methods("GET")
	router.HandleFunc("/updateNews", handlers.UpdateNews).Methods("POST")
	router.HandleFunc("/removenews", handlers.RemoveNews).Methods("GET")
	router.HandleFunc("/publishNews", handlers.PublishNotification).Methods("GET")
	router.HandleFunc("/updateNotification", handlers.UpdateNotification).Methods("GET")
	router.HandleFunc("/removecapitalstructure", handlers.RemoveCapitalStructure).Methods("GET")
	router.HandleFunc("/removeinvestmentstructure", handlers.RemoveInvestmentStructure).Methods("GET")
	router.HandleFunc("/removefinancialresults", handlers.RemoveFinancialResults).Methods("GET")
	router.HandleFunc("/docs", handlers.Docs).Methods("GET")
	router.HandleFunc("/addDoc", handlers.AddDoc).Methods("POST")
	router.HandleFunc("/removeDoc", handlers.RemoveDoc).Methods("GET")
	router.PathPrefix("/files/").Handler(http.StripPrefix("/files/", http.FileServer(http.Dir("./docs"))))

	router.Handle("/users/{id:[0-9]+}", MustLogin(http.HandlerFunc(handlers.PostPutDeleteUsersID))).Methods("POST", "PUT", "DELETE")

	// Path of static files must be last!
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("static")))

	return router
}
