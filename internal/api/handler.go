package api

import (
	_ "embed"
	"encoding/json"
	"net/http"
	"parser/internal/db"
	"strconv"
)

//go:embed index.html
var indexHTML []byte

type Handler struct {
	db *db.DB
}

func NewHandler(db *db.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) GetItems(w http.ResponseWriter, r *http.Request) {
	amount := r.URL.Query().Get("amount")
	floatAmount, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	topItems, err := h.db.GetTopItems(floatAmount)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(topItems)
}

func (h *Handler) GeneralPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	_, err := w.Write(indexHTML)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
