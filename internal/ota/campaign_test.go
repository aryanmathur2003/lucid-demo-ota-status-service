package ota

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSummarizeCountsAndLabels(t *testing.T) {
	s := Summarize(SeedCampaigns()[0])
	if s.Total != 6 || s.CompletePercent != 50 || s.Halted {
		t.Fatalf("unexpected summary: %+v", s)
	}
	if len(s.Phases) != 5 || s.Phases[3].Label != "Complete" || s.Phases[3].Count != 3 {
		t.Fatalf("unexpected phases: %+v", s.Phases)
	}
}

func TestSummarizeHaltsOnFailures(t *testing.T) {
	s := Summarize(SeedCampaigns()[1])
	if !s.Halted {
		t.Fatalf("campaign with 1 of 4 failed should be halted")
	}
	if len(s.FailedVINs) != 1 || s.FailedVINs[0] != "5LDEMAA0000000102" {
		t.Fatalf("unexpected failed vins: %v", s.FailedVINs)
	}
}

func TestSummarizeEmptyCampaign(t *testing.T) {
	s := Summarize(Campaign{ID: "empty"})
	if s.Total != 0 || s.CompletePercent != 0 || s.Halted {
		t.Fatalf("unexpected summary: %+v", s)
	}
}

func TestRoutes(t *testing.T) {
	h := NewServer(SeedCampaigns()).Routes()

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v1/campaigns", nil))
	var list []Summary
	if err := json.Unmarshal(rec.Body.Bytes(), &list); err != nil || len(list) != 2 {
		t.Fatalf("list: err=%v len=%d", err, len(list))
	}

	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v1/campaigns/cmp-2026-11-b", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("get: status = %d", rec.Code)
	}

	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v1/campaigns/missing", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("missing: status = %d", rec.Code)
	}
}
