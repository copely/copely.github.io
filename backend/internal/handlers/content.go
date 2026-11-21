package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"copely-backend/internal/store"
)

// HandleGetContent serves the complete seeded content bundle for the frontend.
func (h *APIHandlers) HandleGetContent(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.Store.GetContent())
}

// HandleGetServices returns all Copely services.
func (h *APIHandlers) HandleGetServices(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.Store.GetContent().Services)
}

// HandleGetServiceBySlug returns a single service detail by slug.
func (h *APIHandlers) HandleGetServiceBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	svc, err := h.Store.GetServiceBySlug(slug)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "service not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to load service")
		return
	}
	writeJSON(w, http.StatusOK, svc)
}

// HandleGetReviews returns verified customer reviews and featured testimonials.
func (h *APIHandlers) HandleGetReviews(w http.ResponseWriter, r *http.Request) {
	bundle := h.Store.GetContent()
	writeJSON(w, http.StatusOK, map[string]any{
		"testimonials":  bundle.Testimonials,
		"googleReviews": bundle.GoogleReviews,
	})
}

// HandleGetFaqs returns categorized FAQs and glossary definitions.
func (h *APIHandlers) HandleGetFaqs(w http.ResponseWriter, r *http.Request) {
	bundle := h.Store.GetContent()
	writeJSON(w, http.StatusOK, map[string]any{
		"homeFaqs":      bundle.HomeFaqs,
		"faqCategories": bundle.FaqCategories,
		"glossaryTerms": bundle.GlossaryTerms,
	})
}

// HandleGetPosts returns all editorial blog posts.
func (h *APIHandlers) HandleGetPosts(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.Store.GetContent().BlogPosts)
}

// HandleGetPostBySlug returns a single blog post by slug.
func (h *APIHandlers) HandleGetPostBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	post, err := h.Store.GetBlogPostBySlug(slug)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "blog post not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to load blog post")
		return
	}
	writeJSON(w, http.StatusOK, post)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]any{
		"success": false,
		"error":   message,
	})
}
