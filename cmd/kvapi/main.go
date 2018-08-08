package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/christianalexander/kvapi/auth0"
	"github.com/christianalexander/kvapi/inmemory"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

var version = "LOCAL"

var config struct {
	Verbose      bool
	Address      string   `default:":8080"`
	AuthAudience []string `required:"true" split_words:"true"`
	AuthIssuer   string   `required:"true" split_words:"true"`
}

func init() {
	envconfig.MustProcess("", &config)
}

func main() {
	logrus.Printf("KV API %s", version)
	kv := inmemory.NewKV()

	r := mux.NewRouter()

	a0 := auth0.NewService(config.AuthAudience, config.AuthIssuer)
	r.Use(a0.Middleware())

	r.HandleFunc("/{key}", func(w http.ResponseWriter, r *http.Request) {
		uid := context.Get(r, "uid").(string)
		pathVars := mux.Vars(r)

		if r.Method == http.MethodGet {
			value := kv.Get(uid, pathVars["key"])
			if value == "" {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("Not Found"))
			}
			_, err := w.Write([]byte(value))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error sending value: %v", err)
			}
		} else if r.Method == http.MethodPost {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Bad Request"))
			}

			r.Body.Close()
			kv.Set(uid, pathVars["key"], string(body))
			w.WriteHeader(http.StatusCreated)
		}
	}).Methods(http.MethodGet, http.MethodPost)

	log.Fatalln("Server terminated", http.ListenAndServe(":8080", r))
}
