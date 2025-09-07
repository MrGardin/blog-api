package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	users, err := h.dbPool.GetAllUsers(r.Context())
	//log error
	if err != nil {
		log.Printf("error getting users: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(users)

	if err != nil {
		log.Printf("Error encoding JSON: %v", err)
		http.Error(w, "Error creating response", http.StatusInternalServerError)
	}
}
func (h *Handler) GetUserById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	userID, _ := strconv.Atoi(id)
	user, err := h.dbPool.GetUserById(r.Context(), userID)
	if err != nil {
		http.Error(w, "Invalid get User", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
}
