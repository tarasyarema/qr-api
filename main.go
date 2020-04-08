package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/rs/cors"
	qrcode "github.com/skip2/go-qrcode"
)

func encode(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Query().Get("size")
	size := 256

	if s != "" {
		tmp, err := strconv.Atoi(s)

		if err == nil && tmp >= 128 && tmp <= 2048 {
			size = tmp
		}
	}

	q := r.URL.Query().Get("quality")
	quality := qrcode.Medium

	if q != "" {
		tmp, err := strconv.Atoi(s)

		if err == nil {
			queryLevel := qrcode.RecoveryLevel(tmp)

			if queryLevel >= qrcode.Low && queryLevel <= qrcode.Highest {
				quality = queryLevel
			}
		}
	}

	d := r.URL.Query().Get("data")
	var png []byte

	png, err := qrcode.Encode(d, quality, size)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Write(png)

	log.Printf("encoded %d bytes with quality %d and size %d", len(png), quality, size)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/encode", encode)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS", "PUT", "PATCH"},
		AllowCredentials: true,
	})

	handler := c.Handler(mux)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("listening on addr %s", addr)

	if err := http.ListenAndServe(":8080", handler); err != nil {
		panic(err)
	}
}
