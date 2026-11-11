package ota

import (
	"encoding/json"
	"net/http"
	"strings"
)

// Server holds the HTTP handlers for the service.
type Server struct {
	campaigns map[string]Campaign
	order     []string
}

func NewServer(campaigns []Campaign) *Server {
	s := &Server{campaigns: make(map[string]Campaign, len(campaigns))}
	for _, c := range campaigns {
		s.campaigns[c.ID] = c
		s.order = append(s.order, c.ID)
	}
	return s
}

// Routes returns the HTTP handler for the service.
func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("/v1/campaigns", s.handleList)
	mux.HandleFunc("/v1/campaigns/", s.handleGet)
	return mux
}

func (s *Server) handleList(w http.ResponseWriter, _ *http.Request) {
	out := make([]Summary, 0, len(s.order))
	for _, id := range s.order {
		out = append(out, Summarize(s.campaigns[id]))
	}
	writeJSON(w, http.StatusOK, out)
}

func (s *Server) handleGet(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/v1/campaigns/")
	c, ok := s.campaigns[id]
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "campaign not found"})
		return
	}
	writeJSON(w, http.StatusOK, Summarize(c))
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
