package handlers

import (
	"encoding/json"
	"go-back/service"
	"net/http"
	"strconv"
	"strings"
)

type TerminalHandler struct {
	terminalService *service.TerminalService
}

type TerminalRequest struct {
	ID           int    `json:"id"`
	SerialNumber string `json:"serial_number"`
	Address      string `json:"address"`
	Name         string `json:"name"`
}

func NewTerminalHandler(terminalService *service.TerminalService) *TerminalHandler {
	return &TerminalHandler{
		terminalService: terminalService,
	}
}

func (h *TerminalHandler) GetAllTerminals(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	terminals, err := h.terminalService.GetAllTerminals()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(terminals)
}

func (h *TerminalHandler) GetTerminalByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idRaw := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idRaw)
	if err != nil {
		http.Error(w, "Invalid terminal ID", http.StatusBadRequest)
		return
	}

	terminal, err := h.terminalService.GetTerminalByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(terminal)
}

func (h *TerminalHandler) CreateTerminal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tokenString, err := extractBearerToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var req TerminalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.terminalService.CreateTerminal(tokenString, req.SerialNumber, req.Address, req.Name)
	if err != nil {
		h.writeTerminalError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Terminal created"})
}

func (h *TerminalHandler) UpdateTerminal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tokenString, err := extractBearerToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var req TerminalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.terminalService.UpdateTerminal(tokenString, req.ID, req.SerialNumber, req.Address, req.Name)
	if err != nil {
		h.writeTerminalError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Terminal updated"})
}

func (h *TerminalHandler) DeleteTerminal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tokenString, err := extractBearerToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	idRaw := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idRaw)
	if err != nil {
		http.Error(w, "Invalid terminal ID", http.StatusBadRequest)
		return
	}

	err = h.terminalService.DeleteTerminal(tokenString, id)
	if err != nil {
		h.writeTerminalError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Terminal deleted"})
}

func (h *TerminalHandler) writeTerminalError(w http.ResponseWriter, err error) {
	errText := err.Error()

	switch {
	case strings.Contains(errText, "forbidden"):
		http.Error(w, errText, http.StatusForbidden)
	case strings.Contains(errText, "token"):
		http.Error(w, errText, http.StatusUnauthorized)
	default:
		http.Error(w, errText, http.StatusBadRequest)
	}
}