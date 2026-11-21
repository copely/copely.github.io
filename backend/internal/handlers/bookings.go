package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/mail"
	"strings"

	"copely-backend/internal/models"
	"copely-backend/internal/store"
)

// HandleBookingAvailability returns open dates and 30-min slots for ?month=YYYY-MM.
func (h *APIHandlers) HandleBookingAvailability(w http.ResponseWriter, r *http.Request) {
	month := strings.TrimSpace(r.URL.Query().Get("month"))
	if month == "" {
		month = "2026-11"
	}

	resp, err := h.Store.GetBookingAvailability(month)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// HandleCreateBooking validates and reserves an on-site consultation slot.
func (h *APIHandlers) HandleCreateBooking(w http.ResponseWriter, r *http.Request) {
	var payload models.BookingRequestPayload
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON payload: "+err.Error())
		return
	}

	payload.Date = strings.TrimSpace(payload.Date)
	payload.TimeSlot = strings.TrimSpace(payload.TimeSlot)
	payload.Town = strings.TrimSpace(payload.Town)
	payload.CustomerName = strings.TrimSpace(payload.CustomerName)
	payload.CustomerPhone = strings.TrimSpace(payload.CustomerPhone)
	payload.CustomerEmail = strings.TrimSpace(payload.CustomerEmail)
	payload.Notes = strings.TrimSpace(payload.Notes)

	if payload.Date == "" {
		writeError(w, http.StatusBadRequest, "date is required (YYYY-MM-DD)")
		return
	}
	if payload.TimeSlot == "" {
		writeError(w, http.StatusBadRequest, "timeSlot is required (HH:MM)")
		return
	}
	if payload.Town == "" {
		writeError(w, http.StatusBadRequest, "town is required")
		return
	}
	if payload.Service != "" && !models.IsValidService(payload.Service) {
		writeError(w, http.StatusBadRequest, "unsupported service category")
		return
	}
	if payload.Timeline != "" && !models.IsValidTimeline(payload.Timeline) {
		writeError(w, http.StatusBadRequest, "unsupported project timeline")
		return
	}
	if payload.CustomerEmail != "" {
		if _, err := mail.ParseAddress(payload.CustomerEmail); err != nil {
			writeError(w, http.StatusBadRequest, "invalid customerEmail address")
			return
		}
	}

	rec, err := h.Store.BookSlot(payload)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrSlotAlreadyBooked):
			writeError(w, http.StatusConflict, err.Error())
		case errors.Is(err, store.ErrInvalidBookingDate):
			writeError(w, http.StatusBadRequest, err.Error())
		default:
			writeError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusCreated, models.BookingResponse{
		Success:   true,
		BookingID: rec.ID,
		Date:      rec.Payload.Date,
		TimeSlot:  rec.Payload.TimeSlot,
		Town:      rec.Payload.Town,
		Service:   rec.Payload.Service,
		Timeline:  rec.Payload.Timeline,
		CreatedAt: rec.CreatedAt,
	})
}
