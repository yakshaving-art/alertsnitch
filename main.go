package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.com/yakshaving.art/alertsnitch/version"
	"gitlab.com/yakshaving.art/alertsnitch/webhook"
)

// Args are the arguments that can be passed to alertsnitch
type Args struct {
	Version bool
	Address string
}

func main() {
	args := Args{}

	flag.BoolVar(&args.Version, "version", false, "print the version and exit")
	flag.StringVar(&args.Address, "listen.address", ":8080", "address in which to listen for http requests")

	flag.Parse()

	if args.Version {
		fmt.Println(version.GetVersion())
		os.Exit(0)
	}

	r := mux.NewRouter()
	r.HandleFunc("/webhook", webhookPost).Methods("POST")
	r.HandleFunc("/-/ready", readyProbe).Methods("GET")
	r.HandleFunc("/-/health", healthyProbe).Methods("GET")
	r.Handle("/metrics", promhttp.Handler())

	log.Println("Starting listener on", args.Address)
	log.Fatal(http.ListenAndServe(args.Address, r))
}

func webhookPost(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read payload: %s", err), http.StatusBadRequest)
		return
	}

	d, err := webhook.Parse(b)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid payload: %s", err), http.StatusBadRequest)
		return
	}

	log.Printf("Webhook Payload:\n%#v\n\nBody:\n%s", d, string(b))
}

func healthyProbe(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Implemented", http.StatusNotImplemented)
}

func readyProbe(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Implemented", http.StatusNotImplemented)
}
