package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	kitlog "github.com/go-kit/kit/log"

	"xsurl/api/shortening"
)

type shorteningHanderler struct {
	s      shortening.Service
	logger kitlog.Logger
}

func (h *shorteningHanderler) router() chi.Router {
	r := chi.NewRouter()

	r.Post("/", h.shortenURL)
	// r.Get("/", h.shortenURL)

	return r
}

func (h *shorteningHanderler) shortenURL(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var request struct {
		OriginURL string `json:"originURL"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Log("process", "decoding", "error", err, "request_msg", request.OriginURL)
		encodeError(ctx, err, w)
		return
	}

	surl, err := h.s.ShortenURL(request.OriginURL)
	if err != nil {
		encodeError(ctx, err, w)
		return
	}

	var response = struct {
		ShortenURL string `json:"shortenurl"`
	}{
		ShortenURL: surl,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Log("process", "encoding", "error", err, "request_msg", request)
		encodeError(ctx, err, w)
		return
	}
}
