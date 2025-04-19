package api

import (
	"encoding/json"
	"net/http"

	"github.com/smafjal/goredis-mini/internal/core"
)

type Handler struct {
	eng *core.Engine
}

func NewHandler(eng *core.Engine) *Handler {
	return &Handler{
		eng: eng,
	}
}

// GET /get?key=foo
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value, ok := h.eng.DB.Get(key)
	if !ok {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"key":   key,
		"value": value,
	})
}

// POST /set { "key": "foo", "value": "bar" }
func (h *Handler) Set(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	h.eng.DB.Set(body.Key, body.Value)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}
