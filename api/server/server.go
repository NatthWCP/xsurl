package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	kitlog "github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"xsurl/api/shortening"
	"xsurl/api/shortenurl"
)

type Server struct {
	Shortening shortening.Service
	Logger     kitlog.Logger
	router     chi.Router
}

type MyResponseWriter struct {
	http.ResponseWriter
	buf *bytes.Buffer
}

func New(ss shortening.Service, logger kitlog.Logger) *Server {
	s := &Server{
		Shortening: ss,
		Logger:     logger,
	}

	r := chi.NewRouter()
	r.Use(accessControl)
	r.Route("/v0", func(r chi.Router) {
		h := shorteningHanderler{s.Shortening, s.Logger}
		r.Mount("/shortening", h.router())
	})

	r.Method("GET", "/metrics", promhttp.Handler())

	s.router = r

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case shortenurl.ErrUnknownURL:
		w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
