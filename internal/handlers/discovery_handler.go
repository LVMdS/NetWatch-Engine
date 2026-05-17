package handlers

import (
	"encoding/json"
	"net"
	"net/http"

	"github.com/leonardo/netwatch/internal/domain"
	"github.com/leonardo/netwatch/internal/monitoring"
	"github.com/leonardo/netwatch/internal/repositories"
)

type DiscoveryHandler struct {
	devRepo *repositories.DeviceRepository
}

func NewDiscoveryHandler(devRepo *repositories.DeviceRepository) *DiscoveryHandler {
	return &DiscoveryHandler{devRepo: devRepo}
}

func (h *DiscoveryHandler) ScanNetwork(w http.ResponseWriter, r *http.Request) {
	var input struct {
		StartIP string `json:"start_ip"`
		EndIP   string `json:"end_ip"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	ips := generateIPRange(input.StartIP, input.EndIP)
	if len(ips) == 0 || len(ips) > 256 {
		http.Error(w, "Intervalo de IP inválido ou superior a 256 alvos", http.StatusBadRequest)
		return
	}

	results := monitoring.RunDiscovery(ips)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func (h *DiscoveryHandler) ImportHosts(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		GroupID *uint `json:"group_id"`
		Hosts   []struct {
			IP       string `json:"ip"`
			Hostname string `json:"hostname"`
			MAC      string `json:"mac"`
		} `json:"hosts"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Payload de importação fragmentado", http.StatusBadRequest)
		return
	}

	for _, item := range payload.Hosts {
		name := item.Hostname
		if name == "Host Desconhecido" {
			name = "Host Genérico (" + item.IP + ")"
		}

		desc := ""
		if item.MAC != "N/A" {
			desc = "MAC: " + item.MAC
		}

		dev := domain.Device{
			Name:        name,
			IPAddress:   item.IP,
			Description: desc,
			Status:      "ONLINE",
			GroupID:     payload.GroupID,
		}
		_ = h.devRepo.Create(&dev)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success"}`))
}

func generateIPRange(startStr, endStr string) []string {
	start := net.ParseIP(startStr)
	end := net.ParseIP(endStr)
	if start == nil || end == nil {
		return nil
	}

	start = start.To4()
	end = end.To4()
	if start == nil || end == nil {
		return nil
	}

	var ips []string
	for i := ip2Long(start); i <= ip2Long(end); i++ {
		ips = append(ips, long2IP(i).String())
	}
	return ips
}

func ip2Long(ip net.IP) uint32 {
	return uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
}

func long2IP(n uint32) net.IP {
	return net.IPv4(byte(n>>24), byte(n>>16), byte(n>>8), byte(n))
}