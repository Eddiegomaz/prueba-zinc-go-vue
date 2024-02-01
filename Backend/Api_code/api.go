package main

import (
	"Project/zincShare"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {

	r := chi.NewRouter()

	cors := cors.New(cors.Options{
		// AllowedOrigins:   []string{"https://example.com"},  // <-- Restrict allowed origins if needed
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum age of preflight requests
	})
	r.Use(cors.Handler)

	r.Use(middleware.Logger)
	r.Get("/query", func(w http.ResponseWriter, r *http.Request) {

		query := r.URL.Query().Get("query")
		fromStr := r.URL.Query().Get("from")
		sizeStr := r.URL.Query().Get("size")

		from, err := strconv.Atoi(fromStr)
		if err != nil {
			http.Error(w, "'desde' must be an integer number", http.StatusBadRequest)
			return
		}

		size, err := strconv.Atoi(sizeStr)
		if err != nil {
			http.Error(w, "'cantidad' must be an integer number", http.StatusBadRequest)
			return
		}

		queryResponse, err := zincShare.Query(query, from, size)
		if err != nil {
			// Maneja el error y responde con un código de estado 500 si es necesario
			http.Error(w, fmt.Sprintf("Error en la consulta: %s", err), http.StatusInternalServerError)
			return
		}

		jsonResponse, err := json.Marshal(queryResponse)
		// fmt.Printf("%T", jsonResponse)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error al convertir a JSON: %s", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	})
	http.ListenAndServe(":3000", r)
}
