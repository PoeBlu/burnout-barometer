// Copyright 2019 Lester James V. Miranda. All rights reserved.
// Licensed under the MIT License. See LICENSE in the project root
// for license information.

package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

// Server implements the Optserve server to be run inside the cluster.
type Server struct {
	Port   int
	Router *httprouter.Router
	Config *Configuration
}

// Routes contain all handler functions that responds to GET or POST requests.
func (s *Server) Routes() {
	log.Debug("serving routes")
	s.Router.HandlerFunc(http.MethodPost, "/log", s.handleLog())
	s.Router.HandlerFunc(http.MethodGet, "/", s.handleIndex())
}

// Start command starts a server on the specific port.
func (s *Server) Start() error {
	log.Infof("listening to port %d", s.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", s.Port), s.Router)
	return nil
}

func (s *Server) handleIndex() http.HandlerFunc {
	type response struct {
		Data string `json:"data"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{"path": "/"}).Trace("received request")
		res := response{Data: "PONG"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&res)
	}
}

func (s *Server) handleLog() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{"path": "/log"}).Trace("received request")

		// config, err := NewConfiguration(s.CfgPath)
		// if err != nil {
		// 	log.WithFields(log.Fields{"err": err}).Fatal("NewConfiguration")
		// }

		if err := r.ParseForm(); err != nil {
			http.Error(w, "couldn't parse form", 400)
			log.WithFields(log.Fields{"err": err}).Error("http.Request.ParseForm")
		}

		// Check if authentication token is correct
		if err := VerifyWebhook(r.Form, s.Config.Token); err != nil {
			http.Error(w, "webhook may be empty, missing, or unauthorized", 400)
			log.WithFields(log.Fields{"err": err}).Error("VerifyWebhook")
		}

		// Check if text exists
		if len(r.Form["text"]) == 0 {
			http.Error(w, "empty text in form", 400)
			log.Error("empty text in form")
		}

		// Identify the scheme
		var db Database
		switch scheme := s.getScheme(s.Config.Table); scheme {
		case "bigquery", "bq":
			db = &BigQuery{URL: s.Config.Table}
		default:
			msg := fmt.Sprintf("unknown database scheme: %s", scheme)
			http.Error(w, msg, 400)
			log.Fatal(msg)
		}

		// Process the request
		req := &Request{
			Text:      r.Form["text"][0],
			UserID:    r.Form["user_id"][0],
			Timestamp: r.Header.Get("X-Slack-Request-Timestamp"),
			Area:      s.Config.Area,
			DB:        db,
		}

		resp, err := req.Process()
		if err != nil {
			http.Error(w, "error in processing request", 400)
			log.WithFields(log.Fields{"err": err}).Error("Request.Process")
		}

		// Send reply back to Slack
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func (s *Server) getScheme(h string) string {
	u, err := url.Parse(h)
	if err != nil {
		log.Fatal("invalid database URL")
	}
	return u.Scheme
}

// VerifyWebhook checks if the submitted request matches the token provided by Slack
func VerifyWebhook(form url.Values, token string) error {
	t := form.Get("token")
	if len(t) == 0 {
		return fmt.Errorf("empty form token")
	}

	if t != token {
		return fmt.Errorf("invalid request/credentials: %q", t[0])
	}

	return nil
}
