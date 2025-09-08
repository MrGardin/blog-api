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
	if requestUser.FirstName == "" || requestUser.SecondName == "" || requestUser.Email == "" {
		http.Error(w, "Error validate name and email", http.StatusBadRequest)
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

// Обновить юзера
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	userId, _ := strconv.Atoi(id)
	var requestUser struct {
		FirstName  string `json:"first_name"`
		SecondName string `json:"second_name"`
		Email      string `json:"email"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestUser)
	if err != nil {
		http.Error(w, "Invalid get update user", http.StatusNotFound)
		return
	}
	userId, err = h.dbPool.UpdateUser(r.Context(), userId, requestUser.FirstName, requestUser.SecondName, requestUser.Email)
	if err != nil {
		http.Error(w, "Invalid update user", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(map[string]any{
		"id":      userId,
		"message": "User updated successfully",
	})
}

//Удалить юзера

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}
	err = h.dbPool.DeleteUser(r.Context(), userId)
	if err != nil {
		http.Error(w, "Error with delete User", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent) // 204 No Content
}

// Создать пост
func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var request struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		UserID  int    `json:"user_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Error with Decode json", http.StatusBadRequest)
		return
	}
	if request.Title == "" || request.UserID == 0 || request.Content == "" {
		http.Error(w, "Title, content and user_id are required", http.StatusBadRequest)
		return
	}
	id, err := h.dbPool.CreateNewPost(r.Context(), request.Title, request.Content, request.UserID)
	if err != nil {
		log.Printf("Error creating post: %v", err)
		http.Error(w, "Invalid create post", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"id":      id,
		"message": "Post created successfully",
	})
}
func (h *Handler) GetAllPosts(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	posts, err := h.dbPool.GetAllPosts(r.Context())
	if err != nil {
		log.Printf("Error getting posts: %v", err)
		http.Error(w, "error get all posts", http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(posts)
}
func (h *Handler) GetPostById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	convId, _ := strconv.Atoi(id)
	post, err := h.dbPool.GetPostById(r.Context(), convId)
	if err != nil {
		http.Error(w, "error hanlder get post", http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(post)
}
func (h *Handler) GetUsersPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userIdStr := ps.ByName("id")
	userId, _ := strconv.Atoi(userIdStr)
	posts, err := h.dbPool.GetUsersPost(r.Context(), userId)
	if err != nil {
		http.Error(w, "error get users post", http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(posts)
}
