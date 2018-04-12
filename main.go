package main

import (
	"encoding/gob"
	"github.com/spf13/viper"
	"github.com/tylerb/graceful"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/kunapuli09/3linesweb/application"
	"github.com/kunapuli09/3linesweb/models"
)

func init() {
	gob.Register(&models.UserRow{})
}

func newConfig() (*viper.Viper, error) {
	defaultDSN := strings.Replace("root:kk@starpath@tcp(localhost:3306)/3linesweb?parseTime=true", "-", "_", -1)

	c := viper.New()
	c.SetDefault("dsn", defaultDSN)
	c.SetDefault("cookie_secret", "z5mOYQcyv3KQHe3W")
	c.SetDefault("http_addr", ":8888")
	c.SetDefault("http_cert_file", "")
	c.SetDefault("http_key_file", "")
	c.SetDefault("http_drain_interval", "1s")

	c.AutomaticEnv()

	return c, nil
}

func main() {
	config, err := newConfig()
	if err != nil {
		log.Fatal(err)
	}

	app, err := application.New(config)
	if err != nil {
		log.Fatal(err)
	}

	middle, err := app.MiddlewareStruct()
	if err != nil {
		log.Fatal(err)
	}

	serverAddress := config.Get("http_addr").(string)

	certFile := config.Get("http_cert_file").(string)
	keyFile := config.Get("http_key_file").(string)
	drainIntervalString := config.Get("http_drain_interval").(string)

	drainInterval, err := time.ParseDuration(drainIntervalString)
	if err != nil {
		log.Fatal(err)
	}

	srv := &graceful.Server{
		Timeout: drainInterval,
		Server:  &http.Server{Addr: serverAddress, Handler: middle},
	}

	log.Println("Running HTTP server on " + serverAddress)

	if certFile != "" && keyFile != "" {
		err = srv.ListenAndServeTLS(certFile, keyFile)
	} else {
		err = srv.ListenAndServe()
	}

	if err != nil {
		log.Fatal(err)
	}
}
