package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// Получить всех юзеров
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

// Получить юзера по айдишнику
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

// Создаем Юзера
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//Проверяем что это именно метод пост
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method post", http.StatusMethodNotAllowed)
	}
	//Тип данных для декода
	var requestUser struct {
		FirstName  string `json:"first_name"`
		SecondName string `json:"second_name"`
		Email      string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestUser); err != nil {
		http.Error(w, "Invalid json decoder", http.StatusInternalServerError)
		return
	}

	id, err := h.dbPool.CreateNewUser(r.Context(), requestUser.FirstName, requestUser.SecondName, requestUser.Email)
	if err != nil {
		log.Printf("error creating user %v", err)
		http.Error(w, "Error creating user", http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(map[string]any{
		"id":      id,
		"message": "User created successfully",
	})
}
