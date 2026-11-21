package handlers

import (
	"encoding/json"
	"net/http"
	"net/mail"
	"strings"

	"copely-backend/internal/models"
	"copely-backend/internal/store"
)

// APIHandlers groups HTTP handlers backed by MemoryStore.
type APIHandlers struct {
	Store *store.MemoryStore
}

// NewAPIHandlers creates an APIHandlers instance.
func NewAPIHandlers(s *store.MemoryStore) *APIHandlers {
	return &APIHandlers{Store: s}
}

// HandleHealth responds with basic API service health telemetry.
func (h *APIHandlers) HandleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "copely-api",
	})
}

// HandleCreateEstimate validates and persists an estimate submission.
func (h *APIHandlers) HandleCreateEstimate(w http.ResponseWriter, r *http.Request) {
	var payload models.EstimateRequestPayload
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON payload: "+err.Error())
		return
	}

	payload.FirstName = strings.TrimSpace(payload.FirstName)
	payload.LastName = strings.TrimSpace(payload.LastName)
	payload.Phone = strings.TrimSpace(payload.Phone)
	payload.Email = strings.TrimSpace(payload.Email)
	payload.Address = strings.TrimSpace(payload.Address)
	payload.Town = strings.TrimSpace(payload.Town)
	payload.Notes = strings.TrimSpace(payload.Notes)

	if !models.IsValidService(payload.Service) {
		writeError(w, http.StatusBadRequest, "unsupported service category")
		return
	}
	if !models.IsValidTimeline(payload.Timeline) {
		writeError(w, http.StatusBadRequest, "unsupported project timeline")
		return
	}
	if payload.Email != "" {
		if _, err := mail.ParseAddress(payload.Email); err != nil {
			writeError(w, http.StatusBadRequest, "invalid email address")
			return
		}
	}

	// If contact details are submitted from /contact, ensure core contact info exists.
	if payload.FirstName != "" || payload.LastName != "" || payload.Address != "" {
		if payload.Phone == "" && payload.Email == "" {
			writeError(w, http.StatusBadRequest, "please provide a phone number or email address")
			return
		}
	}

	rec, err := h.Store.SaveEstimate(payload)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to record estimate request")
		return
	}

	writeJSON(w, http.StatusCreated, models.EstimateResponse{
		Success:        true,
		EstimateID:     rec.ID,
		ResponseWindow: rec.ResponseWindow,
		Service:        rec.Payload.Service,
		Timeline:       rec.Payload.Timeline,
		CreatedAt:      rec.CreatedAt,
	})
}
