package handlers

import (
	"encoding/json"
	"fmt"
	"go-back/service"
	"net/http"
	"strings"
)

type AuthHandler struct {
	authService *service.AuthService
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
	UserID  int    `json:"user_id"`
}

type ValidateTokenRequest struct {
	Token string `json:"token"`
}

type ValidateTokenResponse struct {
	Valid   bool   `json:"valid"`
	Message string `json:"message"`
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

	// fmt.Printf("Received login request: %+v\n", req)

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

	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(LoginResponse{
	// 	Message: "Authorization successful",
	// 	UserID:  user.ID,
	// })

	token, err := h.authService.GenerateJWT(user)

	response := LoginResponse{
		Token:   token,
		Message: "Authorization successful",
		UserID:  user.ID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header is required", http.StatusBadRequest)
		return
	}

	const bearerScheme = "Bearer "
	if !strings.HasPrefix(authHeader, bearerScheme) {
		http.Error(w, "Invalid Authorization header format. Expected: Bearer <token>", http.StatusBadRequest)
		return
	}

	tokenString := strings.TrimPrefix(authHeader, bearerScheme)
	if tokenString == "" {
		http.Error(w, "Token is required", http.StatusBadRequest)
		return
	}

	_, err := h.authService.ValidateJWT(tokenString)
	if err != nil {
		fmt.Printf("Token validation error: %v\n", err)
		response := ValidateTokenResponse{
			Valid:   false,
			Message: "Invalid or expired token",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := ValidateTokenResponse{
		Valid:    true,
		Message:  "Token is valid",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
