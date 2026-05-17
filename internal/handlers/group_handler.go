package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/leonardo/netwatch/internal/domain"
	"github.com/leonardo/netwatch/internal/repositories"
)

type GroupHandler struct {
	repo *repositories.GroupRepository
}

func NewGroupHandler(repo *repositories.GroupRepository) *GroupHandler {
	return &GroupHandler{repo: repo}
}

func (h *GroupHandler) Create(w http.ResponseWriter, r *http.Request) {
	var group domain.Group
	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		http.Error(w, "Dados de grupo inválidos", http.StatusBadRequest)
		return
	}

	if err := h.repo.Create(&group); err != nil {
		http.Error(w, "Erro ao criar grupo", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *GroupHandler) List(w http.ResponseWriter, r *http.Request) {
	groups, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, "Erro ao listar grupos", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groups)
}

func (h *GroupHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "ID de grupo inválido", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		http.Error(w, "Erro ao deletar grupo", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}