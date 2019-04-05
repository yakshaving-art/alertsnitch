package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"gitlab.com/yakshaving.art/alertsnitch/internal"
	"gitlab.com/yakshaving.art/alertsnitch/internal/metrics"
	"gitlab.com/yakshaving.art/alertsnitch/internal/webhook"
)

// SupportedWebhookVersion is the alert webhook data version that is supported
// by this app
const SupportedWebhookVersion = "4"

// Server represents a web server that processes webhooks
type Server struct {
	db internal.Storer
	r  *mux.Router

	debug bool
}

// New returns a new web server
func New(db internal.Storer, debug bool) Server {
	r := mux.NewRouter()

	s := Server{
		db: db,
		r:  r,

		debug: debug,
	}

	r.HandleFunc("/webhook", s.webhookPost).Methods("POST")
	r.HandleFunc("/-/ready", s.readyProbe).Methods("GET")
	r.HandleFunc("/-/health", s.healthyProbe).Methods("GET")
	r.Handle("/metrics", promhttp.Handler())

	return s
}

// Start starts a new server on the given address
func (s Server) Start(address string) {
	log.Println("Starting listener on", address, "using", s.db)
	log.Fatal(http.ListenAndServe(address, s.r))
}

func (s Server) webhookPost(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	metrics.WebhooksReceivedTotal.Inc()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		metrics.InvalidWebhooksTotal.Inc()
		log.Printf("Failed to read payload: %s\n", err)
		http.Error(w, fmt.Sprintf("Failed to read payload: %s", err), http.StatusBadRequest)
		return
	}

	if s.debug {
		log.Println("Received webhook payload", string(body))
	}

	data, err := webhook.Parse(body)
	if err != nil {
		metrics.InvalidWebhooksTotal.Inc()

		log.Printf("Invalid payload: %s\n", err)
		http.Error(w, fmt.Sprintf("Invalid payload: %s", err), http.StatusBadRequest)
		return
	}

	if data.Version != SupportedWebhookVersion {
		metrics.InvalidWebhooksTotal.Inc()

		log.Printf("Invalid payload: webhook version %s is not supported\n", data.Version)
		http.Error(w, fmt.Sprintf("Invalid payload: webhook version %s is not supported", data.Version), http.StatusBadRequest)
		return
	}

	metrics.AlertsReceivedTotal.WithLabelValues(data.Receiver, data.Status).Add(float64(len(data.Alerts)))

	if err = s.db.Save(data); err != nil {
		metrics.AlertsSavingFailuresTotal.WithLabelValues(data.Receiver, data.Status).Add(float64(len(data.Alerts)))

		log.Printf("failed to save alerts: %s\n", err)
		http.Error(w, fmt.Sprintf("failed to save alerts: %s", err), http.StatusInternalServerError)
		return
	}
	metrics.AlertsSavedTotal.WithLabelValues(data.Receiver, data.Status).Add(float64(len(data.Alerts)))
}

func (s Server) healthyProbe(w http.ResponseWriter, r *http.Request) {
	if err := s.db.Ping(); err != nil {
		log.Printf("failed to ping database server: %s", err)
		http.Error(w, fmt.Sprintf("failed to ping database server: %s", err), http.StatusServiceUnavailable)
		return
	}
}

func (s Server) readyProbe(w http.ResponseWriter, r *http.Request) {
	if err := s.db.Ping(); err != nil {
		log.Printf("database is not reachable: %s", err)
		http.Error(w, fmt.Sprintf("database is not reachable: %s", err), http.StatusServiceUnavailable)
		return
	}
	if err := s.db.CheckModel(); err != nil {
		log.Printf("invalid model: %s", err)
		http.Error(w, fmt.Sprintf("invalid model: %s", err), http.StatusServiceUnavailable)
		return
	}
}
