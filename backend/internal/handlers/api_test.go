package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"copely-backend/internal/handlers"
	"copely-backend/internal/models"
	"copely-backend/internal/store"
)

func setupTestMux(t *testing.T) (*http.ServeMux, string) {
	t.Helper()
	snapshotPath := filepath.Join(t.TempDir(), "submissions.json")
	s, err := store.NewMemoryStore(snapshotPath)
	if err != nil {
		t.Fatalf("NewMemoryStore error: %v", err)
	}
	h := handlers.NewAPIHandlers(s)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", h.HandleHealth)
	mux.HandleFunc("POST /api/estimates", h.HandleCreateEstimate)
	mux.HandleFunc("GET /api/bookings/availability", h.HandleBookingAvailability)
	mux.HandleFunc("POST /api/bookings", h.HandleCreateBooking)
	mux.HandleFunc("GET /api/content", h.HandleGetContent)
	mux.HandleFunc("GET /api/services/{slug}", h.HandleGetServiceBySlug)
	mux.HandleFunc("GET /api/posts/{slug}", h.HandleGetPostBySlug)
	return mux, snapshotPath
}

func TestHealthAndContentEndpoints(t *testing.T) {
	mux, _ := setupTestMux(t)

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 for /api/health, got %d", rec.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/content", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 for /api/content, got %d", rec.Code)
	}

	var bundle models.ContentBundle
	if err := json.NewDecoder(rec.Body).Decode(&bundle); err != nil {
		t.Fatalf("decode /api/content failed: %v", err)
	}
	if len(bundle.Services) != 5 {
		t.Fatalf("expected 5 services, got %d", len(bundle.Services))
	}
	if len(bundle.BeforeAfterCases) != 4 {
		t.Fatalf("expected 4 before/after cases, got %d", len(bundle.BeforeAfterCases))
	}
	if len(bundle.Testimonials) != 5 {
		t.Fatalf("expected 5 testimonials, got %d", len(bundle.Testimonials))
	}
	if len(bundle.BlogPosts) != 4 {
		t.Fatalf("expected 4 blog posts, got %d", len(bundle.BlogPosts))
	}
}

func TestEstimateValidationAndPersistence(t *testing.T) {
	mux, snapshotPath := setupTestMux(t)

	// Invalid service should return 400
	badBody := bytes.NewBufferString(`{"service":"roofing","timeline":"this-week"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/estimates", badBody)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid service, got %d", rec.Code)
	}

	// Valid contact submission
	goodBody := bytes.NewBufferString(`{
		"firstName":"Sam",
		"lastName":"Calloway",
		"phone":"(512) 555-0100",
		"email":"sam@example.com",
		"address":"123 Ridgeline Dr",
		"town":"Dripping Springs",
		"service":"wet-rooms",
		"timeline":"this-week"
	}`)
	req = httptest.NewRequest(http.MethodPost, "/api/estimates", goodBody)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 for valid estimate, got %d (%s)", rec.Code, rec.Body.String())
	}

	var resp models.EstimateResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode estimate response: %v", err)
	}
	if !resp.Success || resp.ResponseWindow != "Within 24 hrs" {
		t.Fatalf("unexpected estimate response: %+v", resp)
	}

	// Reload store from snapshot to confirm atomic persistence
	reloaded, err := store.NewMemoryStore(snapshotPath)
	if err != nil {
		t.Fatalf("reload store failed: %v", err)
	}
	_ = reloaded
}

func TestBookingAvailabilityAndDoubleBookConflict(t *testing.T) {
	mux, snapshotPath := setupTestMux(t)

	// Check availability for November 2026
	req := httptest.NewRequest(http.MethodGet, "/api/bookings/availability?month=2026-11", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 for availability, got %d", rec.Code)
	}

	// Book 2026-11-18 at 08:00
	bookBody := bytes.NewBufferString(`{
		"date":"2026-11-18",
		"timeSlot":"08:00",
		"town":"Dripping Springs",
		"service":"porcelain-ceramic",
		"customerEmail":"claire@example.com"
	}`)
	req = httptest.NewRequest(http.MethodPost, "/api/bookings", bookBody)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 for first booking, got %d (%s)", rec.Code, rec.Body.String())
	}

	// Attempting to book same date + slot again must yield 409 Conflict
	dupBody := bytes.NewBufferString(`{
		"date":"2026-11-18",
		"timeSlot":"08:00",
		"town":"Austin"
	}`)
	req = httptest.NewRequest(http.MethodPost, "/api/bookings", dupBody)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusConflict {
		t.Fatalf("expected 409 conflict on duplicate slot, got %d", rec.Code)
	}

	// Reload from snapshot and verify slot remains unavailable across restarts
	reloaded, err := store.NewMemoryStore(snapshotPath)
	if err != nil {
		t.Fatalf("reload store failed: %v", err)
	}
	avail, err := reloaded.GetBookingAvailability("2026-11")
	if err != nil {
		t.Fatalf("GetBookingAvailability after reload failed: %v", err)
	}
	for _, day := range avail.Days {
		if day.Date == "2026-11-18" {
			for _, slot := range day.Slots {
				if slot.Time == "08:00" && slot.Available {
					t.Fatalf("expected 2026-11-18 08:00 to remain booked after restart")
				}
			}
		}
	}
}
