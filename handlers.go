package main

import (
	"encoding/json"
	"net/http"
)

type CreateRequest struct {
	LongURL string `json:"long_url"`
}

type CreateResponse struct {
	ShortURL string `json:"short_url"`
}

func createShortURL(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.LongURL == "" {
		http.Error(w, "url required", http.StatusBadRequest)
		return
	}

	var id int64

	err := db.QueryRow(
		"INSERT INTO urls (long_url) VALUES ($1) returning id",
		req.LongURL,
	).Scan(&id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	shortKey := encodeBase62(id)

	_, err = db.Exec(
		"UPDATE urls SET short_key=$1 WHERE id=$2",
		shortKey,
		id,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := CreateResponse{
		ShortURL: "http://localhost:8080/" + shortKey,
	}

	w.Header().Set("Content-Type", "applicaiton/json")
	json.NewEncoder(w).Encode(resp)
}

func redirect(w http.ResponseWriter, r *http.Request) {
	shortKey := r.URL.Path[1:]

	if shortKey == "" {
		http.NotFound(w, r)
		return
	}

	longURL, err := rdb.Get(ctx, shortKey).Result()

	if err == nil {
		http.Redirect(w, r, longURL, http.StatusFound)
	}

	err = db.QueryRow(
		"SELECT long_url FROM urls WHERE short_key=$1",
		shortKey,
	).Scan(&longURL)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	rdb.Set(ctx, shortKey, longURL, 0)

	go db.Exec(
		"UPDATE urls SET clicks = clicks + 1 WHERE short_key=$1",
		shortKey,
	)

	http.Redirect(w, r, longURL, http.StatusFound)
}
