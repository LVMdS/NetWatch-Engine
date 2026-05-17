package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/leonardo/netwatch/internal/domain"
	"github.com/leonardo/netwatch/internal/repositories"
)

type DeviceHandler struct {
	repo *repositories.DeviceRepository
}

func NewDeviceHandler(repo *repositories.DeviceRepository) *DeviceHandler {
	return &DeviceHandler{repo: repo}
}

func (h *DeviceHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dev domain.Device
	if err := json.NewDecoder(r.Body).Decode(&dev); err != nil {
		http.Error(w, "Payload de host inválido", http.StatusBadRequest)
		return
	}

	dev.Status = "OFFLINE"
	if err := h.repo.Create(&dev); err != nil {
		http.Error(w, "Erro ao salvar host", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *DeviceHandler) List(w http.ResponseWriter, r *http.Request) {
	devices, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, "Erro ao listar hosts", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(devices)
}

func (h *DeviceHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Dados de atualização incorretos", http.StatusBadRequest)
		return
	}

	delete(updates, "id")
	delete(updates, "created_at")

	if err := h.repo.Update(uint(id), updates); err != nil {
		http.Error(w, "Erro ao atualizar host", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *DeviceHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		http.Error(w, "Erro ao deletar host", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}