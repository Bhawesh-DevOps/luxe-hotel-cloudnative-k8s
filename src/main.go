package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//go:embed index.html
var indexHTML []byte

type Reservation struct {
	BookingID   string  `json:"booking_id"`
	GuestName   string  `json:"guest_name"`
	RoomType    string  `json:"room_type"`
	Nights      int     `json:"nights"`
	TotalAmount float64 `json:"total_amount"`
	Status      string  `json:"status"`
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(indexHTML)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "HEALTHY",
		"service": "luxe-hotel-booking-engine",
		"tier":    "production-5star",
	})
}

func reservationsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	reservations := []Reservation{
		{
			BookingID:   "LX-90821",
			GuestName:   "Alexander Wright",
			RoomType:    "Presidential Ocean Suite",
			Nights:      4,
			TotalAmount: 4800.00,
			Status:      "CONFIRMED_KAFKA_DISPATCHED",
		},
		{
			BookingID:   "LX-90822",
			GuestName:   "Sophia Laurent",
			RoomType:    "Royal Penthouse Villa",
			Nights:      3,
			TotalAmount: 3600.00,
			Status:      "PROCESSING_EVENT_STREAM",
		},
	}
	json.NewEncoder(w).Encode(reservations)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/healthz", healthHandler)
	mux.HandleFunc("/readyz", healthHandler)
	mux.HandleFunc("/api/v1/reservations", reservationsHandler)

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("👑 The Grand Bhawesh Palace Luxury Hotel Engine running on http://localhost:%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}
