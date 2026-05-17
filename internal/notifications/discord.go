package notifications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// SendDiscordAlert agora aceita 6 parâmetros para incluir Nome e Grupo
func SendDiscordAlert(webhookURL, name, ip, group, status string, latency time.Duration) {
	if webhookURL == "" {
		return
	}

	color := 15158332 // Vermelho para OFFLINE
	emoji := "🔴"
	if status == "ONLINE" {
		color = 3066993 // Verde para ONLINE
		emoji = "🟢"
	}

	payload := map[string]interface{}{
		"embeds": []map[string]interface{}{{
			"title":       fmt.Sprintf("%s Alerta de Conectividade", emoji),
			"description": fmt.Sprintf("O dispositivo **%s** mudou de estado.", name),
			"color":       color,
			"fields": []map[string]interface{}{
				{"name": "Status", "value": status, "inline": true},
				{"name": "IP", "value": fmt.Sprintf("`%s`", ip), "inline": true},
				{"name": "Grupo / Local", "value": group, "inline": true},
				{"name": "Latência", "value": latency.String(), "inline": true},
			},
			"footer": map[string]string{
				"text": "NetWatch Pro - Enterprise Monitoring",
			},
			"timestamp": time.Now().Format(time.RFC3339),
		}},
	}

	data, _ := json.Marshal(payload)
	client := &http.Client{Timeout: 5 * time.Second}
	_, _ = client.Post(webhookURL, "application/json", bytes.NewBuffer(data))
}