package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"copely-backend/internal/handlers"
	"copely-backend/internal/store"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	snapshotPath := os.Getenv("COPELY_SNAPSHOT_PATH")
	if snapshotPath == "" {
		snapshotPath = filepath.Join("data", "submissions.json")
	}

	memStore, err := store.NewMemoryStore(snapshotPath)
	if err != nil {
		log.Fatalf("failed to initialize Copely store: %v", err)
	}

	h := handlers.NewAPIHandlers(memStore)
	mux := NewRouter(h)

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           withMiddleware(mux),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("Copely Go API listening on http://localhost:%s (snapshot: %s)", port, snapshotPath)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}

// NewRouter registers all Copely REST endpoints on a standard library mux.
func NewRouter(h *handlers.APIHandlers) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", h.HandleHealth)
	mux.HandleFunc("POST /api/estimates", h.HandleCreateEstimate)
	mux.HandleFunc("GET /api/bookings/availability", h.HandleBookingAvailability)
	mux.HandleFunc("GET /api/bookings/slots", h.HandleBookingAvailability)
	mux.HandleFunc("POST /api/bookings", h.HandleCreateBooking)
	mux.HandleFunc("GET /api/content", h.HandleGetContent)
	mux.HandleFunc("GET /api/services", h.HandleGetServices)
	mux.HandleFunc("GET /api/services/{slug}", h.HandleGetServiceBySlug)
	mux.HandleFunc("GET /api/reviews", h.HandleGetReviews)
	mux.HandleFunc("GET /api/faqs", h.HandleGetFaqs)
	mux.HandleFunc("GET /api/posts", h.HandleGetPosts)
	mux.HandleFunc("GET /api/posts/{slug}", h.HandleGetPostBySlug)
	return mux
}

func withMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "http://localhost:5173" || origin == "http://127.0.0.1:5173" || origin == "" {
			if origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			} else {
				w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
			}
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
