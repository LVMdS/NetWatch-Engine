package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/leonardo/netwatch/internal/domain"
	"github.com/leonardo/netwatch/internal/repositories"
	"github.com/leonardo/netwatch/internal/services"
)

type AuthHandler struct {
	service *services.AuthService
	repo    *repositories.UserRepository
}

func NewAuthHandler(service *services.AuthService, repo *repositories.UserRepository) *AuthHandler {
	return &AuthHandler{service: service, repo: repo}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name           string `json:"name"`
		Email          string `json:"email"`
		Password       string `json:"password"`
		Company        string `json:"company"`
		DiscordWebhook string `json:"discord_webhook"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Estrutura de dados inválida", http.StatusBadRequest)
		return
	}

	user := domain.User{
		Name:           input.Name,
		Email:          input.Email,
		Password:       input.Password,
		Company:        input.Company,
		DiscordWebhook: input.DiscordWebhook,
	}

	if err := h.service.Register(&user); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"status":"success","message":"Conta criada com sucesso"}`))
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Payload incorreto", http.StatusBadRequest)
		return
	}

	token, err := h.service.Login(creds.Email, creds.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *AuthHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, "Erro ao listar usuários", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *AuthHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		http.Error(w, "Erro ao remover usuário", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}