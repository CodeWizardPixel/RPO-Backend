package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"go-back/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	// Token   string `json:"token"`
	Message string `json:"message"`
	UserID  int    `json:"user_id"`
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) GetToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	fmt.Printf("Received login request: %+v\n", req)

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Login == "" || req.Password == "" {
		http.Error(w, "Login and password are required", http.StatusBadRequest)
		return
	}

	user, err := h.authService.AuthenticateUser(req.Login, req.Password)
	if err != nil {
		http.Error(w, "Invalid login or password", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(LoginResponse{
		Message: "Authorization successful",
		UserID:  user.ID,
	})

	// token, err := h.authService.GenerateToken(user)
	// if err != nil {
	// 	http.Error(w, "Error generating token", http.StatusInternalServerError)
	// 	return
	// }

	// response := LoginResponse{
	// 	Token:   token,
	// 	Message: "Authorization successful",
	// 	UserID:  user.ID,
	// }

	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(response)
}

