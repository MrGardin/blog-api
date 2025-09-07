package handlers

import (
	"blog-api/storage"
	"context"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Handler struct {
	dbPool *storage.Pool
}

func NewHandler(pool *storage.Pool) *Handler {
	return &Handler{dbPool: pool}
}
func (h *Handler) CheckVersionDataBase(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	version, err := h.dbPool.GetVersion(context.Background())
	if err != nil {
		http.Error(w, "Database unavailable", http.StatusInternalServerError)
		log.Printf("Database version error: %v", err) // Просто логируем
		return
	}
	w.Write([]byte(version))
}
